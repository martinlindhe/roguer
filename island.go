package rogue

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/ojrac/opensimplex-go"
)

type Island struct {
	Width     int
	Height    int
	Seed      int64
	HeightMap [][]byte
}

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

			m[x][y] = b
		}
	}

	island.HeightMap = m
	return island
}

func make2DByteSlice(width int, height int) [][]byte {
	// allocate 2d slice
	m := make([][]byte, width)
	for i := range m {
		m[i] = make([]byte, height)
	}
	return m
}

func (i *Island) WriteHeightMapAsImage(outFileName string) {

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{i.Width, i.Height}})

	for x := 0; x < i.Width; x++ {
		for y := 0; y < i.Height; y++ {
			b := i.HeightMap[x][y]
			c := color.RGBA{b, b, b, 255}
			img.Set(x, y, c)
		}
	}

	myfile, _ := os.Create(outFileName)

	png.Encode(myfile, img)
}
