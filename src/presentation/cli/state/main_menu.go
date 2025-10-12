package state

import (
	"fmt"
	"gogue/model"
	"gogue/model/signal"
	"gogue/view/action"
	viewcli "gogue/view/cli"

	gc "github.com/rthornton128/goncurses"
)

type MainMenuRenderer interface {
	Render(options []string, active int) error
}

type MainMenu struct {
	model    *model.Menu
	renderer MainMenuRenderer
}

func NewMainMenu(parent *gc.Window) *MainMenu {
	renderer, err := viewcli.NewMainMenu(parent)
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
		renderer: renderer,
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
	err := m.renderer.Render(m.model.GetOptionsLabels(), m.model.GetActive())
	if err != nil {
		fmt.Println("Error rendering main menu:", err)
		panic(err)
	}
}
