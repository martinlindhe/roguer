package rogue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpawnsByType(t *testing.T) {

	g := testNewGame()
	assert.Equal(t, 0, len(g.Island.Spawns))

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())
	assert.Equal(t, 1, len(g.Island.Spawns))

	npc := g.Island.Spawns[0]

	assert.Equal(t, 0, len(npc.spawnsByName("small fireplace", 0)))
	assert.Equal(t, 0, len(npc.spawnsByName("small fireplace", 30)))
	assert.Equal(t, 0, len(npc.spawnsByType("fireplace", 0)))
	assert.Equal(t, 0, len(npc.spawnsByType("fireplace", 30)))

	g.Island.addNpcFromName("small fireplace", npc.Position)
	assert.Equal(t, 2, len(g.Island.Spawns))

	assert.Equal(t, 1, len(npc.spawnsByName("small fireplace", 0)))
	assert.Equal(t, 1, len(npc.spawnsByName("small fireplace", 30)))
	assert.Equal(t, 1, len(npc.spawnsByType("fireplace", 0)))
	assert.Equal(t, 1, len(npc.spawnsByType("fireplace", 30)))
}

func TestRandomNearby(t *testing.T) {
	// should never get the input point

	g := testNewGame()

	g.Island.addNpcFromRace("dwarf", g.Island.RandomPointAboveWater())
	npc := g.Island.Spawns[0]

	pos := g.Island.RandomPointAboveWater()

	for i := 0; i < 100; i++ {
		p2, _ := npc.randomNearby()
		assert.Equal(t, false, p2.Equals(pos))
	}
}
