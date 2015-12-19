package main

import (
	"image/png"
	"os"

	"github.com/martinlindhe/rogue"
)

func main() {

	seed := int64(123456)

	island := rogue.GenerateIsland(seed, 220, 160)

	var dwarf rogue.Dwarf
	dwarf.Name = "Gimli"

	island.Add(dwarf)

	//fmt.Println(island.Spawns)

	// debug output:
	islandImage := island.HeightMapAsImage()
	islandColImage := island.ColoredHeightMapAsImage()

	islandImgFile, _ := os.Create("island.png")
	png.Encode(islandImgFile, islandImage)

	islandColImgFile, _ := os.Create("island_col.png")
	png.Encode(islandColImgFile, islandColImage)
}
