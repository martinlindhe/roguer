package rogue

import "fmt"

// time constants
const (
	Minute = 1
	Hour   = Minute * 60
	Day    = Hour * 24
	Month  = Day * 30
	Year   = Month * 12
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

// Set sets the current time
func (t *GameTime) Set(i int64) {
	t.time = i
}

// Tick progress the time
func (t *GameTime) Tick() {

	t.time++
}

// Current returns the current time
func (t *GameTime) Current() int64 {

	return t.time
}

// TimeOfDay returns time as "19:30"
func (t *GameTime) TimeOfDay() string {

	minute, hour, _, _, _ := t.date()

	return fmt.Sprintf("%02d:%02d", hour, minute)
}

func (t *GameTime) date() (int64, int64, int64, int64, int64) {

	rest := t.time

	year := rest / Year
	if year > 0 {
		rest -= year * Year
	}

	month := rest / Month
	if month > 0 {
		rest -= month * Month
	}

	day := rest / Day
	if day > 0 {
		rest -= day * Day
	}

	hour := rest / Hour
	if hour > 0 {
		rest -= hour * Hour
	}

	minute := rest

	return minute, hour, day, month, year
}

// Minute returns the minute
func (t *GameTime) Minute() int64 {

	minute, _, _, _, _ := t.date()
	return minute
}

// Hour returns the hour
func (t *GameTime) Hour() int64 {

	_, hour, _, _, _ := t.date()
	return hour
}

// Day returns the day
func (t *GameTime) Day() int64 {

	_, _, day, _, _ := t.date()
	return day
}

// Month returns the month
func (t *GameTime) Month() int64 {

	_, _, _, month, _ := t.date()
	return month
}

// Year returns the year
func (t *GameTime) Year() int64 {

	_, _, _, _, year := t.date()
	return year
}

// PassedSinceStart describes distance of t.time to 0
func (t *GameTime) PassedSinceStart() string {

	minute, hour, day, month, year := t.date()

	if year > 0 {
		return plural(year, "year")
	}
	if month > 0 {
		return plural(month, "month")
	}
	if day > 0 {
		return plural(day, "day")
	}
	if hour > 0 {
		return plural(hour, "hour")
	}
	return plural(minute, "minute")
}

// DayOfYear ...
func (t *GameTime) DayOfYear() string {

	_, _, day, month, year := t.date()

	return fmt.Sprintf("day %d of month %d in year %d", day, month, year)
}

// Plural returns "1 item" or "2 items"
func plural(t int64, base string) string {

	if t == 1 {
		return fmt.Sprintf("%d %s", t, base)
	}
	return fmt.Sprintf("%d %ss", t, base)
}
