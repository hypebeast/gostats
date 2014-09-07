package controllers

import (
	"net/http"

	"github.com/hypebeast/gostats/web/models"
	"github.com/hypebeast/gostats/web/utils"
)

func Home(w http.ResponseWriter, req *http.Request) {
	templates := utils.BaseTemplates()
	templates = append(templates, "views/index.html")

	data := map[string]interface{}{
		"Title":                "GoStats - Daily statistics for the Go programming language",
		"DailyTrendingRepos":   models.DailyTrendingRepos(),
		"WeeklyTrendingRepos":  models.WeeklyTrendingRepos(),
		"MonthlyTrendingRepos": models.MonthlyTrendingRepos(),
		"MostStarredRepos":     models.MostStarredRepos(),
	}

	err := utils.RenderTemplate(w, templates, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
