package state

import (
	"fmt"
	"gogue/model"
	"gogue/model/signal"
	"gogue/view/action"
	viewcli "gogue/view/cli"

	gc "github.com/rthornton128/goncurses"
)

type MainMenu struct {
	model *model.Menu
	view  *viewcli.MainMenu
}

func NewMainMenu(parent *gc.Window) *MainMenu {
	menu, err := viewcli.NewMainMenu(parent)
	if err != nil {
		fmt.Println("Error creating main menu view:", err)
		return nil
	}

	return &MainMenu{
		model: model.NewMenu([]model.MenuOption{
			{Label: "New Game", Signal: signal.NewGame},
			{Label: "Load Game", Signal: signal.LoadGame},
			{Label: "Scoreboard", Signal: signal.ShowScoreboard},
			{Label: "Exit", Signal: signal.Stop},
		}),
		view: menu,
	}
}

func (m *MainMenu) Input(a action.Type) signal.Type {
	switch a {
	case action.MoveUp:
		m.model.Previous()
	case action.MoveDown:
		m.model.Next()
	case action.Select:
		return m.model.Select().Signal
	}

	return signal.NoSignal
}

func (m *MainMenu) Update() signal.Type {
	return signal.NoSignal
}

func (m *MainMenu) Render() {
	err := m.view.Render(m.model.GetOptionsLabels(), m.model.GetActive())
	if err != nil {
		fmt.Println("Error rendering main menu:", err)
	}
}
