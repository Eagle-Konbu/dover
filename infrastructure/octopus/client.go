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
	Email         string
	Password      string
	APIURL        string
	AccountNumber string
	token         string
}

type ObtainJSONWebTokenInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Client) ensureToken(ctx context.Context) error {
	if c.token != "" {
		return nil
	}

	unauthClient := graphql.NewClient(c.APIURL, nil)

	var m struct {
		ObtainKrakenToken struct {
			Token string
		} `graphql:"obtainKrakenToken(input: $input)"`
	}
	vars := map[string]interface{}{
		"input": ObtainJSONWebTokenInput{
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
	c.gqlClient = graphql.NewClient(c.APIURL, nil).WithRequestModifier(func(r *http.Request) {
		r.Header.Set("Authorization", "JWT "+c.token)
	})
	return nil
}

func (c *Client) FetchDailyReadings(ctx context.Context, date time.Time) ([]domain.Reading, error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, err
	}

	from := date
	to := date.Add(24 * time.Hour)

	var q struct {
		Account struct {
			Properties []struct {
				ElectricitySupplyPoints []struct {
					HalfHourlyReadings []struct {
						StartAt time.Time `scalar:"true"`
						EndAt   time.Time `scalar:"true"`
						Value   string
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
				readings = append(readings, domain.Reading{
					StartAt: hr.StartAt,
					EndAt:   hr.EndAt,
					Value:   v,
				})
			}
		}
	}
	return readings, nil
}

func (c *Client) FetchDailyCost(ctx context.Context, _ time.Time) (float64, error) {
	if err := c.ensureToken(ctx); err != nil {
		return 0, err
	}
	// costOfCharge query shape is not yet confirmed from the API.
	return 0, fmt.Errorf("octopus: fetch cost: not yet supported")
}
