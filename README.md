# 🚀 stock-investing (KIS Hybrid Stock-Investing Bot)

한국투자증권(KIS) API를 활용한 **70% 안정적 DCA ETF + 30% 고변동성 중소형주** 하이브리드 자동매매 봇

[![Go](https://img.shields.io/badge/Go-1.25-green.svg)](https://golang.org)
[![SQLite](https://img.shields.io/badge/SQLite-완전_오프라인-green.svg)](https://sqlite.org)
[![KIS API](https://img.shields.io/badge/KIS-API_v1-orange.svg)](https://apiportal.koreainvestment.com)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## 🎯 핵심 하이브리드 전략

```
총 자본금
├── 70% 안정적 전략 (DCA ETF)
│   ├── KODEX 200 (069500): 35%
│   ├── TIGER S&P500 (360750): 35%
│   └── 매일 7만원 자동 적립 + 분기 리밸런싱
└── 30% 공격적 전략 (고변동성 중소형)
├── 대상: 시총 2천억~2조 KOSDAQ·중소형 성장주
├── 테마: 2차전지소재/반도체장비/AI/로봇/바이오
├── 신호: 10일 +15%↑ + 거래량 2배↑ + 변동성 4%↑
├── 익절: +12~15% / 손절: -7% / 최대 15일 보유
└── 최대 8종목 (종목당 3~5%)
```

---

## 📁 프로젝트 구조

```text
stock-investing/
├── cmd/
│   └── stock-investing/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── kis/
│   │   ├── client.go
│   │   └── auth.go
│   ├── models/
│   │   └── types.go
│   ├── screener/
│   │   └── screener.go
│   ├── strategy/
│   │   ├── hybrid.go
│   │   ├── stable.go
│   │   └── aggressive.go
│   ├── risk/
│   │   └── manager.go
│   └── storage/
│       ├── sqlite.go
│       └── repository.go
├── pkg/
│   └── logger/
│       └── logger.go
├── scheduler/
│   └── scheduler.go
├── data/
│   ├── kosdaq_universe.csv
│   └── .gitkeep
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── main.go

## 폴더별로 들어갈 파일들

### 루트 (`stock-investing/`)

- `main.go`
    - `cmd/stock-investing/main.go`를 그냥 래핑하거나, 바로 `internal/config`, `internal/strategy`를 불러서 실행 엔트리로 사용.[1][2]
- `.env.example`
    - `KIS_APP_KEY`, `KIS_APP_SECRET`, `KIS_ACCOUNT_NO`, `MOCK_TRADING`, `STABLE_ALLOC` 등 환경변수 샘플 정의.[1][2]
- `.gitignore`
    - `.env`, `*.db`, `*.log`, `data/*` 등 민감/런타임 파일 제외.[1][2]
- `go.mod`, `go.sum`
    - Go 모듈 및 의존성.[1][2]

### `cmd/stock-investing/`

- `main.go`
    - `flag`/`cobra`로 `--mode=hybrid|stable|aggressive`, `--init-db`, `--backtest` 등을 파싱하고 `internal` 패키지를 호출하는 진짜 엔트리 포인트.[1][2]

### `internal/config/`

- `config.go`
    - `.env` 로드, `os.Getenv`/`godotenv`로 환경변수 파싱.
    - 구조체 예: `AppConfig{KIS, StableConfig, AggressiveConfig, RiskConfig}`.[1][2]

### `internal/kis/`

- `client.go`
    - KIS REST 호출 클라이언트, 주문, 시세, 계좌조회 메서드.[1][2]
- `auth.go`
    - 토큰 발급/갱신, 헤더 생성 등 인증 관련 helper.[2]

### `internal/models/`

- `types.go`
    - `Stock`, `Quote`, `Order`, `Trade`, `Position`, `Portfolio` 등 도메인 모델 정의.[1][2]

### `internal/screener/`

- `screener.go`
    - KOSDAQ/중소형 종목을 불러와서 변동성, 모멘텀, 거래량 필터를 적용해 상위 종목 리스트 반환.[1][2]

### `internal/strategy/`

- `stable.go`
    - KODEX 200, TIGER S&P500 DCA 매수 로직, 리밸런싱 로직.[1][2]
- `aggressive.go`
    - 10일 +15%, 변동성 4% 이상, 거래량 2배 등 조건 기반 모멘텀 전략 구현.[2]
- `hybrid.go`
    - `stable` + `aggressive`를 `STABLE_ALLOC`, `AGGRESSIVE_ALLOC` 비율로 호출하는 메인 전략.[1][2]

### `internal/risk/`

- `manager.go`
    - 최대 손실률, 포지션당 비중, 테마 집중도, 스탑로스/익절 로직 구현.[1][2]

### `internal/storage/`

- `sqlite.go`
    - SQLite 연결, 마이그레이션, 테이블 생성 (`trades`, `positions`, `daily_pnl` 등).[1][2]
- `repository.go`
    - 트레이드 저장/조회, 포지션 조회용 함수들 (서비스 레이어에서 사용).[1][2]

### `pkg/logger/`

- `logger.go`
    - 공통 로거 설정 (레벨, 파일 출력, JSON 포맷 등).

### `scheduler/`

- `scheduler.go`
    - cron-like 스케줄 정의: 09:00 스크리닝, 10:00 DCA, 18:30 청산 체크 등.[2]

### `data/`

- `kosdaq_universe.csv`
    - 종목코드, 종목명, 섹터 등의 기본 유니버스.[1][2]
- `.gitkeep`
    - 빈 폴더라도 Git에 남기기 위한 더미 파일.
```

---

## 🚀 빠른 시작

### 1. 환경 설정

```
git clone <repo>
cd stock-investing

go mod tidy
mkdir -p data

cp .env.example .env
# 에디터로 APP_KEY, APP_SECRET, ACCOUNT_NO 등 수정
```

### 2. KOSDAQ/중소형 데이터 생성 예시

```
cat > data/kosdaq_universe.csv << EOF
종목코드,종목명
300070,에이비온
300720,빅텍
317850,쎄트
EOF
```

### 3. KIS Developers / 모의투자 계좌

1. KIS Developers 회원가입 및 앱 등록 (APP_KEY / APP_SECRET 발급)
2. 모의투자 계좌 개설 후 ACCOUNT_NO 확인
3. `.env` 파일에 아래와 같이 설정

```
APP_KEY=your_app_key
APP_SECRET=your_app_secret
ACCOUNT_NO=your_mock_account_no
MOCK_TRADING=true        # 모의: true, 실전: false

# 서버 설정
KIS_SERVER=fixi           # 모의: fixi, 실전: v1
```

---

## ⚙️ 주요 환경변수 (.env)

### 하이브리드 비중

```
STABLE_ALLOC=0.7          # 안정적 전략 70%
AGGRESSIVE_ALLOC=0.3      # 공격적 전략 30%
```

### 안정적 전략 (70%) - DCA ETF

```
STABLE_ETFS="069500,360750"  # KODEX 200, TIGER S&P500
DAILY_STABLE_AMT=70000       # 하루/매수액 또는 환산 월 DCA 금액
REBALANCE_PERIOD=90          # 분기(90일) 리밸런싱
```

### 공격적 전략 (30%) - 중소형 고변동성

```
MIN_MARKETCAP=2e11           # 2천억 이상
MAX_MARKETCAP=2e12           # 2조 미만
MIN_VOLATILITY=0.04          # 변동성 4% 이상
MIN_MOMENTUM=0.15            # 10일 +15% 이상
MIN_VOLUME_RATIO=2.0         # 거래량 2배 이상
MAX_POSITIONS=8
POSITION_SIZE=0.04           # 종목당 4%
MAX_HOLDING_DAYS=15
TARGET_SECTORS="2차전지소재,반도체장비,AI솔루션,로봇,바이오신약"
```

### 공통 리스크 파라미터

```
DCA_AMOUNT=100000            # (옵션) 월 DCA 금액
MAX_RISK=0.1                 # 최대 손실 10% 기준
```

---

## 🕒 실행 및 스케줄

### 기본 실행

```
# DB 초기화
go run main.go --init-db

# 하이브리드 모드
go run main.go --mode=hybrid

# 전략별 단독 테스트
go run main.go --mode=stable     # DCA ETF
go run main.go --mode=aggressive # 고변동성
```

### 스케줄 예시 (cron/스케줄러)

```
09:00 장전   → 중소형 고변동성 스크리닝 (200종목 병렬)
10:00 장중   → 안정적 ETF 적립식 매수
16:00 장중   → 모멘텀 신호 실시간 확인
18:30 장후   → 익절/손절 및 포지션 점검
09:00 월요일 → 테마·비중 리밸런싱
```

---

## 📊 전략 및 성과 (백테스트 기준 예시)

### 전략 구성

- Stable DCA ETF (70%)
  - KODEX 200 (069500): 장기 보유, 시장 평균 노출
  - TIGER S&P500 (360750): 글로벌 분산, 달러 자산 노출

- Aggressive KOSDAQ/중소형 (30%)
  - 2차전지 소재, 반도체 장비, AI, 로봇, 바이오 등 성장 테마
  - 모멘텀·거래량·변동성·섹터 필터로 상위 종목 자동 선별

### 예상 성과 (예시)

| 전략           | 연평균/연수익률 | 최대 손실 | 샤프레시오 | 승률 |
|----------------|----------------|-----------|-----------|------|
| 안정적 (70%)   | 8–12%          | -10%      | 1.2       | 95%  |
| 공격적 (30%)   | 20–45%         | -25%      | 1.1       | 60%  |
| 하이브리드     | 13–22%         | -15%      | 1.15      | 78%  |

> 위 수치는 과거 데이터 기반 시뮬레이션 예시이며, 미래 수익을 보장하지 않습니다.

---

## 🛡️ 리스크 관리

```
✅ 최대 손실: 계좌 기준 -10~15% 시 강제 청산 로직
✅ 테마 집중도: 특정 테마 비중 50% 제한
✅ 포지션당: 총 자산 3~5% / 일일 신규 진입 3개 제한
✅ 보유 기간: 최대 15일 이후 자동 청산
✅ 최소 현금: 20% 이상 유지
✅ 스탑로스: ATR × 2 또는 -7% 도달 시 자동 손절
✅ SQLite로 모든 거래 이력 기록
```

---

## 🧪 테스트 체크리스트

| 단계 | 명령어                                                | 기대 결과          |
|------|-------------------------------------------------------|--------------------|
| 1    | `go mod tidy`                                         | 의존성 설치        |
| 2    | `go run main.go --init-db`                            | DB 초기화/생성     |
| 3    | `go run main.go --mode=hybrid`                        | 봇 하이브리드 실행 |
| 4    | `sqlite3 stock_investing.db "SELECT * FROM trades;"`  | 거래 기록 확인     |

---

## 🔧 개발자용 팁 (Go / GoLand)

```
✅ Goroutine + errgroup: 200종목 동시 스크리닝
✅ Channel: 병렬 신호 처리
✅ Context: Graceful Shutdown
✅ Cron: 자동화 스케줄링
✅ Decimal: 주가 정밀 계산
✅ SQLite: 거래 히스토리 영속화
```

- Debug Config 예시
  1. Mock Trading: `MOCK_TRADING=true`
  2. Live Trading: `MOCK_TRADING=false`
  3. Screener Test: `go test ./internal/screener -v`

---

## ⚠️ 중요 경고

```
🚨 고위험 하이브리드 전략 (공격적 30% 포함)!

- 실투자 전 최소 3개월 모의투자 필수
- 실제 투자 시 전체 금융자산의 10% 이내로 시작
- 대출/레버리지 자금 사용 금지
- 초기 1개월은 매일 포트폴리오·로그 모니터링
- 시장 레짐 변화 시 파라미터 재점검 및 리밸런싱
```

---

## 📚 다음 단계

1. KIS Developers 계정 및 모의투자 계좌 개설
2. `.env` 구성 후, 소액/모의 환경에서 3일 이상 연속 실행 테스트
3. 백테스트/리포트 모듈 (`--backtest`, 성과 리포트 자동 생성) 추가 개발
4. 실전 전환 시 자산 10% 이내에서 점진적 증액

---

## 🏷 메타 정보

- **프로젝트명**: stock-investing
- **대상**: 한국 CTO·개발자 초보 투자자
- **목적**: 장기 자산 증식용 자동화 + 리스크 관리 학습
- **라이선스**: MIT (교육·개인 투자용)
- **최종 업데이트**: 2025-12-08
```
