package cli

import (
	"gogue/model"
	view "gogue/view"
	viewcli "gogue/view/cli"

	gc "github.com/rthornton128/goncurses"
)

type MainMenu struct {
	model model.Menu
	view  viewcli.MainMenu
}

func NewMainMenu(menu *gc.Window) *MainMenu {
	return &MainMenu{
		model: *model.NewMenu([]model.MenuOption{
			{Label: "New Game", Signal: model.Signal(model.NewGameSignal)},
			{Label: "Load Game", Signal: model.Signal(model.LoadGameSignal)},
			{Label: "Scoreboard", Signal: model.Signal(model.ShowScoreboardSignal)},
			{Label: "Exit", Signal: model.Signal(model.StopSignal)},
		}),
		view: viewcli.MainMenu{W: menu},
	}
}

func (m *MainMenu) Input(a view.ActionType) model.Signal {
	switch a {
	case view.MoveUp:
		m.model.Previous()
	case view.MoveDown:
		m.model.Next()
	case view.Select:
		return m.model.Select().Signal
	}

	return model.NoSignal
}

func (m *MainMenu) Update() model.Signal {
	return model.NoSignal
}

func (m *MainMenu) Render() {
	m.view.Render(m.model.GetOptionsLabels(), m.model.GetActive())
}
