// Package period contains extensions for time.Duration to represent days and months.
package period

import "time"

// The Period type represents larger time.Duration constants than those
// provided in the standard library.
type Period time.Duration

//nolint: gomnd
const (
	// Day represents the amount of time in a single day.
	Day = Period(time.Hour * 24)
	// Week represents the amount of time in a single week.
	Week = Day * 7
	// Month represents the amount of time in 30 days.
	Month = Day * 30
	// Year represents the amount of time in a year.
	Year = Day * 365
)

var periodNames = map[string]Period{
	"day":   Day,
	"week":  Week,
	"month": Month,
	"year":  Year,
}

var periodValues = map[Period]string{
	Day:   "day",
	Week:  "week",
	Month: "month",
	Year:  "year",
}

// Parse the given string as a Period. Returns false
// if the string does not match a period value.
func Parse(str string) (Period, bool) {
	p, ok := periodNames[str]
	return p, ok
}

// Duration returns the Period as a time.Duration.
func (p Period) Duration() time.Duration {
	return time.Duration(p)
}

func (p Period) String() string {
	str, ok := periodValues[p]
	if !ok {
		return p.Duration().String()
	}
	return str
}
