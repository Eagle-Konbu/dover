package main

import (
	"context"
	"fmt"

	"github.com/Eagle-Konbu/dover/domain"
)

type stdoutNotifier struct{}

func (s *stdoutNotifier) SendDailyReport(_ context.Context, usage domain.DailyUsage) error {
	fmt.Printf("Date:     %s\n", usage.Date.Format("2006-01-02"))
	fmt.Printf("Usage:    %.2f kWh\n", usage.TotalKWh)
	if usage.HasCost {
		fmt.Printf("Cost:     %.0f JPY\n", usage.CostYen)
	} else {
		fmt.Println("Cost:     N/A")
	}
	return nil
}
