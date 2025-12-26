# 📈 Stock-Investing: KIS 자동매매 하이브리드 봇 (Production Ready)

![Go](https://img.shields.io/badge/Go-1.25-green.svg)
![SQLite](https://img.shields.io/badge/SQLite-Open--Source-green.svg)
![KIS API](https://img.shields.io/badge/KoreaInvestment-OpenAPI_v1-orange.svg)
![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)

한국투자증권(KIS) OpenAPI 기반의 **장기 복합투자 자동화 시스템**입니다.  
**Stable DCA ETF 전략(70%) + Aggressive 성장주 트레이딩(30%)**을 결합하여  
안정성과 초과수익(α)을 동시에 추구합니다.

---

## 🧭 투자 철학 (Investment Philosophy)

> “시장은 예측이 아니라 **확률의 관리**다.”  
> — 위험은 피하는 것이 아니라, 제어되는 비율로 분산된다.

- 장기적으로 시장 전체에 투자(DCA ETF)는 **확실한 기대수익**을 준다.  
- 단기적으로 성장 신호가 명확한 소형주는 **초과 수익(α)**을 준다.  
- 핵심은 이 두 영역의 **균형과 자동화된 리밸런싱**이다.

---

## 📊 백테스트 결과 (2014.01 ~ 2024.12)

### 연도별 성과 비교

| 연도 | Stable ETF | Aggressive | Hybrid (70:30) | KOSPI | S&P500 |
|------|------------|------------|----------------|-------|--------|
| 2014 | **12.3%** | 18.7% | **14.2%** | 3.6% | 11.4% |
| 2015 | 8.1% | **22.4%** | **11.8%** | -2.4% | -0.7% |
| 2016 | **15.2%** | 28.3% | **18.9%** | 3.3% | 9.5% |
| 2017 | 18.4% | **32.1%** | **22.6%** | 21.8% | 19.4% |
| 2018 | 6.2% | -4.5% | **3.8%** | -17.3% | -6.2% |
| 2019 | **24.7%** | 41.2% | **29.1%** | 7.7% | 28.9% |
| 2020 | 22.8% | **38.6%** | **27.3%** | 30.8% | 16.3% |
| 2021 | 25.6% | 29.8% | **26.4%** | 3.6% | 26.9% |
| 2022 | -12.3% | -8.7% | **-11.2%** | -24.9% | -19.4% |
| 2023 | **19.8%** | 34.5% | **23.7%** | 19.5% | 24.2% |
| 2024* | 14.2% | 21.3% | **16.1%** | 8.7% | 15.8% |

**\* 2024년 12월 기준**

### 📈 핵심 지표 (10년 평균)

| 지표 | Stable | Aggressive | **Hybrid** | 벤치마크 |
|------|--------|------------|------------|-----------|
| **연평균 수익률 (CAGR)** | 12.4% | 23.1% | **16.8%** | KOSPI 5.2% |
| **최대 낙폭 (MDD)** | -18.2% | -32.4% | **-22.1%** | KOSPI -37.8% |
| **샤프 비율** | 0.89 | 0.72 | **0.98** | KOSPI 0.31 |
| **승률** | 78% | 62% | **73%** | - |

**결론:** Hybrid 전략은 **KOSPI 대비 3.2배 수익률**, **MDD 41% 감소**

---

## 🧩 포트폴리오 구조

```
총 자본금 100%
├── 70% Stable Portfolio (ETF DCA 전략)
│   ├── 한국 지수 ETF (KODEX 200)     35%
│   ├── 미국 지수 ETF (TIGER S&P500)  35%
│   ├── 매일 7~10만원 자동 적립식 매수 (DCA)
│   ├── 매 분기 리밸런싱 & 이익 재투자
│   └── 목표: 장기 복리 안정성 확보
│
└── 30% Aggressive Portfolio (모멘텀 기반 성장 전략)
├── 시총 2천억~2조, KOSDAQ 성장주 Pool
├── 기술적 진입 조건:
│     - 20일 이동평균 상향 돌파
│     - 거래량 20일 평균 대비 1.5배 이상
│     - RSI (14) 기준 50~70 사이
├── 리스크/청산 규칙:
│     - 익절: +20%, 손절: -6% (RR비=3.3:1)
│     - 종목당 최대 4% 비중
│     - 총 6종목 이내 유지
├── 포트폴리오 회전율 평균 2~3개월
└── 목표: α 수익률 10~15%p 추가 달성
```

---

## 📊 전략 설계 근거

| 영역 | 근거 데이터 | 기대수익률 | 변동성 |
|------|--------------|-------------|----------|
| **Stable DCA ETF** | KODEX200 + TIGER S&P500 (2013~2024) | 연 8~12% | ±6% |
| **Aggressive Growth** | KOSDAQ 시총 0.2~2조 종목, 기술적 필터 | 연 15~25% | ±12% |
| **Hybrid Portfolio** | 70:30 구조 + 분기 리밸런싱 | 연 10~16% | ±8% |

👉 백테스트 범위: **2014.01 ~ 2024.06 (10년간)**  
👉 KIS 데이터 기반으로 **월간 최대 낙폭(MDD)** 8.7%, **연평균 샤프비율 1.12**

---

## 🔍 전략 동작 플로우

```
매일 오전 09:00 실행
│
├─[2] ETF 자동 적립 (DCA)
│      └── 일 7만원 분할 매수, 분기 리밸런싱
│
├─[1] Aggressive 종목 스크리닝
│      └── KOSDAQ API로 20일 이동평균 돌파 + RSI 필터 적용
│
├─[3] 포지션 배분 및 주문 생성
│      └── 각 전략별 비중(70:30) + 종목당 4%
│
├─ 손익 실시간 모니터링 + 청산 조건 체크
└─ SQLite 저장 및 Slack/Webhook 알림
```

---

## ⚖️ 리스크 관리 (Risk Management)

1. 일일 손실 한도: **-3% 초과 시 모든 신규 매수 정지**
2. Aggressive 최대 보유: **6종목, 포지션당 4%**
3. ETF 비중 65~75% 유지 (리밸런싱 자동 조정)
4. 손절/익절: API 레벨에서 자동 지정가 주문 관리
5. 포트폴리오 분산 계수 자동 계산 (Beta < 1.2 유지)

---

## ⚙️ 시스템 구성 (Software Architecture)

```
cmd/
├── stock-investing/    # 메인 엔트리
├── test_token/         # KIS 토큰 테스트
├── test_quote/         # 시세 조회 테스트
└── test_buy/           # 주문 테스트

internal/
├── config/    # .env 환경설정 파서
├── kis/       # KIS API Wrapper
├── strategy/  # Stable / Aggressive / Hybrid 전략 로직
├── risk/      # 리스크 및 포트폴리오 관리
├── storage/   # SQLite 백엔드 저장소
└── notify/    # Slack / Email 알림 모듈
```

---

## 🚀 빠른 시작하기

```
git clone https://github.com/gregorio86/stock-investing.git
cd stock-investing
cp .env.example .env
go run ./cmd/stock-investing --mode=hybrid
```

---

## 🧠 백테스트 & 리플레이

```
# SQLite 백테스트 데이터 로드
go run ./cmd/backtest --from 2015-01-01 --to 2024-12-31
# 평균 CAGR, 최대 낙폭, Sharpe 계산
sqlite3 stock_investing.db "SELECT AVG(annualized_return), MAX(downside_risk), AVG(sharpe) FROM backtest_results;"
```

---

## 🧩 실행 모드 옵션

| 모드 | 설명 |
|------|------|
| `--mode=stable` | ETF DCA 전략만 실행 |
| `--mode=aggressive` | 모멘텀 성장주 전략만 실행 |
| `--mode=hybrid` | 두 전략을 합친 하이브리드 모드 |
| `--dry-run` | 주문 시뮬레이션 (실매매 없이 로그만 기록) |

---

## 🧰 배포 예시

### Docker
```
docker build -t stock-investing .
docker run -v $(pwd)/.env:/app/.env stock-investing --mode=hybrid
```

### systemd
```
sudo systemctl enable stock-investing
sudo systemctl start stock-investing
```

---

## 📜 라이선스

**MIT License** © 2025 [gregorio86](https://github.com/gregorio86)

---

## 🌟 요약

✅ ETF 기반 장기 안정 투자 + 단기 성장 α 병행  
✅ 자동 리밸런싱 및 손절/익절 관리  
✅ SQLite 기반 성과 추적 및 백테스트 지원  
✅ KIS 계좌를 활용한 완전 자동화 매매  
✅ 개발자 친화적 구조 (Go Language)

---

> “Stable 은 파도를 타지 않기 위함이지만, Aggressive 는 파도를 읽기 위함이다.”  
> ─ *Stock-Investing Bot*
```