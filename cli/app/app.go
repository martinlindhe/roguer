package main

import (
	"image/png"
	"os"

	"github.com/martinlindhe/rogue"
	"github.com/martinlindhe/rogue/rollingparticle"
)

func main() {

	seed := int64(123456)

	island := rogue.GenerateIsland(seed, 120, 60)

	islandImage := island.HeightMapAsImage()

	particleLength := 50
	outerBlur := 0.65
	innerBlur := 0.90
	particle := rollingparticle.New(island.Width, island.Height, particleLength, innerBlur, outerBlur)

	rollerImage := rogue.Slice2DAsImage(&particle, island.Width, island.Height)

	islandImgFile, _ := os.Create("island.png")
	png.Encode(islandImgFile, islandImage)

	rollerImgFile, _ := os.Create("roller.png")
	png.Encode(rollerImgFile, rollerImage)
}
