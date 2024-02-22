package models

import (
	"github.com/gopxl/pixel"
	"github.com/svidlak/go-space-invaders/constants"
)

type Alien struct {
	Direction constants.Direction
	Position  pixel.Vec
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
