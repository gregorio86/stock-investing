package strategy

import (
	"context"

	"stock-investing/internal/kis"
	"stock-investing/internal/risk"
	"stock-investing/internal/screener"
	"stock-investing/internal/storage"
)

// Stable용 설정만 우선 넣고, 나중에 Aggressive용도 추가 가능
type StableConfig struct {
	ETFs        []string
	DailyAmount int64
}

type Deps struct {
	Ctx      context.Context
	KIS      *kis.Client
	Risk     risk.Manager
	Screener screener.Screener
	Repo     storage.Repository

	Stable StableConfig
}
