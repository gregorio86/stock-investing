package strategy

import (
	"context"

	"stock-investing/pkg/logger"
)

type HybridStrategy struct {
	stable     *StableStrategy
	aggressive *AggressiveStrategy
}

func NewHybridStrategy(deps Deps) *HybridStrategy {
	return &HybridStrategy{
		stable:     NewStableStrategy(deps),
		aggressive: NewAggressiveStrategy(deps),
	}
}

func (h *HybridStrategy) Run(ctx context.Context) error {
	logger.Info.Println("[hybrid] start")
	if err := h.stable.Run(ctx); err != nil {
		return err
	}
	if err := h.aggressive.Run(ctx); err != nil {
		return err
	}
	logger.Info.Println("[hybrid] done")
	return nil
}
