package controllers

import (
	"github.com/hypebeast/gostats/web/utils"

	"net/http"
)

func Home(w http.ResponseWriter, req *http.Request) {
	templates := utils.BaseTemplates()
	templates = append(templates, "views/index.html")
	err := utils.RenderTemplate(w, templates, "base", map[string]string{"Title": "GoStats"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
