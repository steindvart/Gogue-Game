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

	isPlayerReadyToInteract bool
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

	game.updateGameView()
	game.initInput()

	return &game, nil
}

func (g *Game) updateGameView() {
	g.view.UpdateGameField(g.level.MakeCurrentField(MapWidth, MapHeight))

	player := g.level.Player
	if player != nil {
		g.view.UpdatePlayerInfo(&viewcli.PlayerInfo{
			Health:           player.Attributes.Health,
			MaxHealth:        player.Attributes.MaxHealth,
			Strength:         player.Attributes.Strength,
			Agility:          player.Attributes.Agility,
			TemporaryEffects: g.convertEffectsToView(player.TemporaryEffects),
		})
	}
}

func (g *Game) convertEffectsToView(effects []*primitives.Effect) []viewcli.EffectInfo {
	viewEffects := make([]viewcli.EffectInfo, 0, len(effects))
	for _, effect := range effects {
		viewEffects = append(viewEffects, viewcli.EffectInfo{
			DurationSteps:  int(effect.Duration.Steps),
			HealthModify:   effect.Attributes.Health,
			StrengthModify: effect.Attributes.Strength,
			AgilityModify:  effect.Attributes.Agility,
		})
	}
	return viewEffects
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
			g.level.MovePlayerWithBorderControl(movementRegistry[g.eventToAction(event)])

			switch g.level.CheckEntityCollision(g.level.Player.GetPosition()) {
			case world.CollisionTypeEnemy:
				// @todo - обработка столкновения с врагом (атака на врага)
			case world.CollisionTypeItem, world.CollisionTypeTeleport:
				g.isPlayerReadyToInteract = true
			case world.CollisionTypeNone:
				g.isPlayerReadyToInteract = false
			}
		case action.Select:
			g.handleSelectAction()
		case action.Exit:
			g.signal = signals.Stop
			return nil
		default:
			return event
		}
		return nil
	})
}

func (g *Game) handleSelectAction() {
	if g.level.Player == nil {
		return
	}

	if g.isPlayerReadyToInteract {
		pos := g.level.Player.GetPosition()

		switch g.level.CheckEntityCollision(pos) {
		case world.CollisionTypeItem:
			g.level.PlayerUseItemAtPosition(pos)
		case world.CollisionTypeTeleport:
			if err := g.level.GenerateWithExistingPlayer(g.level.Player); err != nil {
				panic("an error occurred when generating next level: " + err.Error())
			}
		}
	}
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
	case 'e', 'у':
		return action.Select
	}

	return action.NoAction
}

func (g *Game) Update(float64) signals.Type {
	g.updateGameView()

	sig := g.signal
	g.signal = signals.NoSignal
	return sig
}

func (g *Game) Primitive() tview.Primitive {
	return g.view.GetContainer()
}
