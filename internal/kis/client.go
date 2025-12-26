package kis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"stock-investing/pkg/logger"
)

type Client struct {
	auth      *AuthClient
	baseURL   string
	accountNo string

	httpClient *http.Client
}

func NewClient(appKey, appSecret, baseURL, accountNo string) *Client {
	return &Client{
		auth:      NewAuthClient(appKey, appSecret, baseURL, accountNo),
		baseURL:   baseURL,
		accountNo: accountNo,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// 공통 HTTP GET 호출 래퍼 (시세조회용)
func (c *Client) doGet(ctx context.Context, path string, query string, trID string, out interface{}) error {
	tok, err := c.auth.GetToken(ctx)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	if query != "" {
		url = url + "?" + query
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("authorization", "Bearer "+tok.AccessToken)
	req.Header.Set("appkey", c.auth.appKey)
	req.Header.Set("appsecret", c.auth.appSecret)
	req.Header.Set("tr_id", trID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error.Printf("[kis] GET %s status=%d\n", path, resp.StatusCode)
		return fmt.Errorf("GET %s failed: status %d", path, resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

// ==== 국내주식 현재가 시세 예시 ====

type quoteOutput struct {
	// 실제 문서 기준 필드 이름/타입 확인 필요
	StockPrice string `json:"stck_prpr"` // 문자열로 내려오면 string으로 받고 나중에 float 변환
}

type quoteResponse struct {
	Output quoteOutput `json:"output"`
}

func (c *Client) GetQuote(ctx context.Context, code string) (float64, error) {
	// 문서 기준 실제 path/TR_ID로 교체 필요
	path := "/uapi/domestic-stock/v1/quotations/inquire-price"
	query := fmt.Sprintf("fid_cond_mrkt_div_code=J&fid_input_iscd=%s", code)
	trID := "FHKST01010100" // 예시: 국내주식 현재가 TR_ID (문서 확인)

	var resp quoteResponse
	if err := c.doGet(ctx, path, query, trID, &resp); err != nil {
		return 0, err
	}

	// stck_prpr를 float으로 변환 (문서에 따라 형식 확인)
	price, err := parsePrice(resp.Output.StockPrice)
	if err != nil {
		return 0, err
	}
	return price, nil
}

// ===== 주문 예시 (현금 매수) =====

type orderRequest struct {
	// 문서 기준으로 필드 정의
	// 예: 계좌번호, 종목코드, 주문수량, 주문가격, 매수/매도 구분 등
	//   CANO: 계좌번호 앞 8자리
	//   ACNT_PRDT_CD: 계좌 상품 코드(뒤 2자리, 예: "01")
	//   PDNO: 종목코드
	//   ORD_DVSN: 주문구분 (00: 지정가 등)
	//   ORD_QTY: 수량
	//   ORD_UNPR: 주문단가
}

type orderResponse struct {
	// rescode, resmsg, output 등 문서 기준으로 정의
}

func (c *Client) getHashKey(ctx context.Context, body []byte) (string, error) {
	tok, err := c.auth.GetToken(ctx)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/uapi/hashkey", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("authorization", "Bearer "+tok.AccessToken)
	req.Header.Set("appkey", c.auth.appKey)
	req.Header.Set("appsecret", c.auth.appSecret)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("hashkey failed: status %d", resp.StatusCode)
	}

	var hk struct {
		Hash string `json:"HASH"` // 문서에서 필드명 확인
	}
	if err := json.NewDecoder(resp.Body).Decode(&hk); err != nil {
		return "", err
	}
	return hk.Hash, nil
}

func (c *Client) Buy(ctx context.Context, code string, quantity int64) error {
	tok, err := c.auth.GetToken(ctx)
	if err != nil {
		return err
	}

	// 계좌번호 분리 (8자리 계좌번호라면 "01" 고정)
	cano := c.accountNo // 8자리 그대로
	acntPrdtCd := "01"  // 종합계좌 고정

	logger.Info.Printf("[kis] account: CANO=%s ACNT_PRDT_CD=%s", cano, acntPrdtCd)

	// KIS 공식 현금주문 Body (TR_ID: TTTC0802U 모의투자용)
	reqBody := map[string]interface{}{
		"CANO":         cano,       // 계좌번호 앞 8자리
		"ACNT_PRDT_CD": acntPrdtCd, // "01"
		"PDNO":         code,       // 종목코드 (6자리)
		"ORD_DVSN":     "01",       // 주문구분: 지정가
		"ORD_QTY":      fmt.Sprintf("%d", quantity),
		"ORD_UNPR":     "0", // 시장가
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// 1) hashkey 먼저 생성
	hash, err := c.getHashKey(ctx, bodyBytes)
	if err != nil {
		logger.Error.Printf("[kis] hashkey failed: %v", err)
		return err
	}

	// 2) 주문 요청 (모의투자 TR_ID: TTTC0802U)
	path := "/uapi/domestic-stock/v1/trading/order-cash"
	url := c.baseURL + path

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	httpReq.Header.Set("authorization", fmt.Sprintf("Bearer %s", tok.AccessToken))
	httpReq.Header.Set("appkey", c.auth.appKey)
	httpReq.Header.Set("appsecret", c.auth.appSecret)
	httpReq.Header.Set("tr_id", "VTTC0802U") // 모의투자 현금매수 TR_ID
	httpReq.Header.Set("hashkey", hash)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json;v=1.0")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyResp, _ := io.ReadAll(resp.Body)
	logger.Info.Printf("[kis] Buy response: status=%d body=%s", resp.StatusCode, string(bodyResp))

	if resp.StatusCode != http.StatusOK {
		logger.Error.Printf("[kis] Buy failed: status=%d body=%s", resp.StatusCode, string(bodyResp))
		return fmt.Errorf("Buy failed: status %d", resp.StatusCode)
	}

	logger.Info.Printf("[kis] Buy SUCCESS: %s x %d", code, quantity)
	return nil
}

func (c *Client) Sell(ctx context.Context, code string, quantity int64) error {
	// Buy와 거의 동일, TR_ID/주문구분만 매도용으로 변경
	logger.Info.Printf("[kis] Sell stub for %s x %d (TODO: implement)\n", code, quantity)
	return nil
}

func parsePrice(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
