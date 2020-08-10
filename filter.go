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

func (f Filter) Handle(u Update, handleFunc func(int, Update)) {
	var fromID int

	switch v := u.(type) {
	case Message:
		// TODO: disallow groups...? if v.Chat.Type != "private" (for all update types?)
		fromID = v.From.ID
	case CallbackQuery:
		fromID = v.From.ID
	case InlineQuery:
		fromID = v.From.ID
	case ChosenInlineResult:
		fromID = v.From.ID
	}

	if fromID == 0 {
		// TODO: error handling
		log.Printf("Couldn't parse fromID from Telegram update: %+v", u)
		return
	}

	if len(f.allowedChatIDs) > 0 && !gslice.ContainsInt(f.allowedChatIDs, fromID) {
		return // filtered
	}

	handleFunc(fromID, u)
}
