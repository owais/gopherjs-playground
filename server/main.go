package main

import (
	"encoding/json"
	"html/template"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/owais/demo/components"
	"github.com/owais/demo/models"
)

type connection struct {
	*models.Player
	ws *websocket.Conn
}

var connections = map[int]connection{}

var colours = []string{
	"red", "green", "blue", "magenta", "orange",
}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func gameHandler(w http.ResponseWriter, r *http.Request) {

	// establish new websocket connection
	conn, _ := upgrader.Upgrade(w, r, nil)

	id := len(connections)

	// listen for position update requests from current player

	// create new player
	player := &models.Player{
		ID:     id,
		X:      rand.Intn(100),
		Y:      rand.Intn(100),
		Colour: colours[rand.Intn(len(colours))],
	}

	connections[id] = connection{player, conn}

	// send new players position to all clients
	playerJSON, _ := json.Marshal(player)
	for _, conn := range connections {
		conn.ws.WriteMessage(websocket.TextMessage, playerJSON)
	}

	for {

		messageType, data, err := conn.ReadMessage()

		if err != nil {
			break
		}

		p := models.Player{}
		json.Unmarshal(data, &p)

		player.X = p.X
		player.Y = p.Y

		data, _ = json.Marshal(player)

		for _, conn := range connections {
			err := conn.ws.WriteMessage(messageType, data)
			if err != nil {
				break
			}
		}

	}

	connections[id].ws.Close()

	delete(connections, id)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	app := components.NewApp("rendered by server")
	for _, conn := range connections {
		app.Players[conn.Player.ID] = *conn.Player
	}

	t, _ := template.ParseFiles("index.html")

	t.Execute(w, template.HTML(app.Render().Text()))
}

func stateHandler(w http.ResponseWriter, r *http.Request) {

	players := []models.Player{}
	for _, conn := range connections {
		players = append(players, *conn.Player)
	}

	data, _ := json.Marshal(players)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {

	http.HandleFunc("/game", gameHandler)

	http.HandleFunc("/state", stateHandler)
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":9100", nil)
}
