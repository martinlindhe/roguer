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

func prepareIsland() {
	if island == nil {
		seed := int64(123)
		island = generateIsland(seed, 200, 100)

		islandColImgFile, _ := os.Create("island_test.png")
		png.Encode(islandColImgFile, island.ColoredHeightMapAsImage())
	}

	// clear spawns between tests
	island.Spawns = nil
}

func TestFindFood(t *testing.T) {

	prepareIsland()

	island.addNpcFromType("dwarf")
	assert.Equal(t, true, len(island.Spawns) == 1)
	dw := island.Spawns[0]

	// make npc hungry
	dw.Hunger = dw.hungerCap() + 1
	island.Tick()

	// make sure that npc has aged
	assert.Equal(t, true, dw.Age > 0)

	// make sure planned action: find food
	assert.Equal(t, "find-food", dw.CurrentAction)
	assert.Equal(t, false, dw.hasItemTypeInInventory("food"))

	// progress until npc found food
	for i := 0; i < 5; i++ { // XXXX need to find action "find food".duration, from actions.yml
		island.Tick()
	}

	// make sure food was found
	assert.Equal(t, true, len(dw.Inventory) > 0)
	assert.Equal(t, true, dw.hasItemTypeInInventory("food"))
	assert.Equal(t, false, dw.hasItemTypeInInventory("drink"))

	oldHunger := dw.Hunger
	island.Tick()

	// make sure npc consumed food
	assert.Equal(t, false, dw.hasItemTypeInInventory("food"))

	// make sure hunger went down
	assert.Equal(t, true, dw.Hunger < oldHunger)
	assert.Equal(t, false, dw.isHungry())
}

func TestFindWater(t *testing.T) {

	prepareIsland()

	island.addNpcFromType("dwarf")
	assert.Equal(t, true, len(island.Spawns) == 1)
	dw := island.Spawns[0]

	// make npc thirsty
	dw.Thirst = dw.thirstCap() + 1
	island.Tick()

	// make sure npc planned action: find water
	assert.Equal(t, "find-water", dw.CurrentAction)
	assert.Equal(t, false, dw.hasItemTypeInInventory("drink"))

	// progress until npc found food
	for i := 0; i < 5; i++ { // XXXX need to find action "find food".duration, from actions.yml
		island.Tick()
	}

	// make sure water was found
	assert.Equal(t, true, len(dw.Inventory) > 0)
	assert.Equal(t, true, dw.hasItemTypeInInventory("drink"))

	oldThirst := dw.Thirst
	island.Tick()

	// make sure npc consumed water
	assert.Equal(t, false, dw.hasItemTypeInInventory("drink"))

	// make sure thirst went down
	assert.Equal(t, true, dw.Thirst < oldThirst)
	assert.Equal(t, false, dw.isThirsty())
}

func TestSleep(t *testing.T) {

	prepareIsland()

	island.addNpcFromType("dwarf")
	assert.Equal(t, true, len(island.Spawns) == 1)
	dw := island.Spawns[0]

	// make npc tired
	dw.Tiredness = dw.tirednessCap() + 1
	island.Tick()

	// make sure npc planned action: find water
	assert.Equal(t, "sleep", dw.CurrentAction)

	oldTiredness := dw.Tiredness

	island.Tick()
	assert.Equal(t, true, dw.isSleeping())

	// progress until npc found food
	for i := 0; i < 10; i++ { // XXXX need to find action "find food".duration, from actions.yml
		island.Tick()
	}

	island.Tick()

	// make sure tiredness went down
	assert.Equal(t, true, dw.Tiredness < oldTiredness)
	assert.Equal(t, false, dw.isTired())
	assert.Equal(t, false, dw.isSleeping())
}
