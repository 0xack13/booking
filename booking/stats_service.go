package booking

type StatsService struct{}

func NewStatsService() (*StatsService, error) {
	return &StatsService{}, nil
}

func (s *StatsService) ProfitPerNight(bookings []Booking) ProfitPerNight {
	var (
		profitPerNightList = make([]float64, 0, len(bookings))
		sumProfitPerNight  float64
	)

	for _, booking := range bookings {
		if booking.valid() {
			p := profitPerNight(booking.SellingRate, booking.Margin, booking.Nights)
			profitPerNightList = append(profitPerNightList, p)
			sumProfitPerNight += p
		}
	}

	var avgProfitPerNight float64
	if len(profitPerNightList) > 0 {
		avgProfitPerNight = sumProfitPerNight / float64(len(profitPerNightList))
	}

	return NewProfitPerNight(avgProfitPerNight, s.min(profitPerNightList), s.max(profitPerNightList))
}

func (s *StatsService) min(list []float64) float64 {
	if len(list) == 0 {
		return 0
	}
	min := list[0]
	for _, v := range list[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

func (s *StatsService) max(list []float64) float64 {
	if len(list) == 0 {
		return 0
	}
	max := list[0]
	for _, v := range list[1:] {
		if v > max {
			max = v
		}
	}
	return max
}
