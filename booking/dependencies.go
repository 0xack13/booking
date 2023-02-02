package booking

import "log"

var (
	statsService    = initStatService()
	maximizeService = initMaximizeService()
	HandleR         = initHandler()
)

func initStatService() *StatsService {
	s, err := NewStatsService()
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func initMaximizeService() *MaximizeService {
	s, err := NewMaximizeService(statsService)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func initHandler() *Handler {
	h, err := NewHandler(statsService, maximizeService)
	if err != nil {
		log.Fatal()
	}
	return h
}
