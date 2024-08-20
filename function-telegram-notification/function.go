package p

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/MaxRazen/functions/pkg/telegram"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type Message struct {
	Header  string `json:"header"`
	Content string `json:"content"`
	Channel string `json:"channel"`
	Urgency string `json:"urgency"`
}

type MessageSender interface {
	SendMessage(context.Context, telegram.Message) (*telegram.Response, error)
}

const (
	// deliver without sound notifications
	Info = "info"
	// deliver with sound notifications
	Warning = "warning"
	Alert   = "alert"
)

var telegramClient MessageSender

// Prepares heap variables
func init() {
	telegramClient = telegram.New(os.Getenv("BOT_TOKEN"))
}

// Handler consumes a Pub/Sub message.
func Handler(ctx context.Context, m PubSubMessage) error {
	log.Printf("[INFO] function is invoked with data: %s", string(m.Data))

	var msg Message
	if err := json.Unmarshal(m.Data, &msg); err != nil {
		// mailformed json
		return err
	}

	if err := Validate(&msg); err != nil {
		// validtion error
		return err
	}

	err := Send(ctx, &msg)
	if err != nil {
		return err
	}

	return nil
}

var (
	ErrValidation = errors.New("validation error")
)

// Validates the given payload
func Validate(msg *Message) error {
	if msg.Content == "" {
		return errors.Join(ErrValidation, fmt.Errorf("property content must be present and not be empty"))
	}
	return nil
}

// Sends a message using SenderMessage interface instance
func Send(ctx context.Context, msg *Message) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	content := msg.Content
	if msg.Header != "" {
		content = fmt.Sprintf("*%s*\n\n%s", msg.Header, content)
	}

	tgMsg := telegram.Message{
		ChatId:    GetChatId(msg.Channel),
		Message:   content,
		ParseMode: "Markdown",
		Silent:    msg.Urgency == Info,
	}

	res, err := telegramClient.SendMessage(ctx, tgMsg)
	if err != nil {
		return err
	}

	log.Printf("[INFO] message is sent: %s", string(res.Body))

	return nil
}

// Resolves valid chat id by the channel name using environment variables
func GetChatId(channel string) string {
	chatId := os.Getenv(strings.ToUpper(channel) + "_CHAT_ID")
	if chatId != "" {
		return chatId
	}

	return os.Getenv("DEFAULT_CHAT_ID")
}
