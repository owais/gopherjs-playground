package components

import (
	"strconv"

	c "github.com/owais/rendr/components"

	"github.com/owais/demo/models"
)

func newPlayer(player models.Player) c.Renderer {
	style := map[string]string{
		"width":      "10px",
		"height":     "10px",
		"position":   "absolute",
		"background": player.Colour,
		"left":       strconv.Itoa(player.X) + "px",
		"top":        strconv.Itoa(player.Y) + "px",
	}
	return c.Div(c.Styles(style)...)
}
