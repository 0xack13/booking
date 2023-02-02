//go:build unit

package booking

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandler_HandlerStats(t *testing.T) {
	tests := map[string]struct {
		payload      []byte
		method       string
		expectedCode int
		expected     ProfitPerNight
	}{
		"pdf first example payload": {
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
  						"check_in": "2020-01-04",
  						"nights": 4,
  						"selling_rate": 156,
  						"margin": 22
					}
				]
			`),
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			expected: ProfitPerNight{
				AvgNight: 8.29,
				MinNight: 8,
				MaxNight: 8.58,
			},
		},
		"pdf second example payload": {
			payload: []byte(`
				[
					{
  						"request_id": "bookata_XY123",
  						"check_in": "2020-01-01",
  						"nights": 1,
  						"selling_rate": 50,
  						"margin": 20
					},
					{
  						"request_id": "kayete_PP234",
  						"check_in": "2020-01-04",
  						"nights": 1,
  						"selling_rate": 55,
  						"margin": 22
					},
					{
  						"request_id": "trivoltio_ZX69",
  						"check_in": "2020-01-07",
  						"nights": 1,
  						"selling_rate": 49,
  						"margin": 21
					}
				]
			`),
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			expected: ProfitPerNight{
				AvgNight: 10.80,
				MinNight: 10,
				MaxNight: 12.1,
			},
		},
		"invalid payload": {
			payload: []byte(`
				[
					{
  						"request_id": "bookata_XY123",
  						"check_in": "2020-01-01"
					}
				]
			`),
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expected:     ProfitPerNight{},
		},
		"payload with incorrect format": {
			payload: []byte(`
				[
					{
  						"request_id": "bookata_XY123",
  						"check_in": "2020-01-01",
					},
				]
			`),
			method:       http.MethodPost,
			expectedCode: http.StatusInternalServerError,
			expected:     ProfitPerNight{},
		},
		"invalid method": {
			payload: []byte(`
				[
					{
  						"request_id": "trivoltio_ZX69",
  						"check_in": "2020-01-07",
  						"nights": 1,
  						"selling_rate": 49,
  						"margin": 21
					}
				]
			`),
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			expected:     ProfitPerNight{},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			wr := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, "/stats", bytes.NewBuffer([]byte(tt.payload)))

			HandleR.HandlerStats(wr, req)
			if wr.Code != tt.expectedCode {
				t.Errorf("got HTTP status code %d, expected %d", wr.Code, tt.expectedCode)
			}

			if wr.Code == http.StatusOK {
				var got ProfitPerNight
				err := json.Unmarshal(wr.Body.Bytes(), &got)
				if err != nil {
					t.Error(err)
				}

				if got != tt.expected {
					t.Errorf("got: %v, expected: %v", got, tt.expected)
				}
			}
		})
	}
}

func TestHandler_HandlerMaximize(t *testing.T) {
	tests := map[string]struct {
		payload      []byte
		expectedCode int
		expected     MaximizeProfit
	}{
		"pdf example ( A + C )": {
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
			expectedCode: http.StatusOK,
			expected: MaximizeProfit{
				RequestIDS:  []string{"A", "C"},
				TotalProfit: 140,
				ProfitPerNight: ProfitPerNight{
					AvgNight: 7,
					MinNight: 4,
					MaxNight: 10,
				},
			},
		},
		"pdf example ( bookata_XY123 + acme_AAAAA )": {
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
  						"check_in": "2020-01-04",
  						"nights": 4,
  						"selling_rate": 156,
  						"margin": 5
					},
					{
  						"request_id": "atropote_AA930",
  						"check_in": "2020-01-04",
  						"nights": 4,
  						"selling_rate": 150,
  						"margin": 6
					},
					{
  						"request_id": "acme_AAAAA",
  						"check_in": "2020-01-10",
  						"nights": 4,
  						"selling_rate": 160,
  						"margin": 30
					}
				]
			`),
			expectedCode: http.StatusOK,
			expected: MaximizeProfit{
				RequestIDS:  []string{"bookata_XY123", "acme_AAAAA"},
				TotalProfit: 88,
				ProfitPerNight: ProfitPerNight{
					AvgNight: 10,
					MinNight: 8,
					MaxNight: 12,
				},
			},
		},
		"best combo: A": {
			payload: []byte(`
				[
					{
  						"request_id": "A",
  						"check_in": "2023-01-31",
  						"nights": 5,
  						"selling_rate": 2000,
  						"margin": 20
					},
					{
  						"request_id": "B",
  						"check_in": "2023-02-01",
  						"nights": 3,
  						"selling_rate": 156,
  						"margin": 22
					},
					{
  						"request_id": "C",
  						"check_in": "2023-02-05",
  						"nights": 4,
  						"selling_rate": 156,
  						"margin": 22
					}
				]
			`),
			expectedCode: http.StatusOK,
			expected: MaximizeProfit{
				RequestIDS:  []string{"A"},
				TotalProfit: 400,
				ProfitPerNight: ProfitPerNight{
					AvgNight: 80,
					MinNight: 80,
					MaxNight: 80,
				},
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			wr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/maximize", bytes.NewBuffer([]byte(tt.payload)))

			HandleR.HandlerMaximize(wr, req)
			if wr.Code != tt.expectedCode {
				t.Errorf("got HTTP status code %d, expected %d", wr.Code, tt.expectedCode)
			}

			var got MaximizeProfit
			err := json.Unmarshal(wr.Body.Bytes(), &got)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(got.RequestIDS, tt.expected.RequestIDS) {
				t.Errorf("got: %v, expected: %v", got.RequestIDS, tt.expected.RequestIDS)
			}

			if got.TotalProfit != tt.expected.TotalProfit {
				t.Errorf("got: %d, expected: %d", got.TotalProfit, tt.expected.TotalProfit)
			}

			if got.ProfitPerNight != tt.expected.ProfitPerNight {
				t.Errorf("got: %v, expected: %v", got.ProfitPerNight, tt.expected.ProfitPerNight)
			}
		})
	}
}
