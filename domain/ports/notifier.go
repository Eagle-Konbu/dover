package ports

import (
	"context"

	"github.com/Eagle-Konbu/dover/domain"
)

type Notifier interface {
	SendDailyReport(ctx context.Context, usage domain.DailyUsage) error
}
