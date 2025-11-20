package cli

import (
	"gogue/internal/model/signals"
	"gogue/internal/presentation/cli/state"
	"time"

	"github.com/rivo/tview"
)

const DefaultFPSLimit = 60

// @todo - переименовать в StateMachine
type Game struct {
	States   []state.State
	App      *tview.Application
	FPSLimit int
}

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
	var lastTime = time.Now()
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			dt := now.Sub(lastTime).Seconds()
			lastTime = now

			state := g.CurrentState()
			if state == nil {
				g.App.Stop()
				return
			}

			sig := state.Update(dt)

			if sig != signals.NoSignal {
				g.App.QueueUpdateDraw(func() {
					g.handleSignal(sig)
					if s := g.CurrentState(); s != nil {
						g.App.SetRoot(s.Primitive(), true)
					} else {
						g.App.Stop()
					}
				})
			} else {
				// Для анимации: обновляем UI даже если сигнала нет
				g.App.QueueUpdateDraw(func() {})
			}
		}
	}()
}

func (g *Game) handleSignal(s signals.Type) {
	switch s {
	case signals.Stop:
		g.PopState()
	case signals.NewGame:
		game, err := state.NewGame()
		if err != nil {
			panic(err)
		}
		g.PushState(game)
	case signals.LoadGame:
		// @todo push load game state
	case signals.ShowScoreboard:
		// @todo push scoreboard state
	}
}
