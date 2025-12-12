package state

import (
	"gogue/internal/model/primitives"
	"gogue/internal/model/signals"
	"gogue/internal/model/world"
	"gogue/internal/presentation/action"
	viewcli "gogue/internal/view/cli"
	"math/rand"
	"time"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MapHeight = 30
	MapWidth  = 90
)

type Game struct {
	level  *world.Level
	view   *viewcli.Game
	signal signals.Type
}

func NewGame() (*Game, error) {
	source := rand.New(rand.NewSource(time.Now().UnixNano()))

	game := Game{
		level: world.NewLevelWithDefaults(source, primitives.Size2D[uint]{Height: MapHeight, Width: MapWidth}),
		view:  viewcli.NewGame(),
	}

	err := game.level.Generate()
	if err != nil {
		return nil, err
	}

	game.initRender()
	game.initInput()

	return &game, nil
}

func (g *Game) initRender() {
	g.view.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		// -2: учёт рамки
		fw, fh := width-2, height-2
		field := g.level.MakeCurrentMap(fw, fh)
		g.view.SetFieldToScreen(screen, field, x+1, y+1)
		return x, y, width, height
	})
}

func (g *Game) initInput() {
	movementRegistry := map[action.Type]primitives.Point2D[int]{
		action.MoveUp:               {X: 0, Y: -1},
		action.MoveDown:             {X: 0, Y: 1},
		action.MoveLeft:             {X: -1, Y: 0},
		action.MoveRight:            {X: 1, Y: 0},
		action.MoveLeftUpperCorner:  {X: -1, Y: -1},
		action.MoveRightUpperCorner: {X: 1, Y: -1},
		action.MoveLefLowerCorner:   {X: -1, Y: 1},
		action.MoveRightLowerCorner: {X: 1, Y: 1},
	}

	g.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch g.eventToAction(event) {
		case action.MoveUp,
			action.MoveDown,
			action.MoveLeft,
			action.MoveRight,
			action.MoveLeftUpperCorner,
			action.MoveRightUpperCorner,
			action.MoveLefLowerCorner,
			action.MoveRightLowerCorner:
			g.level.MovePlayerWithCheckCollision(movementRegistry[g.eventToAction(event)])
		case action.Exit:
			g.signal = signals.Stop
			return nil
		default:
			return event
		}
		return nil
	})
}

func (g *Game) eventToAction(event *tcell.EventKey) action.Type {
	switch event.Key() {
	case tcell.KeyUp:
		return action.MoveUp
	case tcell.KeyDown:
		return action.MoveDown
	case tcell.KeyLeft:
		return action.MoveLeft
	case tcell.KeyRight:
		return action.MoveRight
	case tcell.KeyEsc:
		return action.Exit
	}

	ch := unicode.ToLower(event.Rune())
	switch ch {
	case 'w', 'ц':
		return action.MoveUp
	case 's', 'ы':
		return action.MoveDown
	case 'a', 'ф':
		return action.MoveLeft
	case 'd', 'в':
		return action.MoveRight
	case 'y', 'н':
		return action.MoveLeftUpperCorner
	case 'u', 'г':
		return action.MoveRightUpperCorner
	case 'b', 'и':
		return action.MoveLefLowerCorner
	case 'n', 'т':
		return action.MoveRightLowerCorner
	}

	return action.NoAction
}

func (g *Game) Update(float64) signals.Type {
	// Нет lastField, всё строится на лету
	sig := g.signal
	g.signal = signals.NoSignal
	return sig
}

func (g *Game) Primitive() tview.Primitive {
	return g.view
}
