package gtelegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

type Action int

const (
	Send = iota
	Recv
	Dele
	Call
)

type Dialog []Line

func (d Dialog) String() string {
	var s strings.Builder

	for i, line := range d {
		if i > 0 {
			s.WriteString("\n")
		}

		s.WriteString(line.String())
	}

	return s.String()
}

type Line struct {
	Action  Action
	ChatID  int
	Message string
}

func (l Line) String() string {
	var s strings.Builder

	s.WriteString(strconv.Itoa(l.ChatID))

	switch l.Action {
	case Send:
		s.WriteString(">")
	case Recv:
		s.WriteString("<")
	case Dele:
		s.WriteString("X")
	case Call:
		s.WriteString("!")
	}

	s.WriteString(" ")

	s.WriteString(l.Message)

	return s.String()
}

type HistLine struct {
	ID int
	Line
}

func (h HistLine) Raw() string {
	return fmt.Sprintf("[%d] %s", h.ID, h.String())
}

type Mock struct {
	capacity int

	mu sync.Mutex
	i  int

	updates chan Update
	sent    chan HistLine
	history chan HistLine
}

func NewMock(capacity int) *Mock {
	return &Mock{
		capacity: capacity,
		updates:  make(chan Update),
		sent:     make(chan HistLine),
	}
}

func (m *Mock) Play(dialog []Line) {
	m.i = 0
	m.history = make(chan HistLine, m.capacity)

	for _, line := range dialog {
		log.Printf("dialog line %+v", line)

		switch line.Action {
		case Send:
			ID := m.nextID()

			msg := Message{
				MessageID: ID,
				From: From{
					ID: line.ChatID,
				},
				Text: line.Message,
			}

			log.Printf("waiting for hist")
			m.history <- HistLine{
				ID:   ID,
				Line: line,
			}
			log.Printf("hist done")
			log.Printf("waiting for updates")
			m.updates <- msg
			log.Printf("updates done")
		case Recv:
			h := <-m.sent
			log.Printf("waiting for hist")
			m.history <- h
			log.Printf("hist done")
		case Dele:
			// TODO
		case Call:
			ID := m.nextID()
			data := line.Message[1:]

			callback := CallbackQuery{
				ID: strconv.Itoa(ID), // TODO: ID?
				From: From{
					ID: line.ChatID,
				},
				Data: data,
			}

			log.Printf("waiting for hist")
			m.history <- HistLine{
				ID:   ID,
				Line: line,
			}
			log.Printf("hist done")
			log.Printf("waiting for updates")
			m.updates <- callback
			log.Printf("updates done")
		}
	}

	close(m.updates)
}

func (m *Mock) SendMessage(chatID int, text string, disableNotification bool) (messageID int, err error) {
	h := HistLine{
		ID: m.nextID(),
		Line: Line{
			Action:  Recv,
			ChatID:  chatID,
			Message: text,
		},
	}

	m.sent <- h

	return 0, nil // TODO: message ID
}

func (m Mock) DeleteMessage(chatID int, messageID int) (err error) {
	h := HistLine{
		ID: 0,
		Line: Line{
			Action:  Dele,
			ChatID:  chatID,
			Message: strconv.Itoa(messageID),
		},
	}

	m.sent <- h

	return nil
}

func (m *Mock) ShowKeyboard(chatID int, text string, disableNotification bool, buttons []InlineKeyboardButton, layout KeyboardLayout) (messageID int, err error) {
	log.Printf("ShowKeyboard")

	var keyboard strings.Builder

	for i, button := range buttons {
		if i > 0 {

			keyboard.WriteString(",")
		}
		keyboard.WriteString("!")
		keyboard.WriteString(button.CallbackData)
		keyboard.WriteString(":")
		keyboard.WriteString(button.Text)
	}

	// TODO

	h := HistLine{
		ID: m.nextID(),
		Line: Line{
			Action:  Recv,
			ChatID:  chatID,
			Message: keyboard.String(),
		},
	}

	log.Printf("before")
	m.sent <- h
	log.Printf("after")

	return 0, nil
}

func (m Mock) Updates() <-chan Update {
	return m.updates
}

func (m Mock) Result() (string, string) {
	close(m.history)

	var s, raw strings.Builder

	var i int

	for line := range m.history {
		if i > 0 {
			s.WriteString("\n")
			raw.WriteString("\n")
		}

		s.WriteString(line.String())
		raw.WriteString(line.Raw())
		i++
	}

	return s.String(), raw.String()
}

func (m *Mock) nextID() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.i++
	log.Printf("nextID() %d", m.i)
	return m.i
}
