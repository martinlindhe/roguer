package rogue

import (
	"image"
	"image/png"
	"math"
	"os"

	"github.com/martinlindhe/rogue/rollingparticle"
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

	particleLength := 50
	outerBlur := 0.65
	innerBlur := 0.90
	roller := rollingparticle.New(seed, island.Width, island.Height, particleLength, innerBlur, outerBlur)

	rollerImage := Slice2DAsImage(&roller, island.Width, island.Height)
	rollerImgFile, _ := os.Create("roller.png")
	png.Encode(rollerImgFile, rollerImage)

	m := make2DByteSlice(width, height)

	noise := opensimplex.NewWithSeed(seed)

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

			// XXX combine with rolling particle
			opacity := 0.5

			b = byte((1-opacity)*float64(b) + opacity*float64(roller[y][x]))

			// 566883 ns/op benchmark with [x][y]
			m[y][x] = b
		}
	}

	island.HeightMap = m
	return island
}

// HeightMapAsImage renders height map to an Image
func (i *Island) HeightMapAsImage() image.Image {

	return Slice2DAsImage(&i.HeightMap, i.Width, i.Height)
}
