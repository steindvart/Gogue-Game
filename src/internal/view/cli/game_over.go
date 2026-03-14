package cli

import (
	"fmt"
	"gogue/internal/presentation/dto"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// GameOverType определяет тип окончания игры (победа или поражение).
type GameOverType int

const (
	GameOverTypeDeath GameOverType = iota
	GameOverTypeWin
)

// GameOver - экран завершения игры с отображением статистики сессии.
type GameOver struct {
	flex *tview.Flex
	text *tview.TextView
}

// NewGameOver создаёт экран завершения игры с заголовком и статистикой.
func NewGameOver(gameOverType GameOverType, stats *dto.GameOverStats) *GameOver {
	g := &GameOver{
		flex: tview.NewFlex(),
		text: tview.NewTextView(),
	}

	title, borderColor := g.buildTitle(gameOverType)
	content := g.buildContent(title, stats)

	g.text.SetText(content).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorWhite).
		SetDynamicColors(true).
		SetWordWrap(true)

	g.text.SetBorder(true).
		SetBorderColor(borderColor).
		SetTitleColor(borderColor)

	gapBox := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)

	textFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(gapBox, 0, 1, false).
		AddItem(g.text, 50, 0, true).
		AddItem(gapBox, 0, 1, false)

	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(gapBox, 0, 1, false).
		AddItem(textFlex, 20, 0, true).
		AddItem(gapBox, 0, 1, false)

	g.flex = rootFlex
	return g
}

func (g *GameOver) buildTitle(gameOverType GameOverType) (string, tcell.Color) {
	switch gameOverType {
	case GameOverTypeWin:
		return "\n\n🎉 Game over... You won! 🎉", tcell.ColorGreen
	case GameOverTypeDeath:
		return "\n\n💀 Game over... You die... 💀", tcell.ColorRed
	default:
		return "\n\nGame Over", tcell.ColorWhite
	}
}

func (g *GameOver) buildContent(title string, stats *dto.GameOverStats) string {
	if stats == nil {
		return title + "\n\n[gray]Press Enter to return to menu"
	}

	return fmt.Sprintf(
		"%s\n\n"+
			"[gold]Treasures collected:[white]   %d\n"+
			"[red]Enemies killed:[white]         %d\n"+
			"[skyblue]Level reached:[white]      %d\n"+
			"[green]Food eaten:[white]          %d\n"+
			"[magenta]Elixirs drunk:[white]      %d\n"+
			"[cyan]Scrolls read:[white]          %d\n"+
			"[yellow]Hits dealt:[white]          %d\n"+
			"[orange]Hits missed:[white]         %d\n"+
			"[lime]Cells moved:[white]           %d\n\n"+
			"[gray]Press Enter to return to menu",
		title,
		stats.Treasures,
		stats.EnemiesKilled,
		stats.LevelReached,
		stats.FoodEaten,
		stats.ElixirsDrunk,
		stats.ScrollsRead,
		stats.HitsDealt,
		stats.HitsMissed,
		stats.CellsMoved,
	)
}

// Primitive возвращает корневой примитив для отображения.
func (g *GameOver) Primitive() tview.Primitive {
	return g.flex
}

// SetInputCapture устанавливает обработчик ввода.
func (g *GameOver) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) {
	g.flex.SetInputCapture(capture)
}
