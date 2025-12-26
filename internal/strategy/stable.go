package strategy

import (
	"context"
	"math"
	"time"

	"stock-investing/internal/models"
	"stock-investing/pkg/logger"
)

type StableStrategy struct {
	deps Deps
}

func NewStableStrategy(deps Deps) *StableStrategy {
	return &StableStrategy{deps: deps}
}

func (s *StableStrategy) Run(ctx context.Context) error {
	logger.Info.Println("[stable] running DCA ETF strategy")

	// 0) 간단한 equity 가정 (TODO: 실제 계좌 평가액으로 대체)
	equity := 10_000_000.0

	// 0-1) 최대 손실 한도 체크
	if err := s.deps.Risk.CheckMaxLoss(ctx, equity); err != nil {
		logger.Error.Printf("[stable] max loss check failed: %v\n", err)
		return err
	}

	// 각 ETF에 대해 동일 금액으로 분할 매수 (간단 버전)
	if len(s.deps.Stable.ETFs) == 0 || s.deps.Stable.DailyAmount <= 0 {
		logger.Info.Println("[stable] no ETFs or DailyAmount configured, skipping")
		return nil
	}

	allocPerETF := float64(s.deps.Stable.DailyAmount) / float64(len(s.deps.Stable.ETFs))

	for _, code := range s.deps.Stable.ETFs {
		select {
		case <-ctx.Done():
			logger.Info.Println("[stable] context canceled, aborting")
			return ctx.Err()
		default:
		}

		// 1) 현재가 조회 (KIS stub)
		price, err := s.deps.KIS.GetQuote(ctx, code)
		if err != nil {
			logger.Error.Printf("[stable] failed to get quote for %s: %v\n", code, err)
			continue
		}
		if price <= 0 {
			logger.Error.Printf("[stable] invalid price for %s: %.2f\n", code, price)
			continue
		}

		// 2) 수량 계산
		qty := int64(math.Floor(allocPerETF / price))
		if qty <= 0 {
			logger.Info.Printf("[stable] amount too small for %s (price=%.2f, alloc=%.2f)\n", code, price, allocPerETF)
			continue
		}

		// 3) 포지션 사이즈 리스크 체크
		newPosValue := float64(qty) * price
		if err := s.deps.Risk.CheckPositionSize(ctx, equity, newPosValue); err != nil {
			logger.Error.Printf("[stable] position risk check failed for %s: %v\n", code, err)
			continue
		}

		// 4) 매수 주문 (stub)
		if err := s.deps.KIS.Buy(ctx, code, qty); err != nil {
			logger.Error.Printf("[stable] buy failed for %s: %v\n", code, err)
			continue
		}

		// 5) 트레이드 기록
		trade := &models.Trade{
			Code:     code,
			Side:     "BUY",
			Quantity: qty,
			Price:    price,
			Time:     time.Now(),
			Strategy: "stable",
		}
		if err := s.deps.Repo.InsertTrade(ctx, trade); err != nil {
			logger.Error.Printf("[stable] failed to insert trade for %s: %v\n", code, err)
			continue
		}

		logger.Info.Printf("[stable] DCA buy %s x %d @ %.2f\n", code, qty, price)
	}

	logger.Info.Println("[stable] DCA ETF strategy completed")
	return nil
}
