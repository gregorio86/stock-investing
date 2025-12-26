package strategy

import (
	"context"
	"math"
	"time"

	"stock-investing/internal/models"
	"stock-investing/pkg/logger"
)

type AggressiveStrategy struct {
	deps Deps
}

func NewAggressiveStrategy(deps Deps) *AggressiveStrategy {
	return &AggressiveStrategy{deps: deps}
}

func (s *AggressiveStrategy) Run(ctx context.Context) error {
	logger.Info.Println("[aggressive] running high-volatility strategy")

	// 0) 간단한 equity 가정 (TODO: 실제 계좌 평가액으로 대체)
	equity := 10_000_000.0

	// 0-1) 최대 손실 한도 체크
	if err := s.deps.Risk.CheckMaxLoss(ctx, equity); err != nil {
		logger.Error.Printf("[aggressive] max loss check failed: %v\n", err)
		return err
	}

	// 1) 스크리너로 후보 종목 리스트 얻기
	stocks, err := s.deps.Screener.Screen(ctx)
	if err != nil {
		logger.Error.Printf("[aggressive] screener error: %v\n", err)
		return err
	}
	if len(stocks) == 0 {
		logger.Info.Println("[aggressive] no screened stocks, skipping")
		return nil
	}

	// 간단한 예: 상위 3개 종목만 매수 시도
	maxCandidates := 3
	if len(stocks) < maxCandidates {
		maxCandidates = len(stocks)
	}
	candidates := stocks[:maxCandidates]

	for _, stock := range candidates {
		select {
		case <-ctx.Done():
			logger.Info.Println("[aggressive] context canceled, aborting")
			return ctx.Err()
		default:
		}

		// 2) 현재가 조회
		price, err := s.deps.KIS.GetQuote(ctx, stock.Code)
		if err != nil {
			logger.Error.Printf("[aggressive] failed to get quote for %s: %v\n", stock.Code, err)
			continue
		}
		if price <= 0 {
			logger.Error.Printf("[aggressive] invalid price for %s: %.2f\n", stock.Code, price)
			continue
		}

		// 3) 포지션당 목표 비중 (예: 4%)
		targetValue := equity * 0.04
		qty := int64(math.Floor(targetValue / price))
		if qty <= 0 {
			logger.Info.Printf("[aggressive] amount too small for %s (price=%.2f, target=%.2f)\n", stock.Code, price, targetValue)
			continue
		}

		// 3-1) 포지션 사이즈 리스크 체크
		newPosValue := float64(qty) * price
		if err := s.deps.Risk.CheckPositionSize(ctx, equity, newPosValue); err != nil {
			logger.Error.Printf("[aggressive] position risk check failed for %s: %v\n", stock.Code, err)
			continue
		}

		// 4) 매수 주문 (stub)
		if err := s.deps.KIS.Buy(ctx, stock.Code, qty); err != nil {
			logger.Error.Printf("[aggressive] buy failed for %s: %v\n", stock.Code, err)
			continue
		}

		// 5) 트레이드 기록
		trade := &models.Trade{
			Code:     stock.Code,
			Side:     "BUY",
			Quantity: qty,
			Price:    price,
			Time:     time.Now(),
			Strategy: "aggressive",
		}
		if err := s.deps.Repo.InsertTrade(ctx, trade); err != nil {
			logger.Error.Printf("[aggressive] failed to insert trade for %s: %v\n", stock.Code, err)
			continue
		}

		logger.Info.Printf("[aggressive] buy %s x %d @ %.2f\n", stock.Code, qty, price)
	}

	logger.Info.Println("[aggressive] high-volatility strategy completed")
	return nil
}
