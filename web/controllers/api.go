package controllers

import (
	"net/http"

	"github.com/hypebeast/gostats/web/models"
	"github.com/hypebeast/gostats/web/utils"
)

func GodocStats(w http.ResponseWriter, req *http.Request) {
	stats := models.GetGodocStats()
	utils.RenderJSON(w, http.StatusOK, stats)
}

func GithubRepoStarsHistory(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query()["name"][0]
	stats := models.RepoStarsHistory(name)
	utils.RenderJSON(w, http.StatusOK, stats)
}
