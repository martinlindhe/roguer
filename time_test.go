package rogue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeOne(t *testing.T) {

	t1 := newTime(5)
	assert.Equal(t, int64(5), t1.Current())
	assert.Equal(t, int64(0), t1.Day())
	t1.Tick()
	assert.Equal(t, int64(6), t1.Current())
	assert.Equal(t, "6 minutes", t1.PassedSinceStart())

	t2 := newTime(3000)
	assert.Equal(t, int64(2), t2.Day())
	assert.Equal(t, "2 days", t2.PassedSinceStart())
}
