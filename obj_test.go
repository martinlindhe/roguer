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

	pos := island.randomPointAboveWater()

	assert.Equal(t, 0, len(island.withinRadiusOfName("small fireplace", 0, pos)))
	assert.Equal(t, 0, len(island.withinRadiusOfName("small fireplace", 30, pos)))
	assert.Equal(t, 0, len(island.withinRadiusOfType("fireplace", 0, pos)))
	assert.Equal(t, 0, len(island.withinRadiusOfType("fireplace", 30, pos)))

	island.addNpcFromName("small fireplace", pos)
	assert.Equal(t, true, len(island.Spawns) == 1)

	assert.Equal(t, 1, len(island.withinRadiusOfName("small fireplace", 0, pos)))
	assert.Equal(t, 1, len(island.withinRadiusOfName("small fireplace", 30, pos)))
	assert.Equal(t, 1, len(island.withinRadiusOfType("fireplace", 0, pos)))
	assert.Equal(t, 1, len(island.withinRadiusOfType("fireplace", 30, pos)))

	pos2 := pos
	pos2.Y++

	assert.Equal(t, 0, len(island.withinRadiusOfName("small fireplace", 0, pos2)))
	assert.Equal(t, 1, len(island.withinRadiusOfName("small fireplace", 1, pos2)))
	assert.Equal(t, 1, len(island.withinRadiusOfName("small fireplace", 30, pos2)))
	assert.Equal(t, 0, len(island.withinRadiusOfType("fireplace", 0, pos2)))
	assert.Equal(t, 1, len(island.withinRadiusOfType("fireplace", 1, pos2)))
	assert.Equal(t, 1, len(island.withinRadiusOfType("fireplace", 30, pos2)))
}

func TestCanBuildAt(t *testing.T) {

	prepareIsland()
	assert.Equal(t, 0, len(island.Spawns))

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))

	assert.Equal(t, true, island.canBuildAt(island.Spawns[0].Position))

	island.addNpcFromName("small fireplace", island.Spawns[0].Position)

	assert.Equal(t, false, island.canBuildAt(island.Spawns[0].Position))
}

func TestFindFoodAndEat(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))
	dw := island.Spawns[0]
	dw.addToInventory("small branch")
	assert.Equal(t, 1, len(dw.Inventory))

	// place food nearby
	nextTo := dw.Position
	nextTo.Y += 2
	island.addNpcFromName("carrot", nextTo)

	// make npc hungry
	dw.Hunger = dw.hungerCap() + 1
	island.Tick()

	// make sure planned action: find food
	assert.Equal(t, "find food", dw.CurrentAction.Name)
	assert.Equal(t, false, dw.hasItemTypeInInventory("food"))

	// progress until npc found food
	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for {
		island.Tick()
		if len(dw.Inventory) == 2 {
			// food + small branch
			break
		}
	}

	// make sure that npc has aged
	assert.Equal(t, true, dw.Age > 0)
	// make sure food was found
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
	assert.Equal(t, 1, len(island.Spawns))
	dw := island.Spawns[0]
	dw.addToInventory("small branch")
	assert.Equal(t, 1, len(dw.Inventory)) // firewood

	// place water nearby
	nextTo := dw.Position
	nextTo.Y -= 2
	island.addNpcFromName("pouch of water", nextTo)

	// make npc thirsty
	dw.Thirst = dw.thirstCap() + 1
	island.Tick()

	// make sure npc planned action: find water
	assert.Equal(t, "find water", dw.CurrentAction.Name)
	assert.Equal(t, false, dw.hasItemTypeInInventory("drink"))

	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for {
		island.Tick()
		if len(dw.Inventory) == 2 {
			// water + small branch
			break
		}
	}

	assert.Equal(t, 2, len(dw.Inventory)) // water + firewood
	assert.Equal(t, true, dw.hasItemTypeInInventory("drink"))

	oldThirst := dw.Thirst
	island.Tick()

	// make sure npc consumed water
	assert.Equal(t, false, dw.hasItemTypeInInventory("drink"))

	// make sure thirst went down
	assert.Equal(t, true, dw.Thirst < oldThirst)
	assert.Equal(t, false, dw.isThirsty())
}

func TestFindFirewood(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())

	// add nessecities, so they dont need to be built
	island.addNpcFromName("small fireplace", island.Spawns[0].Position)
	island.addNpcFromName("small shelter", island.Spawns[0].Position)
	island.addNpcFromName("apple tree", island.Spawns[0].Position)
	island.addNpcFromName("farmland", island.Spawns[0].Position)

	assert.Equal(t, 5, len(island.Spawns))
	dw := island.Spawns[0]

	// make dwarf wanna find firewood
	dw.Coldness = dw.coldnessCap() + 1
	assert.Equal(t, true, dw.isCold())

	// place firewood nearby
	nextTo := dw.Position
	nextTo.Y += 1
	island.addNpcFromName("small branch", nextTo)

	// tick so npc decides to pick firewood
	island.Tick()

	assert.Equal(t, true, dw.hasPlanned("find fire wood"))
	// tick so fire wood is picked up
	island.Tick()

	// make sure it was picked up
	assert.Equal(t, true, dw.hasItemTypeInInventory("wood"))
}

func TestSleep(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))
	dw := island.Spawns[0]
	dw.addToInventory("small branch")

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
	assert.Equal(t, 1, len(island.Spawns))
	dw := island.Spawns[0]
	dw.addToInventory("small branch")

	island.Tick()
	assert.Equal(t, "dig small hole", dw.CurrentAction.Name)

	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	assert.Equal(t, 1, len(island.withinRadiusOfName("small hole", 0, dw.Position)))
	assert.Equal(t, 1, len(island.withinRadiusOfType("shelter", 0, dw.Position)))
}

func TestBuildFireplace(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())

	dw := island.Spawns[0]
	dw.addToInventory("small branch")

	// add nessecities nearby, so they dont need to be built
	nextTo := island.Spawns[0].Position
	nextTo.Y++
	// make sure nextTo is changed
	assert.Equal(t, false, island.Spawns[0].Position.Y == nextTo.Y)

	island.addNpcFromName("small shelter", nextTo)
	island.addNpcFromName("farmland", nextTo)
	island.addNpcFromName("apple tree", nextTo)

	assert.Equal(t, 1, len(island.withinRadiusOfType("shelter", 30, dw.Position)))

	assert.Equal(t, true, len(island.Spawns) == 4)

	island.Tick()

	assert.Equal(t, false, dw.CurrentAction == nil)
	assert.Equal(t, "build small fireplace", dw.CurrentAction.Name)

	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	assert.Equal(t, 1, len(island.withinRadiusOfName("small fireplace", 0, dw.Position)))
	assert.Equal(t, 1, len(island.withinRadiusOfType("fireplace", 0, dw.Position)))
}

func TestBuildShelter(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())

	// add nessecities nearby, so they dont need to be built
	nextTo := island.Spawns[0].Position
	nextTo.Y++

	// add nessecities, so they dont need to be built
	island.addNpcFromName("small fireplace", nextTo)
	island.addNpcFromName("farmland", nextTo)
	island.addNpcFromName("apple tree", nextTo)
	assert.Equal(t, 4, len(island.Spawns))

	dw := island.Spawns[0]
	dw.addToInventory("small branch")

	assert.Equal(t, 1, len(island.withinRadiusOfType("fireplace", 1, dw.Position)))

	island.Tick()
	assert.Equal(t, false, dw.CurrentAction == nil)
	assert.Equal(t, "build small shelter", dw.CurrentAction.Name)

	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	assert.Equal(t, 1, len(island.withinRadiusOfType("shelter", 0, dw.Position)))
}

func TestBuildFarmland(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())

	// add nessecities nearby, so they dont need to be built
	nextTo := island.Spawns[0].Position
	nextTo.Y++

	// add nessecities, so they dont need to be built
	island.addNpcFromName("small fireplace", nextTo)
	island.addNpcFromName("small shelter", nextTo)
	island.addNpcFromName("apple tree", nextTo)
	assert.Equal(t, 4, len(island.Spawns))

	dw := island.Spawns[0]
	dw.addToInventory("small branch")

	island.Tick()
	assert.Equal(t, false, dw.CurrentAction == nil)
	assert.Equal(t, "build farmland", dw.CurrentAction.Name)

	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	assert.Equal(t, 1, len(island.withinRadiusOfName("farmland", 0, dw.Position)))
	assert.Equal(t, 1, len(island.withinRadiusOfType("food producer", 0, dw.Position)))
}

func TestTree(t *testing.T) {

	prepareIsland()

	pos := island.randomPointAboveWater()
	island.addNpcFromName("oak tree", pos)

	assert.Equal(t, 1, len(island.Spawns))

	// let tree drop some spawns
	for i := 0; i <= 100; i++ {
		island.Tick()
	}

	assert.Equal(t, true, len(island.Spawns) > 1)
	assert.Equal(t, true, len(island.withinRadiusOfType("wood", 0, pos)) > 0)
}

func TestNpcDiesOfOldAge(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))

	dw := island.Spawns[0]
	dw.Age = dw.ageCap() + 1

	island.Tick()

	// dwarf should have died of old age
	assert.Equal(t, 0, len(island.Spawns))
}

func TestSpawnGravel(t *testing.T) {

	prepareIsland()

	assert.Equal(t, 0, len(island.Spawns))
	island.spawnGravel()
	assert.Equal(t, true, len(island.Spawns) > 10000)
}

func TestNpcMovesToFireplace(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))

	dw := island.Spawns[0]
	dw.addToInventory("small branch")

	nextTo := island.Spawns[0].Position
	nextTo.X -= 8
	nextTo.Y += 20

	assert.Equal(t, false, dw.Position == nextTo)

	island.addNpcFromName("small fireplace", nextTo)

	// add nessecities, so they dont need to be built
	island.addNpcFromName("small shelter", nextTo)
	island.addNpcFromName("apple tree", nextTo)
	island.addNpcFromName("farmland", nextTo)
	assert.Equal(t, 5, len(island.Spawns))

	// make dwarf wanna move to shelter
	dw.Coldness = dw.coldnessCap() + 1
	assert.Equal(t, true, dw.isCold())

	// let them travel to destination
	for {
		island.Tick()
		if dw.Position.intMatches(&nextTo) {
			break
		}
	}

	assert.Equal(t, true, dw.Position.intMatches(&nextTo))

	// let npc start the fire
	island.Tick()

	// let fire burn
	island.Tick()

	// let them get warm by the fire
	assert.Equal(t, false, dw.isCold())
}

func TestNpcFindFirewoodThenMovesToFireplace(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.randomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))

	dw := island.Spawns[0]

	// NOTE: similar to TestNpcMovesToFireplace, but now also make sure dwarf finds a firewood

	nextTo := island.Spawns[0].Position
	nextTo.X -= 3
	nextTo.Y += 3

	assert.Equal(t, false, dw.Position == nextTo)

	island.addNpcFromName("farmland", island.Spawns[0].Position)
	island.addNpcFromName("small fireplace", nextTo)
	island.addNpcFromName("small shelter", nextTo)
	island.addNpcFromName("farmland", nextTo)
	island.addNpcFromName("apple tree", nextTo)
	assert.Equal(t, 6, len(island.Spawns))

	nextTo2 := island.Spawns[0].Position
	nextTo2.X += 2
	island.addNpcFromName("branch", nextTo)
	assert.Equal(t, 7, len(island.Spawns))

	// make dwarf wanna move to shelter
	dw.Coldness = dw.coldnessCap() + 1
	assert.Equal(t, true, dw.isCold())

	for {
		island.Tick()

		// have dwarf find branch
		if dw.hasItemTypeInInventory("wood") {
			break
		}
	}

	// let them travel to destination
	for {
		island.Tick()

		// have dwarf find branch
		if dw.Position.intMatches(&nextTo) {
			break
		}
	}

	assert.Equal(t, true, dw.Position.intMatches(&nextTo))

	// let npc start the fire, wait and get warmed up
	island.Tick()

	// let them get warm by the fire
	assert.Equal(t, false, dw.isCold())
}
