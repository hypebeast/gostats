package controllers

import (
	"net/http"

	"github.com/hypebeast/gostats/web/models"
)

func GodocStats(w http.ResponseWriter, req *http.Request) {
	stats := models.GetGodocStats()
	RenderJSON(w, http.StatusOK, stats)
}

func GithubRepoStarsHistory(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query()["name"][0]
	stats := models.RepoStarsHistory(name)
	RenderJSON(w, http.StatusOK, stats)
}
