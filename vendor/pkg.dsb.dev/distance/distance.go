// Package distance contains methods for working with distances, primarily in lat/long formats.
package distance

import "math"

// Between calculates the distance between two latitude and longitude values.
func Between(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	const radius = 6371 // Radius of Earth according to google.
	const degrees = 180.0

	a1 := lat1 * math.Pi / degrees
	b1 := lon1 * math.Pi / degrees
	a2 := lat2 * math.Pi / degrees
	b2 := lon2 * math.Pi / degrees

	x := math.Sin(a1)*math.Sin(a2) + math.Cos(a1)*math.Cos(a2)*math.Cos(b2-b1)
	return radius * math.Acos(x)
}
