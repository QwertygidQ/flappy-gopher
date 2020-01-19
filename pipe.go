package main

import "github.com/faiface/pixel"

type Pipe struct {
	rect     pixel.Rect
	inverted bool
}

func (p *Pipe) update(dt float64) {
	const moveSpeed float64 = 200
	p.rect = p.rect.Moved(pixel.V(-moveSpeed*dt, 0))
}

func (p *Pipe) draw(target *pixel.Target, sprite *pixel.Sprite) {
	mat := pixel.IM
	if p.inverted {
		mat = mat.ScaledXY(pixel.ZV, pixel.V(1, -1))
	}
	mat = mat.Moved(p.rect.Center())
	sprite.Draw(*target, mat)
}
