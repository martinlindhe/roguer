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
	island.addNpcFromType("dwarf")

	islandColImgFile, _ := os.Create("island_test.png")
	png.Encode(islandColImgFile, island.ColoredHeightMapAsImage())

	// make sure npcs was created
	assert.Equal(t, true, len(island.Spawns) == 2)

	dw := island.Spawns[0]
	dw2 := island.Spawns[1]

	// make npc 1 hungry
	dw.Hunger = dw.hungerCap() + 1

	// make npc 2 thirsty
	dw2.Thirst = dw.thirstCap() + 1

	island.Tick()

	// make sure that npc has aged
	assert.Equal(t, true, dw.Age > 0)
	assert.Equal(t, true, dw2.Age > 0)

	// make sure npc 1 planned action: find food
	assert.Equal(t, "*rogue.lookForFood", fmt.Sprintf("%s", reflect.TypeOf(dw.CurrentAction)))
	assert.Equal(t, false, dw.hasItemTypeInInventory("food"))

	// make sure npc 2 planned action: find water
	assert.Equal(t, "*rogue.lookForWater", fmt.Sprintf("%s", reflect.TypeOf(dw2.CurrentAction)))

	// progress until npc found food
	for i := 0; i < 5; i++ { // XXXX need to find action "find food".duration, from actions.yml
		island.Tick()
	}

	// make sure food was found for npc 1
	assert.Equal(t, true, len(dw.Inventory) > 0)
	assert.Equal(t, true, dw.hasItemTypeInInventory("food"))
	assert.Equal(t, false, dw.hasItemTypeInInventory("drink"))

	// make sure water was found for npc 2
	assert.Equal(t, true, len(dw2.Inventory) > 0)
	assert.Equal(t, true, dw2.hasItemTypeInInventory("drink"))
	assert.Equal(t, false, dw2.hasItemTypeInInventory("food"))

	oldHunger := dw.Hunger
	island.Tick()

	// make sure npc 1 consumed food
	assert.Equal(t, false, dw.hasItemTypeInInventory("food"))

	// make sure npc 2 consumed water
	assert.Equal(t, false, dw.hasItemTypeInInventory("drink"))

	// make sure hunger went down
	assert.Equal(t, true, dw.Hunger < oldHunger)

	assert.Equal(t, false, dw.isHungry())

	assert.Equal(t, false, dw2.isThirsty())
}

// XXX need tests for behaviour now: need 1 dwarf to live on the map, and be hungry

// XXXX: make dwarf thirsty and give it water, make sure it drinks it.. etc
