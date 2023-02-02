//go:build unit

package booking

import (
	"testing"
	"time"
)

func TestValid(t *testing.T) {
	tests := map[string]struct {
		booking  Booking
		expected bool
	}{
		"valid": {
			booking: Booking{
				RequestID:   "made-up",
				CheckIn:     time.Time{},
				Nights:      10,
				SellingRate: 55,
				Margin:      15,
			},
			expected: true,
		},
		"invalid: no nights": {
			booking: Booking{
				RequestID:   "made-up",
				CheckIn:     time.Time{},
				Nights:      0,
				SellingRate: 13,
				Margin:      5,
			},
			expected: false,
		},
		"invalid: no selling rate": {
			booking: Booking{
				RequestID:   "made-up",
				CheckIn:     time.Time{},
				Nights:      10,
				SellingRate: 0,
				Margin:      5,
			},
			expected: false,
		},
		"invalid: no margin": {
			booking: Booking{
				RequestID:   "made-up",
				CheckIn:     time.Time{},
				Nights:      10,
				SellingRate: 13,
				Margin:      0,
			},
			expected: false,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tt.booking.valid()
			if got != tt.expected {
				t.Errorf("got: %t, expected: %t", got, tt.expected)
			}
		})
	}
}
