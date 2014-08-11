package main

import (
	"github.com/hypebeast/gojistatic"
	"github.com/zenazn/goji"

	"github.com/hypebeast/gostats/web/routes"
)

func main() {
	// Serve static files
	goji.Use(gojistatic.Static("public", gojistatic.StaticOptions{SkipLogging: true}))

	// Add routes
	routes.Include()

	// Run Goji
	goji.Serve()
}
