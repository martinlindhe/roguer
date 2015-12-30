package rogue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSpritesetDefinitionCharacters(t *testing.T) {

	ss, err := parseSpritesetDefinition("resources/assets/tilesets/oddball/characters.yml")
	assert.Equal(t, nil, err)
	assert.Equal(t, true, len(ss) > 2)
}

func TestParseSpritesetDefinitionItems(t *testing.T) {

	ss, err := parseSpritesetDefinition("resources/assets/tilesets/oddball/items.yml")
	assert.Equal(t, nil, err)
	assert.Equal(t, true, len(ss) > 2)
}
