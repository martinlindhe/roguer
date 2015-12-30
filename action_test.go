package rogue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseActionsDefinition(t *testing.T) {

	_, err := parseActionsDefinition("data/actions.yml")
	assert.Equal(t, nil, err)
}
