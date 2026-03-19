package ports

import (
	"context"
	"time"

	"github.com/Eagle-Konbu/dover/domain"
)

type EnergyClient interface {
	FetchDailyReadings(ctx context.Context, date time.Time) ([]domain.Reading, error)
	FetchDailyCost(ctx context.Context, date time.Time) (float64, error)
}
