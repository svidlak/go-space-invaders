package models

type GameState struct {
	Score  int
	Aliens []Alien
	Player Player
}
