package main

import (
	"image"
	"os"
	"image/png"
	"image/draw"
	"image/color"
	"fmt"
)

func Render(id Identicon, key string)  {
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))

	setBackgroundWhite(img)

	for i, v := range id.bitmap {
		if v == 1 {
			drawRect(img, i, id.color)
		}
	}

	f, _ := os.OpenFile(fmt.Sprintf("%s.png", key), os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}

func setBackgroundWhite(img *image.RGBA) {
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
}

func drawRect(rgba *image.RGBA, i int, c color.Color) {
	sizeSquare := 10
	maxRow := 5

	r := image.Rect(
		(i%maxRow)*sizeSquare,
		(i/maxRow)*sizeSquare,
		(i%maxRow)*sizeSquare+sizeSquare,
		(i/maxRow)*sizeSquare+sizeSquare,
	)

	draw.Draw(rgba, r, &image.Uniform{c}, image.ZP, draw.Src)
}
