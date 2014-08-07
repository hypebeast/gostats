package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	Info    *log.Logger
	Error   *log.Logger
	Periods []string
)

type GoDocPackage struct {
	Path string
}

type GoDocPackages struct {
	Results []GoDocPackage
}

type GitHubTrendingRepo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Stars       string `json:"stars"`
	Date        int    `json:"date"`
	Since       string `json:"since"`
}

type GitHubTotalStarsRepo struct {
	Name  string
	Url   string
	Stars int
}

func init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
	Periods = append(Periods, "daily", "weekly", "monthly")
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
		return
	}

	_, err = f.Write(outData)
	if err != nil {
		Error.Println(err)
		return
	}
}

func scrapeGithubTrendingRepos(language string) {
	var doc *goquery.Document
	var err error

	filename := dateFilename("github_trending_repos", ".json")
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

	for _, period := range Periods {
		if doc, err = goquery.NewDocument(fmt.Sprintf("https://github.com/trending?l=%s&since=%s", language, period)); err != nil {
			Error.Println(err)
		}

		repos := getRepos(doc, period)
		err = writeRepos(f, repos)
		if err != nil {
			Error.Println(err)
			return
		}
	}
}

func getRepos(doc *goquery.Document, since string) []GitHubTrendingRepo {
	var repos []GitHubTrendingRepo

	doc.Find("li.repo-leaderboard-list-item").Each(func(i int, s *goquery.Selection) {
		title := s.Find("div h2 a").Text()
		description := s.Find(".repo-leaderboard-description").Text()
		url, _ := s.Find("div h2 a").Attr("href")
		url = "https://github.com" + url
		stars := s.Find("span.collection-stat").First().Text()
		stars = strings.Replace(stars, ",", "", -1)

		repo := GitHubTrendingRepo{
			Title:       title,
			Description: description,
			Url:         url,
			Stars:       stars,
			Date:        int(time.Now().Unix()),
			Since:       since,
		}

		repos = append(repos, repo)
	})
	return repos
}

func writeRepos(file *os.File, repos []GitHubTrendingRepo) error {
	for _, repo := range repos {
		outData, err := json.Marshal(repo)
		if err != nil {
			return err
		}
		outData = append(outData, '\n')

		_, err = file.Write(outData)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	Info.Println("Getting packages from GoDoc...")
	scrapeGoDocPackages()

	Info.Println("Getting trending repos from GtiHub...")
	scrapeGithubTrendingRepos("go")

	Info.Println("Done")
}
