package gtelegram

type KeyboardLayout interface {
	Render(buttons []InlineKeyboardButton) []InlineKeyboardRow
}

type TwoColumnLayout struct{}

func (t TwoColumnLayout) Render(buttons []InlineKeyboardButton) []InlineKeyboardRow {
	var keyboardButtons []InlineKeyboardRow

	var row InlineKeyboardRow

	for i, button := range buttons {
		row = append(row, button)

		if i%2 != 0 {
			keyboardButtons = append(keyboardButtons, row)
			row = InlineKeyboardRow{}
		}
	}

	if len(buttons)%2 != 0 {
		keyboardButtons = append(keyboardButtons, row)
	}

	return keyboardButtons
}
