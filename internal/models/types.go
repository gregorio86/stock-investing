package models

import "time"

type Stock struct {
	Code   string
	Name   string
	Market string
}

type Trade struct {
	ID       int64
	Code     string
	Side     string // "BUY" or "SELL"
	Quantity int64
	Price    float64
	Time     time.Time
	Strategy string // "stable", "aggressive", ...
}

type Position struct {
	Code     string
	Quantity int64
	AvgPrice float64
}

type DailyPnL struct {
	Date     time.Time
	Equity   float64
	Profit   float64
	Drawdown float64
}
