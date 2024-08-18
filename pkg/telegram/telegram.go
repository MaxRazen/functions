package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type TelegramMessage struct {
	ChatId    string `json:"chat_id"`
	Message   string `json:"text"`
	ParseMode string `json:"parse_mode"`
	Silent    bool   `json:"disable_notification"`
}

type TelegramClient struct {
	HTTPClient *http.Client
	BaseURL    string
	token      string
}

type Response struct {
	Status int
	Body   []byte
}

func New(token string) *TelegramClient {
	return &TelegramClient{
		HTTPClient: http.DefaultClient,
		BaseURL:    "https://api.telegram.org/bot",
		token:      token,
	}
}

func (tc *TelegramClient) SendMessage(ctx context.Context, tm TelegramMessage) (*Response, error) {
	return tc.doRequest(ctx, "sendMessage", tm)
}

func (tc *TelegramClient) constructEndpointUrl(path string) string {
	return tc.BaseURL + tc.token + "/" + path
}

func (tc *TelegramClient) doRequest(ctx context.Context, path string, data any) (*Response, error) {
	endpoint := tc.constructEndpointUrl(path)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	httpResp, err := tc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	resp := Response{
		Status: httpResp.StatusCode,
		Body:   []byte{},
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return &resp, err
	}

	resp.Body = body

	return &resp, nil
}
