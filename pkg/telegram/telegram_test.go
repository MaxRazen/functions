package telegram

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	token = "testtoken"
)

func TestNew(t *testing.T) {
	tc := New(token)
	if tc.token != token {
		t.Errorf("token is not equal to expected %s != %s", token, tc.token)
	}
	if tc.BaseURL == "" {
		t.Errorf("BaseURL is not set automatically")
	}
}

func TestSendMessage(t *testing.T) {
	expectedBodyData := `{"chat_id":"-1001","text":"testmessage","parse_mode":"HTML","disable_notification":true}`
	expectedPath := fmt.Sprintf("/bot%s/sendMessage", token)
	responseBody := `{"message_id":20001,"ok":true}`

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		reqPath := req.URL.String()
		body, _ := io.ReadAll(req.Body)

		if reqPath != expectedPath {
			t.Errorf("expected body is not equal to actual:\n%s\n%s", expectedPath, reqPath)
		}

		if string(body) != expectedBodyData {
			t.Errorf("expected body is not equal to actual:\n%s\n%s", expectedBodyData, string(body))
		}

		rw.Write([]byte(responseBody))
	}))
	defer server.Close()

	tm := Message{
		ChatId:    "-1001",
		Message:   "testmessage",
		ParseMode: "HTML",
		Silent:    true,
	}

	tc := New(token)
	tc.BaseURL = server.URL + "/bot"
	tc.HTTPClient = server.Client()
	ctx := context.Background()

	resp, err := tc.SendMessage(ctx, tm)

	if err != nil {
		t.Errorf("SendMessage method returned unexpected error %v", err)
	}
	if resp.Status != 200 {
		t.Errorf("SendMessage method returned unexpected status code %v", resp.Status)
	}
	if string(resp.Body) != responseBody {
		t.Errorf("SendMessage method returned unexpected body %s", string(resp.Body))
	}
}
