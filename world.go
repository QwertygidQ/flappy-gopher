package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type World struct {
	targetWin *pixelgl.Window

	pipeSprite *pixel.Sprite
	pipes      []*Pipe

	playerSprite *pixel.Sprite
	player       *Player
}

type Pipe struct {
	rect pixel.Rect
}

func newWorld(
	pipeSprite *pixel.Sprite,
	playerSprite *pixel.Sprite,
	targetWin *pixelgl.Window,
) *World {
	player := newPlayer(
		playerSprite.Frame(),
		pixel.V(targetWin.Bounds().W()/4, targetWin.Bounds().Center().Y),
	)
	return &World{targetWin: targetWin, pipeSprite: pipeSprite, playerSprite: playerSprite, player: player}
}

func (w *World) draw() {
	w.player.draw(w.targetWin, w.playerSprite)
}
