package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"stock-investing/internal/config"
	"stock-investing/internal/kis"
	"stock-investing/internal/risk"
	"stock-investing/internal/screener"
	"stock-investing/internal/storage"
	"stock-investing/internal/strategy"
	"stock-investing/pkg/logger"
	"stock-investing/scheduler"
)

func main() {
	modeFlag := flag.String("mode", "hybrid", "trading mode: hybrid|stable|aggressive")
	initDB := flag.Bool("init-db", false, "initialize database")
	flag.Parse()

	// 1) 로거 초기화
	logger.Init()

	// 2) 환경변수 로드
	cfg := config.Load()
	logger.Info.Printf("config loaded, mock=%v, server=%s\n", cfg.MockTrading, cfg.KIS)

	// 3) DB 초기화 모드
	if *initDB {
		logger.Info.Println("initializing SQLite DB...")
		store, err := storage.NewSQLiteStore("stock_investing.db")
		if err != nil {
			logger.Error.Fatalf("failed to open sqlite: %v", err)
		}
		defer store.Close()

		if err := store.Migrate(); err != nil {
			logger.Error.Fatalf("failed to migrate sqlite: %v", err)
		}

		logger.Info.Println("SQLite DB initialized successfully")
		return
	}

	// 4) 공통 의존성 초기화
	store, err := storage.NewSQLiteStore("stock_investing.db")
	if err != nil {
		logger.Error.Fatalf("failed to open sqlite: %v", err)
	}
	defer store.Close()

	repo := storage.NewRepository(store)
	kisClient := kis.NewClient(
		cfg.KIS.AppKey,
		cfg.KIS.AppSecret,
		cfg.KIS.BaseURL,
		cfg.KIS.AccountNo,
	)

	riskMgr := risk.NewManager(risk.Config{
		MaxRiskRatio:     cfg.Risk.MaxRisk, // .env의 MAX_RISK
		MaxPositionRatio: 0.05,             // 종목당 5% (임시)
		MaxThemeRatio:    0.5,              // 테마당 50% (임시)
		MinCashRatio:     0.2,              // 현금 20% 유지 (임시)
	})

	scr := screener.NewKosdaqScreener()

	deps := strategy.Deps{
		KIS:      kisClient,
		Risk:     riskMgr,
		Screener: scr,
		Repo:     repo,
		Stable: strategy.StableConfig{
			ETFs:        cfg.Stable.ETFs,
			DailyAmount: cfg.Stable.DailyAmount,
		},
	}

	// 5) 전략 선택
	var runner strategy.Runner
	switch strategy.Mode(*modeFlag) {
	case strategy.ModeHybrid:
		runner = strategy.NewHybridStrategy(deps)
	case strategy.ModeStable:
		runner = strategy.NewStableStrategy(deps)
	case strategy.ModeAggressive:
		runner = strategy.NewAggressiveStrategy(deps)
	default:
		log.Fatalf("unknown mode: %s", *modeFlag)
	}

	// 6) 스케줄러 시작 (현재는 stub)
	s := scheduler.New()
	s.Start()
	defer s.Stop()

	// 7) 종료 시그널 처리 + 컨텍스트
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 8) 선택된 전략 실행
	if err := runner.Run(ctx); err != nil {
		logger.Error.Printf("strategy run error: %v\n", err)
	}
}
