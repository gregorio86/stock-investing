package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	KIS         KISConfig
	Stable      StableConfig
	Aggressive  AggressiveConfig
	Risk        RiskConfig
	MockTrading bool
}

type KISConfig struct {
	BaseURL   string
	AppKey    string
	AppSecret string
	AccountNo string
}

type StableConfig struct {
	Alloc         float64
	ETFs          []string
	DailyAmount   int64
	RebalanceDays int
}

type AggressiveConfig struct {
	Alloc float64
}

type RiskConfig struct {
	MaxRisk float64
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env %s is not set", key)
	}
	return v
}

func getEnvFloat(key string, def float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		log.Fatalf("invalid float env %s: %v", key, err)
	}
	return f
}

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid int env %s: %v", key, err)
	}
	return i
}

func Load() *AppConfig {
	rawMock := strings.TrimSpace(os.Getenv("MOCK_TRADING"))
	mock := strings.EqualFold(rawMock, "true")

	log.Printf("[config] MOCK_TRADING=%q -> mock=%v\n", rawMock, mock)

	var kisCfg KISConfig
	if mock {
		kisCfg = KISConfig{
			BaseURL:   mustEnv("KIS_BASE_URL"),
			AppKey:    mustEnv("APP_KEY_PAPER"),
			AppSecret: mustEnv("APP_SECRET_PAPER"),
			AccountNo: mustEnv("ACCOUNT_NO_PAPER"),
		}
	} else {
		kisCfg = KISConfig{
			BaseURL:   mustEnv("KIS_BASE_URL_LIVE"),
			AppKey:    mustEnv("APP_KEY_LIVE"),
			AppSecret: mustEnv("APP_SECRET_LIVE"),
			AccountNo: mustEnv("ACCOUNT_NO_LIVE"),
		}
	}

	etfsEnv := os.Getenv("STABLE_ETFS")
	etfs := []string{"069500", "360750"}
	if etfsEnv != "" {
		parts := strings.Split(etfsEnv, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		etfs = parts
	}

	return &AppConfig{
		KIS: kisCfg,
		Stable: StableConfig{
			Alloc:         getEnvFloat("STABLE_ALLOC", 0.7),
			ETFs:          etfs,
			DailyAmount:   70000,
			RebalanceDays: getEnvInt("REBALANCE_PERIOD", 90),
		},
		Aggressive: AggressiveConfig{
			Alloc: getEnvFloat("AGGRESSIVE_ALLOC", 0.3),
		},
		Risk: RiskConfig{
			MaxRisk: getEnvFloat("MAX_RISK", 0.1),
		},
		MockTrading: mock,
	}
}
