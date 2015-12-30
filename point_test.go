package rogue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpawnsByType(t *testing.T) {

	prepareIsland()
	assert.Equal(t, true, len(island.Spawns) == 0)

	pos := island.RandomPointAboveWater()

	assert.Equal(t, 0, len(pos.spawnsByName("small fireplace", 0)))
	assert.Equal(t, 0, len(pos.spawnsByName("small fireplace", 30)))
	assert.Equal(t, 0, len(pos.spawnsByType("fireplace", 0)))
	assert.Equal(t, 0, len(pos.spawnsByType("fireplace", 30)))

	island.addNpcFromName("small fireplace", pos)
	assert.Equal(t, true, len(island.Spawns) == 1)
	fmt.Println("pos1", pos)

	assert.Equal(t, 1, len(pos.spawnsByName("small fireplace", 0)))
	assert.Equal(t, 1, len(pos.spawnsByName("small fireplace", 30)))
	assert.Equal(t, 1, len(pos.spawnsByType("fireplace", 0)))
	assert.Equal(t, 1, len(pos.spawnsByType("fireplace", 30)))

	pos2, _ := pos.randomNearby()
	fmt.Println("pos2", pos2)
	assert.Equal(t, 0, len(pos2.spawnsByName("small fireplace", 0)))
	assert.Equal(t, 1, len(pos2.spawnsByName("small fireplace", 1)))
	assert.Equal(t, 1, len(pos2.spawnsByName("small fireplace", 20)))
	assert.Equal(t, 0, len(pos2.spawnsByType("fireplace", 0)))
	assert.Equal(t, 1, len(pos2.spawnsByType("fireplace", 1)))
	assert.Equal(t, 1, len(pos2.spawnsByType("fireplace", 20)))
}
