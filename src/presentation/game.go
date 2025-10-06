package presentation

type GameStateSignal int

const (
	StopStateSignal GameStateSignal = iota
)

type GameState interface {
	Input()
	Render()
	Update() GameStateSignal
}

type Game struct {
	States []GameState
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
		state.Input()
		sign := state.Update()
		state.Render()

		if sign == StopStateSignal {
			g.PopState()
		}
	}
}
