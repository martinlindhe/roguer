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

func TestGenerateIslandOneDwarf(t *testing.T) {

	seed := int64(123)
	generateIsland(seed, 200, 100)
	island.addNpcFromType("dwarf")

	islandColImgFile, _ := os.Create("island_test.png")
	png.Encode(islandColImgFile, island.ColoredHeightMapAsImage())

	// make sure spawns was created
	assert.Equal(t, true, len(island.Spawns) == 1)

	dw := island.Spawns[0]

	// make npc hungry
	island.Spawns[0].Hunger = dw.hungerCap() + 1

	island.Tick()

	// make sure that first critter has aged
	assert.Equal(t, true, dw.Age > 0)

	// make sure planned action: find food
	v := reflect.ValueOf(dw.CurrentAction)
	assert.Equal(t, "*rogue.lookForFood", fmt.Sprintf("%s", v.Type()))

	assert.Equal(t, false, dw.hasItemTypeInInventory("food"))

	// progress until npc found food
	for i := 0; i < 5; i++ { // XXXX need to find action "find food".duration, from actions.yml
		island.Tick()
	}

	// make sure food was found
	assert.Equal(t, true, len(dw.Inventory) > 0)
	assert.Equal(t, true, dw.hasItemTypeInInventory("food"))
	assert.Equal(t, false, dw.hasItemTypeInInventory("water"))

	oldHunger := dw.Hunger
	// let npc consume food
	island.Tick()
	assert.Equal(t, false, dw.hasItemTypeInInventory("food"))

	// make sure hunger went down
	assert.Equal(t, true, dw.Hunger < oldHunger)

	assert.Equal(t, false, dw.isHungry())
}

// XXX need tests for behaviour now: need 1 dwarf to live on the map, and be hungry

// XXXX: make dwarf thirsty and give it water, make sure it drinks it.. etc
