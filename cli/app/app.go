package main

import (
	"image/png"
	"os"

	"github.com/martinlindhe/rogue"
)

func main() {

	seed := int64(123456)

	island := rogue.GenerateIsland(seed, 220, 160)

	islandImage := island.HeightMapAsImage()

	islandImgFile, _ := os.Create("island.png")
	png.Encode(islandImgFile, islandImage)
}
