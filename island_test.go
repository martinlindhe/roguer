package rogue

import (
	"fmt"
	"image/png"
	"os"
	"reflect"
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

	// make npc hungry
	island.Spawns[0].Hunger = island.Spawns[0].hungerCap() + 1

	island.Tick()

	// make sure that first critter has aged
	assert.Equal(t, true, island.Spawns[0].Age > 0)

	// make sure planned action: find food
	v := reflect.ValueOf(island.Spawns[0].CurrentAction)
	assert.Equal(t, "*rogue.lookForFood", fmt.Sprintf("%s", v.Type()))

	assert.Equal(t, false, island.Spawns[0].hasItemTypeInInventory("food"))

	// progress until npc found food
	for i := 0; i < 5; i++ { // XXXX need to find action "find food".duration, from actions.yml
		island.Tick()
	}

	// make sure food was found
	assert.Equal(t, true, len(island.Spawns[0].Inventory) > 0)
	assert.Equal(t, true, island.Spawns[0].hasItemTypeInInventory("food"))
	assert.Equal(t, false, island.Spawns[0].hasItemTypeInInventory("water"))
}

// XXX need tests for behaviour now: need 1 dwarf to live on the map, and be hungry

// XXXX: make dwarf thirsty and give it water, make sure it drinks it.. etc
