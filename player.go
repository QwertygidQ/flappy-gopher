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
	rect = rect.Resized(pixel.ZV, rect.Size().Scaled(.15))
	rect = rect.Moved(startPos)
	return &Player{rect: rect}
}

func (p *Player) update(dt float64, spacePressed bool) {
	const g float64 = 980
	const jumpSpeed float64 = g / 1.7

	p.ySpeed -= dt * g
	if spacePressed {
		p.ySpeed = jumpSpeed
	}

	p.rect = p.rect.Moved(pixel.V(0, dt*p.ySpeed))
}

func (p *Player) draw(target *pixel.Target, sprite *pixel.Sprite) {
	const scaleFactor float64 = .15

	playerMat := pixel.IM.Scaled(pixel.ZV, scaleFactor)
	playerMat = playerMat.Moved(p.rect.Center())
	sprite.Draw(*target, playerMat)
}
