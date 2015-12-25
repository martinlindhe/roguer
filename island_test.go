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
	island := generateIsland(seed, 200, 100)
	island.fillWithCritters()

	// XXX currently failing:
	assert.Equal(t, true, len(island.Spawns) > 0)

	islandColImgFile, _ := os.Create("island_test.png")
	png.Encode(islandColImgFile, island.ColoredHeightMapAsImage())

	for i := 0; i < 3; i++ {
		island.Tick()
	}

	//spew.Dump(island.Spawns)
}
