package rogue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGroundTilesetDefinition(t *testing.T) {

	ts, err := parseGroundTilesetDefinition("resources/assets/tilesets/oddball/ground.yml")

	assert.Equal(t, 10, len(ts.DeepWater))
	assert.Equal(t, 14, len(ts.ShallowWater))
	assert.Equal(t, 8, len(ts.Beach))
	assert.Equal(t, 14, len(ts.Grass))

	assert.Equal(t, 14, len(ts.Lava))
	assert.Equal(t, 14, len(ts.Dirt))
	assert.Equal(t, 6, len(ts.Wall))
	assert.Equal(t, nil, err)
}
