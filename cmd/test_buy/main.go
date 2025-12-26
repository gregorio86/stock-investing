package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"stock-investing/internal/config"
	"stock-investing/internal/kis"
	"stock-investing/pkg/logger"
)

func main() {
	_ = godotenv.Load(".env")

	logger.Init()
	cfg := config.Load()

	if len(os.Args) < 3 {
		fmt.Println("usage: test_buy <종목코드> <수량>")
		os.Exit(1)
	}
	code := os.Args[1]
	qtyStr := os.Args[2]

	qty, err := strconv.ParseInt(qtyStr, 10, 64)
	if err != nil || qty <= 0 {
		fmt.Fprintf(os.Stderr, "invalid quantity: %s\n", qtyStr)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := kis.NewClient(
		cfg.KIS.AppKey,
		cfg.KIS.AppSecret,
		cfg.KIS.BaseURL,
		cfg.KIS.AccountNo,
	)

	fmt.Printf("TRY BUY (mock): code=%s, qty=%d\n", code, qty)

	if err := client.Buy(ctx, code, qty); err != nil {
		fmt.Fprintf(os.Stderr, "Buy error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Buy order sent. 모의투자 HTS/앱에서 체결 내역 확인해봐.")
}
