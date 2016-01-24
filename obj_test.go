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

func testNewGame() *Game {

	g := Game{
		Island: NewIsland(),
	}

	islandColImgFile, _ := os.Create("island_test.png")
	png.Encode(islandColImgFile, g.Island.ColoredHeightMapAsImage())

	return &g
}

func TestParseObjectsDefinition(t *testing.T) {
	_, err := parseObjectsDefinition("data/objs.yml")
	assert.Equal(t, nil, err)
}

func TestCanBuildAt(t *testing.T) {

	g := testNewGame()

	spawnCnt := len(g.Island.Spawns)
	assert.Equal(t, 0, spawnCnt)

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())
	assert.Equal(t, spawnCnt+1, len(g.Island.Spawns))

	assert.Equal(t, true, g.Island.canBuildAt(g.Island.Spawns[0].Position))

	g.Island.addNpcFromName("small fireplace", g.Island.Spawns[0].Position)

	assert.Equal(t, false, g.Island.canBuildAt(g.Island.Spawns[0].Position))
}

func TestFindFoodAndEat(t *testing.T) {

	g := testNewGame()

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())
	assert.Equal(t, true, len(g.Island.Spawns) > 0)
	npc := g.Island.Spawns[0]
	npc.addToInventory("small branch")
	assert.Equal(t, 1, len(npc.Inventory))

	// add nessecities, so they dont need to be built
	g.Island.addNpcFromName("small fireplace", npc.Position)
	g.Island.addNpcFromName("small shelter", npc.Position)
	g.Island.addNpcFromName("apple tree", npc.Position)
	g.Island.addNpcFromName("farmland", npc.Position)

	// place food nearby
	nextTo, err := npc.randomNearby()
	assert.Equal(t, nil, err)
	g.Island.addNpcFromName("carrot", nextTo)

	// make npc hungry
	npc.Hunger = npc.hungerCap() + 1
	g.Island.Tick()

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
		g.Island.Tick()
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
	g.Island.Tick()

	// make sure npc consumed food
	assert.Equal(t, false, npc.hasItemTypeInInventory("food"))

	// make sure hunger went down
	assert.Equal(t, true, npc.Hunger < oldHunger)
	assert.Equal(t, false, npc.isHungry())
}

func TestFindWaterAndDrink(t *testing.T) {

	g := testNewGame()

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())
	assert.Equal(t, 1, len(g.Island.Spawns))
	npc := g.Island.Spawns[0]
	npc.addToInventory("small branch")
	assert.Equal(t, 1, len(npc.Inventory)) // firewood

	// place water nearby
	nextTo, err := npc.randomNearby()
	assert.Equal(t, nil, err)
	g.Island.addNpcFromName("pouch of water", nextTo)

	// make npc thirsty
	npc.Thirst = npc.thirstCap() + 1
	g.Island.Tick()

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
		g.Island.Tick()
		if len(npc.Inventory) == 2 {
			// water + small branch
			break
		}
	}

	assert.Equal(t, 2, len(npc.Inventory)) // water + firewood
	assert.Equal(t, true, npc.hasItemTypeInInventory("drink"))

	oldThirst := npc.Thirst
	g.Island.Tick()

	// make sure npc consumed water
	assert.Equal(t, false, npc.hasItemTypeInInventory("drink"))

	// make sure thirst went down
	assert.Equal(t, true, npc.Thirst < oldThirst)
	assert.Equal(t, false, npc.isThirsty())
}

func TestFindFirewood(t *testing.T) {

	g := testNewGame()

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())

	// add nessecities, so they dont need to be built
	g.Island.addNpcFromName("small fireplace", g.Island.Spawns[0].Position)
	g.Island.addNpcFromName("small shelter", g.Island.Spawns[0].Position)
	g.Island.addNpcFromName("apple tree", g.Island.Spawns[0].Position)
	g.Island.addNpcFromName("farmland", g.Island.Spawns[0].Position)

	assert.Equal(t, 5, len(g.Island.Spawns))
	npc := g.Island.Spawns[0]

	// make dwarf wanna find firewood
	npc.Coldness = npc.coldnessCap() + 1
	assert.Equal(t, true, npc.isCold())

	// place firewood nearby
	nextTo, err := npc.randomNearby()
	assert.Equal(t, nil, err)
	g.Island.addNpcFromName("small branch", nextTo)

	// tick so npc decides to pick firewood
	g.Island.Tick()

	assert.Equal(t, true, npc.hasPlanned("find fire wood"))
	// tick so fire wood is picked up
	g.Island.Tick()

	// make sure it was picked up
	assert.Equal(t, true, npc.hasItemTypeInInventory("wood"))
}

func TestSleep(t *testing.T) {

	g := testNewGame()

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())
	assert.Equal(t, 1, len(g.Island.Spawns))
	npc := g.Island.Spawns[0]
	npc.addToInventory("small branch")

	// make npc tired
	npc.Tiredness = npc.tirednessCap() + 1
	g.Island.Tick()

	// make sure npc planned action: sleep
	assert.Equal(t, true, npc.hasPlanned("sleep"))

	oldTiredness := npc.Tiredness

	g.Island.Tick()
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
		g.Island.Tick()
	}

	// make sure tiredness went down
	assert.Equal(t, true, npc.Tiredness < oldTiredness)
	assert.Equal(t, false, npc.isTired())
	assert.Equal(t, false, npc.isSleeping())
}

func TestSleepAtShelter(t *testing.T) {

	g := testNewGame()

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())
	assert.Equal(t, 1, len(g.Island.Spawns))
	npc := g.Island.Spawns[0]

	// add nessecities nearby
	nextTo, err := npc.randomNearby()
	assert.Equal(t, nil, err)
	g.Island.addNpcFromName("small shelter", nextTo)

	// make npc tired
	npc.Tiredness = npc.tirednessCap() + 1
	g.Island.Tick()

	assert.Equal(t, false, npc.isSleeping())

	// XXXXXxxx make sure npc moves to shelter and sleeps there

}

func TestRabbitDigHole(t *testing.T) {

	g := testNewGame()

	assert.Equal(t, 0, len(g.Island.Spawns))

	g.Island.addNpcFromRace("rabbit", g.Island.RandomPointAboveWater())
	assert.Equal(t, 1, len(g.Island.Spawns))

	npc := g.Island.Spawns[0]

	g.Island.Tick()

	// make sure npc is now digging
	if npc.CurrentAction == nil {
		assert.Fail(t, "current action is nil")
		return
	}
	assert.Equal(t, "dig small hole", npc.CurrentAction.Name)

	duration := npc.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		g.Island.Tick()
	}

	assert.Equal(t, 1, len(npc.spawnsByName("small hole", 0)))
	assert.Equal(t, 1, len(npc.spawnsByType("burrow", 0)))

	// XXX make sure rabbit uses the burrow to sleep

	// make npc tired
	npc.Tiredness = npc.tirednessCap() + 1
	g.Island.Tick()

	assert.Equal(t, false, npc.isSleeping())
	assert.Equal(t, true, npc.hasPlanned("sleep"))

	// make npc sleep
	g.Island.Tick()
}

func TestBuildFireplace(t *testing.T) {

	g := testNewGame()

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())

	npc := g.Island.Spawns[0]
	npc.addToInventory("small branch")

	// add nessecities nearby, so they dont need to be built
	nextTo, err := npc.randomNearby()
	assert.Equal(t, nil, err)

	// make sure nextTo is changed
	assert.Equal(t, false, npc.Position == nextTo)

	g.Island.addNpcFromName("small shelter", nextTo)
	g.Island.addNpcFromName("farmland", nextTo)
	g.Island.addNpcFromName("apple tree", nextTo)

	assert.Equal(t, 1, len(npc.spawnsByType("shelter", 30)))

	assert.Equal(t, true, len(g.Island.Spawns) == 4)

	g.Island.Tick()

	if npc.CurrentAction == nil {
		assert.Fail(t, "current action is nil")
		return
	}
	assert.Equal(t, "build small fireplace", npc.CurrentAction.Name)

	duration := npc.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		g.Island.Tick()
	}

	assert.Equal(t, 1, len(npc.spawnsByName("small fireplace", 0)))
	assert.Equal(t, 1, len(npc.spawnsByType("fireplace", 0)))
}

func TestBuildShelter(t *testing.T) {

	g := testNewGame()

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())

	npc := g.Island.Spawns[0]

	// add nessecities nearby, so they dont need to be built
	nextTo, err := npc.randomNearby()
	assert.Equal(t, nil, err)

	// add nessecities, so they dont need to be built
	g.Island.addNpcFromName("farmland", nextTo)
	g.Island.addNpcFromName("small fireplace", nextTo)
	g.Island.addNpcFromName("apple tree", nextTo)
	assert.Equal(t, 4, len(g.Island.Spawns))

	npc.addToInventory("small branch")

	assert.Equal(t, 1, len(npc.spawnsByType("fireplace", 2)))

	g.Island.Tick()
	assert.Equal(t, true, npc.hasPlanned("build small shelter"))

	if npc.CurrentAction == nil {
		assert.Fail(t, "current action is nil")
		return
	}
	duration := npc.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		g.Island.Tick()
	}

	shelters := npc.spawnsByType("shelter", 0)
	assert.Equal(t, 1, len(shelters))

	// make sure npc made this their home
	assert.Equal(t, true, npc.Home == shelters[0])

	// make npc tired
	npc.Tiredness = npc.tirednessCap() + 1
	prevTiredness := npc.Tiredness
	g.Island.Tick()

	// make sure npc planned action: sleep
	assert.Equal(t, true, npc.hasPlanned("sleep"))

	// make npc fall asleep
	g.Island.Tick()

	// make sure they get the shelter bonus
	assert.Equal(t, 48, prevTiredness-npc.Tiredness) // 50 for the bonus, -2 for the ticks before starting to sleep
}

func TestBuildFarmland(t *testing.T) {

	g := testNewGame()

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())
	npc := g.Island.Spawns[0]

	// add nessecities nearby, so they dont need to be built
	nextTo, err := g.Island.Spawns[0].randomNearby()
	assert.Equal(t, nil, err)

	// add nessecities, so they dont need to be built
	g.Island.addNpcFromName("small fireplace", nextTo)
	assert.Equal(t, 1, len(npc.spawnsByType("fireplace", 30)))

	g.Island.addNpcFromName("small shelter", nextTo)
	g.Island.addNpcFromName("apple tree", nextTo)
	g.Island.addNpcFromName("cooking pit", nextTo)
	g.Island.addNpcFromName("small hut", nextTo)
	assert.Equal(t, 6, len(g.Island.Spawns))

	npc.addToInventory("small branch")

	g.Island.Tick()

	if npc.CurrentAction == nil {
		assert.Fail(t, "current action is nil")
		return
	}
	assert.Equal(t, "build farmland", npc.CurrentAction.Name)

	duration := npc.CurrentAction.Duration
	assert.Equal(t, true, duration > 0)

	for i := 0; i <= duration; i++ {
		g.Island.Tick()
	}

	assert.Equal(t, 1, len(npc.spawnsByName("farmland", 0)))
	assert.Equal(t, 1, len(npc.spawnsByType("food producer", 0)))
}

func TestTree(t *testing.T) {

	g := testNewGame()

	pos := g.Island.RandomPointAboveWater()
	g.Island.addNpcFromName("oak tree", pos)
	npc := g.Island.Spawns[0]

	assert.Equal(t, 1, len(g.Island.Spawns))

	// let tree drop some spawns
	for i := 0; i <= 100; i++ {
		g.Island.Tick()
	}

	assert.Equal(t, true, len(g.Island.Spawns) > 1)
	assert.Equal(t, true, len(npc.spawnsByType("wood", 1)) > 0)
}

func TestNpcDiesOfOldAge(t *testing.T) {

	g := testNewGame()

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())
	assert.Equal(t, 1, len(g.Island.Spawns))

	npc := g.Island.Spawns[0]
	npc.Age.Set(npc.ageCap() + 1)

	g.Island.Tick()

	// dwarf should have died of old age
	assert.Equal(t, 0, len(g.Island.Spawns))
}

func TestSpawnGravel(t *testing.T) {

	g := testNewGame()

	assert.Equal(t, 0, len(g.Island.Spawns))
	g.Island.spawnGravel()
	assert.Equal(t, true, len(g.Island.Spawns) > 1)
}

func TestNpcMovesToFireplace(t *testing.T) {

	g := testNewGame()

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())
	assert.Equal(t, 1, len(g.Island.Spawns))

	npc := g.Island.Spawns[0]
	npc.addToInventory("small branch")

	nextTo, err := npc.randomNearby()
	assert.Equal(t, nil, err)

	assert.Equal(t, false, npc.Position == nextTo)

	g.Island.addNpcFromName("small fireplace", nextTo)

	// add nessecities, so they dont need to be built
	g.Island.addNpcFromName("small shelter", nextTo)
	g.Island.addNpcFromName("apple tree", nextTo)
	g.Island.addNpcFromName("farmland", nextTo)
	assert.Equal(t, 5, len(g.Island.Spawns))

	// make dwarf wanna move to shelter
	npc.Coldness = npc.coldnessCap() + 1
	assert.Equal(t, true, npc.isCold())

	// let them travel to destination
	for {
		g.Island.Tick()
		if npc.Position.intMatches(&nextTo) {
			break
		}
	}

	assert.Equal(t, true, npc.Position.intMatches(&nextTo))

	// let npc start the fire
	g.Island.Tick()

	// let fire burn
	g.Island.Tick()

	// let them get warm by the fire
	assert.Equal(t, false, npc.isCold())
}

func TestNpcFindFirewoodThenMovesToFireplace(t *testing.T) {

	g := testNewGame()

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())
	assert.Equal(t, 1, len(g.Island.Spawns))

	npc := g.Island.Spawns[0]

	// NOTE: similar to TestNpcMovesToFireplace, but now also make sure dwarf finds a firewood

	nextTo, err := npc.randomNearby()
	assert.Equal(t, nil, err)

	assert.Equal(t, false, npc.Position == nextTo)

	g.Island.addNpcFromName("farmland", npc.Position)
	g.Island.addNpcFromName("small fireplace", nextTo)
	g.Island.addNpcFromName("small shelter", nextTo)
	g.Island.addNpcFromName("cooking pit", nextTo)
	g.Island.addNpcFromName("farmland", nextTo)
	g.Island.addNpcFromName("apple tree", nextTo)
	assert.Equal(t, 7, len(g.Island.Spawns))

	nextTo2, err := npc.randomNearby()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, nextTo.Equals(nextTo2))

	g.Island.addNpcFromName("branch", nextTo)
	assert.Equal(t, 8, len(g.Island.Spawns))

	// make dwarf wanna move to shelter
	npc.Coldness = npc.coldnessCap() + 1
	assert.Equal(t, true, npc.isCold())

	assert.Equal(t, false, npc.hasItemTypeInInventory("wood"))

	// make dwarf plan to get firewood
	g.Island.Tick()

	assert.Equal(t, false, npc.hasItemTypeInInventory("wood"))
	assert.Equal(t, true, npc.hasPlanned("find fire wood"))

	for {
		g.Island.Tick()

		// have dwarf find branch
		if npc.hasItemTypeInInventory("wood") {
			break
		}
	}

	assert.Equal(t, true, npc.hasItemTypeInInventory("wood"))

	// make dwarf plan to get warm by fireplace
	g.Island.Tick()
	assert.Equal(t, true, npc.hasPlannedType("wait"))

	// let dwarf get warm
	g.Island.Tick()

	// let them get warm by the fire
	assert.Equal(t, false, npc.isCold())
}

func TestBuildCookingPit(t *testing.T) {

	g := testNewGame()

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())

	assert.Equal(t, 1, len(g.Island.Spawns))
	npc := g.Island.Spawns[0]

	// add nessecities nearby, so they dont need to be built
	nextTo, err := npc.randomNearby()
	assert.Equal(t, nil, err)

	g.Island.addNpcFromName("small fireplace", nextTo)
	g.Island.addNpcFromName("small shelter", nextTo)
	g.Island.addNpcFromName("apple tree", nextTo)
	g.Island.addNpcFromName("farmland", nextTo)
	assert.Equal(t, 5, len(g.Island.Spawns))

	// make sure npc decides to build cooking pit
	g.Island.Tick()
	assert.Equal(t, true, npc.hasPlanned("build cooking pit"))

	// wait until done
	for {
		g.Island.Tick()

		if !npc.isPerforming("build cooking pit") {
			break
		}
	}

	assert.Equal(t, 1, len(npc.spawnsByName("cooking pit", 0)))
}

func TestBuildSmallHut(t *testing.T) {

	g := testNewGame()

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())
	assert.Equal(t, 1, len(g.Island.Spawns))
	npc := g.Island.Spawns[0]

	// add nessecities nearby, so they dont need to be built
	nextTo, err := npc.randomNearby()
	assert.Equal(t, nil, err)
	g.Island.addNpcFromName("small fireplace", nextTo)

	home := g.Island.addNpcFromName("small shelter", nextTo)
	npc.Home = home

	g.Island.addNpcFromName("apple tree", nextTo)
	g.Island.addNpcFromName("farmland", nextTo)
	g.Island.addNpcFromName("cooking pit", nextTo)
	assert.Equal(t, 6, len(g.Island.Spawns))

	// make sure npc decides to build cooking pit
	g.Island.Tick()
	assert.Equal(t, true, npc.hasPlanned("build small hut"))

	// wait until done
	for {
		g.Island.Tick()

		if !npc.isPerforming("build small hut") {
			break
		}
	}

	assert.Equal(t, 1, len(npc.spawnsByName("small hut", 0)))
}
