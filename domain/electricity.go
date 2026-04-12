package domain

import "time"

type Reading struct {
	StartAt      time.Time
	EndAt        time.Time
	Value        float64
	CostEstimate float64
}

type DailyUsage struct {
	Date     time.Time
	TotalKWh float64
	CostYen  float64
	HasCost  bool
}

func AggregateDailyUsage(date time.Time, readings []Reading) DailyUsage {
	var totalKWh, totalCost float64
	for _, r := range readings {
		totalKWh += r.Value
		totalCost += r.CostEstimate
	}
	return DailyUsage{
		Date:     date,
		TotalKWh: totalKWh,
		CostYen:  totalCost,
		HasCost:  len(readings) > 0,
	}
}
