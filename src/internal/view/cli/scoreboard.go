package cli

import (
	"fmt"

	"gogue/internal/presentation/dto"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	scoreboardTableLength = 90
	scoreboardBarLength   = 60
	scoreboardBarSpeed    = 1
)

type Scoreboard struct {
	flex      *tview.Flex
	table     *tview.Table
	runnerBar *RunnerBar
}

func NewScoreboard(entries []dto.ScoreEntry) *Scoreboard {
	sb := &Scoreboard{
		flex:  tview.NewFlex(),
		table: tview.NewTable(),
	}

	// Заголовки столбцов
	headers := []string{
		"#",
		"Treasures",
		"Level",
		"Enemies",
		"Food",
		"Elixirs",
		"Scrolls",
		"Hits",
		"Missed",
		"Moves",
	}

	for col, header := range headers {
		cell := tview.NewTableCell(fmt.Sprintf(" %s ", header)).
			SetTextColor(tcell.ColorGold).
			SetAlign(tview.AlignCenter).
			SetSelectable(false).
			SetAttributes(tcell.AttrBold)
		sb.table.SetCell(0, col, cell)
	}

	// Данные
	for i, entry := range entries {
		var medal string
		switch i {
		case 0:
			medal = "🥇"
		case 1:
			medal = "🥈"
		case 2:
			medal = "🥉"
		default:
			medal = fmt.Sprintf("%d", i+1)
		}

		rowData := []string{
			medal,
			fmt.Sprintf("%d", entry.Treasures),
			fmt.Sprintf("%d", entry.LevelReached),
			fmt.Sprintf("%d", entry.EnemiesKilled),
			fmt.Sprintf("%d", entry.FoodEaten),
			fmt.Sprintf("%d", entry.ElixirsDrunk),
			fmt.Sprintf("%d", entry.ScrollsRead),
			fmt.Sprintf("%d", entry.HitsDealt),
			fmt.Sprintf("%d", entry.HitsMissed),
			fmt.Sprintf("%d", entry.CellsMoved),
		}

		for col, text := range rowData {
			color := tcell.ColorWhite
			if col == 0 && i < 3 {
				color = tcell.ColorGold
			}
			cell := tview.NewTableCell(fmt.Sprintf(" %s ", text)).
				SetTextColor(color).
				SetAlign(tview.AlignCenter).
				SetSelectable(false)
			sb.table.SetCell(i+1, col, cell)
		}
	}

	if len(entries) == 0 {
		cell := tview.NewTableCell("  No records yet  ").
			SetTextColor(tcell.ColorGray).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
		sb.table.SetCell(1, 0, cell)
	}

	sb.table.SetBorder(true).
		SetTitle(" 🌟 SCOREBOARD 🌟 ").
		SetBorderColor(tcell.ColorGold).
		SetTitleColor(tcell.ColorGold).
		SetBackgroundColor(tcell.ColorBlack)

	sb.table.SetSelectable(false, false).
		SetBackgroundColor(tcell.ColorBlack)

	sb.runnerBar = NewRunnerBar(scoreboardBarLength, " 🏆  🏆 ", scoreboardBarSpeed)

	gapBox := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)

	tableFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(gapBox, 0, 1, false).
		AddItem(sb.table, scoreboardTableLength, 0, true).
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
		AddItem(tableFlex, 24, 0, true).
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
