package floatrounder

import (
	"math"
)

func ToNearest(f float64) float64 {
	return math.Round(f*100) / 100
}

func Down(f float64) float64 {
	return math.Floor(f*100) / 100
}

func Up(f float64) float64 {
	return math.Ceil(f*100) / 100
}
