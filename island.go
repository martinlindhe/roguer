package rogue

import (
	"image"
	"image/color"
	"math"

	"github.com/ojrac/opensimplex-go"
)

// Island ...
type Island struct {
	Width     int
	Height    int
	Seed      int64
	HeightMap [][]byte
}

// GenerateIsland returns a new island
func GenerateIsland(seed int64, width int, height int) Island {

	var island Island

	island.Width = width
	island.Height = height
	island.Seed = seed

	m := make2DByteSlice(width, height)

	noise := opensimplex.NewWithSeed(seed)

	//$roller = new RollingParticle($world->width, $world->height, $world->seed);
	//$rolls = $roller->roll();

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			noiseX := float64(x) * 0.04
			noiseY := float64(y) * 0.04

			f := noise.Eval2(noiseX, noiseY)

			// scale from -1.0-1.0 to 0.0-1.0
			f = (f + 1.0) / 2.0

			// scale from 0.0-1.0 to 0-255
			b := byte(0)
			if f == 1.0 {
				b = 255
			} else {
				b = byte(math.Floor(f * 256.0))
			}
			// 566883 ns/op benchmark with [x][y]
			m[y][x] = b
		}
	}

	island.HeightMap = m
	return island
}

// HeightMapAsImage renders height map to an Image
func (i *Island) HeightMapAsImage() image.Image {

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{i.Width, i.Height}})

	for x := 0; x < i.Width; x++ {
		for y := 0; y < i.Height; y++ {
			b := i.HeightMap[x][y]
			c := color.RGBA{b, b, b, 255}
			img.Set(x, y, c)
		}
	}

	return img
}
