package octopus_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Eagle-Konbu/dover/infrastructure/octopus"
)

const tokenResponse = `{
	"data": {
		"obtainKrakenToken": {
			"token": "test-jwt-token"
		}
	}
}`

func readingsResponse(t *testing.T, readings []map[string]string) string {
	t.Helper()
	type gqlReading struct {
		StartAt string `json:"startAt"`
		EndAt   string `json:"endAt"`
		Value   string `json:"value"`
	}
	var rs []gqlReading
	for _, r := range readings {
		rs = append(rs, gqlReading{StartAt: r["startAt"], EndAt: r["endAt"], Value: r["value"]})
	}
	resp := map[string]interface{}{
		"data": map[string]interface{}{
			"account": map[string]interface{}{
				"properties": []interface{}{
					map[string]interface{}{
						"electricitySupplyPoints": []interface{}{
							map[string]interface{}{
								"halfHourlyReadings": rs,
							},
						},
					},
				},
			},
		},
	}
	b, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal readings response: %v", err)
	}
	return string(b)
}

func writeJSON(t *testing.T, w http.ResponseWriter, body string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(body)); err != nil {
		t.Fatalf("failed to write response: %v", err)
	}
}

func newFakeServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

func TestClient_FetchDailyReadings(t *testing.T) {
	t.Run("returns parsed readings", func(t *testing.T) {
		callCount := 0
		srv := newFakeServer(t, func(w http.ResponseWriter, r *http.Request) {
			callCount++
			if callCount == 1 {
				writeJSON(t, w, tokenResponse)
				return
			}

			auth := r.Header.Get("Authorization")
			if auth != "JWT test-jwt-token" {
				t.Errorf("Authorization = %q, want %q", auth, "JWT test-jwt-token")
			}

			writeJSON(t, w, readingsResponse(t, []map[string]string{
				{"startAt": "2025-01-15T00:00:00Z", "endAt": "2025-01-15T00:30:00Z", "value": "0.5"},
				{"startAt": "2025-01-15T00:30:00Z", "endAt": "2025-01-15T01:00:00Z", "value": "1.2"},
			}))
		})

		c := &octopus.Client{
			Email:         "test@example.com",
			Password:      "pass",
			APIURL:        srv.URL,
			AccountNumber: "A-12345",
		}
		date := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
		readings, err := c.FetchDailyReadings(context.Background(), date)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(readings) != 2 {
			t.Fatalf("got %d readings, want 2", len(readings))
		}
		if readings[0].Value != 0.5 {
			t.Errorf("readings[0].Value = %f, want 0.5", readings[0].Value)
		}
		if readings[1].Value != 1.2 {
			t.Errorf("readings[1].Value = %f, want 1.2", readings[1].Value)
		}
	})

	t.Run("empty readings", func(t *testing.T) {
		callCount := 0
		srv := newFakeServer(t, func(w http.ResponseWriter, _ *http.Request) {
			callCount++
			if callCount == 1 {
				writeJSON(t, w, tokenResponse)
				return
			}
			writeJSON(t, w, readingsResponse(t, nil))
		})

		c := &octopus.Client{
			Email:         "test@example.com",
			Password:      "pass",
			APIURL:        srv.URL,
			AccountNumber: "A-12345",
		}
		readings, err := c.FetchDailyReadings(context.Background(), time.Now())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(readings) != 0 {
			t.Errorf("got %d readings, want 0", len(readings))
		}
	})

	t.Run("invalid value returns error", func(t *testing.T) {
		callCount := 0
		srv := newFakeServer(t, func(w http.ResponseWriter, _ *http.Request) {
			callCount++
			if callCount == 1 {
				writeJSON(t, w, tokenResponse)
				return
			}
			writeJSON(t, w, readingsResponse(t, []map[string]string{
				{"startAt": "2025-01-15T00:00:00Z", "endAt": "2025-01-15T00:30:00Z", "value": "not-a-number"},
			}))
		})

		c := &octopus.Client{
			Email:         "test@example.com",
			Password:      "pass",
			APIURL:        srv.URL,
			AccountNumber: "A-12345",
		}
		_, err := c.FetchDailyReadings(context.Background(), time.Now())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "parse reading value") {
			t.Errorf("error = %q, want it to contain %q", err.Error(), "parse reading value")
		}
	})
}

func TestClient_EnsureToken(t *testing.T) {
	t.Run("auth failure returns error", func(t *testing.T) {
		srv := newFakeServer(t, func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(t, w, `{"errors":[{"message":"invalid credentials"}]}`)
		})

		c := &octopus.Client{
			Email:    "bad@example.com",
			Password: "wrong",
			APIURL:   srv.URL,
		}
		_, err := c.FetchDailyReadings(context.Background(), time.Now())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "obtain token") {
			t.Errorf("error = %q, want it to contain %q", err.Error(), "obtain token")
		}
	})

	t.Run("empty token returns error", func(t *testing.T) {
		srv := newFakeServer(t, func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(t, w, `{"data":{"obtainKrakenToken":{"token":""}}}`)
		})

		c := &octopus.Client{
			Email:    "test@example.com",
			Password: "pass",
			APIURL:   srv.URL,
		}
		_, err := c.FetchDailyReadings(context.Background(), time.Now())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "empty token") {
			t.Errorf("error = %q, want it to contain %q", err.Error(), "empty token")
		}
	})

	t.Run("token is reused across calls", func(t *testing.T) {
		tokenCalls := 0
		callCount := 0
		srv := newFakeServer(t, func(w http.ResponseWriter, _ *http.Request) {
			callCount++
			if callCount == 1 {
				tokenCalls++
				writeJSON(t, w, tokenResponse)
				return
			}
			writeJSON(t, w, readingsResponse(t, nil))
		})

		c := &octopus.Client{
			Email:         "test@example.com",
			Password:      "pass",
			APIURL:        srv.URL,
			AccountNumber: "A-12345",
		}
		ctx := context.Background()
		if _, err := c.FetchDailyReadings(ctx, time.Now()); err != nil {
			t.Fatalf("first call: %v", err)
		}
		if _, err := c.FetchDailyReadings(ctx, time.Now()); err != nil {
			t.Fatalf("second call: %v", err)
		}
		if tokenCalls != 1 {
			t.Errorf("token mutations = %d, want 1", tokenCalls)
		}
	})
}

func TestClient_FetchDailyCost(t *testing.T) {
	t.Run("returns not yet supported error", func(t *testing.T) {
		callCount := 0
		srv := newFakeServer(t, func(w http.ResponseWriter, _ *http.Request) {
			callCount++
			if callCount == 1 {
				writeJSON(t, w, tokenResponse)
				return
			}
		})

		c := &octopus.Client{
			Email:         "test@example.com",
			Password:      "pass",
			APIURL:        srv.URL,
			AccountNumber: "A-12345",
		}
		_, err := c.FetchDailyCost(context.Background(), time.Now())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "not yet supported") {
			t.Errorf("error = %q, want it to contain %q", err.Error(), "not yet supported")
		}
	})
}
