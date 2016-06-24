package components

import (
	c "github.com/owais/rendr/components"
	"honnef.co/go/js/dom"

	"github.com/owais/demo/api"
	"github.com/owais/demo/models"
)

const (
	left  int = 37
	right int = 39
	up    int = 38
	down  int = 40
)

type root struct {
	c.Component

	Title   string
	Players map[int]models.Player
	id      int
	client  api.Client
}

func (r root) onKeyDown(ev dom.Event) {
	offset := 10

	switch ev.Underlying().Get("keyCode").Int() {
	case left:
		r.client.Out <- [2]int{

			r.Players[r.id].X - offset,
			r.Players[r.id].Y,
		}
	case right:
		r.client.Out <- [2]int{
			r.Players[r.id].X + offset,
			r.Players[r.id].Y,
		}
	case up:
		r.client.Out <- [2]int{
			r.Players[r.id].X,
			r.Players[r.id].Y - offset,
		}
	case down:
		r.client.Out <- [2]int{
			r.Players[r.id].X,
			r.Players[r.id].Y + offset,
		}
	}
}

func (r root) Render() c.Renderer {

	children := []c.Renderer{}
	for _, player := range r.Players {
		children = append(children, newPlayer(player))
	}

	children = append(children, c.H3(
		c.Text(r.Title)),
		c.Style("color", "red"),
	)

	return r.Append(children...)
}

func (r root) Start(client api.Client, player models.Player) {
	r.client = client
	r.id = player.ID
	r.Players[r.id] = player

	document := dom.GetWindow().Document()

	document.AddEventListener("keydown", false, r.onKeyDown)
}

func NewApp(title string) root {
	return root{
		Title:     title,
		Component: c.Div(nil).(c.Component),
		Players:   map[int]models.Player{},
	}
}
