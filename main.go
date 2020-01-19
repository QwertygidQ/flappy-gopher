package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const debug bool = true

func loadSprite(path string) *pixel.Sprite {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	picture := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(picture, picture.Bounds())

	return sprite
}

func setupWindow() *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, 1024, 768),
		//VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	return win
}

func run() {
	win := setupWindow()

	var (
		basicAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
		scoreText  = text.New(pixel.V(15, win.Bounds().H()-30), basicAtlas)
	)

	var (
		pipeSprite   = loadSprite("img/pipe.png")
		playerSprite = loadSprite("img/gopher.png")
		drawTarget   = pixel.Target(win)
		world        = newWorld(&drawTarget, win.Bounds(), pipeSprite, playerSprite)
	)

	for i := 0; i < 7; i++ {
		world.makePipe()
	}

	var (
		lastTime = time.Now()
		frames   = 0
		ticker   = time.Tick(time.Second)
	)

	for !win.Closed() {
		win.Update()

		// Delta time calculations
		dt := time.Since(lastTime).Seconds()
		lastTime = time.Now()

		// FPS calculations
		frames++
		select {
		case <-ticker:
			win.SetTitle(fmt.Sprintf("Flappy gopher| FPS: %d ", frames))
			frames = 0
		default:
		}

		// Game logic
		gameOver := world.update(dt, win.JustPressed(pixelgl.KeySpace))
		if gameOver {
			break
		}

		// Drawing
		win.Clear(colornames.Skyblue)

		world.draw(debug)

		scoreText.Clear()
		fmt.Fprintf(scoreText, "Score: %d", world.player.score)
		scoreText.Draw(win, pixel.IM.Scaled(scoreText.Orig, 2))
	}
}

func main() {
	pixelgl.Run(run)
}
