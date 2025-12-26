package kis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"stock-investing/pkg/logger"
)

type Token struct {
	AccessToken string
	ExpiresAt   time.Time
}

type AuthClient struct {
	appKey    string
	appSecret string
	baseURL   string
	accountNo string

	httpClient *http.Client

	mu    sync.Mutex
	token *Token
}

func NewAuthClient(appKey, appSecret, baseURL string, accountNo string) *AuthClient {
	return &AuthClient{
		appKey:    appKey,
		appSecret: appSecret,
		baseURL:   baseURL,
		accountNo: accountNo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type tokenPRequest struct {
	GrantType string `json:"grant_type"`
	AppKey    string `json:"appkey"`
	AppSecret string `json:"appsecret"`
	AccountNo string `json:"account_no"`
	// 문서에 따라 필요한 필드 추가
}

type tokenPResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"` // 초 단위라면
	// 에러코드/메시지 필드도 문서 보고 추가
}

func (a *AuthClient) GetToken(ctx context.Context) (*Token, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 1) 캐시된 토큰이 유효하면 재사용
	if a.token != nil && time.Now().Before(a.token.ExpiresAt) {
		return a.token, nil
	}

	// 2) 새 토큰 발급
	reqBody := tokenPRequest{
		GrantType: "client_credentials",
		AppKey:    a.appKey,
		AppSecret: a.appSecret,
		AccountNo: a.accountNo,
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/oauth2/tokenP", a.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	// 응답 바디를 먼저 읽어서 로깅
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Printf("[kis] failed to read response body: %v", err)
		return nil, fmt.Errorf("tokenP failed: status %d", resp.StatusCode)
	}
	logger.Error.Printf("[kis] tokenP status=%d body=%s\n", resp.StatusCode, string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tokenP failed: status %d body=%s", resp.StatusCode, string(bodyBytes))
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error.Printf("[kis] tokenP status=%d body=%s\n", resp.StatusCode, resp.Body)
		return nil, fmt.Errorf("tokenP failed: status %d", resp.StatusCode)
	}

	// 성공 시 bodyBytes를 다시 사용해서 JSON 파싱
	var tr tokenPResponse
	if err := json.Unmarshal(bodyBytes, &tr); err != nil {
		return nil, err
	}
	if tr.AccessToken == "" {
		return nil, fmt.Errorf("tokenP: empty access_token")
	}

	// expires_in 몇 초 전에 미리 갱신할지 여유를 둔다 (예: 60초)
	exp := time.Now().Add(time.Duration(tr.ExpiresIn-60) * time.Second)

	a.token = &Token{
		AccessToken: tr.AccessToken,
		ExpiresAt:   exp,
	}
	logger.Info.Printf("[kis] token acquired, expires at %s\n", exp.Format(time.RFC3339))

	return a.token, nil
}
