package state

import (
	"gogue/internal/model/signals"
	viewcli "gogue/internal/view/cli"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Message struct {
	view   *viewcli.Message
	signal signals.Type
}

func NewMessage(msgType viewcli.MessageType) *Message {
	m := &Message{
		view:   viewcli.NewMessage(msgType),
		signal: signals.NoSignal,
	}
	m.initInput()
	return m
}

func (m *Message) Update(float64) signals.Type {
	sig := m.signal
	m.signal = signals.NoSignal // сброс после чтения
	return sig
}

func (m *Message) Primitive() tview.Primitive {
	return m.view.Primitive()
}

func (m *Message) initInput() {
	m.view.SetInputCapture(m.handleEvent)
}

func (m *Message) handleEvent(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyEnter {
		m.signal = signals.ReturnToMenu
		return nil
	}
	return event
}
