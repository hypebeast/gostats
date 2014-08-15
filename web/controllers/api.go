package controllers

import (
	"github.com/hypebeast/gostats/web/models"
	"github.com/hypebeast/gostats/web/utils"
	"net/http"
)

func GodocStats(w http.ResponseWriter, req *http.Request) {
	stats := models.GetGodocStats()
	utils.RenderJSON(w, http.StatusOK, stats)
}
