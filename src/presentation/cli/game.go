package cli

import (
	"gogue/model/signal"
	"gogue/presentation/cli/state"

	"github.com/rivo/tview"
)

type Game struct {
	States []state.State
	App    *tview.Application
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
	state := g.CurrentState()
	if state != nil {
		g.App.SetRoot(g.CurrentState().Primitive(), true)
	}

	if err := g.App.Run(); err != nil {
		panic(err)
	}

	for {
		state := g.CurrentState()
		if state == nil {
			break
		}

		g.handleSignal(state.Update())
		state = g.CurrentState()
		if state == nil {
			break
		}
	}

	g.App.Stop()
}

func (g *Game) handleSignal(s signal.Type) {
	switch s {
	case signal.Stop:
		g.PopState()
	case signal.NewGame:
		// g.PushState(state.NewGame(g.App))
	case signal.LoadGame:
		// @todo push load game state
	case signal.ShowScoreboard:
		// @todo push scoreboard state
	}
}
