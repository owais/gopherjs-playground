package main

import (
	"github.com/owais/rendr/dom"

	"github.com/owais/demo/api"
	"github.com/owais/demo/components"
)

func main() {
	client := api.New()

	app := components.NewApp("Rendered by browser")

	for _, player := range client.LoadInitialState() {
		app.Players[player.ID] = player
	}

	player := <-client.In

	app.Start(client, player)

	for {
		dom.Render("#app-container", app)
		player := <-client.In
		app.Players[player.ID] = player
	}
}
