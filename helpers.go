package rogue

import (
	"image"
	"image/color"
)

// returns a 2d slice in [height][width]
func make2DByteSlice(width int, height int) [][]byte {
	// allocate 2d slice
	m := make([][]byte, height)
	for i := range m {
		m[i] = make([]byte, width)
	}
	return m
}

// returns a 2d slice in [height][width]
func make2DUintSlice(width int, height int) [][]uint {
	// allocate 2d slice
	m := make([][]uint, height)
	for i := range m {
		m[i] = make([]uint, width)
	}
	return m
}

func Slice2DAsImage(data *[][]byte, width int, height int) image.Image {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

	p := *data

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			b := p[y][x]
			c := color.RGBA{b, b, b, 255}
			img.Set(x, y, c)
		}
	}

	return img
}
