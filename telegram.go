package gtelegram

type Telegram interface {
	// TODO: SendMessagef (for conveniece?)

	// TODO: ShowOrUpdateKeyboard
	// TODO: ShowKeyboard

	// TODO: "Say"Inline

	SendMessage(chatID int, text string, disableNotification bool) (messageID int, err error)

	DeleteMessage(chatID int, messageID int) (err error)

	ShowKeyboard(chatID int, text string, disableNotification bool, buttons []InlineKeyboardButton, layout KeyboardLayout) (messageID int, err error)

	Updates() <-chan Update
}
