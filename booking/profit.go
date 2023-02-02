package booking

import (
	"github.com/carlos/s4l/floatrounder"
)

type ProfitPerNight struct {
	AvgNight float64 `json:"avg_night"`
	MinNight float64 `json:"min_night"`
	MaxNight float64 `json:"max_night"`
}

func NewProfitPerNight(avg, min, max float64) ProfitPerNight {
	return ProfitPerNight{
		AvgNight: floatrounder.ToNearest(avg),
		MinNight: floatrounder.ToNearest(min),
		MaxNight: floatrounder.ToNearest(max),
	}
}

type MaximizeProfit struct {
	RequestIDS  []string `json:"request_ids"`
	TotalProfit int32    `json:"total_profit"`
	ProfitPerNight
}

func NewMaximizeProfit(ids []string, total int32, night ProfitPerNight) MaximizeProfit {
	return MaximizeProfit{
		RequestIDS:     ids,
		TotalProfit:    total,
		ProfitPerNight: night,
	}
}

func profit(sellingRate, margin int32) int32 {
	return sellingRate * margin / 100
}

// it'll be wise to call profit to no replicate code, but dealing with floats could lead to wrong accuracy.
func profitPerNight(sellingRate, margin, nights int32) float64 {
	return (float64(sellingRate) * float64(margin) / 100) / float64(nights)
}
