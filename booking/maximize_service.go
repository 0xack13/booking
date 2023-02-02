package booking

import (
	"sort"
)

type MaximizeService struct {
	statsService *StatsService
}

func NewMaximizeService(svc *StatsService) (*MaximizeService, error) {
	return &MaximizeService{
		statsService: svc,
	}, nil
}

func (m *MaximizeService) MaximTotalProfits(bookings []Booking) MaximizeProfit {
	m.sortBookingsByCheckIn(bookings)

	var (
		combinations          = make(map[string][]Booking)
		checked               = make(map[string]struct{})
		bestIndividualBooking MaximizeProfit
		bestBookingCombo      MaximizeProfit
	)

	for i, booking := range bookings {
		individual := m.calculateIndividual(booking)
		if m.hasBetterTotalProfit(bestIndividualBooking, individual) {
			bestIndividualBooking = individual
		}

		if _, ok := checked[booking.RequestID]; ok {
			continue
		}

		combinations[booking.RequestID] = append(combinations[booking.RequestID], booking)
	Combo:
		for j := len(bookings[i+1:]); j > i; j-- {
			b := bookings[j]
			if !m.noOverlap(booking, b) {
				break Combo
			}
			combinations[booking.RequestID] = append(combinations[booking.RequestID], b)
			checked[b.RequestID] = struct{}{}
		}

		if len(combinations[booking.RequestID]) > 1 {
			combo := m.calculateCombination(combinations[booking.RequestID]...)
			if m.hasBetterTotalProfit(bestBookingCombo, combo) {
				bestBookingCombo = combo
			}
		}
	}

	return m.bestOne(bestIndividualBooking, bestBookingCombo)
}

func (m *MaximizeService) calculateIndividual(b Booking) MaximizeProfit {
	p := profit(b.SellingRate, b.Margin)
	perNight := m.statsService.ProfitPerNight([]Booking{b})
	return NewMaximizeProfit([]string{b.RequestID}, p, perNight)
}

func (m *MaximizeService) calculateCombination(bookings ...Booking) MaximizeProfit {
	var (
		requestID   []string
		totalProfit int32
	)

	for _, b := range bookings {
		requestID = append(requestID, b.RequestID)
		totalProfit += profit(b.SellingRate, b.Margin)
	}

	perNight := m.statsService.ProfitPerNight(bookings)
	return NewMaximizeProfit(requestID, totalProfit, perNight)
}

func (m *MaximizeService) hasBetterTotalProfit(a, b MaximizeProfit) bool {
	return a.TotalProfit == 0 || a.TotalProfit < b.TotalProfit
}

func (m *MaximizeService) sortBookingsByCheckIn(bookings []Booking) {
	sort.Slice(bookings, func(i, j int) bool {
		//a := bookings[i].CheckIn.AddDate(0, 0, int(bookings[i].Nights))
		//b := bookings[j].CheckIn.AddDate(0, 0, int(bookings[j].Nights))
		return bookings[i].CheckIn.Before(bookings[j].CheckIn)
	})
}

func (m *MaximizeService) noOverlap(a, b Booking) bool {
	checkOut := a.CheckIn.AddDate(0, 0, int(a.Nights))
	return b.CheckIn.After(checkOut)
}

func (m *MaximizeService) bestOne(a, b MaximizeProfit) MaximizeProfit {
	if a.TotalProfit > b.TotalProfit {
		return a
	}
	return b
}
