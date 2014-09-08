package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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
	fetchMostStarred()
	fetchStarsHistory()
}

func UpdateGithubTrendingRepos() {
	scrapeTrendingRepos("go")
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

func scrapeTrendingRepos(language string) {
	var doc *goquery.Document
	var err error

	periods := [3]string{"daily", "weekly", "monthly"}
	stores := [3]*[]GithubRepo{&database.DailyTrending, &database.WeeklyTrending, &database.MonthlyTrending}
	for i, period := range periods {
		if doc, err = goquery.NewDocument(fmt.Sprintf("https://github.com/trending?l=%s&since=%s", language, period)); err != nil {
			log.Println(err)
		}

		*stores[i] = parseTrendingRepos(doc)
	}
}

func parseTrendingRepos(doc *goquery.Document) []GithubRepo {
	var repos []GithubRepo
	var regStars = regexp.MustCompile("[0-9]+")

	doc.Find("li.repo-list-item").Each(func(i int, s *goquery.Selection) {
		title := strings.Trim(s.Find("h3.repo-list-name a").Text(), "\n\t ")
		title = strings.Replace(title, " ", "", -1)
		title = strings.Replace(title, "\n", "", -1)
		description := strings.Trim(s.Find("p.repo-list-description").Text(), "\n\t ")
		url, _ := s.Find("h3.repo-list-name a").Attr("href")
		url = "https://github.com" + url
		starsString := s.Find("p.repo-list-meta").Text()
		starsString = strings.Replace(starsString, ",", "", -1)
		starsString = regStars.FindString(starsString)
		if starsString == "" {
			starsString = "0"
		}
		stars, _ := strconv.Atoi(starsString)

		repo := GithubRepo{
			Title:       title,
			Description: description,
			Url:         url,
			Stars:       stars,
			Forks:       0,
			Date:        time.Now().UTC().Unix(),
		}

		repos = append(repos, repo)
	})

	return repos
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
