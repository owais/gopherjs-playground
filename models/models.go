package models

type Player struct {
	ID     int    `json:"id"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Colour string `json:"colour"`
}
