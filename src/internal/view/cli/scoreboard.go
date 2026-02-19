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

	for _, opt := range options {
		sb.list.AddItem(opt, "", 0, nil)
	}

	// Внешний вид списка
	sb.list.SetBorder(true).
		SetTitle("🌟RECORDS🌟").
		SetBorderColor(tcell.ColorWhite).
		SetTitleColor(tcell.ColorGold)

	// Создаём пустой primitives.Box с тёмным фоном для выравнивания элементов
	gapBox := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)

	// Центрируем список по горизонтали с помощью Flex
	listFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(gapBox, 0, 1, false).  // отступ от левой стенки терминала
		AddItem(sb.list, 22, 0, true). // ширина окошка
		AddItem(gapBox, 0, 1, false)   // отступ от правой стенки терминала

	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(gapBox, 2, 0, false).   // верхний отступ
		AddItem(listFlex, 22, 0, true). // высота окошка 22
		AddItem(gapBox, 0, 1, false)    // закраска пространства под окошком

	sb.flex = rootFlex

	return sb
}

func (s *Scoreboard) Primitive() tview.Primitive {
	return s.flex
}

func (s *Scoreboard) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) {
	s.flex.SetInputCapture(capture)
}
