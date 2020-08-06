package gtelegram

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/thomasheller/ghttp"
)

type Service struct {
	token      string
	lastUpdate int
	updates    chan Update
	console    *Console
}

func NewService(token string) *Service {
	return &Service{
		token:   token,
		updates: make(chan Update),
	}
}

func (s *Service) SendMessage(chatID int, text string, disableNotification bool) (messageID int, err error) {
	if s.console != nil && chatID == ConsoleChatID {
		s.console.sendMessage(text)
		return 0, nil // TODO: chat ID
	}

	td := telegramSendMessageData{
		Text:                text,
		ChatID:              chatID,
		DisableNotification: disableNotification,
		ReplyMarkup: telegramReplyMarkup{
			InlineKeyboard: []InlineKeyboardRow{},
		},
	}

	return s.send("sendMessage", td)
}

func (s *Service) ShowKeyboard(chatID int, text string, disableNotification bool, buttons []InlineKeyboardButton, layout KeyboardLayout) (messageID int, err error) {
	if s.console != nil && chatID == ConsoleChatID {
		s.console.showKeyboard(text, buttons)
		return 0, nil // TODO: chat ID
	}

	keyboardMarkup := telegramReplyMarkup{
		InlineKeyboard: layout.Render(buttons),
	}

	td := telegramSendMessageData{
		Text:                text,
		ChatID:              chatID,
		DisableNotification: disableNotification,
		ReplyMarkup:         keyboardMarkup,
	}

	return s.send("sendMessage", td)
}

func (s *Service) DeleteMessage(chatID int, messageID int) (err error) {
	if s.console != nil && chatID == ConsoleChatID {
		// TODO
		return nil
	}

	td := telegramDeleteMessageData{
		MessageID: messageID,
		ChatID:    chatID,
	}

	_, err = s.send("deleteMessage", td) // TODO: can't return message ID?
	return
}

func (s *Service) send(action string, data interface{}) (messageID int, err error) {
	url := s.buildURL("%s", action)

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

func (s *Service) Updates() <-chan Update {
	go s.pollLoop()
	return s.updates
}

func (s *Service) Console() *Console {
	if s.console == nil {
		s.console = &Console{updates: s.updates}
	}

	return s.console
}

func (s *Service) pollLoop() {
	ticker := time.NewTicker(1000 * time.Millisecond)
	for range ticker.C { // TODO: cancel loop
		s.poll()
	}
}

func (s *Service) poll() {
	url := s.buildURL("getUpdates?offset=%d", s.lastUpdate+1)

	req, err := ghttp.NewRequest(http.MethodGet, url)
	if err != nil {
		log.Printf("Failed to build HTTP request: %s", err)
		return
	}

	updates := telegramUpdates{}

	if err := ghttp.JSON(req, &updates); err != nil {
		log.Printf("Failed poll messages from Telegram API: %v", err)
		return
	}

	if !updates.OK {
		log.Print("Failed to query updates from Telegram: ok is not true")
		// TODO: error handling
		return
	}

	for _, update := range updates.Result {
		// TODO: debug logging?
		log.Printf("%+v", update)

		if update.UpdateID > s.lastUpdate {
			s.lastUpdate = update.UpdateID
		}

		s.updates <- update.Update()
	}
}

func (s *Service) buildURL(format string, a ...interface{}) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/", s.token) + fmt.Sprintf(format, a...)
}
