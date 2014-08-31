package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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
	Forks       string `json:"forks"`
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
		Date:  int(time.Now().UTC().Unix()),
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

		repos := readTrendingRepos(doc, period)
		err = writeRepos(f, repos)
		if err != nil {
			Error.Println(err)
			return
		}
	}
}

func scrapeMostStarredRepos(language string, outDir string) {
	var doc *goquery.Document
	var err error

	filename := dateFilename("github_most_starred", ".json")
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

	for i := 1; i <= 5; i++ {
		if doc, err = goquery.NewDocument(fmt.Sprintf("https://github.com/search?q=stars:>1&type=Repositories&l=%s&p=%d", language, i)); err != nil {
			Error.Println(err)
		}
		repos := readMostStarredRepos(doc)
		err = writeRepos(f, repos)
		if err != nil {
			Error.Println(err)
			return
		}
	}
}

func readMostStarredRepos(doc *goquery.Document) []GithubRepo {
	var repos []GithubRepo

	doc.Find("li.repo-list-item").Each(func(i int, s *goquery.Selection) {
		title := strings.Trim(s.Find("h3.repo-list-name a").Text(), "\n ")
		description := strings.Trim(s.Find("p.repo-list-description").Text(), "\n ")
		url, _ := s.Find("h3.repo-list-name a").Attr("href")
		url = "https://github.com" + url
		stars := strings.Trim(s.Find("a.repo-list-stat-item[aria-label=\"Stargazers\"]").Text(), "\n\t ")
		stars = strings.Replace(stars, " ", "", -1)
		stars = strings.Replace(stars, "\n", "", -1)
		stars = strings.Replace(stars, ",", "", -1)
		if stars == "" {
			stars = "0"
		}
		forks := strings.Trim(s.Find("a.repo-list-stat-item[aria-label=\"Forks\"]").Text(), "\n ")
		forks = strings.Replace(forks, " ", "", -1)
		forks = strings.Replace(forks, "\n", "", -1)
		forks = strings.Replace(forks, ",", "", -1)
		if forks == "" {
			forks = "0"
		}

		repo := GithubRepo{
			Title:       title,
			Description: description,
			Url:         url,
			Stars:       stars,
			Forks:       forks,
			Date:        int(time.Now().UTC().Unix()),
		}

		repos = append(repos, repo)
	})
	return repos
}

func readTrendingRepos(doc *goquery.Document, since string) []GithubRepo {
	var repos []GithubRepo
	var regStars = regexp.MustCompile("[0-9]+")

	doc.Find("li.repo-list-item").Each(func(i int, s *goquery.Selection) {
		title := strings.Trim(s.Find("h3.repo-list-name a").Text(), "\n\t ")
		title = strings.Replace(title, " ", "", -1)
		title = strings.Replace(title, "\n", "", -1)
		description := strings.Trim(s.Find("p.repo-list-description").Text(), "\n\t ")
		url, _ := s.Find("h3.repo-list-name a").Attr("href")
		url = "https://github.com" + url
		stars := s.Find("p.repo-list-meta").Text()
		stars = regStars.FindString(stars)
		if stars == "" {
			stars = "0"
		}

		repo := GithubRepo{
			Title:       title,
			Description: description,
			Url:         url,
			Stars:       stars,
			Forks:       "0",
			Date:        int(time.Now().UTC().Unix()),
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

	Info.Println("Getting trending repos from GitHub...")
	scrapeTrendingRepos("go", *out)

	Info.Println("Getting most starred repos from GitHub...")
	scrapeMostStarredRepos("Go", *out)

	Info.Println("Done")
}
