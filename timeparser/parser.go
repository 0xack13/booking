package timeparser

import "time"

const layout = "2006-01-02"

func ToTime(date string) (time.Time, error) {
	return time.Parse(layout, date)
}
