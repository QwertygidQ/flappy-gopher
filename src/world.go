package main

import (
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type World struct {
	target     *pixel.Target
	targetRect pixel.Rect

	soundHandler *SoundHandler

	pipeSprite *pixel.Sprite
	pipes      []*Pipe
	pipeTicker <-chan time.Time

	playerSprite *pixel.Sprite
	player       *Player
}

func newWorld(
	target *pixel.Target,
	targetRect pixel.Rect,
	soundHandler *SoundHandler,
	pipeSprite *pixel.Sprite,
	playerSprite *pixel.Sprite,
) *World {
	player := newPlayer(
		playerSprite.Frame().Moved(playerSprite.Frame().Size().Scaled(-.5)),
		pixel.V(targetRect.W()/4, targetRect.Center().Y),
	)
	return &World{
		target:       target,
		targetRect:   targetRect,
		soundHandler: soundHandler,
		pipeSprite:   pipeSprite,
		playerSprite: playerSprite,
		pipeTicker:   time.Tick(time.Second * 3),
		player:       player,
	}
}

func (w *World) makePipe() {
	inverted := rand.Intn(2) == 1 // Randomly picks either 0 or 1 and compares to 1 to convert to bool

	rect := w.pipeSprite.Frame().Moved(w.pipeSprite.Frame().Size().Scaled(-.5))
	if !inverted {
		rect = rect.Moved(pixel.V(0, w.targetRect.H()-w.pipeSprite.Frame().H()/2))
	} else {
		rect = rect.Moved(pixel.V(0, w.pipeSprite.Frame().H()/2))
	}
	rect = rect.Moved(pixel.V(w.targetRect.W()+w.pipeSprite.Frame().W()/2, 0))

	pipe := Pipe{rect: rect, inverted: inverted}
	w.pipes = append(w.pipes, &pipe)
}

func (w *World) update(dt float64, spaceJustPressed bool) (gameOver bool) {
	w.player.update(dt, spaceJustPressed)
	if w.player.rect.Min.Y <= 0 || w.player.rect.Max.Y >= w.targetRect.H() {
		gameOver = true
		return
	}

	newPipes := make([]*Pipe, 0)
	for _, pipe := range w.pipes {
		pipe.update(dt)

		if w.player.rect.Intersects(pipe.rect) {
			gameOver = true
			return
		} else if !pipe.passed && w.player.rect.Center().X >= pipe.rect.Center().X {
			w.soundHandler.playSound("pipe_passed", false)
			w.player.score++
			pipe.passed = true
		}

		if pipe.rect.Max.X >= 0 {
			newPipes = append(newPipes, pipe)
		}
	}
	w.pipes = newPipes

	select {
	case <-w.pipeTicker:
		w.makePipe()
	default:
	}

	return
}

func (w *World) draw(debug bool) {
	if debug {
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(1, 0, 0)

		for _, vert := range w.player.rect.Vertices() {
			imd.Push(vert)
		}
		imd.Polygon(3)

		for _, pipe := range w.pipes {
			for _, vert := range pipe.rect.Vertices() {
				imd.Push(vert)
			}
			imd.Polygon(3)
		}

		imd.Draw(*w.target)
	}

	w.player.draw(w.target, w.playerSprite)
	for _, pipe := range w.pipes {
		pipe.draw(w.target, w.pipeSprite)
	}
}
