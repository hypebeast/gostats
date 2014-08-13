package models

type GithubRepo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	Date        int    `json:"date"`
	Since       string `json:"since"`
}

var (
	trendingRepos    []GithubRepo
	mostStarredRepos []GithubRepo
)

func UpdateGithubStats() {
	// TODO: Update trending and most starred repos from BigQuery
}

func TrendingRepos() []GithubRepo {
	var repos []GithubRepo
	return repos
}

func MostStarredRepos() []GithubRepo {
	var repos []GithubRepo
	return repos
}
