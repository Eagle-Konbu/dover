package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Eagle-Konbu/dover/application"
	"github.com/Eagle-Konbu/dover/domain"
)

type mockEnergyClient struct {
	costErr  error
	readings []domain.Reading
	cost     float64
}

func (m *mockEnergyClient) FetchDailyReadings(_ context.Context, _ time.Time) ([]domain.Reading, error) {
	return m.readings, nil
}

func (m *mockEnergyClient) FetchDailyCost(_ context.Context, _ time.Time) (float64, error) {
	return m.cost, m.costErr
}

type mockNotifier struct {
	usage  domain.DailyUsage
	called bool
}

func (m *mockNotifier) SendDailyReport(_ context.Context, usage domain.DailyUsage) error {
	m.called = true
	m.usage = usage
	return nil
}

func TestReportService_Run(t *testing.T) {
	t.Run("aggregates and notifies", func(t *testing.T) {
		energy := &mockEnergyClient{
			readings: []domain.Reading{{Value: 1.0}, {Value: 2.0}},
			cost:     300.0,
		}
		notifier := &mockNotifier{}
		svc := &application.ReportService{Energy: energy, Notifier: notifier}

		if err := svc.Run(context.Background()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !notifier.called {
			t.Fatal("notifier was not called")
		}
		if notifier.usage.TotalKWh != 3.0 {
			t.Errorf("TotalKWh = %f, want 3.0", notifier.usage.TotalKWh)
		}
		if !notifier.usage.HasCost {
			t.Error("HasCost = false, want true")
		}
	})

	t.Run("cost error results in HasCost false", func(t *testing.T) {
		energy := &mockEnergyClient{
			readings: []domain.Reading{{Value: 1.5}},
			costErr:  errors.New("not available"),
		}
		notifier := &mockNotifier{}
		svc := &application.ReportService{Energy: energy, Notifier: notifier}

		if err := svc.Run(context.Background()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if notifier.usage.HasCost {
			t.Error("HasCost = true, want false")
		}
		if notifier.usage.CostYen != 0 {
			t.Errorf("CostYen = %f, want 0", notifier.usage.CostYen)
		}
	})
}
