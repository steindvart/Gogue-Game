package cli

import (
	"gogue/model/signal"
	"gogue/presentation/cli/state"
	"time"

	"gogue/view/action"
	viewcli "gogue/view/cli"

	gc "github.com/rthornton128/goncurses"
)

const FPS_DEFAULT = 60

type Game struct {
	States []state.State
	Window *gc.Window
}

func (g *Game) PushState(state state.State) {
	g.States = append(g.States, state)
}

func (g *Game) PopState() {
	if len(g.States) == 0 {
		return
	}
	g.States = g.States[:len(g.States)-1]
}

func (g *Game) CurrentState() state.State {
	if len(g.States) == 0 {
		return nil
	}
	return g.States[len(g.States)-1]
}

func (g *Game) Run() {
	actions := g.runInputActionsRoutine()
	ticker := time.NewTicker(time.Second / FPS_DEFAULT)
	defer ticker.Stop()

	for {
		<-ticker.C // ограничение FPS

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
