package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Message struct {
	ChatId    string `json:"chat_id"`
	Message   string `json:"text"`
	ParseMode string `json:"parse_mode"`
	Silent    bool   `json:"disable_notification"`
}

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	token      string
}

type Response struct {
	Status int
	Body   []byte
}

// Returns new instance of Telegram client with base URL and given BOT security token
func New(token string) *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    "https://api.telegram.org/bot",
		token:      token,
	}
}

// Calls sendMessage method of API
func (tc *Client) SendMessage(ctx context.Context, tm Message) (*Response, error) {
	return tc.doRequest(ctx, "sendMessage", tm)
}

// Constructs endpoint URL with the proper base URL, token and path
func (tc *Client) constructEndpointUrl(path string) string {
	return tc.BaseURL + tc.token + "/" + path
}

// Performs request to the API
func (tc *Client) doRequest(ctx context.Context, path string, data any) (*Response, error) {
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
