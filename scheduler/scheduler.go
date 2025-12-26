package scheduler

import (
	"context"
	"time"

	"stock-investing/pkg/logger"
)

// 매우 단순한 틱 기반 스케줄러 (실제 cron 대체는 이후 단계에서 개선)
type Scheduler struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func New() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Scheduler) Start() {
	logger.Info.Println("scheduler started (stub)")
	// TODO: 이후 cron 기반으로 시간별 작업 연결
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-s.ctx.Done():
				logger.Info.Println("scheduler stopped")
				return
			case <-ticker.C:
				logger.Info.Println("scheduler tick (stub)")
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	s.cancel()
}
