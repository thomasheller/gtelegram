package gtelegram

type Update interface{}

func (t telegramUpdate) Update() Update {
	if t.Message.MessageID > 0 {
		return t.Message
	} else if t.CallbackQuery.ID != "" {
		return t.CallbackQuery
	} else if t.InlineQuery.ID != "" {
		return t.InlineQuery
	} else if t.ChosenInlineResult.ResultID != "" {
		return t.ChosenInlineResult
	}

	return nil
}
