package rogue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTilesetDefinition(t *testing.T) {

	_, err := parseTilesetDefinition("resources/assets/tilesets/oddball/tiles.yml")
	//spew.Dump(ts)
	assert.Equal(t, nil, err)
}
