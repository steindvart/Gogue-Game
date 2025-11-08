package state

import (
	"gogue/internal/model"
	"gogue/internal/model/signals"
	"gogue/internal/presentation/action"
	viewcli "gogue/internal/view/cli"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MainMenu struct {
	model  *model.Menu
	view   *viewcli.MainMenu
	signal signals.Type
}

func NewMainMenu() *MainMenu {
	menu := model.NewMenu([]model.MenuOption{
		{Label: "  New Game  ", Signal: signals.NewGame},
		{Label: "    Load    ", Signal: signals.LoadGame},
		{Label: " Scoreboard ", Signal: signals.ShowScoreboard},
		{Label: "    Exit    ", Signal: signals.Stop},
	})
	view := viewcli.NewMainMenu(menu.GetOptionsLabels(), menu.GetActive())

	m := &MainMenu{
		model:  menu,
		view:   view,
		signal: signals.NoSignal,
	}

	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch m.eventToAction(event) {
		case action.MoveUp:
			m.model.Previous()
			m.view.SetActive(m.model.GetActive())
			return nil
		case action.MoveDown:
			m.model.Next()
			m.view.SetActive(m.model.GetActive())
			return nil
		case action.Select:
			m.signal = m.model.Select().Signal
			return nil
		case action.Exit:
			m.signal = signals.Stop
			return nil
		default:
			return event
		}
	})

	return m
}

func (m *MainMenu) eventToAction(event *tcell.EventKey) action.Type {
	switch event.Key() {
	case tcell.KeyUp:
		return action.MoveUp
	case tcell.KeyDown:
		return action.MoveDown
	case tcell.KeyEnter:
		return action.Select
	case tcell.KeyEsc:
		return action.Exit
	}

	switch event.Rune() {
	case 'w', 'W', 'ц', 'Ц':
		return action.MoveUp
	case 's', 'S', 'ы', 'Ы':
		return action.MoveDown
	}

	return action.NoAction
}

func (m *MainMenu) Primitive() tview.Primitive {
	return m.view.Primitive()
}

func (m *MainMenu) Update(dt float64) signals.Type {
	const animationSpeed = 10.0

	m.view.Update(dt, animationSpeed)

	sig := m.signal
	m.signal = signals.NoSignal // сброс сигнала после чтения
	return sig
}
