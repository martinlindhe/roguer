package rogue

import (
	"fmt"
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
	// clear spawns between tests
	island.Spawns = nil

	seed := int64(780)
	log.Info("Creating island with seed ", seed)
	generateIsland(seed, 200, 100)
	island.spawnGravel()
	island.spawnTrees()

	islandColImgFile, _ := os.Create("island_test.png")
	png.Encode(islandColImgFile, island.ColoredHeightMapAsImage())
}

func TestParseObjectsDefinition(t *testing.T) {
	_, err := parseObjectsDefinition("data/objs.yml")
	assert.Equal(t, nil, err)
}

func TestCanBuildAt(t *testing.T) {

	prepareIsland()
	spawnCnt := len(island.Spawns)
	assert.Equal(t, true, spawnCnt > 0)

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, spawnCnt+1, len(island.Spawns))

	assert.Equal(t, true, island.canBuildAt(island.Spawns[0].Position))

	island.addNpcFromName("small fireplace", island.Spawns[0].Position)

	assert.Equal(t, false, island.canBuildAt(island.Spawns[0].Position))
}

func TestFindFoodAndEat(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, true, len(island.Spawns) > 0)
	npc := island.Spawns[0]
	npc.addToInventory("small branch")

	assert.Equal(t, 1, len(npc.Inventory))

	// place food nearby
	nextTo, err := npc.Position.randomNearby()
	assert.Equal(t, nil, err)
	island.addNpcFromName("carrot", nextTo)

	// make npc hungry
	npc.Hunger = npc.hungerCap() + 1
	island.Tick()

	// make sure planned action: find food
	if npc.CurrentAction == nil {
		assert.Fail(t, "current action is nil")
		return
	}

	assert.Equal(t, "find food", npc.CurrentAction.Name)
	assert.Equal(t, false, npc.hasItemTypeInInventory("food"))

	// progress until npc found food
	duration := npc.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for {
		island.Tick()
		if len(npc.Inventory) == 2 {
			// food + small branch
			break
		}
	}

	// make sure that npc has aged
	assert.Equal(t, true, npc.Age.Current() > 0)
	// make sure food was found
	assert.Equal(t, true, npc.hasItemTypeInInventory("food"))
	assert.Equal(t, false, npc.hasItemTypeInInventory("drink"))

	oldHunger := npc.Hunger
	island.Tick()

	// make sure npc consumed food
	assert.Equal(t, false, npc.hasItemTypeInInventory("food"))

	// make sure hunger went down
	assert.Equal(t, true, npc.Hunger < oldHunger)
	assert.Equal(t, false, npc.isHungry())
}

func TestFindWaterAndDrink(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))
	npc := island.Spawns[0]
	npc.addToInventory("small branch")
	assert.Equal(t, 1, len(npc.Inventory)) // firewood

	// place water nearby
	nextTo, err := npc.Position.randomNearby()
	assert.Equal(t, nil, err)
	island.addNpcFromName("pouch of water", nextTo)

	// make npc thirsty
	npc.Thirst = npc.thirstCap() + 1
	island.Tick()

	if npc.CurrentAction == nil {
		assert.Fail(t, "current action is nil")
		return
	}

	// make sure npc planned action: find water
	assert.Equal(t, "find water", npc.CurrentAction.Name)
	assert.Equal(t, false, npc.hasItemTypeInInventory("drink"))

	duration := npc.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for {
		island.Tick()
		if len(npc.Inventory) == 2 {
			// water + small branch
			break
		}
	}

	assert.Equal(t, 2, len(npc.Inventory)) // water + firewood
	assert.Equal(t, true, npc.hasItemTypeInInventory("drink"))

	oldThirst := npc.Thirst
	island.Tick()

	// make sure npc consumed water
	assert.Equal(t, false, npc.hasItemTypeInInventory("drink"))

	// make sure thirst went down
	assert.Equal(t, true, npc.Thirst < oldThirst)
	assert.Equal(t, false, npc.isThirsty())
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
	npc := island.Spawns[0]

	// make dwarf wanna find firewood
	npc.Coldness = npc.coldnessCap() + 1
	assert.Equal(t, true, npc.isCold())

	// place firewood nearby
	nextTo, err := npc.Position.randomNearby()
	assert.Equal(t, nil, err)
	island.addNpcFromName("small branch", nextTo)

	// tick so npc decides to pick firewood
	island.Tick()

	assert.Equal(t, true, npc.hasPlanned("find fire wood"))
	// tick so fire wood is picked up
	island.Tick()

	// make sure it was picked up
	assert.Equal(t, true, npc.hasItemTypeInInventory("wood"))
}

func TestSleep(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))
	npc := island.Spawns[0]
	npc.addToInventory("small branch")

	// make npc tired
	npc.Tiredness = npc.tirednessCap() + 1
	island.Tick()

	// make sure npc planned action: sleep
	assert.Equal(t, true, npc.hasPlanned("sleep"))

	oldTiredness := npc.Tiredness

	island.Tick()
	// make sure npc is now sleeping
	assert.Equal(t, true, npc.isSleeping())
	if npc.CurrentAction == nil {
		assert.Fail(t, "current action is nil")
		return
	}

	duration := npc.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	// progress until npc wakes up
	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	// make sure tiredness went down
	assert.Equal(t, true, npc.Tiredness < oldTiredness)
	assert.Equal(t, false, npc.isTired())
	assert.Equal(t, false, npc.isSleeping())
}

func TestSleepAtShelter(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))
	npc := island.Spawns[0]

	// add nessecities nearby
	nextTo, err := island.Spawns[0].Position.randomNearby()
	assert.Equal(t, nil, err)
	island.addNpcFromName("small shelter", nextTo)

	// make npc tired
	npc.Tiredness = npc.tirednessCap() + 1
	island.Tick()

	assert.Equal(t, false, npc.isSleeping())

	// XXXXXxxx make sure npc moves to shelter and sleeps there

}

func TestRabbitDigHole(t *testing.T) {

	prepareIsland()
	assert.Equal(t, 0, len(island.Spawns))

	island.addNpcFromRace("rabbit", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))

	fmt.Println(island.Spawns)

	npc := island.Spawns[0]

	island.Tick()

	// make sure npc is now digging
	if npc.CurrentAction == nil {
		assert.Fail(t, "current action is nil")
		return
	}
	assert.Equal(t, "dig small hole", npc.CurrentAction.Name)

	duration := npc.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	assert.Equal(t, 1, len(npc.Position.spawnsByName("small hole", 0)))
	assert.Equal(t, 1, len(npc.Position.spawnsByType("burrow", 0)))

	// XXX make sure rabbit uses the burrow to sleep

	// make npc tired
	npc.Tiredness = npc.tirednessCap() + 1
	island.Tick()

	assert.Equal(t, false, npc.isSleeping())
	assert.Equal(t, true, npc.hasPlanned("sleep"))

	// make npc sleep
	island.Tick()
}

func TestBuildFireplace(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())

	npc := island.Spawns[0]
	npc.addToInventory("small branch")

	// add nessecities nearby, so they dont need to be built
	nextTo, err := island.Spawns[0].Position.randomNearby()
	assert.Equal(t, nil, err)

	// make sure nextTo is changed
	assert.Equal(t, false, island.Spawns[0].Position == nextTo)

	island.addNpcFromName("small shelter", nextTo)
	island.addNpcFromName("farmland", nextTo)
	island.addNpcFromName("apple tree", nextTo)

	assert.Equal(t, 1, len(npc.Position.spawnsByType("shelter", 30)))

	assert.Equal(t, true, len(island.Spawns) == 4)

	island.Tick()

	if npc.CurrentAction == nil {
		assert.Fail(t, "current action is nil")
		return
	}
	assert.Equal(t, "build small fireplace", npc.CurrentAction.Name)

	duration := npc.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	assert.Equal(t, 1, len(npc.Position.spawnsByName("small fireplace", 0)))
	assert.Equal(t, 1, len(npc.Position.spawnsByType("fireplace", 0)))
}

func TestBuildShelter(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())

	// add nessecities nearby, so they dont need to be built
	nextTo, err := island.Spawns[0].Position.randomNearby()
	assert.Equal(t, nil, err)

	// add nessecities, so they dont need to be built
	island.addNpcFromName("farmland", nextTo)
	island.addNpcFromName("small fireplace", nextTo)
	island.addNpcFromName("apple tree", nextTo)
	assert.Equal(t, 4, len(island.Spawns))

	npc := island.Spawns[0]
	npc.addToInventory("small branch")

	assert.Equal(t, 1, len(npc.Position.spawnsByType("fireplace", 2)))

	island.Tick()
	assert.Equal(t, true, npc.hasPlanned("build small shelter"))

	if npc.CurrentAction == nil {
		assert.Fail(t, "current action is nil")
		return
	}
	duration := npc.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	shelters := npc.Position.spawnsByType("shelter", 0)
	assert.Equal(t, 1, len(shelters))

	// make sure npc made this their home
	assert.Equal(t, true, npc.Home == &shelters[0])

	// make npc tired
	npc.Tiredness = npc.tirednessCap() + 1
	prevTiredness := npc.Tiredness
	island.Tick()

	// make sure npc planned action: sleep
	assert.Equal(t, true, npc.hasPlanned("sleep"))

	// make npc fall asleep
	island.Tick()

	// make sure they get the shelter bonus
	assert.Equal(t, 48, prevTiredness-npc.Tiredness) // 50 for the bonus, -2 for the ticks before starting to sleep
}

func TestBuildFarmland(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())

	fmt.Printf("xxxx dwarf at %s\n", island.Spawns[0].Position)

	// add nessecities nearby, so they dont need to be built
	nextTo, err := island.Spawns[0].Position.randomNearby()
	assert.Equal(t, nil, err)

	fmt.Printf("xxxx nextTo = %s\n", nextTo)

	// add nessecities, so they dont need to be built
	island.addNpcFromName("small fireplace", nextTo)
	assert.Equal(t, 1, len(nextTo.spawnsByType("fireplace", 30)))

	island.addNpcFromName("small shelter", nextTo)
	island.addNpcFromName("apple tree", nextTo)
	island.addNpcFromName("cooking pit", nextTo)
	island.addNpcFromName("small hut", nextTo)
	assert.Equal(t, 6, len(island.Spawns))

	npc := island.Spawns[0]
	npc.addToInventory("small branch")

	island.Tick()

	if npc.CurrentAction == nil {
		assert.Fail(t, "current action is nil")
		return
	}
	assert.Equal(t, "build farmland", npc.CurrentAction.Name)

	duration := npc.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		island.Tick()
	}

	assert.Equal(t, 1, len(npc.Position.spawnsByName("farmland", 0)))
	assert.Equal(t, 1, len(npc.Position.spawnsByType("food producer", 0)))
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

	npc := island.Spawns[0]
	npc.Age.Set(npc.ageCap() + 1)

	island.Tick()

	// dwarf should have died of old age
	assert.Equal(t, 0, len(island.Spawns))
}

func TestSpawnGravel(t *testing.T) {

	prepareIsland()

	assert.Equal(t, 0, len(island.Spawns))
	island.spawnGravel()
	assert.Equal(t, true, len(island.Spawns) > 1)
}

func TestNpcMovesToFireplace(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))

	npc := island.Spawns[0]
	npc.addToInventory("small branch")

	nextTo, err := island.Spawns[0].Position.randomNearby()
	assert.Equal(t, nil, err)

	assert.Equal(t, false, npc.Position == nextTo)

	island.addNpcFromName("small fireplace", nextTo)

	// add nessecities, so they dont need to be built
	island.addNpcFromName("small shelter", nextTo)
	island.addNpcFromName("apple tree", nextTo)
	island.addNpcFromName("farmland", nextTo)
	assert.Equal(t, 5, len(island.Spawns))

	// make dwarf wanna move to shelter
	npc.Coldness = npc.coldnessCap() + 1
	assert.Equal(t, true, npc.isCold())

	// let them travel to destination
	for {
		island.Tick()
		if npc.Position.intMatches(&nextTo) {
			break
		}
	}

	assert.Equal(t, true, npc.Position.intMatches(&nextTo))

	// let npc start the fire
	island.Tick()

	// let fire burn
	island.Tick()

	// let them get warm by the fire
	assert.Equal(t, false, npc.isCold())
}

func TestNpcFindFirewoodThenMovesToFireplace(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))

	npc := island.Spawns[0]

	// NOTE: similar to TestNpcMovesToFireplace, but now also make sure dwarf finds a firewood

	nextTo, err := npc.Position.randomNearby()
	assert.Equal(t, nil, err)

	assert.Equal(t, false, npc.Position == nextTo)

	island.addNpcFromName("farmland", island.Spawns[0].Position)
	island.addNpcFromName("small fireplace", nextTo)
	island.addNpcFromName("small shelter", nextTo)
	island.addNpcFromName("cooking pit", nextTo)
	island.addNpcFromName("farmland", nextTo)
	island.addNpcFromName("apple tree", nextTo)
	assert.Equal(t, 7, len(island.Spawns))

	nextTo2, err := npc.Position.randomNearby()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, nextTo.Equals(nextTo2))

	island.addNpcFromName("branch", nextTo)
	assert.Equal(t, 8, len(island.Spawns))

	// make dwarf wanna move to shelter
	npc.Coldness = npc.coldnessCap() + 1
	assert.Equal(t, true, npc.isCold())

	assert.Equal(t, false, npc.hasItemTypeInInventory("wood"))

	// make dwarf plan to get firewood
	island.Tick()

	assert.Equal(t, false, npc.hasItemTypeInInventory("wood"))
	assert.Equal(t, true, npc.hasPlanned("find fire wood"))

	for {
		island.Tick()

		// have dwarf find branch
		if npc.hasItemTypeInInventory("wood") {
			break
		}
	}

	assert.Equal(t, true, npc.hasItemTypeInInventory("wood"))

	// make dwarf plan to get warm by fireplace
	island.Tick()
	assert.Equal(t, true, npc.hasPlannedType("wait"))

	// let dwarf get warm
	island.Tick()

	// let them get warm by the fire
	assert.Equal(t, false, npc.isCold())
}

func TestBuildCookingPit(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())

	assert.Equal(t, 1, len(island.Spawns))
	npc := island.Spawns[0]

	// add nessecities nearby, so they dont need to be built
	nextTo, err := npc.Position.randomNearby()
	assert.Equal(t, nil, err)

	island.addNpcFromName("small fireplace", nextTo)
	island.addNpcFromName("small shelter", nextTo)
	island.addNpcFromName("apple tree", nextTo)
	island.addNpcFromName("farmland", nextTo)
	assert.Equal(t, 5, len(island.Spawns))

	// make sure npc decides to build cooking pit
	island.Tick()
	assert.Equal(t, true, npc.hasPlanned("build cooking pit"))

	// wait until done
	for {
		island.Tick()

		if !npc.isPerforming("build cooking pit") {
			break
		}
	}

	assert.Equal(t, 1, len(npc.Position.spawnsByName("cooking pit", 0)))
}

func TestBuildSmallHut(t *testing.T) {

	prepareIsland()

	island.addNpcFromRace("dwarf", island.RandomPointAboveWater())
	assert.Equal(t, 1, len(island.Spawns))
	npc := island.Spawns[0]

	// add nessecities nearby, so they dont need to be built
	nextTo, err := npc.Position.randomNearby()
	assert.Equal(t, nil, err)
	island.addNpcFromName("small fireplace", nextTo)

	home := island.addNpcFromName("small shelter", nextTo)
	npc.Home = &home

	island.addNpcFromName("apple tree", nextTo)
	island.addNpcFromName("farmland", nextTo)
	island.addNpcFromName("cooking pit", nextTo)
	assert.Equal(t, 6, len(island.Spawns))

	// make sure npc decides to build cooking pit
	island.Tick()
	assert.Equal(t, true, npc.hasPlanned("build small hut"))

	// wait until done
	for {
		island.Tick()

		if !npc.isPerforming("build small hut") {
			break
		}
	}

	assert.Equal(t, 1, len(npc.Position.spawnsByName("small hut", 0)))
}
