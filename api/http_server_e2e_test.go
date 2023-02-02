//go:build e2e

package api

import (
	"bytes"
	"encoding/json"
	"github.com/carlos/s4l/booking"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAPI_Booking(t *testing.T) {
	tests := map[string]struct {
		payload            []byte
		endpoint           string
		handlerFunc        func(w http.ResponseWriter, req *http.Request)
		validateResponse   func([]byte, interface{}, *testing.T)
		expected           interface{}
		expectedStatusCode int
	}{
		"stats e2e call with pdf example": {
			payload: []byte(`
				[
					{
  						"request_id": "bookata_XY123",
  						"check_in": "2020-01-01",
  						"nights": 5,
  						"selling_rate": 200,
  						"margin": 20
					},
					{
  						"request_id": "kayete_PP234",
  						"check_in": "2018-01-04",
  						"nights": 4,
  						"selling_rate": 156,
  						"margin": 22
					}
				]
			`),
			handlerFunc:      booking.HandleR.HandlerStats,
			validateResponse: validateStatsResponse,
			expected: booking.ProfitPerNight{
				AvgNight: 8.29,
				MinNight: 8,
				MaxNight: 8.58,
			},
			expectedStatusCode: http.StatusOK,
		},
		"stats e2e call invalid booking payload": {
			payload: []byte(`
				[
					{
  						"request_id": "bookata_XY123",
  						"check_in": "2020-01-01",
  						"nights": 5,
  						"selling_rate": 200
					},
					{
  						"request_id": "kayete_PP234",
  						"check_in": "2018-01-04",
  						"nights": 4,
  						"selling_rate": 156,
  						"margin": 22
					}
				]
			`),
			handlerFunc:        booking.HandleR.HandlerStats,
			validateResponse:   nil,
			expected:           nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		"stats e2e call invalid payload": {
			payload:            []byte(""),
			handlerFunc:        booking.HandleR.HandlerStats,
			validateResponse:   nil,
			expected:           nil,
			expectedStatusCode: http.StatusInternalServerError,
		},
		"maximize e2e call with pdf example": {
			payload: []byte(`
				[
					{
  						"request_id": "A",
  						"check_in": "2018-01-01",
  						"nights": 10,
  						"selling_rate": 1000,
						"margin": 10
					},
					{
  						"request_id": "B",
  						"check_in": "2018-01-06",
  						"nights": 10,
  						"selling_rate": 700,
  						"margin": 10
					},
					{
  						"request_id": "C",
  						"check_in": "2018-01-12",
  						"nights": 10,
  						"selling_rate": 400,
  						"margin": 10
					}
				]
			`),
			handlerFunc:      booking.HandleR.HandlerMaximize,
			validateResponse: validateMaximizeResponse,
			expected: booking.MaximizeProfit{
				RequestIDS:  []string{"A", "C"},
				TotalProfit: 140,
				ProfitPerNight: booking.ProfitPerNight{
					AvgNight: 7,
					MinNight: 4,
					MaxNight: 10,
				},
			},
			expectedStatusCode: http.StatusOK,
		},
		"maximize e2e call with invalid booking": {
			payload: []byte(`
				[
					{
  						"request_id": "bookata_XY123",
  						"check_in": "2020-01-01",
  						"nights": 5,
  						"selling_rate": 200,
						"margin": 20
					},
					{
  						"request_id": "acme_AAAAA",
  						"check_in": "2020-01-10",
  						"nights": 4,
  						"margin": 30
					}
				]
			`),
			handlerFunc:        booking.HandleR.HandlerMaximize,
			validateResponse:   nil,
			expected:           booking.MaximizeProfit{},
			expectedStatusCode: http.StatusBadRequest,
		},
		"maximize e2e call with invalid payload": {
			payload:            []byte(""),
			handlerFunc:        booking.HandleR.HandlerMaximize,
			validateResponse:   nil,
			expected:           booking.MaximizeProfit{},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(tt.handlerFunc))
			resp, err := http.Post(server.URL+tt.endpoint, "", bytes.NewBuffer(tt.payload))
			if err != nil {
				t.Error(err)
			}

			defer func() {
				err = resp.Body.Close()
				if err != nil {
					t.Error(err)
				}
			}()

			if resp.StatusCode != tt.expectedStatusCode {
				t.Error(err)
			}

			if resp.StatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
				}
				tt.validateResponse(body, tt.expected, t)
			}
		})
	}
}

func validateStatsResponse(b []byte, e interface{}, t *testing.T) {
	var got booking.ProfitPerNight
	err := json.Unmarshal(b, &got)
	if err != nil {
		t.Error(err)
	}
	if _, ok := e.(booking.ProfitPerNight); !ok {
		t.Error(err)
	}
	expected := e.(booking.ProfitPerNight)
	if got != expected {
		t.Errorf("got: %v, expected: %v", got, expected)
	}
}

func validateMaximizeResponse(b []byte, e interface{}, t *testing.T) {
	var got booking.MaximizeProfit
	err := json.Unmarshal(b, &got)
	if err != nil {
		t.Error(err)
	}
	if _, ok := e.(booking.MaximizeProfit); !ok {
		t.Error(err)
	}
	expected := e.(booking.MaximizeProfit)
	if !reflect.DeepEqual(got.RequestIDS, expected.RequestIDS) {
		t.Errorf("got: %v, expected: %v", got.RequestIDS, expected.RequestIDS)
	}
	if got.TotalProfit != expected.TotalProfit {
		t.Errorf("got: %d, expected: %d", got.TotalProfit, expected.TotalProfit)
	}
	if got.ProfitPerNight != expected.ProfitPerNight {
		t.Errorf("got: %v, expected: %v", got.ProfitPerNight, expected.ProfitPerNight)
	}
}
