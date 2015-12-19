package main

import (
	"image/png"
	"os"

	"github.com/martinlindhe/rogue"
)

func main() {

	seed := int64(123456)

	island := rogue.GenerateIsland(seed, 120, 60)

	img := island.HeightMapAsImage()

	myfile, _ := os.Create("test.png")
	png.Encode(myfile, img)
}
