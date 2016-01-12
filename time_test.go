package rogue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeAssumptions(t *testing.T) {

	// NOTE: 1 tick = 1 minute
	t1 := newTime(1 * Minute)
	assert.Equal(t, int64(1), t1.Current())

	t1.Set(1 * Hour)
	assert.Equal(t, int64(60), t1.Current())

	// NOTE: 60 * 24 = 1440 ticks in a day
	t1.Set(1 * Day)
	assert.Equal(t, int64(1440), t1.Current())

	// NOTE: 1440 * 30 = 43200 ticks in a month
	t1.Set(1 * Month)
	assert.Equal(t, int64(43200), t1.Current())

	// NOTE: 43200 * 12 = 518400 ticks in a year
	t1.Set(1 * Year)
	assert.Equal(t, int64(518400), t1.Current())
}

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
	t1.Set(Day * 1)
	assert.Equal(t, "jan 1", t1.DayOfYear())

	t1.Set(Month*1 + Day*20)
	assert.Equal(t, "feb 20", t1.DayOfYear())
}

func TestDateString(t *testing.T) {

	t1 := newTime(Month*1 + Day*20 + Hour*18 + Minute*20)
	assert.Equal(t, "18:20\nfeb 20", t1.DateString())
}

func TestSeason(t *testing.T) {

	t1 := newTime(Month*1 + Day*20 + Hour*18 + Minute*20)
	assert.Equal(t, "spring", t1.Season())
}

func TestIsDaytime(t *testing.T) {

	// NOTE: daytime is 06:00 to 17:59

	t1 := newTime(Hour*05 + Minute*59)
	assert.Equal(t, false, t1.IsDaytime())

	t1.Set(Hour * 06)
	assert.Equal(t, true, t1.IsDaytime())

	t1.Set(Hour*17 + Minute*59)
	assert.Equal(t, true, t1.IsDaytime())

	t1.Set(Hour * 18)
	assert.Equal(t, false, t1.IsDaytime())
}
