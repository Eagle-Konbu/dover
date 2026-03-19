package octopus

import (
	"context"
	"errors"
	"time"

	"github.com/Eagle-Konbu/dover/domain"
)

var ErrNotImplemented = errors.New("octopus client: not implemented")

type Client struct {
	Email         string
	Password      string
	APIURL        string
	AccountNumber string
}

func (c *Client) FetchDailyReadings(_ context.Context, _ time.Time) ([]domain.Reading, error) {
	return nil, ErrNotImplemented
}

func (c *Client) FetchDailyCost(_ context.Context, _ time.Time) (float64, error) {
	return 0, ErrNotImplemented
}
