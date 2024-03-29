package models

import "github.com/gopxl/pixel"

type Player struct {
	Position pixel.Vec
	Bounds   pixel.Rect
	Health   int
}

func (p *Player) DetectCollision(bullet *Bullet) bool {
	playerCenterPoint := p.Bounds.Max.X / 2

	minWidthPoint := p.Position.X - playerCenterPoint
	maxWidthPoint := p.Position.X + playerCenterPoint

	if bullet.Position.X >= minWidthPoint && bullet.Position.X <= maxWidthPoint && bullet.Position.Y <= p.Position.Y {
		p.Health -= 1
		return true
	}
	return false
}
