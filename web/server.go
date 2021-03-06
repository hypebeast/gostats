package main

import (
	"time"

	"github.com/hypebeast/gojistatic"
	"github.com/zenazn/goji"

	"github.com/hypebeast/gostats/web/models"
	"github.com/hypebeast/gostats/web/routes"
)

func main() {
	// Serve static files
	goji.Use(gojistatic.Static("public", gojistatic.StaticOptions{SkipLogging: true}))

	// Add routes
	routes.Include()

	trigger := make(chan bool)
	startTrending := make(chan bool)
	tickerChan := time.NewTicker(time.Minute * 10).C
	tickerTrending := time.NewTicker(time.Minute * 10).C

	// Update data from BigQuery
	go func() {
		for {
			select {
			case <-trigger:
				models.Update()
			case <-tickerChan:
				models.Update()
			}
		}
	}()

	// Update Github trending repos
	go func() {
		for {
			select {
			case <-startTrending:
				models.UpdateGithubTrendingRepos()
			case <-tickerTrending:
				models.UpdateGithubTrendingRepos()
			}
		}
	}()

	trigger <- true
	startTrending <- true

	// Run Goji
	goji.Serve()
}
