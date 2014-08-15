package routes

import (
	"github.com/zenazn/goji"

	"github.com/hypebeast/gostats/web/controllers"
)

func Include() {
	goji.Get("/", controllers.Home)

	// API
	goji.Get("/api/godocstats", controllers.GodocStats)
	goji.Get("/api/github/trending", controllers.GodocStats)
	goji.Get("/api/github/stars", controllers.GodocStats)
}
