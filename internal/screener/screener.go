package screener

import (
	"context"

	"stock-investing/internal/models"
	"stock-investing/pkg/logger"
)

type Screener interface {
	Screen(ctx context.Context) ([]*models.Stock, error)
}

type KosdaqScreener struct{}

func NewKosdaqScreener() *KosdaqScreener {
	return &KosdaqScreener{}
}

func (s *KosdaqScreener) Screen(ctx context.Context) ([]*models.Stock, error) {
	logger.Info.Println("[screener] running dummy screener: returning 1 test stock")
	// TODO: CSV + KIS 시세 기반으로 실제 필터링 구현
	return []*models.Stock{
		{
			Code:   "123456",
			Name:   "DummyHighVol",
			Market: "KOSDAQ",
		},
	}, nil
}
