package models

type GithubRepo struct {
}

var (
	trendingRepos    []GithubRepo
	mostStarredRepos []GithubRepo
)

func Update() {
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
