package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

type Godoc struct {
	Count int   `json:"count"`
	Date  int64 `json:"date"`
}

var godocQuery = "SELECT count, date FROM godoc.packages WHERE date > '%d' ORDER BY date ASC"

func UpdateGodocStats() {
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, time.UTC).Add(time.Hour * 24 * 365 * -1).Unix()
	query := fmt.Sprintf(godocQuery, date)
	data, _ := runQuery(query)

	jsonData := []struct {
		Count string `json:"count"`
		Date  string `json:"date"`
	}{}

	if ok := json.Unmarshal(data.Bytes(), &jsonData); ok != nil {
		log.Printf("ERROR - %s", ok)
		return
	}

	database.GodocRepos = nil
	for _, x := range jsonData {
		count, _ := strconv.Atoi(x.Count)
		timestamp, _ := time.Parse("2006-01-02 15:04:05", x.Date)
		database.GodocRepos = append(database.GodocRepos, Godoc{Count: count, Date: timestamp.Unix()})
	}
}

func GetGodocStats() []Godoc {
	return database.GodocRepos
}
