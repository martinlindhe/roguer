package rogue

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestParseGroundTilesetDefinition(t *testing.T) {

	ts, err := parseGroundTilesetDefinition("resources/assets/tilesets/oddball/ground.yml")
	spew.Dump(ts)
	assert.Equal(t, nil, err)
}
