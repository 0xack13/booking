//go:build unit

package timeparser

import (
	"testing"
)

func TestToTime(t *testing.T) {
	tests := map[string]struct {
		date string
	}{
		"ok cast": {
			date: "2023-01-29",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := ToTime(tt.date)
			if err != nil {
				t.Errorf("got error %s", err)
			}
			if got.Format(layout) != tt.date {
				t.Errorf("got: %s, expected: %s", got.String(), tt.date)
			}
		})
	}
}
