package main

// Combines all *.png:s in a dir into a new tileset

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	inDir       = kingpin.Arg("indir", "Input directory").Required().String()
	outFile     = kingpin.Arg("outfile", "Output png file").Required().String()
	tilesPerRow = kingpin.Arg("tilesperrow", "Tiles per row").Required().Int()
)

func main() {

	// support -h for --help
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	//inDir := ""
	//outFile := "joined-tiles.png"

	var images []image.Image

	tileWidth := 0
	tileHeight := 0
	tileCount := 0

	// XXX must be alphabetical order!!!
	files, _ := ioutil.ReadDir(*inDir)
	for _, f := range files {

		p := filepath.Join(*inDir, f.Name())

		img, _, err := decodeImage(p)
		if err != nil {
			fmt.Printf("Error decoding: %s", err)
			continue
		}

		b := img.Bounds()

		if tileWidth == 0 && tileHeight == 0 {
			tileWidth = b.Max.X
			tileHeight = b.Max.Y
		} else if b.Max.X != tileWidth || b.Max.Y != tileHeight {
			fmt.Printf("Error: tile %s did not have expected dimensions of %d,%d\n", p, tileWidth, tileHeight)
		}

		tileCount++
		images = append(images, img)
	}

	outWidth := *tilesPerRow * tileWidth
	outHeight := (tileCount / *tilesPerRow) * tileHeight

	fmt.Printf("Creating tileset of %d tiles with %d,%d pixels, %d tiles per row. Output is image is %d,%d pixels\n", tileCount, tileWidth, tileHeight, *tilesPerRow, outWidth, outWidth)

	dst := image.NewRGBA(image.Rect(0, 0, outWidth, outHeight))
	for i, img := range images {
		x0 := (i % tileWidth) * tileWidth
		y0 := (i / tileWidth) * tileHeight

		// fmt.Printf("Writing %d to %d,%d\n", i, x0, y0)

		dr := image.Rect(x0, y0, x0+tileWidth, y0+tileHeight)

		draw.Draw(dst, dr, img, image.Point{0, 0}, draw.Src)
	}

	fmt.Printf("Writing to %s\n", *outFile)
	writePng(*outFile, dst)
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
