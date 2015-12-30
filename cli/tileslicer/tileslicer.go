package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"math"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"image/draw"
	_ "image/jpeg"
	"image/png"
)

var (
	file       = kingpin.Arg("file", "Input png tileset").Required().File()
	outDir     = kingpin.Arg("outdir", "Output dir").Required().String()
	tileWidth  = kingpin.Arg("width", "Tile width").Required().Int()
	tileHeight = kingpin.Arg("height", "Tile height").Required().Int()
)

func main() {

	// support -h for --help
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	inFileName := (*file).Name()

	if pathDontExist(*outDir) {
		err := os.Mkdir(*outDir, 0777)
		if err != nil {
			fmt.Printf("Could not create %s: %s", *outDir, err)
			os.Exit(1)
		}
	}

	sliceImage(inFileName, *outDir, *tileWidth, *tileHeight)
}

func pathDontExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		return true
	}
	return false
}

func sliceImage(imgFile string, outDir string, tileWidth int, tileHeight int) []image.Image {

	var slices []image.Image

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

	//fmt.Printf("%f cols, %f rows\n", cols, rows)

	// slice up image into tiles
	cnt := 0
	for row := 0; row < int(rows); row++ {
		for col := 0; col < int(cols); col++ {
			x0 := col * tileWidth
			y0 := row * tileHeight
			x1 := (col + 1) * tileWidth
			y1 := (row + 1) * tileHeight
			sr := image.Rect(x0, y0, x1, y1)

			dst := image.NewRGBA(image.Rect(0, 0, tileWidth, tileHeight))
			r := sr.Sub(sr.Min).Add(image.Point{0, 0})
			draw.Draw(dst, r, img, sr.Min, draw.Src)

			if isOnlyTransparent(dst) {
				fmt.Printf("Skipping empty tile at row %d, col %d\n", row, col)
				continue
			}

			outFile := fmt.Sprintf("%s/%03d.png", outDir, cnt)
			writePng(outFile, dst)
			cnt++
		}
	}

	fmt.Printf("%d tiles written to %s\n", cnt, outDir)
	return slices
}

// is this an empty tile?
func isOnlyTransparent(img *image.RGBA) bool {

	b := img.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				return false
			}
		}
	}

	return true
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
