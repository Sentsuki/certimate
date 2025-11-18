package telegrambot_test

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/notifier/providers/telegrambot"
)

const (
	mockSubject = "test_subject"
	mockMessage = "test_message"
)

var (
	fApiToken string
	fChatId   string
)

func init() {
	argsPrefix := "CERTIMATE_NOTIFIER_TELEGRAMBOT_"

	flag.StringVar(&fApiToken, argsPrefix+"APITOKEN", "", "")
	flag.StringVar(&fChatId, argsPrefix+"CHATID", "", "")
}

/*
Shell command to run this test:

	go test -v ./telegrambot_test.go -args \
	--CERTIMATE_NOTIFIER_TELEGRAMBOT_APITOKEN="your-api-token" \
	--CERTIMATE_NOTIFIER_TELEGRAMBOT_CHATID="your-chat-id"
*/
func TestNotify(t *testing.T) {
	flag.Parse()

	t.Run("Notify", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("APITOKEN: %v", fApiToken),
			fmt.Sprintf("CHATID: %v", fChatId),
		}, "\n"))

		provider, err := provider.NewNotifier(&provider.NotifierConfig{
			BotToken: fApiToken,
			ChatId:   fChatId,
		})
		if err != nil {
			t.Errorf("err: %+v", err)
			return
		}

		res, err := provider.Notify(context.Background(), mockSubject, mockMessage)
		if err != nil {
			t.Errorf("err: %+v", err)
			return
		}

		t.Logf("ok: %v", res)
	})
}
