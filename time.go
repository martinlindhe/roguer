package rogue

import "fmt"

// time constants
const (
	Minute = 1
	Hour   = Minute * 60
	Day    = Hour * 24
	Month  = Day * 30
	Season = Month * 3
	Year   = Month * 12
)

// GameTime is the object representing game server time
//
// Server ticks every 3 real-world seconds.
// A day has 24 hours.
// There are 4 seasons.
// Each season has 3 months (90 days).
//
// Each tick progresses the in-game time by 1 minutes,
// giving 60 ticks for an hour
//
// So 1 game year = 518400 * 3 = 1555200 real seconds,
// or 25920 real time minutes (432 real hours, 7.2 real days).
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

// DateString returns "evening 19:20\nspring feb 20"
func (t *GameTime) DateString() string {

	return t.PartOfDay() + " " + t.TimeOfDay() + "\n" + t.Season() + " " + t.DayOfYear()
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

	// NOTE: day is stored 0-based, but displayed 1-based
	day++

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

// DayOfYear returns "19 feb in year 5"
func (t *GameTime) DayOfYear() string {

	_, _, day, month, year := t.date()

	months := []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}

	if year > 0 {
		return fmt.Sprintf("%s %d in year %d", months[month], day, year)
	}
	return fmt.Sprintf("%s %d", months[month], day)
}

// Season returns "spring" or "winter"
func (t *GameTime) Season() string {
	seasons := []string{"spring", "summer", "autumn", "winter"}

	season := t.Month() / 4
	return seasons[season]
}

// IsDaytime ...
func (t *GameTime) IsDaytime() bool {
	h := t.Hour()
	if h >= 6 && h <= 17 {
		return true
	}

	return false
}

// IsNighttime ...
func (t *GameTime) IsNighttime() bool {
	return !t.IsDaytime()
}

// PartOfDay returns "morning" or "midnight"
func (t *GameTime) PartOfDay() string {

	hour := t.Hour()
	if hour == 0 {
		return "midnight"
	}
	if hour == 6 {
		return "sunrise"
	}
	if hour == 12 {
		return "midday"
	}
	if hour == 18 {
		return "sunset"
	}
	if hour <= 12 {
		return "morning"
	}
	if hour <= 18 {
		return "afternoon"
	}
	return "evening"
}

// Plural returns "1 item" or "2 items"
func plural(t int64, base string) string {

	if t == 1 {
		return fmt.Sprintf("%d %s", t, base)
	}
	return fmt.Sprintf("%d %ss", t, base)
}
