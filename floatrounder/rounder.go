package floatrounder

import (
	"math"
)

// ToNearest The average profit per night from all the booking requests (two digit precision)
// https://app.swaggerhub.com/apis-docs/BlackfireSFL/BackendChallenge/1.0.1#/StatsResponse
func ToNearest(f float64) float64 {
	return math.Round(f*100) / 100
}

func Down(f float64) float64 {
	return math.Floor(f*100) / 100
}

func Up(f float64) float64 {
	return math.Ceil(f*100) / 100
}
