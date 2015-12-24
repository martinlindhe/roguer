package rogue

import (
	"image/png"
	"os"
	"testing"
)

func BenchmarkGenerateIsland(b *testing.B) {
	for n := 0; n < b.N; n++ {

		GenerateIsland(666, 220, 140)
	}
}

func TestGenerateIsland(t *testing.T) {

	seed := int64(123)
	island := GenerateIsland(seed, 200, 100)
	island.FillWithCritters()

	islandColImgFile, _ := os.Create("island_test.png")
	png.Encode(islandColImgFile, island.ColoredHeightMapAsImage())

	for i := 0; i < 1; i++ {
		island.Tick()
	}

	island.PrintSpawns()
}
