package main

import (
	"context"
	"fmt"
	"os"
	"stock-investing/internal/config"
	"stock-investing/internal/kis"
	"stock-investing/pkg/logger"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// 1) .env 로드
	_ = godotenv.Load(".env")

	// 2) 로거 + 설정
	logger.Init()
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := kis.NewAuthClient(cfg.KIS.AppKey, cfg.KIS.AppSecret, cfg.KIS.BaseURL)

	tok, err := auth.GetToken(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetToken error: %v\n", err)
		os.Exit(1)
	}

	access := tok.AccessToken
	if len(access) > 20 {
		access = access[:20] + "..."
	}
	fmt.Printf("access_token: %s\n", access)
	fmt.Printf("expires_at : %s\n", tok.ExpiresAt.Format(time.RFC3339))
}
