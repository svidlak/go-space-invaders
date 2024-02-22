package models

import (
	"github.com/gopxl/pixel"
	"github.com/svidlak/go-space-invaders/constants"
)

type Alien struct {
	Direction constants.Direction
	Position  pixel.Vec
	Bounds    pixel.Rect
}

func (a *Alien) Move(height float64) {
	if a.Direction == constants.Right {
		a.Position.X += constants.PlayerSpeed / constants.AlienMovementOffset
	} else if a.Direction == constants.Down {
		a.Position.Y -= height
	} else {
		a.Position.X -= constants.PlayerSpeed / constants.AlienMovementOffset
	}
}

func (a *Alien) UpdateDirection(alienSpriteWidthBound, windowWidth float64) {
	if a.Position.X <= alienSpriteWidthBound {
		if a.Direction == constants.Down {
			a.Direction = constants.Right
		} else {
			a.Direction = constants.Down
		}
	}

	if a.Position.X >= windowWidth-((alienSpriteWidthBound*2)*2) {
		if a.Direction == constants.Down {
			a.Direction = constants.Left
		} else {
			a.Direction = constants.Down
		}
	}
}

func (a *Alien) DetectCollision(bullet *Bullet) bool {
	alienCenterPoint := a.Bounds.Max.X / 2

	minWidthPoint := a.Position.X - alienCenterPoint
	maxWidthPoint := a.Position.X + alienCenterPoint

	if bullet.Position.X >= minWidthPoint && bullet.Position.X <= maxWidthPoint && bullet.Position.Y >= a.Position.Y {
		return true
	}
	return false
}
