package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"

	"github.com/owais/demo/models"
)

type Client struct {
	ws  *websocket.WebSocket
	In  chan models.Player
	Out chan [2]int
}

func New() Client {
	c := Client{
		In:  make(chan models.Player),
		Out: make(chan [2]int),
	}
	ws, err := websocket.New("ws://localhost:9100/game")
	if err != nil {
		panic(err)
	}

	c.ws = ws

	c.ws.AddEventListener("message", false, func(ev *js.Object) {
		go func() {
			data := []byte(ev.Get("data").String())
			p := models.Player{}
			json.Unmarshal(data, &p)
			c.In <- p
		}()
	})

	go c.listenToChanges()
	return c
}

func (c Client) listenToChanges() {
	for {
		change := <-c.Out
		player := models.Player{
			X: change[0],
			Y: change[1],
		}
		data, _ := json.Marshal(player)
		c.ws.Send(string(data))
	}
}

func (c Client) LoadInitialState() []models.Player {
	resp, _ := http.Get("http://localhost:9100/state")
	defer resp.Body.Close()

	contents, _ := ioutil.ReadAll(resp.Body)
	players := []models.Player{}
	json.Unmarshal(contents, &players)
	return players
}
