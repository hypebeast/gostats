package routes

import (
	"github.com/zenazn/goji"

	"github.com/hypebeast/gostats/web/controllers"
)

func Include() {
	goji.Get("/", controllers.Home)
	goji.Get("/about", controllers.About)
}
