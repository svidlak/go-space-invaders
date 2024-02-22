package models

import (
	"github.com/gopxl/pixel"
	"github.com/svidlak/go-space-invaders/constants"
)

type Bullet struct {
	Entity   constants.Entity
	Position pixel.Vec
}

func (b *Bullet) MoveBullet(bullet *Bullet) {
	if bullet.Entity == constants.AlienEntity {
		bullet.Position.Y--
	} else {
		bullet.Position.Y++
	}
}
