package rogue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// NOTE: 1 tick = 1 minute

func TestPassedSinceStart(t *testing.T) {

	t1 := newTime(5)
	assert.Equal(t, int64(5), t1.Current())
	assert.Equal(t, int64(0), t1.Day())
	assert.Equal(t, "5 minutes", t1.PassedSinceStart())
	t1.Tick()
	assert.Equal(t, int64(6), t1.Current())
	assert.Equal(t, "6 minutes", t1.PassedSinceStart())

	t2 := newTime(3000)
	assert.Equal(t, int64(2), t2.Day())
	assert.Equal(t, "2 days", t2.PassedSinceStart())
}

func TestTimeOfDay(t *testing.T) {

	t1 := newTime(5)
	assert.Equal(t, "00:05", t1.TimeOfDay())

	t1.Set(2 * Hour)
	assert.Equal(t, "02:00", t1.TimeOfDay())
}

func TestDayOfYear(t *testing.T) {
	t1 := newTime(0)
	assert.Equal(t, "day 0 of month 0 in year 0", t1.DayOfYear())
}
