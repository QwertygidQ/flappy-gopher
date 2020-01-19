package main

import (
	"github.com/faiface/pixel"
)

type Player struct {
	rect   pixel.Rect
	ySpeed float64
	score  int
}

func newPlayer(rect pixel.Rect, startPos pixel.Vec) *Player {
	rect = rect.Moved(startPos)
	return &Player{rect: rect}
}

func (p *Player) updatePosition(dt float64, spacePressed bool) {
	const g float64 = 980
	const jumpSpeed float64 = g / 1.7

	p.ySpeed -= dt * g
	if spacePressed {
		p.ySpeed = jumpSpeed
	}

	p.rect = p.rect.Moved(pixel.V(0, dt*p.ySpeed))
}

func (p *Player) draw(target *pixel.Target, sprite *pixel.Sprite) {
	playerMat := pixel.IM.Scaled(pixel.ZV, .15)
	playerMat = playerMat.Moved(p.rect.Min)
	sprite.Draw(*target, playerMat)
}
