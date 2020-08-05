package gtelegram

// Common types

type Message struct {
	MessageID int          `json:"message_id"`
	From      From         `json:"from"`
	Chat      telegramChat `json:"chat"`
	Date      int          `json:"date"`
	Text      string       `json:"text"`
}

type From struct {
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

// Send types

type telegramSendMessageData struct {
	Text                string `json:"text"`
	ChatID              int    `json:"chat_id"`
	DisableNotification bool   `json:"disable_notification"`
	ReplyMarkup telegramReplyMarkup `json:"reply_markup"`
}

type telegramReplyMarkup struct {
	InlineKeyboard []InlineKeyboardRow `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type InlineKeyboardRow []InlineKeyboardButton

type telegramDeleteMessageData struct {
	ChatID    int `json:"chat_id"`
	MessageID int `json:"message_id"`
}

type telegramResult struct {
	OK          bool    `json:"ok"`
	ErrorCode   int     `json:"error_code"`
	Description string  `json:"description"`
	Result      Message `json:"result"`
}

// Receive types

type telegramUpdates struct {
	OK     bool             `json:"ok"`
	Result []telegramUpdate `json:"result"`
}

type telegramUpdate struct {
	UpdateID           int                `json:"update_id"`
	Message            Message            `json:"message"`
	CallbackQuery      CallbackQuery      `json:"callback_query"`
	InlineQuery        InlineQuery        `json:"inline_query"`
	ChosenInlineResult ChosenInlineResult `json:"chosen_inline_result"`
}

type CallbackQuery struct {
	ID   string `json:"id"`
	From From   `json:"from"`
	Data string `json:"data"`
}

type InlineQuery struct {
	ID     string `json:"id"`
	From   From   `json:"from"`
	Query  string `json:"query"`
	Offset string `json:"offset"`
}

type ChosenInlineResult struct {
	ResultID        string `json:"result_id"`
	From            From   `json:"from"`
	InlineMessageID string `json:"inline_message_id"`
	Query           string `json:"query"`
}
