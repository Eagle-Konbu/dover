package domain_test

import (
	"testing"
	"time"

	"github.com/Eagle-Konbu/dover/domain"
)

func TestAggregateDailyUsage(t *testing.T) {
	date := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

	t.Run("sums readings", func(t *testing.T) {
		readings := []domain.Reading{
			{Value: 0.5},
			{Value: 1.2},
			{Value: 0.3},
		}
		got := domain.AggregateDailyUsage(date, readings, 150.0, true)

		if got.Date != date {
			t.Errorf("Date = %v, want %v", got.Date, date)
		}
		wantKWh := 2.0
		if diff := got.TotalKWh - wantKWh; diff > 0.001 || diff < -0.001 {
			t.Errorf("TotalKWh = %f, want %f", got.TotalKWh, wantKWh)
		}
		if got.CostYen != 150.0 {
			t.Errorf("CostYen = %f, want 150.0", got.CostYen)
		}
		if !got.HasCost {
			t.Error("HasCost = false, want true")
		}
	})

	t.Run("empty readings", func(t *testing.T) {
		got := domain.AggregateDailyUsage(date, nil, 0, false)
		if got.TotalKWh != 0 {
			t.Errorf("TotalKWh = %f, want 0", got.TotalKWh)
		}
		if got.HasCost {
			t.Error("HasCost = true, want false")
		}
	})
}
