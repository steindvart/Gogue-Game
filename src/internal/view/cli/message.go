package cli

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MessageType int

const (
	MsgNoSavedGame MessageType = iota
)

func (mt MessageType) String() string {
	switch mt {
	case MsgNoSavedGame:
		return "⚠  No Saved Game ⚠\n\nNo save file found"
	default:
		return "Message"
	}
}

type Message struct {
	flex *tview.Flex
	text *tview.TextView
}

func NewMessage(msgType MessageType) *Message {
	m := &Message{
		flex: tview.NewFlex(),
		text: tview.NewTextView(),
	}

	m.text.SetText(msgType.String()).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorWhite).
		SetDynamicColors(true).
		SetWordWrap(true)

	borderColor := tcell.ColorRed
	switch msgType {
	case MsgNoSavedGame:
		borderColor = tcell.ColorYellow
	}

	m.text.SetBorder(true).
		SetBorderColor(borderColor).
		SetTitleColor(borderColor)

	gapBox := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)

	textFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(gapBox, 0, 1, false).
		AddItem(m.text, 40, 0, true).
		AddItem(gapBox, 0, 1, false)

	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(gapBox, 0, 1, false).
		AddItem(textFlex, 5, 0, true).
		AddItem(gapBox, 0, 1, false)

	m.flex = rootFlex
	return m
}

func (m *Message) Primitive() tview.Primitive {
	return m.flex
}

func (m *Message) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) {
	m.flex.SetInputCapture(capture)
}
