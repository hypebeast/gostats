package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

type GithubRepo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	Date        int64  `json:"date"`
}

var trendingQuery = "SELECT title, description, url, stars, date FROM github.trending WHERE since='%s' AND date > '%d' ORDER BY stars DESC"
var mostStarredQuery = "SELECT title, description, url, stars, date FROM github.stars WHERE date > '%d' ORDER BY stars DESC"

func UpdateGithubStats() {
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).UTC().Unix()

	times := [3]string{"daily", "weekly", "monthly"}
	stores := [3]*[]GithubRepo{&database.DailyTrending, &database.WeeklyTrending, &database.MonthlyTrending}
	for i, since := range times {
		query := fmt.Sprintf(trendingQuery, since, date)
		mapGithubData(query, stores[i])
	}

	query := fmt.Sprintf(mostStarredQuery, date)
	mapGithubData(query, &database.MostStarred)
}

func mapGithubData(query string, store *[]GithubRepo) {
	jsonData := []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"`
		Stars       string `json:"stars"`
		Date        string `json:"date"`
	}{}

	data, _ := runQuery(query)

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
