package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"path"
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

type GithubRepo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Stars       string `json:"stars"`
	Date        int    `json:"date"`
	Since       string `json:"since"`
}

func init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
	Periods = append(Periods, "daily", "weekly", "monthly")
}

func dateFilename(prefix string, extension string) string {
	dateString := time.Now().Format("2006-01-02")
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

func scrapeGoDocPackages(outDir string) {
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
	if outDir != "" {
		filename = path.Join(outDir, filename)
	}

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

func scrapeTrendingRepos(language string, outDir string) {
	var doc *goquery.Document
	var err error

	filename := dateFilename("github_trending_repos", ".json")
	if outDir != "" {
		filename = path.Join(outDir, filename)
	}

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

func getRepos(doc *goquery.Document, since string) []GithubRepo {
	var repos []GithubRepo

	doc.Find("li.repo-leaderboard-list-item").Each(func(i int, s *goquery.Selection) {
		title := s.Find("div h2 a").Text()
		description := s.Find(".repo-leaderboard-description").Text()
		url, _ := s.Find("div h2 a").Attr("href")
		url = "https://github.com" + url
		stars := s.Find("span.collection-stat").First().Text()
		stars = strings.Replace(stars, ",", "", -1)

		repo := GithubRepo{
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

func writeRepos(file *os.File, repos []GithubRepo) error {
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
	out := flag.String("out", "", "the output directory where the data is written to")
	flag.Parse()

	if *out != "" {
		if _, err := os.Stat(*out); err != nil {
			if os.IsNotExist(err) {
				if err = os.Mkdir(*out, 0755); err != nil {
					Error.Panicln(err)
				}
			}
		}
	}

	Info.Println("Getting packages from GoDoc...")
	scrapeGoDocPackages(*out)

	Info.Println("Getting trending repos from GtiHub...")
	scrapeTrendingRepos("go", *out)

	Info.Println("Done")
}
