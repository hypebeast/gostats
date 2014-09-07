package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	MAX_TRIES = 10
)

// GithubRepo represents a Github repository.
type GithubRepo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	Date        int64  `json:"date"`
}

type RepoStarsSerie struct {
	Name string       `json:"name"`
	Data []StarsPoint `json:"data"`
}

type StarsPoint struct {
	Stars int   `json:"stars"`
	Date  int64 `json:"date"`
}

var trendingQuery = "SELECT title, description, url, stars, date FROM github.trending WHERE since='%s' AND date > '%d' ORDER BY stars DESC"
var mostStarredQuery = "SELECT title, description, url, stars, date FROM github.stars WHERE date > '%d' ORDER BY stars DESC"
var repoHistory = "SELECT stars, date FROM github.stars WHERE title='%s' AND date > '%d' ORDER BY date DESC"

func UpdateGithubStats() {
	fetchTrending()
	fetchMostStarred()
	fetchStarsHistory()
}

func fetchStarsHistory() {
	// Get data for one year
	since := time.Now().Add(time.Hour * 24 * 365 * -1).Unix()

	database.StarsSeries = nil
	for _, repo := range database.MostStarred {
		if repo.Title == "" {
			continue
		}

		query := fmt.Sprintf(repoHistory, repo.Title, since)
		data, err := runQuery(query)
		if err != nil {
			continue
		}

		jsonData := []struct {
			Stars string `json:"stars"`
			Date  string `json:"date"`
		}{}

		if ok := json.Unmarshal(data.Bytes(), &jsonData); ok != nil {
			log.Printf("ERROR - %s", ok)
			continue
		}

		serie := RepoStarsSerie{
			Name: repo.Title,
		}

		// Map the data
		for _, x := range jsonData {
			stars, _ := strconv.Atoi(x.Stars)
			timestamp, _ := time.Parse("2006-01-02 15:04:05", x.Date)
			point := StarsPoint{
				Stars: stars,
				Date:  timestamp.Unix(),
			}
			serie.Data = append(serie.Data, point)
		}

		database.StarsSeries = append(database.StarsSeries, serie)
	}
}

func fetchTrending() {
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	periods := [3]string{"daily", "weekly", "monthly"}
	stores := [3]*[]GithubRepo{&database.DailyTrending, &database.WeeklyTrending, &database.MonthlyTrending}
	data := new(bytes.Buffer)

	var err error
	for i, period := range periods {
		query := fmt.Sprintf(trendingQuery, period, date.Unix())
		data, err = runQuery(query)
		if err != nil {
			continue
		}

		tries := 1
		for data.Len() == 0 && tries <= MAX_TRIES {
			date = date.Add(time.Hour * 24 * -1)
			query := fmt.Sprintf(trendingQuery, period, date.Unix())
			data, _ = runQuery(query)

			if data.Len() > 0 {
				break
			}
			tries += 1
		}

		if tries >= MAX_TRIES && data.Len() == 0 {
			return
		}
		decodeRepoData(data, stores[i])
	}
}

func fetchMostStarred() {
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	data := new(bytes.Buffer)
	query := fmt.Sprintf(mostStarredQuery, date.Unix())
	data, err := runQuery(query)
	if err != nil {
		return
	}

	tries := 1
	for data.Len() == 0 && tries <= MAX_TRIES {
		date = date.Add(time.Hour * 24 * -1)
		query := fmt.Sprintf(mostStarredQuery, date.Unix())
		data, _ = runQuery(query)

		if data.Len() > 0 {
			break
		}
		tries += 1
	}

	if tries >= MAX_TRIES && data.Len() == 0 {
		return
	}
	decodeRepoData(data, &database.MostStarred)
}

func decodeRepoData(data *bytes.Buffer, store *[]GithubRepo) {
	jsonData := []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"`
		Stars       string `json:"stars"`
		Date        string `json:"date"`
	}{}

	if ok := json.Unmarshal(data.Bytes(), &jsonData); ok != nil {
		log.Printf("ERROR - %s", ok)
		return
	}

	*store = nil
	for _, x := range jsonData {
		stars, _ := strconv.Atoi(x.Stars)
		timestamp, _ := time.Parse("2006-01-02 15:04:05", x.Date)
		*store = append(*store, GithubRepo{Title: x.Title, Description: x.Description, Url: x.Url, Stars: stars, Date: timestamp.Unix()})
	}
}

func DailyTrendingRepos() []GithubRepo {
	return database.DailyTrending
}

func WeeklyTrendingRepos() []GithubRepo {
	return database.WeeklyTrending
}

func MonthlyTrendingRepos() []GithubRepo {
	return database.MonthlyTrending
}

func MostStarredRepos() []GithubRepo {
	return database.MostStarred
}

func RepoStarsHistory(name string) RepoStarsSerie {
	for _, serie := range database.StarsSeries {
		if serie.Name == name {
			return serie
		}
	}
	return RepoStarsSerie{}
}
