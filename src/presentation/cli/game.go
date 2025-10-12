package cli

import (
	"gogue/model/signal"
	"gogue/presentation/cli/state"
	"time"

	"github.com/rivo/tview"
)

const DefaultFPSLimit = 60

type Game struct {
	States   []state.State
	App      *tview.Application
	FPSLimit int
}

// NewGame создаёт новый игровой цикл с заданным FPS (по умолчанию 60)
func NewGame(app *tview.Application, initialState state.State) *Game {
	game := &Game{
		States:   []state.State{initialState},
		App:      app,
		FPSLimit: DefaultFPSLimit,
	}
	game.App.SetRoot(game.CurrentState().Primitive(), true)
	return game
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
	g.runUpdateLoop(g.FPSLimit)

	if err := g.App.Run(); err != nil {
		panic(err)
	}
}

// runUpdateLoop запускает обновление игровых состояний с опросом сигналов от них и с ограничением по FPS
func (g *Game) runUpdateLoop(fpsLimit int) {
	ticker := time.NewTicker(time.Second / time.Duration(fpsLimit))
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			state := g.CurrentState()
			if state == nil {
				g.App.Stop()
				return
			}

			sig := state.Update()

			if sig != signal.NoSignal {
				g.App.QueueUpdateDraw(func() {
					g.handleSignal(sig)
					if s := g.CurrentState(); s != nil {
						g.App.SetRoot(s.Primitive(), true)
					} else {
						g.App.Stop()
					}
				})
			}
		}
	}()
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
