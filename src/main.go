package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const debug bool = false

func handleFlags() pixelgl.WindowConfig {
	vSync := flag.Bool("vsync", false, "Enable VSync")
	resolutionStr := flag.String("resolution", "1024x768", "Window resolution")

	flag.Parse()

	resolution := strings.Split(*resolutionStr, "x")
	if len(resolution) != 2 {
		panic("Resolution must be in <WIDTH>x<HEIGHT> format.")
	}

	width, err := strconv.Atoi(resolution[0])
	if err != nil {
		panic(err)
	}

	height, err := strconv.Atoi(resolution[1])
	if err != nil {
		panic(err)
	}

	if width <= 0 || height <= 0 {
		panic("Both width and height must be positive.")
	}

	cfg := pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, float64(width), float64(height)),
		VSync:  *vSync,
	}
	return cfg
}

func setupWindow() *pixelgl.Window {
	cfg := handleFlags()

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	return win
}

func setupScoreText(targetRect pixel.Rect) *text.Text {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	scoreText := text.New(pixel.V(15, targetRect.H()-30), basicAtlas)

	return scoreText
}

func setupSoundHandler() *SoundHandler {
	soundHandler := initBeep(48000)
	soundHandler.loadSound("snd/jetpack.wav", "jetpack")
	soundHandler.loadSound("snd/pipe_passed.wav", "pipe_passed")
	soundHandler.playSound("jetpack", true)

	return soundHandler
}

func setupWorld(target *pixel.Target, targetRect pixel.Rect, soundHandler *SoundHandler) *World {
	var (
		pipeSprite   = loadSprite("img/pipe.png")
		playerSprite = loadSprite("img/gopher.png")
		world        = newWorld(target, targetRect, soundHandler, pipeSprite, playerSprite)
	)
	world.makePipe()

	return world
}

func run() {
	win := setupWindow()

	scoreText := setupScoreText(win.Bounds())
	soundHandler := setupSoundHandler()

	var (
		drawTarget = pixel.Target(win)
		world      = setupWorld(&drawTarget, win.Bounds(), soundHandler)
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
		lastTime = time.Now()

		// FPS calculations
		frames++
		select {
		case <-ticker:
			win.SetTitle(fmt.Sprintf("Flappy Gopher | FPS: %d ", frames))
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
