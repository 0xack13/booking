package booking

import (
	"encoding/json"
	"errors"
	"github.com/xsolrac87/booking/timeparser"
	"time"
)

var (
	errRequestIDMissing   = errors.New("booking payload does not contain request_id property")
	errCheckInMissing     = errors.New("booking payload does not contain check_in property")
	errNightsMissing      = errors.New("booking payload does not contain nights property")
	errSellingRateMissing = errors.New("booking payload does not contain selling_rate property")
	errMarginMissing      = errors.New("booking payload does not contain margin property")
)

type Booking struct {
	RequestID   string    `json:"request_id"`
	CheckIn     time.Time `json:"check_in"`
	Nights      int32     `json:"nights"`
	SellingRate int32     `json:"selling_rate"`
	Margin      int32     `json:"margin"`
}

func (b *Booking) valid() bool {
	return b.Nights > 0 && b.SellingRate > 0 && b.Margin > 0
}

func (b *Booking) UnmarshalJSON(data []byte) error {
	payload := make(map[string]interface{})
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return err
	}

	if _, ok := payload["request_id"]; !ok {
		return errRequestIDMissing
	}
	if _, ok := payload["check_in"]; !ok {
		return errCheckInMissing
	}
	if _, ok := payload["nights"]; !ok {
		return errNightsMissing
	}
	if _, ok := payload["selling_rate"]; !ok {
		return errSellingRateMissing
	}
	if _, ok := payload["margin"]; !ok {
		return errMarginMissing
	}

	t, err := timeparser.ToTime(payload["check_in"].(string))
	if err != nil {
		return err
	}

	b.RequestID = payload["request_id"].(string)
	b.CheckIn = t
	b.Nights = int32(payload["nights"].(float64))
	b.SellingRate = int32(payload["selling_rate"].(float64))
	b.Margin = int32(payload["margin"].(float64))
	return nil
}
