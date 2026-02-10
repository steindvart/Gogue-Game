package cli

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Scoreboard struct {
	flex *tview.Flex
	list *tview.List
}

func NewScoreboard(options []string) *Scoreboard {
	sb := &Scoreboard{
		flex: tview.NewFlex(),
		list: tview.NewList(),
	}

	// @todo - может ли тут быть список не int?
	for _, opt := range options {
		sb.list.AddItem(opt, "", 0, nil)
	}

	// Создаём пустой primitives.Box с тёмным фоном для выравнивания элементов
	gapBox := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)
	// Центрируем список по горизонтали с помощью Flex
	listFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(gapBox, 0, 1, false).
		AddItem(gapBox, 6, 1, false).
		AddItem(sb.list, 14, 0, true).
		AddItem(gapBox, 0, 1, false)
	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(gapBox, 2, 0, false). // верхний отступ
		AddItem(listFlex, 5, 1, true).
		AddItem(gapBox, 3, 0, false).
		AddItem(gapBox, 0, 1, false) // всё оставшееся пространство

	sb.flex = rootFlex

	return sb
}

func (s *Scoreboard) Primitive() tview.Primitive {
	return s.flex
}

func (s *Scoreboard) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) {
	s.flex.SetInputCapture(capture)
}
