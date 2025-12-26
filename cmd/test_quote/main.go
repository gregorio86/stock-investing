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
	_ = godotenv.Load(".env")

	logger.Init()
	cfg := config.Load()

	if len(os.Args) < 2 {
		fmt.Println("usage: test_quote <종목코드>")
		os.Exit(1)
	}
	code := os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := kis.NewClient(
		cfg.KIS.AppKey,
		cfg.KIS.AppSecret,
		cfg.KIS.BaseURL,
		cfg.KIS.AccountNo,
	)

	price, err := client.GetQuote(ctx, code)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetQuote error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("code=%s, price=%.2f\n", code, price)
}
