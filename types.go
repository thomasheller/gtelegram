package gtelegram

type telegramSendMessageData struct {
	Text                string `json:"text"`
	ChatID              int    `json:"chat_id"`
	DisableNotification bool   `json:"disable_notification"`
}

type telegramResult struct {
	OK          bool            `json:"ok"`
	ErrorCode   int             `json:"error_code"`
	Description string          `json:"description"`
	Result      telegramMessage `json:"result"`
}

type telegramMessage struct {
	MessageID int          `json:"message_id"`
	From      telegramFrom `json:"from"`
	Chat      telegramChat `json:"chat"`
	Date      int          `json:"date"`
	Text      string       `json:"text"`
}

type telegramFrom struct {
	ID           int    `json:"id"`
	Bot          bool   `json:"is_boot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type telegramChat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}
