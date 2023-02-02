//go:build unit

package floatrounder

import (
	"testing"
)

func TestRound(t *testing.T) {
	tests := map[string]struct {
		number   float64
		expected float64
		method   func(f float64) float64
	}{
		"to nearest": {
			number:   5.7896,
			expected: 5.79,
			method:   ToNearest,
		},
		"down": {
			number:   2.33291,
			expected: 2.33,
			method:   Down,
		},
		"up": {
			number:   1.55789,
			expected: 1.56,
			method:   Up,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tt.method(tt.number)
			//todo: As I said in the booking package, not a big fan of comparing floats this way
			if got != tt.expected {
				t.Errorf("got: %f, expected: %f", got, tt.expected)
			}
		})
	}
}
