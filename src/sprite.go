package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
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
