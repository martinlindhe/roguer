package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"math"
	"os"

	"github.com/davecgh/go-spew/spew"

	"image/draw"
	_ "image/jpeg"
	"image/png"
)

func main() {

	// XXXX read cli args: file width height

	// XXX 1 read source tileset

	/*
		font.png: 8x8
	*/

	inFile := "resources/assets/tilesets/oddball/font.png"
	tileWidth := 8
	tileHeight := 8

	xx := sliceImage(inFile, tileWidth, tileHeight)

	spew.Dump(xx)
}

func sliceImage(imgFile string, tileWidth int, tileHeight int) []image.Image {

	var slices []image.Image

	fmt.Printf("Source: %s\n", imgFile)

	img, _, err := decodeImage(imgFile)
	if err != nil {
		panic(err)
	}

	b := img.Bounds()
	imgWidth := b.Max.X
	imgHeight := b.Max.Y

	cols := float64(imgWidth) / float64(tileWidth)
	rows := float64(imgHeight) / float64(tileHeight)

	if cols != math.Floor(cols) {
		log.Fatalf("Input image width %d is not evenly divisable by tile width %d", imgWidth, tileWidth)
	}

	if rows != math.Floor(rows) {
		log.Fatalf("Input image height %d is not evenly divisable by tile height %d", imgHeight, tileHeight)
	}

	fmt.Printf("%f cols, %f rows\n", cols, rows)

	// slice up image into tiles
	cnt := -1
	for row := 0; row < int(rows); row++ {
		for col := 0; col < int(cols); col++ {
			cnt++
			x0 := col * tileWidth
			x1 := (col + 1) * tileWidth
			y0 := row * tileHeight
			y1 := (row + 1) * tileHeight
			sr := image.Rect(x0, y0, x1, y1)

			dst := image.NewRGBA(image.Rect(0, 0, tileWidth, tileHeight))
			r := sr.Sub(sr.Min).Add(image.Point{0, 0})
			draw.Draw(dst, r, img, sr.Min, draw.Src)

			if isOnlyTransparent(dst) {
				// XXX identify if this img has no pixels (only transparent ones), show error and skip
				fmt.Printf("Skipping empty tile %d\n", cnt)
			}

			outFile := fmt.Sprintf("tmp/%d.png", cnt)
			writePng(outFile, dst)
		}
	}

	return slices
}

func isOnlyTransparent(img *image.RGBA) bool {
	// XXX check alpha channel for all pixels

	b := img.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			oldPixel := img.At(x, y)
			_, _, _, a := oldPixel.RGBA()
			fmt.Printf("%d", a)
		}
	}
	fmt.Println("")

	return false
}

func decodeImage(filename string) (image.Image, string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	return image.Decode(bufio.NewReader(f))
}

func writePng(fileName string, m image.Image) error {

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = png.Encode(b, m)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}
	return nil
}
