package p

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/MaxRazen/functions/pkg/telegram"
)

const (
	DefaultChatId = "-1001"
	CustomChatId  = "51002"
)

func TestGetChatId(t *testing.T) {
	os.Setenv("DEFAULT_CHAT_ID", DefaultChatId)
	os.Setenv("CUSTOM_CHAT_ID", CustomChatId)
	var actualChatId string

	actualChatId = GetChatId("")
	if actualChatId != DefaultChatId {
		t.Errorf("chat id for empty channel should be default")
	}

	actualChatId = GetChatId("unknown")
	if actualChatId != DefaultChatId {
		t.Errorf("chat id for unknown channel should be default")
	}

	actualChatId = GetChatId("custom")
	if actualChatId != CustomChatId {
		t.Errorf("chat id for custom channel should be custom")
	}

	actualChatId = GetChatId("CUSTOM")
	if actualChatId != CustomChatId {
		t.Errorf("chat id for CUSTOM channel should be custom")
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		msg Message
		err error
	}{
		{
			msg: Message{
				Header:  "test header",
				Content: "test content",
				Channel: "custom",
				Urgency: Info,
			},
			err: nil,
		},
		{
			msg: Message{
				Content: "test content",
				Urgency: "unknown",
			},
			err: nil,
		},
		{
			msg: Message{
				Content: "test content",
			},
			err: nil,
		},
		{
			msg: Message{},
			err: ErrValidation,
		},
	}

	for i, tc := range testCases {
		err := Validate(&tc.msg)
		if err != tc.err && !errors.Is(err, tc.err) {
			t.Errorf("[%d] expected result does not match with actual\n%v\n%v", i, tc.err, err)
		}
	}
}

type MessageSenderMock struct {
	t     *testing.T
	resp  *telegram.Response
	err   error
	calls int
}

func (msm *MessageSenderMock) SendMessage(_ context.Context, tm telegram.TelegramMessage) (*telegram.Response, error) {
	msm.calls++

	expected := telegram.TelegramMessage{
		ChatId:    CustomChatId,
		Message:   "*ALERT*\n\nmessage text",
		ParseMode: "Markdown",
		Silent:    false,
	}

	if tm != expected {
		msm.t.Errorf("arguments are not valid \n%v\n%v", expected, tm)
	}

	return msm.resp, msm.err
}

func TestHandler(t *testing.T) {
	os.Setenv("CUSTOM_CHAT_ID", CustomChatId)

	logBufer := new(bytes.Buffer)
	log.SetOutput(logBufer)

	ctx := context.Background()
	m := PubSubMessage{
		Data: []byte(`{"header":"ALERT","content":"message text","channel":"custom","urgency":"alert"}`),
	}

	client := &MessageSenderMock{
		t:    t,
		resp: &telegram.Response{Body: []byte(`{"message_id":1007001}`)},
		err:  nil,
	}

	telegramClient = client

	err := Handler(ctx, m)

	if err != nil {
		t.Errorf("tested function returned unexpected error: %v", err)
	}

	if client.calls != 1 {
		t.Errorf("the number or calls are not expected: %d", client.calls)
	}

	t.Logf("\n----------- LOGS -----------\n%s----------- LOGS END -----------\n", logBufer.String())
}
