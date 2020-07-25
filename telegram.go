package gtelegram

import (
	"fmt"
	"net/http"

	"github.com/thomasheller/ghttp"
)

type Telegram struct {
	token string
}

func NewTelegram(token string) *Telegram {
	return &Telegram{token: token}
}

func (t *Telegram) SendMessage(chatID int, text string, disableNotification bool) (messageID int, err error) {
	td := telegramSendMessageData{
		Text:                text,
		ChatID:              chatID,
		DisableNotification: disableNotification,
	}

	return t.send("sendMessage", td)
}

func (t *Telegram) send(action string, data interface{}) (messageID int, err error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", t.token, action)

	var req *http.Request
	req, err = ghttp.NewRequest(http.MethodPost, url)
	if err != nil {
		return
	}

	var result telegramResult
	if err = ghttp.JSONJSON(req, data, &result); err != nil {
		return
	}

	if !result.OK {
		err = fmt.Errorf("%d %s", result.ErrorCode, result.Description)
		return
	}

	messageID = result.Result.MessageID
	return
}
