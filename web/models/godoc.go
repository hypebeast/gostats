package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Godoc struct {
	Count int
	Date  int
}

var godocQuery = "SELECT count, date FROM godoc.packages WHERE date > '%d'"

func UpdateGodocStats() {
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, time.UTC).Add(time.Hour * 24 * 365 * -1).Unix()
	query := fmt.Sprintf(godocQuery, date)
	data, _ := runQuery(query)

	if ok := json.Unmarshal(data.Bytes(), &database.GodocRepos); ok != nil {

	}
}

func GetGodocStats() []Godoc {
	return database.GodocRepos
}
