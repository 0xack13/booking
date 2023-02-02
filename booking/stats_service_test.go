//go:build unit

package booking

import (
	"testing"
	"time"
)

func TestStatsService_ProfitPerNight(t *testing.T) {
	tests := map[string]struct {
		bookings []Booking
		expected ProfitPerNight
	}{
		"first use case": {
			bookings: []Booking{
				{
					RequestID:   "bookata_XY123",
					CheckIn:     time.Time{},
					Nights:      5,
					SellingRate: 200,
					Margin:      25,
				},
				{
					RequestID:   "kayete_PP234",
					CheckIn:     time.Time{},
					Nights:      4,
					SellingRate: 156,
					Margin:      22,
				},
			},
			expected: ProfitPerNight{
				AvgNight: 9.29,
				MinNight: 8.58,
				MaxNight: 10.00,
			},
		},
		"second use case": {
			bookings: []Booking{
				{
					RequestID:   "bookata_XY123",
					CheckIn:     time.Time{},
					Nights:      1,
					SellingRate: 50,
					Margin:      20,
				},
				{
					RequestID:   "kayete_PP234",
					CheckIn:     time.Time{},
					Nights:      1,
					SellingRate: 55,
					Margin:      22,
				},
				{
					RequestID:   "bookata_XY123",
					CheckIn:     time.Time{},
					Nights:      1,
					SellingRate: 49,
					Margin:      21,
				},
			},
			expected: ProfitPerNight{
				AvgNight: 10.80,
				MinNight: 10,
				MaxNight: 12.1,
			},
		},
		"my own use case": {
			bookings: []Booking{
				{
					RequestID:   "A",
					CheckIn:     time.Time{},
					Nights:      7,
					SellingRate: 244,
					Margin:      5,
				},
				{
					RequestID:   "B",
					CheckIn:     time.Time{},
					Nights:      5,
					SellingRate: 100,
					Margin:      9,
				},
				{
					RequestID:   "C",
					CheckIn:     time.Time{},
					Nights:      1,
					SellingRate: 79,
					Margin:      11,
				},
				{
					RequestID:   "D",
					CheckIn:     time.Time{},
					Nights:      3,
					SellingRate: 49,
					Margin:      21,
				},
			},
			expected: ProfitPerNight{
				AvgNight: 3.92,
				MinNight: 1.74,
				MaxNight: 8.69,
			},
		},
		"use case with invalid booking": {
			bookings: []Booking{
				{
					RequestID:   "bookata_XY123",
					CheckIn:     time.Time{},
					Nights:      0,
					SellingRate: 50,
					Margin:      20,
				},
				{
					RequestID:   "kayete_PP234",
					CheckIn:     time.Time{},
					Nights:      1,
					SellingRate: 55,
					Margin:      22,
				},
				{
					RequestID:   "bookata_XY123",
					CheckIn:     time.Time{},
					Nights:      1,
					SellingRate: 49,
					Margin:      21,
				},
			},

			expected: ProfitPerNight{
				AvgNight: 11.20,
				MinNight: 10.29,
				MaxNight: 12.1,
			},
		},
		"use case with unique invalid booking": {
			bookings: []Booking{
				{
					RequestID:   "bookata_XY123",
					CheckIn:     time.Time{},
					Nights:      0,
					SellingRate: 50,
					Margin:      20,
				},
			},
			expected: ProfitPerNight{
				AvgNight: 0,
				MinNight: 0,
				MaxNight: 0,
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := statsService.ProfitPerNight(tt.bookings)
			//todo: Using the != operator to compare two floating-point numbers can lead to inaccuracies ( mantissa )
			// Instead, we should compare their difference to see if it is less than some small error value.
			// This can be done with testify testing library and InDelta function. ( https://pkg.go.dev/github.com/stretchr/testify/assert?utm_source=godoc#InDelta )
			if got != tt.expected {
				t.Errorf("got: %f, expected: %f", got, tt.expected)
			}
		})
	}
}

const n = 1_000_000

var statsGlobal ProfitPerNight

func BenchmarkProfitPerNight(b *testing.B) {
	var p ProfitPerNight
	bookings := make([]Booking, 0, n)
	for i := 0; i < n; i++ {
		bookings = append(bookings, Booking{
			RequestID:   "kayete_PP234",
			CheckIn:     time.Time{},
			Nights:      4,
			SellingRate: 156,
			Margin:      22,
		})
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p = statsService.ProfitPerNight(bookings)
	}
	statsGlobal = p
}
