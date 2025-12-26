package risk

import (
	"context"
	"errors"

	"stock-investing/internal/models"
	"stock-investing/pkg/logger"
)

type Manager interface {
	CheckMaxLoss(ctx context.Context, equity float64) error
	CheckPositionSize(ctx context.Context, equity float64, newPositionValue float64) error
	CheckThemeConcentration(ctx context.Context, positions []models.Position) error
}

type Config struct {
	MaxRiskRatio     float64 // e.g. 0.1 = 최대 손실 10%
	MaxPositionRatio float64 // e.g. 0.05 = 종목당 5%
	MaxThemeRatio    float64 // e.g. 0.5 = 테마당 50%
	MinCashRatio     float64 // e.g. 0.2 = 현금 20% 유지
}

type manager struct {
	cfg Config
}

func NewManager(cfg Config) Manager {
	return &manager{cfg: cfg}
}

func (m *manager) CheckMaxLoss(ctx context.Context, equity float64) error {
	logger.Info.Printf("[risk] CheckMaxLoss stub: equity=%.2f, maxRisk=%.2f\n", equity, m.cfg.MaxRiskRatio)
	// TODO: 실제 계좌 손실률 계산 후, 한도를 넘으면 에러 반환
	return nil
}

func (m *manager) CheckPositionSize(ctx context.Context, equity float64, newPositionValue float64) error {
	ratio := newPositionValue / equity
	if ratio > m.cfg.MaxPositionRatio {
		logger.Error.Printf("[risk] position size too big: %.2f > %.2f\n", ratio, m.cfg.MaxPositionRatio)
		return errors.New("position size exceeds max ratio")
	}
	return nil
}

func (m *manager) CheckThemeConcentration(ctx context.Context, positions []models.Position) error {
	logger.Info.Println("[risk] CheckThemeConcentration stub")
	// TODO: 섹터/테마 정보 기반 집중도 계산
	return nil
}
