package models

import (
	"github.com/gopxl/pixel"
	"github.com/svidlak/go-space-invaders/constants"
)

type Bullet struct {
	Entity   constants.Entity
	Position pixel.Vec
}
