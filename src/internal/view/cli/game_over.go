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
func NewGameOver(gameOverType GameOverType, stats *dto.ScoreEntry) *GameOver {
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

func (g *GameOver) buildContent(title string, stats *dto.ScoreEntry) string {
	if stats == nil {
		return title + "\n\n[gray]Press Enter to return to menu"
	}

	return fmt.Sprintf(
		"%s\n\n"+
			"[%s::b]Treasures:[-:-:-]  %d\n"+
			"[%s::b]Level:[-:-:-]      %d\n"+
			"[%s::b]Enemies:[-:-:-]    %d\n"+
			"[%s::b]Food:[-:-:-]       %d\n"+
			"[%s::b]Elixirs:[-:-:-]    %d\n"+
			"[%s::b]Scrolls:[-:-:-]    %d\n"+
			"[%s::b]Hits:[-:-:-]       %d\n"+
			"[%s::b]Missed:[-:-:-]     %d\n"+
			"[%s::b]Moves:[-:-:-]      %d\n\n"+
			"[gray]Press Enter to return to menu",
		title,
		statColorTreasures, stats.Treasures,
		statColorLevel, stats.LevelReached,
		statColorEnemies, stats.EnemiesKilled,
		statColorFood, stats.FoodEaten,
		statColorElixirs, stats.ElixirsDrunk,
		statColorScrolls, stats.ScrollsRead,
		statColorHitsDealt, stats.HitsDealt,
		statColorHitsMissed, stats.HitsMissed,
		statColorCellsMoved, stats.CellsMoved,
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
