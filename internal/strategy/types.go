package strategy

import "context"

type Mode string

const (
	ModeHybrid     Mode = "hybrid"
	ModeStable     Mode = "stable"
	ModeAggressive Mode = "aggressive"
)

// Runner 모든 전략의 공통 인터페이스
type Runner interface {
	Run(ctx context.Context) error
}
