package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

type GoDocPackage struct {
	Path string
}

type GoDocPackages struct {
	Results []GoDocPackage
}

type GitHubTrendingRepo struct {
	Name  string
	Url   string
	Stars int
	Date  time.Time
}

type GitHubTotalStarsRepo struct {
	Name  string
	Url   string
	Stars int
}

func init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
}

func scrapeGoDocPackages() {
	response, err := http.Get("http://api.godoc.org/packages")
	if err != nil {
		Error.Println(err)
		return
	}

	var packages GoDocPackages
	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()
	err = decoder.Decode(&packages)
	if err != nil {
		Error.Println(err)
		return
	}

	filename := dateFilename("godoc_packages", ".json")
	err = createFile(filename)
	if err != nil {
		Error.Println(err)
		return
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()
	if err != nil {
		Error.Println(err)
		return
	}

	data := struct {
		Count int `json:"count"`
		Date  int `json:"date"`
	}{
		Count: len(packages.Results),
		Date:  int(time.Now().Unix()),
	}
	outData, err := json.Marshal(data)
	if err != nil {
		Error.Println(err)
	}

	_, err = f.Write(outData)
	if err != nil {
		Error.Println(err)
		return
	}
}

func dateFilename(prefix string, extension string) string {
	dateString := time.Now().Format("2006-01-02-15")
	return prefix + "-" + dateString + extension
}

func createFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}

func scrapeGithubTrendingRepos() {
	// Get trending Go repositories from Gtihub
	// Get it for today, week and month
	// Get once per day
	// Data: Name, Link, Stars and total stars
	// Save all data in a JSON file on a daily base
}

func main() {
	Info.Println("Getting packages from GoDoc...")
	scrapeGoDocPackages()

	Info.Println("Done")
}
