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
	island.addNpcFromType("dwarf")

	islandColImgFile, _ := os.Create("island_test.png")
	png.Encode(islandColImgFile, island.ColoredHeightMapAsImage())

	// make sure spawns was created
	assert.Equal(t, true, len(island.Spawns) == 1)

	island.Tick()

	// make sure that first critter has aged
	assert.Equal(t, true, island.Spawns[0].Age > 0)

	// inject an idea so it will happen before something is auto-picked

	for i := 0; i < 5; i++ {
		island.Tick()
	}

	assert.Equal(t, true, len(island.Spawns[0].Inventory) > 0)

	//spew.Dump(island.Spawns)
}

// XXX need tests for behaviour now: need 1 dwarf to live on the map, and be thirsty
