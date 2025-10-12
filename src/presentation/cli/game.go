package cli

import (
	"gogue/model/signal"
	"gogue/presentation/cli/state"

	"github.com/rivo/tview"
)

const FPS_DEFAULT = 60

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
	// actions := g.runInputActionsRoutine()
	// ticker := time.NewTicker(time.Second / FPS_DEFAULT)
	// defer ticker.Stop()

	state := g.CurrentState()
	if state != nil {
		g.App.SetRoot(g.CurrentState().Primitive(), true)
	}

	if err := g.App.Run(); err != nil {
		panic(err)
	}

	for {
		// <-ticker.C // ограничение FPS

		state := g.CurrentState()
		if state == nil {
			break
		}

		// g.handleSignal(state.Input(<-actions))
		// state = g.CurrentState()
		// if state == nil {
		// 	break
		// }

		g.handleSignal(state.Update())
		state = g.CurrentState()
		if state == nil {
			break
		}

		// state.Render()

		// err := gc.Update()
		// if err != nil {
		// 	panic(err)
		// }
	}
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
