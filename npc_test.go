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

func TestWithinRadiusOfType(t *testing.T) {

	prepareIsland()
	assert.Equal(t, true, len(island.Spawns) == 0)

	island.addNpcFromName("small fireplace", island.randomPointAboveWater())
	assert.Equal(t, true, len(island.Spawns) == 1)

	assert.Equal(t, 1, len(island.withinRadiusOfName("small fireplace", 0, island.Spawns[0].Position)))
	assert.Equal(t, 1, len(island.withinRadiusOfName("small fireplace", 30, island.Spawns[0].Position)))
	assert.Equal(t, 1, len(island.withinRadiusOfType("fireplace", 0, island.Spawns[0].Position)))
	assert.Equal(t, 1, len(island.withinRadiusOfType("fireplace", 30, island.Spawns[0].Position)))
}

func TestFindFoodAndEat(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())
	assert.Equal(t, true, len(island.Spawns) == 1)
	dw := island.Spawns[0]

	// make npc hungry
	dw.Hunger = dw.hungerCap() + 1
	island.Tick()

	// make sure planned action: find food
	assert.Equal(t, "find food", dw.CurrentAction.Name)
	assert.Equal(t, false, dw.hasItemTypeInInventory("food"))

	// progress until npc found food
	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	// make sure that npc has aged
	assert.Equal(t, true, dw.Age > 0)

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

func TestFindWaterAndDrink(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())
	assert.Equal(t, true, len(island.Spawns) == 1)
	dw := island.Spawns[0]

	// make npc thirsty
	dw.Thirst = dw.thirstCap() + 1
	island.Tick()

	// make sure npc planned action: find water
	assert.Equal(t, "find water", dw.CurrentAction.Name)
	assert.Equal(t, false, dw.hasItemTypeInInventory("drink"))

	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	// progress until npc found food
	for i := 0; i <= duration; i++ {
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

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())
	assert.Equal(t, true, len(island.Spawns) == 1)
	dw := island.Spawns[0]

	// make npc tired
	dw.Tiredness = dw.tirednessCap() + 1
	island.Tick()

	// make sure npc planned action: sleep
	assert.Equal(t, "sleep", dw.CurrentAction.Name)

	oldTiredness := dw.Tiredness

	island.Tick()
	assert.Equal(t, true, dw.isSleeping())

	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	// progress until npc wakes up
	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	// make sure tiredness went down
	assert.Equal(t, true, dw.Tiredness < oldTiredness)
	assert.Equal(t, false, dw.isTired())
	assert.Equal(t, false, dw.isSleeping())
}

func TestRabbitDigHole(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("rabbit", island.randomPointAboveWater())
	assert.Equal(t, true, len(island.Spawns) == 1)
	dw := island.Spawns[0]

	island.Tick()
	assert.Equal(t, "dig small hole", dw.CurrentAction.Name)

	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	assert.Equal(t, true, len(island.withinRadiusOfName("small hole", 0, dw.Position)) == 1)
	assert.Equal(t, true, len(island.withinRadiusOfType("shelter", 0, dw.Position)) == 1)
}

func TestBuildFireplace(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())

	// add nessecities, so they dont need to be built
	island.addNpcFromName("small shelter", island.Spawns[0].Position)
	island.addNpcFromName("farmland", island.Spawns[0].Position)

	assert.Equal(t, true, len(island.Spawns) == 3)
	dw := island.Spawns[0]

	island.Tick()
	assert.Equal(t, "build small fireplace", dw.CurrentAction.Name)

	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	assert.Equal(t, true, len(island.withinRadiusOfName("small fireplace", 0, dw.Position)) == 1)
	assert.Equal(t, true, len(island.withinRadiusOfType("fireplace", 0, dw.Position)) == 1)
}

func TestBuildShelter(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())

	// add nessecities, so they dont need to be built
	island.addNpcFromName("small fireplace", island.Spawns[0].Position)
	island.addNpcFromName("farmland", island.Spawns[0].Position)

	assert.Equal(t, true, len(island.Spawns) == 3)
	dw := island.Spawns[0]

	assert.Equal(t, true, len(island.withinRadiusOfType("fireplace", 0, dw.Position)) == 1)

	island.Tick()
	assert.Equal(t, false, dw.CurrentAction == nil)
	assert.Equal(t, "build small shelter", dw.CurrentAction.Name)

	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	assert.Equal(t, true, len(island.withinRadiusOfType("shelter", 0, dw.Position)) == 1)
}

func TestBuildFarmland(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())

	// add nessecities, so they dont need to be built
	island.addNpcFromName("small fireplace", island.Spawns[0].Position)
	island.addNpcFromName("small shelter", island.Spawns[0].Position)

	assert.Equal(t, true, len(island.Spawns) == 3)
	dw := island.Spawns[0]

	island.Tick()
	assert.Equal(t, false, dw.CurrentAction == nil)
	assert.Equal(t, "build farmland", dw.CurrentAction.Name)

	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	assert.Equal(t, true, len(island.withinRadiusOfName("farmland", 0, dw.Position)) == 1)
	assert.Equal(t, true, len(island.withinRadiusOfType("food producer", 0, dw.Position)) == 1)
}
