package booking

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

var (
	errInvalidHttpMethod  = errors.New("invalid HTTP method only post allowed")
	errRequestBody        = errors.New("error reading request body")
	errInvalidRequestBody = errors.New("invalid request body")
)

type Handler struct {
	statsService    *StatsService
	maximizeService *MaximizeService
}

func NewHandler(stats *StatsService, max *MaximizeService) (*Handler, error) {
	return &Handler{
		statsService:    stats,
		maximizeService: max,
	}, nil
}

func (h *Handler) HandlerStats(w http.ResponseWriter, req *http.Request) {
	defer func() {
		err := req.Body.Close()
		if err != nil {
			log.Printf("failed to close response: %v\n", err)
		}
	}()

	log.Println("processing request from stats handler")
	bookings, err := h.handleRequest(req)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	p := h.statsService.ProfitPerNight(bookings)
	err = h.writeJSON(w, http.StatusOK, p)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}
	return
}

func (h *Handler) HandlerMaximize(w http.ResponseWriter, req *http.Request) {
	defer func() {
		err := req.Body.Close()
		if err != nil {
			log.Printf("failed to close response: %v\n", err)
		}
	}()

	log.Println("processing request from maximize handler")
	bookings, err := h.handleRequest(req)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	p := h.maximizeService.MaximTotalProfits(bookings)
	err = h.writeJSON(w, http.StatusOK, p)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}
	return
}

func (h *Handler) handleRequest(req *http.Request) ([]Booking, error) {
	if req.Method != http.MethodPost {
		return nil, errInvalidHttpMethod
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, errRequestBody
	}

	var bookings []Booking
	err = json.Unmarshal(body, &bookings)
	if err != nil {
		return nil, err
	}

	if len(bookings) == 0 {
		return nil, errInvalidRequestBody
	}
	return bookings, nil
}

func (h *Handler) sendErrorResponse(w http.ResponseWriter, err error) {
	var statusCode int
	switch {
	case errors.Is(err, errInvalidHttpMethod):
		statusCode = http.StatusMethodNotAllowed
	case errors.Is(err, errRequestBody),
		errors.Is(err, errInvalidRequestBody),
		errors.Is(err, errRequestIDMissing),
		errors.Is(err, errCheckInMissing),
		errors.Is(err, errNightsMissing),
		errors.Is(err, errSellingRateMissing),
		errors.Is(err, errMarginMissing):
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}

	http.Error(w, err.Error(), statusCode)
}

func (h *Handler) writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
