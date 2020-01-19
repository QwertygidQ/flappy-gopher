package main

import (
	"fmt"
	"image"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const debug bool = false

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

	soundHandler := initBeep(48000)
	soundHandler.loadSound("snd/jetpack.wav", "jetpack")
	soundHandler.loadSound("snd/pipe_passed.wav", "pipe_passed")
	soundHandler.playSound("jetpack", true)

	var (
		pipeSprite   = loadSprite("img/pipe.png")
		playerSprite = loadSprite("img/gopher.png")
		drawTarget   = pixel.Target(win)
		world        = newWorld(&drawTarget, win.Bounds(), soundHandler, pipeSprite, playerSprite)
	)
	world.makePipe()

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

	fmt.Printf("Thanks for playing! Final score: %d\n", world.player.score)
	for _, decoded := range soundHandler.decoded {
		decoded.streamer.Close()
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	pixelgl.Run(run)
}
