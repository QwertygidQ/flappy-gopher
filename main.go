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

	var (
		lastTime = time.Now()
		frames   = 0
		ticker   = time.Tick(time.Second)
	)

	for !win.Closed() {
		win.Update()

		// Delta time calculations
		dt := time.Since(lastTime).Seconds()
		_ = dt
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
		world.update(dt, win.JustPressed(pixelgl.KeySpace))

		// Drawing
		win.Clear(colornames.Skyblue)
		world.draw()

		pipeMat := pixel.IM.Moved(pixel.V(0, -pipeSprite.Frame().H()/2))
		pipeMat = pipeMat.Moved(pixel.V(win.Bounds().Center().X*3/4, win.Bounds().H()))
		pipeSprite.Draw(win, pipeMat)

		scoreText.Clear()
		fmt.Fprintf(scoreText, "Score: %d", world.player.score)
		scoreText.Draw(win, pixel.IM.Scaled(scoreText.Orig, 2))
	}
}

func main() {
	pixelgl.Run(run)
}
