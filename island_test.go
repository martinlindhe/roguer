package rogue

import (
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkGenerateIsland(b *testing.B) {
	for n := 0; n < b.N; n++ {

		generateIsland(666, 220, 140)
	}
}

func TestGenerateIsland(t *testing.T) {

	seed := int64(123)
	generateIsland(seed, 200, 100)
	island.fillWithCritters()

	islandColImgFile, _ := os.Create("island_test.png")
	png.Encode(islandColImgFile, island.ColoredHeightMapAsImage())

	// make sure spawns was created
	assert.Equal(t, true, len(island.Spawns) > 0)

	island.Tick()
	island.Tick()

	// make sure that first critter has aged! XXX currently failing:
	assert.Equal(t, true, island.Spawns[0].Age > 1)

	//spew.Dump(island.Spawns)
}
