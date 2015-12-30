package rogue

import (
	"image/png"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func BenchmarkGenerateIsland(b *testing.B) {
	for n := 0; n < b.N; n++ {

		generateIsland(666, 220, 140)
	}
}

func prepareIsland() {
	if island == nil {
		seed := int64(779)
		log.Printf("Creating island with seed %d", seed)
		island = generateIsland(seed, 200, 100)
		island.spawnGravel()
		island.spawnTrees()

		islandColImgFile, _ := os.Create("island_test.png")
		png.Encode(islandColImgFile, island.ColoredHeightMapAsImage())
	}

	// clear spawns between tests
	island.Spawns = nil
}

func TestParseObjectsDefinition(t *testing.T) {
	_, err := parseObjectsDefinition("data/objs.yml")
	assert.Equal(t, nil, err)
}

func TestCanBuildAt(t *testing.T) {

	prepareIsland()
	assert.Equal(t, 0, len(island.Spawns))

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))

	assert.Equal(t, true, island.canBuildAt(island.Spawns[0].Position))

	island.addNpcFromName("small fireplace", island.Spawns[0].Position)

	assert.Equal(t, false, island.canBuildAt(island.Spawns[0].Position))
}

func TestFindFoodAndEat(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))
	dw := island.Spawns[0]
	dw.addToInventory("small branch")
	assert.Equal(t, 1, len(dw.Inventory))

	// place food nearby
	nextTo, _ := dw.Position.randomNearby()
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

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))
	dw := island.Spawns[0]
	dw.addToInventory("small branch")
	assert.Equal(t, 1, len(dw.Inventory)) // firewood

	// place water nearby
	nextTo, _ := dw.Position.randomNearby()
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

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())

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
	nextTo, _ := dw.Position.randomNearby()
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

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))
	dw := island.Spawns[0]
	dw.addToInventory("small branch")

	// make npc tired
	dw.Tiredness = dw.tirednessCap() + 1
	island.Tick()

	// make sure npc planned action: sleep
	assert.Equal(t, true, dw.hasPlanned("sleep"))

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

func TestSleepAtShelter(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))
	dw := island.Spawns[0]

	// add nessecities nearby
	nextTo, _ := island.Spawns[0].Position.randomNearby()
	island.addNpcFromName("small shelter", nextTo)

	// make npc tired
	dw.Tiredness = dw.tirednessCap() + 1
	island.Tick()

	assert.Equal(t, false, dw.isSleeping())

	// XXXXXxxx make sure npc moves to shelter and sleeps there

}

func TestRabbitDigHole(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("rabbit", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))
	ra := island.Spawns[0]

	island.Tick()
	assert.Equal(t, "dig small hole", ra.CurrentAction.Name)

	duration := ra.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	assert.Equal(t, 1, len(ra.Position.spawnsByName("small hole", 0)))
	assert.Equal(t, 1, len(ra.Position.spawnsByType("burrow", 0)))

	// XXX make sure rabbit uses the burrow to sleep

	// make npc tired
	ra.Tiredness = ra.tirednessCap() + 1
	island.Tick()

	assert.Equal(t, false, ra.isSleeping())
	assert.Equal(t, true, ra.hasPlanned("sleep"))

	// make npc sleep
	island.Tick()
}

func TestBuildFireplace(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())

	dw := island.Spawns[0]
	dw.addToInventory("small branch")

	// add nessecities nearby, so they dont need to be built
	nextTo, _ := island.Spawns[0].Position.randomNearby()

	// make sure nextTo is changed
	assert.Equal(t, false, island.Spawns[0].Position == nextTo)

	island.addNpcFromName("small shelter", nextTo)
	island.addNpcFromName("farmland", nextTo)
	island.addNpcFromName("apple tree", nextTo)

	assert.Equal(t, 1, len(dw.Position.spawnsByType("shelter", 30)))

	assert.Equal(t, true, len(island.Spawns) == 4)

	island.Tick()

	assert.Equal(t, false, dw.CurrentAction == nil)
	assert.Equal(t, "build small fireplace", dw.CurrentAction.Name)

	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	assert.Equal(t, 1, len(dw.Position.spawnsByName("small fireplace", 0)))
	assert.Equal(t, 1, len(dw.Position.spawnsByType("fireplace", 0)))
}

func TestBuildShelter(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())

	// add nessecities nearby, so they dont need to be built
	nextTo, _ := island.Spawns[0].Position.randomNearby()

	// add nessecities, so they dont need to be built
	island.addNpcFromName("farmland", nextTo)
	island.addNpcFromName("small fireplace", nextTo)
	island.addNpcFromName("apple tree", nextTo)
	assert.Equal(t, 4, len(island.Spawns))

	dw := island.Spawns[0]
	dw.addToInventory("small branch")

	assert.Equal(t, 1, len(dw.Position.spawnsByType("fireplace", 2)))

	island.Tick()
	assert.Equal(t, true, dw.hasPlanned("build small shelter"))

	duration := dw.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	shelters := dw.Position.spawnsByType("shelter", 0)
	assert.Equal(t, 1, len(shelters))

	// make sure npc made this their home
	assert.Equal(t, true, dw.Home == shelters[0])

	// make npc tired
	dw.Tiredness = dw.tirednessCap() + 1
	prevTiredness := dw.Tiredness
	island.Tick()

	// make sure npc planned action: sleep
	assert.Equal(t, true, dw.hasPlanned("sleep"))

	// make npc fall asleep
	island.Tick()

	// make sure they get the shelter bonus
	assert.Equal(t, 48, prevTiredness-dw.Tiredness) // 50 for the bonus, -2 for the ticks before starting to sleep
}

func TestBuildFarmland(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())

	// add nessecities nearby, so they dont need to be built
	nextTo, _ := island.Spawns[0].Position.randomNearby()

	// add nessecities, so they dont need to be built
	island.addNpcFromName("small fireplace", nextTo)
	island.addNpcFromName("small shelter", nextTo)
	island.addNpcFromName("apple tree", nextTo)
	island.addNpcFromName("cooking pit", nextTo)
	island.addNpcFromName("small hut", nextTo)
	assert.Equal(t, 6, len(island.Spawns))

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

	assert.Equal(t, 1, len(dw.Position.spawnsByName("farmland", 0)))
	assert.Equal(t, 1, len(dw.Position.spawnsByType("food producer", 0)))
}

func TestTree(t *testing.T) {

	prepareIsland()

	pos := island.RandomPointAboveWater()
	island.addNpcFromName("oak tree", pos)

	assert.Equal(t, 1, len(island.Spawns))

	// let tree drop some spawns
	for i := 0; i <= 100; i++ {
		island.Tick()
	}

	assert.Equal(t, true, len(island.Spawns) > 1)
	assert.Equal(t, true, len(pos.spawnsByType("wood", 1)) > 0)
}

func TestNpcDiesOfOldAge(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
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
	assert.Equal(t, true, len(island.Spawns) > 1000)
}

func TestNpcMovesToFireplace(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))

	dw := island.Spawns[0]
	dw.addToInventory("small branch")

	nextTo, _ := island.Spawns[0].Position.randomNearby()

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

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))

	dw := island.Spawns[0]

	// NOTE: similar to TestNpcMovesToFireplace, but now also make sure dwarf finds a firewood

	nextTo, _ := dw.Position.randomNearby()

	assert.Equal(t, false, dw.Position == nextTo)

	island.addNpcFromName("farmland", island.Spawns[0].Position)
	island.addNpcFromName("small fireplace", nextTo)
	island.addNpcFromName("small shelter", nextTo)
	island.addNpcFromName("cooking pit", nextTo)
	island.addNpcFromName("farmland", nextTo)
	island.addNpcFromName("apple tree", nextTo)
	assert.Equal(t, 7, len(island.Spawns))

	nextTo2, _ := dw.Position.randomNearby()
	assert.Equal(t, false, nextTo == nextTo2)

	island.addNpcFromName("branch", nextTo)
	assert.Equal(t, 8, len(island.Spawns))

	// make dwarf wanna move to shelter
	dw.Coldness = dw.coldnessCap() + 1
	assert.Equal(t, true, dw.isCold())

	assert.Equal(t, false, dw.hasItemTypeInInventory("wood"))

	// make dwarf plan to get firewood
	island.Tick()

	assert.Equal(t, false, dw.hasItemTypeInInventory("wood"))
	assert.Equal(t, true, dw.hasPlanned("find fire wood"))

	for {
		island.Tick()

		// have dwarf find branch
		if dw.hasItemTypeInInventory("wood") {
			break
		}
	}

	assert.Equal(t, true, dw.hasItemTypeInInventory("wood"))

	// make dwarf plan to get warm by fireplace
	island.Tick()
	assert.Equal(t, true, dw.hasPlannedType("wait"))

	// let dwarf get warm
	island.Tick()

	// let them get warm by the fire
	assert.Equal(t, false, dw.isCold())
}

func TestBuildCookingPit(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))
	dw := island.Spawns[0]

	// add nessecities nearby, so they dont need to be built
	nextTo, _ := dw.Position.randomNearby()
	island.addNpcFromName("small fireplace", nextTo)
	island.addNpcFromName("small shelter", nextTo)
	island.addNpcFromName("apple tree", nextTo)
	island.addNpcFromName("farmland", nextTo)
	assert.Equal(t, 5, len(island.Spawns))

	// make sure npc decides to build cooking pit
	island.Tick()
	assert.Equal(t, true, dw.hasPlanned("build cooking pit"))

	// wait until done
	for {
		island.Tick()

		if !dw.isPerforming("build cooking pit") {
			break
		}
	}

	assert.Equal(t, 1, len(dw.Position.spawnsByName("cooking pit", 0)))
}

func TestBuildSmallHut(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))
	dw := island.Spawns[0]

	// add nessecities nearby, so they dont need to be built
	nextTo, _ := dw.Position.randomNearby()
	island.addNpcFromName("small fireplace", nextTo)

	dw.Home = island.addNpcFromName("small shelter", nextTo)

	island.addNpcFromName("apple tree", nextTo)
	island.addNpcFromName("farmland", nextTo)
	island.addNpcFromName("cooking pit", nextTo)
	assert.Equal(t, 6, len(island.Spawns))

	// make sure npc decides to build cooking pit
	island.Tick()
	assert.Equal(t, true, dw.hasPlanned("build small hut"))

	// wait until done
	for {
		island.Tick()

		if !dw.isPerforming("build small hut") {
			break
		}
	}

	assert.Equal(t, 1, len(dw.Position.spawnsByName("small hut", 0)))
}
