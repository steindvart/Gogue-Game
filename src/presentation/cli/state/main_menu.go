package state

import (
	"gogue/model"
	"gogue/model/signal"
	"gogue/view/action"
	viewcli "gogue/view/cli"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	animationSpeed = 10 // скорость анимации радуги
)

type MainMenu struct {
	model        *model.Menu
	view         *viewcli.MainMenuTView
	signal       signal.Type
	rainbowTimer float64 // для анимации радуги
}

func NewMainMenu() *MainMenu {
	menu := model.NewMenu([]model.MenuOption{
		{Label: "New Game", Signal: signal.NewGame},
		{Label: "Load Game", Signal: signal.LoadGame},
		{Label: "Scoreboard", Signal: signal.ShowScoreboard},
		{Label: "Exit", Signal: signal.Stop},
	})
	view := viewcli.NewMainMenuTView(menu.GetOptionsLabels(), menu.GetActive())

	m := &MainMenu{
		model:  menu,
		view:   view,
		signal: signal.NoSignal,
	}

	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch EventToAction(event) {
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
			m.signal = signal.Stop
			return nil
		default:
			return event
		}
	})

	return m
}

func EventToAction(event *tcell.EventKey) action.Type {
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

func (m *MainMenu) Update(dt float64) signal.Type {
	m.AnimateTitle(dt)

	sig := m.signal
	m.signal = signal.NoSignal // сброс сигнала после чтения
	return sig
}

// Анимация радуги: увеличиваем frame с учётом времени
func (m *MainMenu) AnimateTitle(dt float64) {
	m.rainbowTimer += dt
	frame := int(m.rainbowTimer * animationSpeed) // скорость анимации
	m.view.SetRainbowTitleFrame(frame)
}
