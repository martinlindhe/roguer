package rogue

import (
	"image"
	"image/color"
	"math"

	"github.com/martinlindhe/rogue/rollingparticle"
	"github.com/ojrac/opensimplex-go"
)

// Island ...
type Island struct {
	Width     int
	Height    int
	Seed      int64
	HeightMap [][]uint
	Spawns    []worldObject
}

// Tick ...
func (i *Island) Tick() {
	for _, o := range i.Spawns {
		o.Tick()
	}
}

// FillWithCritters ...
func (i *Island) FillWithCritters() {
	var dwarf dwarf
	dwarf.Defaults()
	dwarf.Name = "Gimli"
	dwarf.Position = Point{5, 5}

	i.Add(&dwarf)
}

// GenerateIsland returns a new island
func GenerateIsland(seed int64, width int, height int) Island {

	particleLength := 8
	innerBlur := 0.85
	outerBlur := 0.60
	roller := rollingparticle.New(seed, width, height, particleLength, innerBlur, outerBlur)

	/*
		rollerImage := Slice2DAsImage(&roller, width, height)
		rollerImgFile, _ := os.Create("roller.png")
		png.Encode(rollerImgFile, rollerImage)
	*/

	m := make2DUintSlice(width, height)

	noise := opensimplex.NewWithSeed(seed)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			// combine some sizes of noise
			fBig := noise.Eval2(float64(x)*0.005, float64(y)*0.005)
			fMid := noise.Eval2(float64(x)*0.01, float64(y)*0.01)
			fSmall := noise.Eval2(float64(x)*0.02, float64(y)*0.02)
			fMini := noise.Eval2(float64(x)*0.04, float64(y)*0.04)

			f := (fBig + fMid + fSmall + fMini) / 4

			// scale from -1.0-1.0 to 0.0-1.0
			f = (f + 1.0) / 2.0

			// scale from 0.0-1.0 to 0-255
			b := uint(0)
			if f == 1.0 {
				b = 255
			} else {
				b = uint(math.Floor(f * 256.0))
			}

			// combine with rolling particle
			opacity := 0.5
			b = uint((1-opacity)*float64(b) + opacity*float64(roller[y][x]))

			m[y][x] = b
		}
	}

	island := Island{
		Width:     width,
		Height:    height,
		Seed:      seed,
		HeightMap: m}

	return island
}

// Add ...
func (i *Island) Add(o worldObject) {
	i.Spawns = append(i.Spawns, o)
}

// ColoredHeightMapAsImage ...
func (i *Island) ColoredHeightMapAsImage() image.Image {

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{i.Width, i.Height}})

	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			b := i.HeightMap[y][x]

			var col color.RGBA
			switch {
			case b < 80:
				col = color.RGBA{0x26, 0x2f, 0x71, 0xff} // deep water

			case b < 90:
				col = color.RGBA{0x46, 0x4D, 0x85, 0xff} // shallow water

			case b < 95:
				col = color.RGBA{0xD4, 0xBC, 0x6A, 0xff} // beach

			case b < 150:
				col = color.RGBA{0x2D, 0x88, 0x2D, 0xff} // grass (green)

			case b < 230:
				col = color.RGBA{0x00, 0x4E, 0x00, 0xff} // forest (dark green)

			case b < 240:
				col = color.RGBA{0x4B, 0x2D, 0x12, 0xff} // hills (brown)

			default:
				col = color.RGBA{0xF2, 0xED, 0xE6, 0xff} // gray (mountains)
			}

			img.Set(x, y, col)
		}
	}

	return img
}
