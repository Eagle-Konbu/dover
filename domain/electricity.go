package domain

import "time"

type Reading struct {
	StartAt time.Time
	EndAt   time.Time
	Value   float64
}

type DailyUsage struct {
	Date     time.Time
	TotalKWh float64
	CostYen  float64
	HasCost  bool
}

func AggregateDailyUsage(date time.Time, readings []Reading, costYen float64, hasCost bool) DailyUsage {
	var total float64
	for _, r := range readings {
		total += r.Value
	}
	return DailyUsage{
		Date:     date,
		TotalKWh: total,
		CostYen:  costYen,
		HasCost:  hasCost,
	}
}
