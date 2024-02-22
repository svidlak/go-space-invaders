package models

import (
	"math/rand"

	"github.com/gopxl/pixel"
	"github.com/svidlak/go-space-invaders/constants"
)

type GameStateController struct {
	Score   int
	Aliens  []Alien
	Player  Player
	Bullets []Bullet
	Health  int
}

func (g *GameStateController) UpdateAlienMovement(dt float64, direction float64, movement float64) {
	for idx, alienVal := range g.Aliens {
		alien := &alienVal

		alien.UpdateDirection(direction, constants.WindowWidth)
		alien.Move(movement)

		g.Aliens[idx] = *alien
	}
}

func (g *GameStateController) ShootBullet(entity constants.Entity) {
	position := g.Player.Position

	if entity == constants.AlienEntity {
		randomAlienNumber := rand.Intn(len(g.Aliens) - 0)
		position = g.Aliens[randomAlienNumber].Position
	}

	g.Bullets = append(g.Bullets, Bullet{Position: position, Entity: entity})
}

func (g *GameStateController) AddAlien(bounds pixel.Rect, direction constants.Direction, position pixel.Vec) {
	alien := Alien{Direction: direction, Position: position, Bounds: bounds}
	g.Aliens = append(g.Aliens, alien)
}
