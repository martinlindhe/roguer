package main

import (
	"image/png"
	"os"
	"time"

	"github.com/martinlindhe/rogue"
)

func main() {

	seed := time.Now().Unix() // int64(12345610)

	island := rogue.GenerateIsland(seed, 220, 160)

	island.FillWithCritters()
	for i := 0; i < 10; i++ {
		island.Tick()
	}

	// debug output:
	islandImage := island.HeightMapAsImage()
	islandColImage := island.ColoredHeightMapAsImage()

	islandImgFile, _ := os.Create("island.png")
	png.Encode(islandImgFile, islandImage)

	islandColImgFile, _ := os.Create("island_col.png")
	png.Encode(islandColImgFile, islandColImage)
}
