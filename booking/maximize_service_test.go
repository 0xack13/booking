//go:build unit

package booking

import (
	"github.com/xsolrac87/booking/timeparser"
	"reflect"
	"testing"
	"time"
)

func TestMaximizeService_MaximTotalProfits(t *testing.T) {
	tests := map[string]struct {
		bookings []Booking
		expected MaximizeProfit
	}{
		"test combo ( B + C + D )": {
			bookings: []Booking{
				{
					RequestID:   "A",
					CheckIn:     parse("2023-01-05"),
					Nights:      5,
					SellingRate: 200,
					Margin:      20,
				},
				{
					RequestID:   "B",
					CheckIn:     parse("2023-01-04"),
					Nights:      4,
					SellingRate: 156,
					Margin:      5,
				},
				{
					RequestID:   "C",
					CheckIn:     parse("2023-01-09"),
					Nights:      4,
					SellingRate: 150,
					Margin:      6,
				},
				{
					RequestID:   "D",
					CheckIn:     parse("2023-01-09"),
					Nights:      1,
					SellingRate: 1600,
					Margin:      30,
				},
			},
			expected: MaximizeProfit{
				RequestIDS:  []string{"B", "D", "C"},
				TotalProfit: 496,
				ProfitPerNight: ProfitPerNight{
					AvgNight: 161.4,
					MinNight: 1.95,
					MaxNight: 480,
				},
			},
		},
		"unsorted test where booking no overlap each other": {
			bookings: []Booking{
				{
					RequestID:   "A",
					CheckIn:     parse("2023-02-05"),
					Nights:      5,
					SellingRate: 55,
					Margin:      6,
				},
				{
					RequestID:   "B",
					CheckIn:     parse("2023-01-17"),
					Nights:      7,
					SellingRate: 231,
					Margin:      12,
				},
				{
					RequestID:   "C",
					CheckIn:     parse("2023-01-04"),
					Nights:      4,
					SellingRate: 150,
					Margin:      6,
				},
				{
					RequestID:   "D",
					CheckIn:     parse("2023-01-22"),
					Nights:      1,
					SellingRate: 1600,
					Margin:      30,
				},
			},
			expected: MaximizeProfit{
				RequestIDS:  []string{"C", "A", "D", "B"},
				TotalProfit: 519,
				ProfitPerNight: ProfitPerNight{
					AvgNight: 121.72,
					MinNight: 0.66,
					MaxNight: 480,
				},
			},
		},
		"no combo just A": {
			bookings: []Booking{
				{
					RequestID:   "A",
					CheckIn:     parse("2023-01-05"),
					Nights:      5,
					SellingRate: 200,
					Margin:      20,
				},
			},
			expected: MaximizeProfit{
				RequestIDS:  []string{"A"},
				TotalProfit: 40,
				ProfitPerNight: ProfitPerNight{
					AvgNight: 8,
					MinNight: 8,
					MaxNight: 8,
				},
			},
		},
		"pdf first use case": {
			bookings: []Booking{
				{
					RequestID:   "A",
					CheckIn:     parse("2018-01-01"),
					Nights:      10,
					SellingRate: 1000,
					Margin:      10,
				},
				{
					RequestID:   "B",
					CheckIn:     parse("2018-01-06"),
					Nights:      10,
					SellingRate: 700,
					Margin:      10,
				},
				{
					RequestID:   "C",
					CheckIn:     parse("2018-01-12"),
					Nights:      10,
					SellingRate: 400,
					Margin:      10,
				},
			},
			expected: MaximizeProfit{
				RequestIDS:  []string{"A", "C"},
				TotalProfit: 140,
				ProfitPerNight: ProfitPerNight{
					AvgNight: 7,
					MinNight: 4,
					MaxNight: 10,
				},
			},
		},
		"pdf second use case": {
			bookings: []Booking{
				{
					RequestID:   "bookata_XY123",
					CheckIn:     parse("2020-01-01"),
					Nights:      5,
					SellingRate: 200,
					Margin:      20,
				},
				{
					RequestID:   "kayete_PP234",
					CheckIn:     parse("2020-01-04"),
					Nights:      4,
					SellingRate: 156,
					Margin:      5,
				},
				{
					RequestID:   "atropote_AA930",
					CheckIn:     parse("2020-01-04"),
					Nights:      4,
					SellingRate: 150,
					Margin:      6,
				},
				{
					RequestID:   "acme_AAAAA",
					CheckIn:     parse("2020-01-10"),
					Nights:      4,
					SellingRate: 160,
					Margin:      30,
				},
			},
			expected: MaximizeProfit{
				RequestIDS:  []string{"bookata_XY123", "acme_AAAAA"},
				TotalProfit: 88,
				ProfitPerNight: ProfitPerNight{
					AvgNight: 10,
					MinNight: 8,
					MaxNight: 12,
				},
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := maximizeService.MaximTotalProfits(tt.bookings)
			if !reflect.DeepEqual(got.RequestIDS, tt.expected.RequestIDS) {
				t.Errorf("got: %v, expected: %v", got.RequestIDS, tt.expected.RequestIDS)
			}

			if got.TotalProfit != tt.expected.TotalProfit {
				t.Errorf("got: %d, expected: %d", got.TotalProfit, tt.expected.TotalProfit)
			}

			if got.ProfitPerNight != tt.expected.ProfitPerNight {
				t.Errorf("got: %v, expected: %v", got.ProfitPerNight, tt.expected.ProfitPerNight)
			}
		})
	}
}

var MaxGlobal MaximizeProfit

func BenchmarkMaximTotalProfits(b *testing.B) {
	test := []Booking{
		{
			RequestID:   "bookata_XY123",
			CheckIn:     parse("2020-01-01"),
			Nights:      5,
			SellingRate: 200,
			Margin:      20,
		},
		{
			RequestID:   "kayete_PP234",
			CheckIn:     parse("2020-01-04"),
			Nights:      4,
			SellingRate: 156,
			Margin:      5,
		},
		{
			RequestID:   "atropote_AA930",
			CheckIn:     parse("2020-01-04"),
			Nights:      4,
			SellingRate: 150,
			Margin:      6,
		},
		{
			RequestID:   "acme_AAAAA",
			CheckIn:     parse("2020-01-10"),
			Nights:      4,
			SellingRate: 160,
			Margin:      30,
		},
	}

	var max MaximizeProfit
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		max = maximizeService.MaximTotalProfits(test)
	}
	MaxGlobal = max
}

func parse(date string) time.Time {
	t, _ := timeparser.ToTime(date)
	return t
}
