package state

import (
	"gogue/model"
	"gogue/model/signal"
	viewcli "gogue/view/cli"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MainMenu struct {
	model    *model.Menu
	view     *viewcli.MainMenuTView
	onSignal func(signal.Type)
}

func NewMainMenu(onSignal func(signal.Type)) *MainMenu {
	menu := model.NewMenu([]model.MenuOption{
		{Label: "New Game", Signal: signal.NewGame},
		{Label: "Load Game", Signal: signal.LoadGame},
		{Label: "Scoreboard", Signal: signal.ShowScoreboard},
		{Label: "Exit", Signal: signal.Stop},
	})
	view := viewcli.NewMainMenuTView(menu.GetOptionsLabels(), menu.GetActive())

	m := &MainMenu{
		model:    menu,
		view:     view,
		onSignal: onSignal,
	}

	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			m.model.Previous()
			m.view.SetActive(m.model.GetActive())
			return nil
		case tcell.KeyDown:
			m.model.Next()
			m.view.SetActive(m.model.GetActive())
			return nil
		case tcell.KeyEnter:
			if m.onSignal != nil {
				m.onSignal(m.model.Select().Signal)
			}
			return nil
		case tcell.KeyEsc:
			if m.onSignal != nil {
				m.onSignal(signal.Stop)
			}
			return nil
		}
		return event
	})

	return m
}

// Primitive возвращает tview-примитив для интеграции с приложением
func (m *MainMenu) Primitive() tview.Primitive {
	return m.view.Primitive()
}

func (m *MainMenu) Update() signal.Type {
	return signal.NoSignal
}
