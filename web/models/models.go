package models

import (
	"bytes"
	"os/exec"
)

var database Database

type Database struct {
	GodocRepos      []Godoc
	DailyTrending   []GithubRepo
	WeeklyTrending  []GithubRepo
	MonthlyTrending []GithubRepo
	MostStarred     []GithubRepo
}

// runQuery executes bq with the given query and returns the result.
func runQuery(query string) (buf *bytes.Buffer, err error) {
	out, err := exec.Command("bq", "-q", "--format=prettyjson", "query", query).Output()
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(out), nil
}

// Update updates data from BigQuery and writes it to the local database.
func Update() {
	UpdateGithubStats()
	UpdateGodocStats()
}
