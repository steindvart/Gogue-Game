package cli

import (
	"gogue/model/signal"

	"gogue/view/action"
	viewcli "gogue/view/cli"

	gc "github.com/rthornton128/goncurses"
)

type GameState interface {
	Input(a action.Type) signal.Type
	Update() signal.Type
	Render()
}

type Game struct {
	States []GameState
	Window *gc.Window
}

func (g *Game) PushState(state GameState) {
	g.States = append(g.States, state)
}

func (g *Game) PopState() {
	if len(g.States) == 0 {
		return
	}
	g.States = g.States[:len(g.States)-1]
}

func (g *Game) CurrentState() GameState {
	if len(g.States) == 0 {
		return nil
	}
	return g.States[len(g.States)-1]
}

func (g *Game) Run() {
	actions := g.runInputActionsRoutine()

	for {
		state := g.CurrentState()
		if state == nil {
			break
		}

		g.handleSignal(state.Input(<-actions))
		state = g.CurrentState()
		if state == nil {
			break
		}

		g.handleSignal(state.Update())
		state = g.CurrentState()
		if state == nil {
			break
		}

		state.Render()
	}
}

func (g *Game) runInputActionsRoutine() <-chan action.Type {
	actions := make(chan action.Type, 1)
	go func(ch chan<- action.Type) {
		for {
			ch <- viewcli.HandleInput(g.Window)
		}
	}(actions)

	return actions
}

func (g *Game) handleSignal(s signal.Type) {
	switch s {
	case signal.Stop:
		g.PopState()
	case signal.NewGame:
		// @todo push new game state
	case signal.LoadGame:
		// @todo push load game state
	case signal.ShowScoreboard:
		// @todo push scoreboard state
	}
}
