package main

import (
	"github.com/faiface/pixel"
)

type World struct {
	target     *pixel.Target
	targetRect pixel.Rect

	pipeSprite *pixel.Sprite
	pipes      []*Pipe

	playerSprite *pixel.Sprite
	player       *Player
}

type Pipe struct {
	rect pixel.Rect
}

func newWorld(
	target *pixel.Target,
	targetRect pixel.Rect,
	pipeSprite *pixel.Sprite,
	playerSprite *pixel.Sprite,
) *World {
	player := newPlayer(
		playerSprite.Frame(),
		pixel.V(targetRect.W()/4, targetRect.Center().Y),
	)
	return &World{
		target:       target,
		targetRect:   targetRect,
		pipeSprite:   pipeSprite,
		playerSprite: playerSprite,
		player:       player,
	}
}

func (w *World) update(dt float64, spaceJustPressed bool) {
	w.player.updatePosition(dt, spaceJustPressed)
}

func (w *World) draw() {
	w.player.draw(w.target, w.playerSprite)
}
