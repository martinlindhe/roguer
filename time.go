package rogue

import "fmt"

const (
	minute = 1
	hour   = minute * 60
	day    = hour * 24
	month  = day * 30
	season = month * 3
	year   = season * 4
)

// GameTime is the object representing game server time
//
// Server ticks every 3 real-world seconds.
//
// Each tick progresses the in-game time by 1 minutes,
// giving 60 ticks for an hour, or 1440 ticks for a day (24 hour days in-game).
//
// There are 4 seasons. Each season has 25 days,
// making 1 year in 1440 * 100 = 144 000 ticks
//
// So 1 game year = 144000 * 3 = 432 000 real seconds,
// or 120 real time hours (5 real time days).
//
// Some examples:
//    day 0, hour 0, min 1: 1 ticks
//    day 0, hour 1, min 0: 60 ticks
//    day 1, hour 0, min 0: 1440 ticks
type GameTime struct {
	time int64
}

func newTime(t int64) *GameTime {

	return &GameTime{time: t}
}

// Tick progress the time
func (t *GameTime) Tick() {

	t.time++
}

// Current returns the current time
func (t *GameTime) Current() int64 {

	return t.time
}

// Day returns the day of the year
func (t *GameTime) Day() int64 {

	return t.time / day
}

// Plural returns "s" to construct english plural forms of words
func (t *GameTime) Plural(base string) string {

	if t.time == 1 {
		return base
	}
	return base + "s"
}

// PassedSinceStart describes distance of t.time to 0
func (t *GameTime) PassedSinceStart() string {

	if t.time < hour {
		return fmt.Sprintf("%d %s", t.time, t.Plural("minute"))
	}
	if t.time < day {
		r := newTime(t.time / hour)
		return fmt.Sprintf("%d %s", r.time, r.Plural("hour"))
	}
	if t.time < month {
		r := newTime(t.time / day)
		return fmt.Sprintf("%d %s", r.time, r.Plural("day"))
	}
	if t.time < year {
		r := newTime(t.time / month)
		return fmt.Sprintf("%d %s", r.time, r.Plural("month"))
	}
	r := newTime(t.time / year)
	return fmt.Sprintf("%d %s", r.time, r.Plural("year"))
}
