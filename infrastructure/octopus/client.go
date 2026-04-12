package octopus

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	graphql "github.com/hasura/go-graphql-client"

	"github.com/Eagle-Konbu/dover/domain"
)

type Client struct {
	gqlClient     *graphql.Client
	HTTPClient    *http.Client
	Email         string
	Password      string
	APIURL        string
	AccountNumber string
	token         string
}

type obtainJSONWebTokenInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{Timeout: 30 * time.Second}
}

func (c *Client) ensureToken(ctx context.Context) error {
	if c.token != "" {
		return nil
	}

	unauthClient := graphql.NewClient(c.APIURL, c.httpClient())

	var m struct {
		ObtainKrakenToken struct {
			Token string
		} `graphql:"obtainKrakenToken(input: $input)"`
	}
	vars := map[string]interface{}{
		"input": obtainJSONWebTokenInput{
			Email:    c.Email,
			Password: c.Password,
		},
	}
	if err := unauthClient.Mutate(ctx, &m, vars); err != nil {
		return fmt.Errorf("octopus: obtain token: %w", err)
	}
	if m.ObtainKrakenToken.Token == "" {
		return fmt.Errorf("octopus: obtain token: empty token returned")
	}

	c.token = m.ObtainKrakenToken.Token
	c.gqlClient = graphql.NewClient(c.APIURL, c.httpClient()).WithRequestModifier(func(r *http.Request) {
		r.Header.Set("Authorization", "JWT "+c.token)
	})
	return nil
}

func (c *Client) FetchDailyReadings(ctx context.Context, date time.Time) ([]domain.Reading, error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, err
	}

	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	from := date
	to := date.AddDate(0, 0, 1)

	var q struct {
		Account struct {
			Properties []struct {
				ElectricitySupplyPoints []struct {
					HalfHourlyReadings []struct {
						StartAt      time.Time `scalar:"true"`
						EndAt        time.Time `scalar:"true"`
						Value        string
						CostEstimate string
					} `graphql:"halfHourlyReadings(fromDatetime: $fromDatetime, toDatetime: $toDatetime)"`
				}
			}
		} `graphql:"account(accountNumber: $accountNumber)"`
	}
	vars := map[string]interface{}{
		"accountNumber": graphql.String(c.AccountNumber),
		"fromDatetime":  from,
		"toDatetime":    to,
	}
	if err := c.gqlClient.Query(ctx, &q, vars); err != nil {
		return nil, fmt.Errorf("octopus: fetch readings: %w", err)
	}

	var readings []domain.Reading
	for _, prop := range q.Account.Properties {
		for _, esp := range prop.ElectricitySupplyPoints {
			for _, hr := range esp.HalfHourlyReadings {
				v, err := strconv.ParseFloat(hr.Value, 64)
				if err != nil {
					return nil, fmt.Errorf("octopus: parse reading value %q: %w", hr.Value, err)
				}
				ce, err := strconv.ParseFloat(hr.CostEstimate, 64)
				if err != nil {
					return nil, fmt.Errorf("octopus: parse cost estimate %q: %w", hr.CostEstimate, err)
				}
				readings = append(readings, domain.Reading{
					StartAt:      hr.StartAt,
					EndAt:        hr.EndAt,
					Value:        v,
					CostEstimate: ce,
				})
			}
		}
	}
	return readings, nil
}
