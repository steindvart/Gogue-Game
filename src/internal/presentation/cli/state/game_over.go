package state

import (
	"gogue/internal/model/signals"
	"gogue/internal/presentation/dto"
	viewcli "gogue/internal/view/cli"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// GameOver - состояние экрана завершения игры.
// Отображает статистику и по нажатию Enter возвращает в главное меню.
type GameOver struct {
	view   *viewcli.GameOver
	signal signals.Type
}

// NewGameOver создаёт состояние экрана завершения игры.
func NewGameOver(gameOverType viewcli.GameOverType, stats *dto.ScoreEntry) *GameOver {
	g := &GameOver{
		view:   viewcli.NewGameOver(gameOverType, stats),
		signal: signals.NoSignal,
	}
	g.initInput()
	return g
}

func (g *GameOver) Update(float64) signals.Type {
	sig := g.signal
	g.signal = signals.NoSignal
	return sig
}

func (g *GameOver) Primitive() tview.Primitive {
	return g.view.Primitive()
}

func (g *GameOver) initInput() {
	g.view.SetInputCapture(g.handleEvent)
}

func (g *GameOver) handleEvent(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEnter {
		g.signal = signals.ReturnToMenu
		return nil
	}
	return event
}
