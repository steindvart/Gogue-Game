package cli

import (
	"gogue/model"
	"gogue/view/action"
	viewcli "gogue/view/cli"
	"log"

	gc "github.com/rthornton128/goncurses"
)

type MainMenu struct {
	model *model.Menu
	view  *viewcli.MainMenu
}

func NewMainMenu(parent *gc.Window) *MainMenu {
	menu, err := viewcli.NewMainMenu(parent)
	if err != nil {
		log.Println("Error creating main menu view:", err)
		return nil
	}

	return &MainMenu{
		model: model.NewMenu([]model.MenuOption{
			{Label: "New Game", Signal: model.Signal(model.NewGameSignal)},
			{Label: "Load Game", Signal: model.Signal(model.LoadGameSignal)},
			{Label: "Scoreboard", Signal: model.Signal(model.ShowScoreboardSignal)},
			{Label: "Exit", Signal: model.Signal(model.StopSignal)},
		}),
		view: menu,
	}
}

func (m *MainMenu) Input(a action.Type) model.Signal {
	switch a {
	case action.MoveUp:
		m.model.Previous()
	case action.MoveDown:
		m.model.Next()
	case action.Select:
		return m.model.Select().Signal
	}

	return model.NoSignal
}

func (m *MainMenu) Update() model.Signal {
	return model.NoSignal
}

func (m *MainMenu) Render() {
	err := m.view.Render(m.model.GetOptionsLabels(), m.model.GetActive())
	if err != nil {
		log.Println("Error rendering main menu:", err)
	}
}
