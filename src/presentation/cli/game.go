package cli

import (
	"gogue/model"
)

type GameState interface {
	Input() model.Signal
	Update() model.Signal
	Render()
}

type Game struct {
	States  []GameState
	Signals chan model.Signal
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
	for {
		state := g.CurrentState()
		if state == nil {
			break
		}

		g.HandleSignal(state.Input())
		state = g.CurrentState()
		if state == nil {
			break
		}

		g.HandleSignal(state.Update())
		state = g.CurrentState()
		if state == nil {
			break
		}

		state.Render()
	}
}

func (g *Game) HandleSignal(s model.Signal) {
	switch s {
	case model.StopSignal:
		g.PopState()
	case model.NewGameSignal:
		// @todo push new game state
	case model.LoadGameSignal:
		// @todo push load game state
	case model.ShowScoreboardSignal:
		// @todo push scoreboard state
	}
}
