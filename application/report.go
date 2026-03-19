package application

import (
	"context"
	"fmt"
	"time"

	"github.com/Eagle-Konbu/dover/domain"
	"github.com/Eagle-Konbu/dover/domain/ports"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

type ReportService struct {
	Energy   ports.EnergyClient
	Notifier ports.Notifier
}

func (s *ReportService) Run(ctx context.Context) error {
	yesterday := time.Now().In(jst).AddDate(0, 0, -1).Truncate(24 * time.Hour)

	readings, err := s.Energy.FetchDailyReadings(ctx, yesterday)
	if err != nil {
		return fmt.Errorf("fetch readings: %w", err)
	}

	costYen, err := s.Energy.FetchDailyCost(ctx, yesterday)
	hasCost := err == nil
	if !hasCost {
		costYen = 0
	}

	usage := domain.AggregateDailyUsage(yesterday, readings, costYen, hasCost)

	if err := s.Notifier.SendDailyReport(ctx, usage); err != nil {
		return fmt.Errorf("send report: %w", err)
	}
	return nil
}
