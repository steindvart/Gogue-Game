package cli

import (
	"gogue/internal/model/signals"
	"gogue/internal/presentation/cli/state"
	"gogue/internal/presentation/dto"
	viewcli "gogue/internal/view/cli"
	"time"

	"github.com/rivo/tview"
)

const DefaultFPSLimit = 60

type StateMachine struct {
	States   []state.State
	App      *tview.Application
	FPSLimit int
}

func NewStateMachine(app *tview.Application, initialState state.State) *StateMachine {
	game := &StateMachine{
		States:   []state.State{initialState},
		App:      app,
		FPSLimit: DefaultFPSLimit,
	}
	game.App.SetRoot(game.CurrentState().Primitive(), true)
	return game
}

func (g *StateMachine) PushState(state state.State) {
	g.States = append(g.States, state)
}

func (g *StateMachine) PopState() {
	if len(g.States) == 0 {
		return
	}

	g.States = g.States[:len(g.States)-1]
}

func (g *StateMachine) rootState() {
	if len(g.States) == 0 {
		return
	}

	g.States = g.States[:1]
}

func (g *StateMachine) CurrentState() state.State {
	if len(g.States) == 0 {
		return nil
	}
	return g.States[len(g.States)-1]
}

func (g *StateMachine) Run() {
	g.runUpdateLoop(g.FPSLimit)

	if err := g.App.Run(); err != nil {
		panic(err)
	}
}

// runUpdateLoop запускает обновление игровых состояний с опросом сигналов от них и с ограничением по FPS
func (g *StateMachine) runUpdateLoop(fpsLimit int) {
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

func (g *StateMachine) handleSignal(s signals.Type) {
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
		game, err := state.LoadGame()
		if err != nil {
			// @todo - архитектурная проблема: машина состояний не должна зависеть от view
			msgState := state.NewMessage(viewcli.MsgNoSavedGame)
			g.PushState(msgState)
			return
		}
		g.PushState(game)
	case signals.ShowScoreboard:
		leaderboard, err := state.GetScoreboard()
		if err != nil {
			panic(err)
		}
		g.PushState(leaderboard)
	case signals.GameWon:
		stats := g.extractGameOverStats()
		gameOverState := state.NewGameOver(viewcli.GameOverTypeWin, stats)
		g.PushState(gameOverState)
	case signals.GameOver:
		stats := g.extractGameOverStats()
		gameOverState := state.NewGameOver(viewcli.GameOverTypeDeath, stats)
		g.PushState(gameOverState)
	case signals.ReturnToMenu:
		g.rootState()
	}
}

// extractGameOverStats извлекает статистику из текущего состояния Game (если оно активно).
func (g *StateMachine) extractGameOverStats() *dto.GameOverStats {
	if gameState, ok := g.CurrentState().(*state.Game); ok {
		return gameState.GameOverStats
	}
	return nil
}
