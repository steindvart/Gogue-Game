package cli

import (
	"gogue/model"
	view "gogue/view"
	viewcli "gogue/view/cli"

	gc "github.com/rthornton128/goncurses"
)

type MainMenuState struct {
	model model.Menu
	view  viewcli.MainMenu
}

func NewMainMenuState(menu *gc.Window) *MainMenuState {
	return &MainMenuState{
		model: *model.NewMenu([]model.MenuOption{
			{Label: "New Game", Signal: model.Signal(model.NewGameSignal)},
			{Label: "Load Game", Signal: model.Signal(model.LoadGameSignal)},
			{Label: "Scoreboard", Signal: model.Signal(model.ShowScoreboardSignal)},
			{Label: "Exit", Signal: model.Signal(model.StopSignal)},
		}),
		view: viewcli.MainMenu{W: menu},
	}
}

func (m *MainMenuState) Input() model.Signal {
	action := viewcli.HandleInput(m.view.W)
	switch action {
	case view.MoveUp:
		m.model.Previous()
	case view.MoveDown:
		m.model.Next()
	case view.Select:
		return m.model.Select().Signal
	}

	return model.NoSignal
}

func (m *MainMenuState) Update() model.Signal {
	if gc.Key(m.view.W.GetChar()) == 'q' {
		return model.StopSignal
	}

	return model.NoSignal
}

func (m *MainMenuState) Render() {
	m.view.Render(m.model.GetOptionsLabels(), m.model.GetActive())
}
