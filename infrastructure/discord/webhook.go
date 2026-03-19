package discord

import (
	"context"
	"errors"

	"github.com/Eagle-Konbu/dover/domain"
)

var ErrNotImplemented = errors.New("discord webhook: not implemented")

type Webhook struct {
	URL string
}

func (w *Webhook) SendDailyReport(_ context.Context, _ domain.DailyUsage) error {
	return ErrNotImplemented
}
