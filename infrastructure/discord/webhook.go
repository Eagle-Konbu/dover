package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Eagle-Konbu/dover/domain"
)

type Webhook struct {
	URL string
}

type embed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Timestamp   string       `json:"timestamp"`
	Fields      []embedField `json:"fields"`
	Color       int          `json:"color"`
}

type embedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type webhookPayload struct {
	Embeds []embed `json:"embeds"`
}

func (w *Webhook) SendDailyReport(ctx context.Context, usage domain.DailyUsage) error {
	fields := []embedField{
		{
			Name:   "🔌 使用量",
			Value:  fmt.Sprintf("%.1f kWh", usage.TotalKWh),
			Inline: true,
		},
	}
	if usage.HasCost {
		fields = append(fields, embedField{
			Name:   "💴 料金",
			Value:  fmt.Sprintf("¥%.0f", usage.CostYen),
			Inline: true,
		})
	}

	payload := webhookPayload{
		Embeds: []embed{
			{
				Title:       "⚡ 本日の電力使用量",
				Description: fmt.Sprintf("%s の電力使用状況", usage.Date.Format("2006-01-02")),
				Color:       5814783,
				Fields:      fields,
				Timestamp:   usage.Date.Format("2006-01-02T15:04:05Z07:00"),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord webhook: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord webhook: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("discord webhook: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, err := io.ReadAll(io.LimitReader(resp.Body, 512))
		if err != nil {
			return fmt.Errorf("discord webhook: unexpected status code: %d", resp.StatusCode)
		}
		return fmt.Errorf("discord webhook: unexpected status code: %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
