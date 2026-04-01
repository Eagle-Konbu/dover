package discord_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Eagle-Konbu/dover/domain"
	"github.com/Eagle-Konbu/dover/infrastructure/discord"
)

type webhookPayload struct {
	Embeds []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Timestamp   string `json:"timestamp"`
		Fields      []struct {
			Name   string `json:"name"`
			Value  string `json:"value"`
			Inline bool   `json:"inline"`
		} `json:"fields"`
		Color int `json:"color"`
	} `json:"embeds"`
}

func TestSendDailyReport_WithCost(t *testing.T) {
	t.Parallel()

	var got webhookPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", r.Header.Get("Content-Type"))
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	wh := &discord.Webhook{URL: srv.URL}
	usage := domain.DailyUsage{
		Date:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		TotalKWh: 12.5,
		CostYen:  350,
		HasCost:  true,
	}

	if err := wh.SendDailyReport(context.Background(), usage); err != nil {
		t.Fatalf("SendDailyReport: %v", err)
	}

	if len(got.Embeds) != 1 {
		t.Fatalf("embeds count = %d, want 1", len(got.Embeds))
	}
	e := got.Embeds[0]

	if e.Title != "⚡ 本日の電力使用量" {
		t.Errorf("title = %q, want ⚡ 本日の電力使用量", e.Title)
	}
	if e.Description != "2025-01-15 の電力使用状況" {
		t.Errorf("description = %q, want '2025-01-15 の電力使用状況'", e.Description)
	}
	if e.Color != 5814783 {
		t.Errorf("color = %d, want 5814783", e.Color)
	}
	if len(e.Fields) != 2 {
		t.Fatalf("fields count = %d, want 2", len(e.Fields))
	}
	if e.Fields[0].Name != "🔌 使用量" || e.Fields[0].Value != "12.5 kWh" || !e.Fields[0].Inline {
		t.Errorf("field[0] = %+v, want 🔌 使用量/12.5 kWh/inline", e.Fields[0])
	}
	if e.Fields[1].Name != "💴 料金" || e.Fields[1].Value != "¥350" || !e.Fields[1].Inline {
		t.Errorf("field[1] = %+v, want 💴 料金/¥350/inline", e.Fields[1])
	}
}

func TestSendDailyReport_WithoutCost(t *testing.T) {
	t.Parallel()

	var got webhookPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	wh := &discord.Webhook{URL: srv.URL}
	usage := domain.DailyUsage{
		Date:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		TotalKWh: 8.3,
		HasCost:  false,
	}

	if err := wh.SendDailyReport(context.Background(), usage); err != nil {
		t.Fatalf("SendDailyReport: %v", err)
	}

	if len(got.Embeds) != 1 {
		t.Fatalf("embeds count = %d, want 1", len(got.Embeds))
	}
	if len(got.Embeds[0].Fields) != 1 {
		t.Fatalf("fields count = %d, want 1 (no cost field)", len(got.Embeds[0].Fields))
	}
	if got.Embeds[0].Fields[0].Name != "🔌 使用量" {
		t.Errorf("field[0].Name = %q, want 🔌 使用量", got.Embeds[0].Fields[0].Name)
	}
}

func TestSendDailyReport_HTTPError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	wh := &discord.Webhook{URL: srv.URL}
	usage := domain.DailyUsage{
		Date:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		TotalKWh: 10.0,
	}

	err := wh.SendDailyReport(context.Background(), usage)
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}
