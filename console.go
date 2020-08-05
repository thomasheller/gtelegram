package gtelegram

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const ConsoleChatID = 1

type Console struct {
	updates chan Update
}

func (c *Console) Loop() {
	log.Println("Console loop started")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		if len(input) == 0 {
			continue
		}

		if strings.HasPrefix(input, "!") {
			callback := CallbackQuery{
				// MessageID: ID, // TODO
				From: From{
					ID: 1,
				},
				Data: input[1:],
			}

			fmt.Printf(">%s\n", callback.Data)
			c.updates <- callback
		} else {
			msg := Message{
				// MessageID: ID, // TODO
				From: From{
					ID: 1,
				},
				Text: input,
			}

			fmt.Printf("> %s\n", msg.Text)
			c.updates <- msg
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from stdin: %s", err)
	}

}

func (c *Console) sendMessage(text string) {
	fmt.Printf("< %s\n", text)
}

func (c *Console) showKeyboard(text string, buttons []InlineKeyboardButton) {
	fmt.Printf("< %s\n", text)

	for _, button := range buttons {
		fmt.Printf("<!%s: %s\n", button.CallbackData, button.Text)
	}
}

func (c Console) deleteMessage(messageID int) {
	fmt.Printf("X deleting message ID %d is not supported on console", messageID)
}
