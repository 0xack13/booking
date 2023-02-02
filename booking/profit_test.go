//go:build unit

package booking

import "testing"

func TestProfit(t *testing.T) {
	tests := map[string]struct {
		sellingRate int32
		margin      int32
		expected    int32
	}{
		"profit with selling rate 200 and margin 20": {
			sellingRate: 200,
			margin:      20,
			expected:    40,
		},
		"profit with selling rate 50 and margin 5": {
			sellingRate: 50,
			margin:      5,
			expected:    2,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := profit(tt.sellingRate, tt.margin)
			if got != tt.expected {
				t.Errorf("got: %d, expected: %d", got, tt.expected)
			}
		})
	}
}

func TestProfitPerNight(t *testing.T) {
	tests := map[string]struct {
		sellingRate int32
		margin      int32
		nights      int32
		expected    float64
	}{
		"profit per night 1": {
			sellingRate: 200,
			margin:      20,
			nights:      5,
			expected:    8.0,
		},
		"profit per night 2": {
			sellingRate: 50,
			margin:      20,
			nights:      1,
			expected:    10.000000,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := profitPerNight(tt.sellingRate, tt.margin, tt.nights)
			// Using the != operator to compare two floating-point numbers can lead to inaccuracies ( mantissa )
			// Instead, we should compare their difference to see if it is less than some small error value.
			// This can be done with testify testing library and InDelta function. ( https://pkg.go.dev/github.com/stretchr/testify/assert?utm_source=godoc#InDelta )
			if got != tt.expected {
				t.Errorf("got: %f, expected: %f", got, tt.expected)
			}
		})
	}
}

func FuzzProfit(f *testing.F) {
	f.Fuzz(func(t *testing.T, rate, margin int32) {
		profit(rate, margin)
	})
}

func FuzzPerNight(f *testing.F) {
	f.Fuzz(func(t *testing.T, rate, margin, nights int32) {
		profitPerNight(rate, margin, nights)
	})
}
