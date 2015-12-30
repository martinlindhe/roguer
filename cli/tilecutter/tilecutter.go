package main

import (
	"bufio"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {

	// TODO: take command line input: allow to choose bottom or top part, allow to choose half, 1/3 or 2/3 to keep

	// loop over input folder, keep bottom third of each image, overwrite

	inDir := "resources/assets/tilesets/oddball/tiles/8x4"
	files, _ := ioutil.ReadDir(inDir)
	for _, f := range files {

		p := filepath.Join(inDir, f.Name())
		img := getBottomThirdOfImage(p)

		writePng(p, img)
	}
}

func getBottomThirdOfImage(fileName string) *image.RGBA {
	img, _, err := decodeImage(fileName)
	if err != nil {
		panic(err)
	}

	b := img.Bounds()

	tX := b.Max.X
	tY := b.Max.Y / 3

	sr := image.Rect(0, tY*2, b.Max.X, b.Max.Y)

	dst := image.NewRGBA(image.Rect(0, 0, tX, tY))
	r := sr.Sub(sr.Min).Add(image.Point{0, 0})
	draw.Draw(dst, r, img, sr.Min, draw.Src)

	return dst
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
