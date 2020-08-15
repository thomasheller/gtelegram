package gtelegram

import (
	"log"

	"github.com/thomasheller/gslice"
)

type Filter struct {
	allowedChatIDs []int
}

func NewFilter(allowedChatIDs []int) *Filter {
	return &Filter{allowedChatIDs: allowedChatIDs}
}

func (f Filter) Handle(u Update, handleFunc func(From, Update)) {
	var from From

	switch v := u.(type) {
	case Message:
		// TODO: disallow groups...? if v.Chat.Type != "private" (for all update types?)
		from = v.From
	case CallbackQuery:
		from = v.From
	case InlineQuery:
		from = v.From
	case ChosenInlineResult:
		from = v.From
	}

	if from.ID == 0 {
		// TODO: error handling
		log.Printf("Couldn't parse from.ID from Telegram update: %+v", u)
		return
	}

	if len(f.allowedChatIDs) > 0 && !gslice.ContainsInt(f.allowedChatIDs, from.ID) {
		return // filtered
	}

	handleFunc(from, u)
}
