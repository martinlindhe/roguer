package rogue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGroundTilesetDefinition(t *testing.T) {

	ts, err := parseGroundTilesetDefinition("resources/assets/tilesets/oddball/ground.yml")
	assert.Equal(t, true, len(ts.DeepWater) > 1)
	assert.Equal(t, true, len(ts.ShallowWater) > 1)
	assert.Equal(t, true, len(ts.Beach) > 1)
	assert.Equal(t, true, len(ts.Grass) > 1)

	assert.Equal(t, true, len(ts.Lava) > 1)
	assert.Equal(t, true, len(ts.Dirt) > 1)
	assert.Equal(t, true, len(ts.Wall) > 1)
	assert.Equal(t, nil, err)
}
