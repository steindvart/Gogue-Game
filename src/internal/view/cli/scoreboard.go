package cli

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	scoreboardBarLength = 30
	scoreboardBarSpeed  = 1
)

type Scoreboard struct {
	flex      *tview.Flex
	list      *tview.List
	runnerBar *RunnerBar
}

func NewScoreboard(options []string) *Scoreboard {
	sb := &Scoreboard{
		flex: tview.NewFlex(),
		list: tview.NewList(),
	}

	for i, opt := range options {
		var medal string
		switch i {
		case 0:
			medal = "🥇 "
		case 1:
			medal = "🥈 "
		case 2:
			medal = "🥉 "
		default:
			medal = fmt.Sprintf("%2d.", i+1)
		}
		sb.list.AddItem(fmt.Sprintf(" %s %s", medal, opt), "", 0, nil)
	}

	if len(options) == 0 {
		sb.list.AddItem("       No records yet", "", 0, nil)
	}

	sb.list.SetBorder(true).
		SetTitle(" 🌟 RECORDS 🌟 ").
		SetBorderColor(tcell.ColorGold).
		SetTitleColor(tcell.ColorGold).
		SetBackgroundColor(tcell.ColorBlack)

	sb.list.SetMainTextColor(tcell.ColorWhite).
		SetSelectedTextColor(tcell.ColorGold).
		SetSelectedBackgroundColor(tcell.ColorBlack)

	sb.runnerBar = NewRunnerBar(scoreboardBarLength, " 🏆  🏆 ", scoreboardBarSpeed)

	gapBox := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)

	listFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(gapBox, 0, 1, false).
		AddItem(sb.list, 30, 0, true).
		AddItem(gapBox, 0, 1, false)

	runnerFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(gapBox, 0, 1, false).
		AddItem(sb.runnerBar.Primitive(), scoreboardBarLength, 0, false).
		AddItem(gapBox, 0, 1, false)

	hintText := tview.NewTextView().
		SetText("[gray]Press Esc to return").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	hintText.SetBackgroundColor(tcell.ColorBlack)

	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(gapBox, 2, 0, false).
		AddItem(listFlex, 22, 0, true).
		AddItem(gapBox, 1, 0, false).
		AddItem(runnerFlex, 1, 0, false).
		AddItem(gapBox, 1, 0, false).
		AddItem(hintText, 1, 0, false).
		AddItem(gapBox, 0, 1, false)

	sb.flex = rootFlex

	return sb
}

// Update обновляет анимацию runner bar.
func (s *Scoreboard) Update(dt float64) {
	s.runnerBar.UpdatePositions(dt)
	s.runnerBar.Update()
}

func (s *Scoreboard) Primitive() tview.Primitive {
	return s.flex
}

func (s *Scoreboard) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) {
	s.flex.SetInputCapture(capture)
}
