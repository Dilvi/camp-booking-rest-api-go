package gigachat

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Client struct {
	authKey            string
	scope              string
	model              string
	oauthURL           string
	chatCompletionsURL string
	httpClient         *http.Client

	mu          sync.Mutex
	accessToken string
	expiresAt   time.Time
}

type Config struct {
	AuthKey            string
	Scope              string
	Model              string
	OAuthURL           string
	ChatCompletionsURL string
	InsecureSkipVerify bool
}

func NewClient(cfg Config) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if cfg.InsecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}

	return &Client{
		authKey:            cfg.AuthKey,
		scope:              cfg.Scope,
		model:              cfg.Model,
		oauthURL:           cfg.OAuthURL,
		chatCompletionsURL: cfg.ChatCompletionsURL,
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}
}

func (c *Client) Complete(ctx context.Context, prompt string) (string, error) {
	if strings.TrimSpace(c.authKey) == "" {
		return "", errors.New("gigachat auth key is not configured")
	}

	result, err := c.complete(ctx, prompt)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (c *Client) complete(ctx context.Context, prompt string) (string, error) {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return "", err
	}

	result, err := c.completeWithToken(ctx, prompt, token)
	if err == nil {
		return result, nil
	}

	var apiErr apiError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusUnauthorized {
		return "", err
	}

	c.clearAccessToken()

	token, err = c.getAccessToken(ctx)
	if err != nil {
		return "", err
	}

	return c.completeWithToken(ctx, prompt, token)
}

func (c *Client) completeWithToken(ctx context.Context, prompt, token string) (string, error) {
	reqBody := chatCompletionRequest{
		Model: c.model,
		Messages: []chatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		N:                 1,
		Stream:            false,
		MaxTokens:         900,
		RepetitionPenalty: 1,
		UpdateInterval:    0,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.chatCompletionsURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", apiError{
			Operation:  "completion",
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	var completion chatCompletionResponse
	if err := json.Unmarshal(respBody, &completion); err != nil {
		return "", err
	}
	if len(completion.Choices) == 0 {
		return "", errors.New("gigachat returned empty choices")
	}

	return completion.Choices[0].Message.Content, nil
}

func (c *Client) clearAccessToken() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.accessToken = ""
	c.expiresAt = time.Time{}
}

func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.accessToken != "" && time.Now().Before(c.expiresAt.Add(-time.Minute)) {
		return c.accessToken, nil
	}

	form := url.Values{}
	form.Set("scope", c.scope)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.oauthURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", newRqUID())
	req.Header.Set("Authorization", normalizeAuthHeader(c.authKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("gigachat token request failed: status %d: %s", resp.StatusCode, string(respBody))
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return "", err
	}
	if tokenResp.AccessToken == "" {
		return "", errors.New("gigachat token response is empty")
	}

	c.accessToken = tokenResp.AccessToken
	c.expiresAt = expiresAtFromGigaChat(tokenResp.ExpiresAt)

	return c.accessToken, nil
}

func expiresAtFromGigaChat(expiresAt int64) time.Time {
	if expiresAt == 0 {
		return time.Now().Add(30 * time.Minute)
	}

	// GigaChat returns expires_at as Unix time in milliseconds.
	if expiresAt > 1_000_000_000_000 {
		return time.UnixMilli(expiresAt)
	}

	return time.Unix(expiresAt, 0)
}

func normalizeAuthHeader(authKey string) string {
	authKey = strings.TrimSpace(authKey)
	if strings.HasPrefix(strings.ToLower(authKey), "basic ") {
		return authKey
	}
	return "Basic " + authKey
}

func newRqUID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return strings.Join([]string{
		hex.EncodeToString(b[0:4]),
		hex.EncodeToString(b[4:6]),
		hex.EncodeToString(b[6:8]),
		hex.EncodeToString(b[8:10]),
		hex.EncodeToString(b[10:16]),
	}, "-")
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

type chatCompletionRequest struct {
	Model             string        `json:"model"`
	Messages          []chatMessage `json:"messages"`
	N                 int           `json:"n"`
	Stream            bool          `json:"stream"`
	MaxTokens         int           `json:"max_tokens"`
	RepetitionPenalty float64       `json:"repetition_penalty"`
	UpdateInterval    int           `json:"update_interval"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

type apiError struct {
	Operation  string
	StatusCode int
	Body       string
}

func (e apiError) Error() string {
	return fmt.Sprintf("gigachat %s failed: status %d: %s", e.Operation, e.StatusCode, e.Body)
}
