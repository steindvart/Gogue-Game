package state

import (
	"encoding/json"
	"gogue/internal/model/primitives"
	"gogue/internal/model/signals"
	"gogue/internal/model/world"
	"gogue/internal/presentation/action"
	"gogue/internal/presentation/dto"
	viewcli "gogue/internal/view/cli"
	"math/rand"
	"os"
	"time"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MapHeight      = 30
	MapWidth       = 90
	maxLevelNumber = 21
)

const SaveFileName = "save.json"

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

func LoadGame() (*Game, error) {
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	game := Game{
		level: world.NewLevelWithDefaults(source, primitives.Size2D[uint]{Height: MapHeight, Width: MapWidth}),
		view:  viewcli.NewGame(),
	}

	file, err := os.Open(SaveFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gameSave := dto.GameSaveDto{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&gameSave); err != nil {
		return nil, err
	}
	game.level.Number = gameSave.LevelNumber
	game.level.Rooms = gameSave.Rooms
	game.level.Passages = gameSave.Passages
	game.level.FinishPortal = gameSave.FinishPortal
	game.level.Player = &gameSave.Player

	fogOfWarMap := make(map[primitives.Point2D[int]]bool, len(gameSave.FogOfWar.ExploredTiles))
	for _, point := range gameSave.FogOfWar.ExploredTiles {
		fogOfWarMap[point] = true
	}
	game.level.FogOfWar = &world.FogOfWar{
		ExploredTiles: fogOfWarMap,
		Width:         gameSave.FogOfWar.Width,
		Height:        gameSave.FogOfWar.Height,
	}

	game.updateGameView()
	game.initInput()

	return &game, nil
}

func (g *Game) updateGameView() {
	g.view.UpdateGameField(g.level.MakeCurrentField(MapWidth, MapHeight))

	player := g.level.Player
	if player != nil {
		g.view.UpdatePlayerInfo(dto.ConvertPlayerToDto(g.level.Player, g.level.Number))
	}
}

func (g *Game) initInput() {
	g.view.SetInputCapture(g.handleEvent)
}

func (g *Game) handleEvent(event *tcell.EventKey) *tcell.EventKey {
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
			g.isPlayerReadyToInteract = false
		case world.CollisionTypeItem:
			g.isPlayerReadyToInteract = true
			g.updateItemInfo(g.level.Player.GetPosition())
		case world.CollisionTypeTeleport:
			g.isPlayerReadyToInteract = true
			g.view.SetPlayerIsOnPortal(true)
		case world.CollisionTypeNone:
			g.isPlayerReadyToInteract = false
			g.resetInteraction()
		}

		g.level.ProcessTurns(1)
	case action.Select:
		g.handleSelectAction()
	case action.Exit:
		g.signal = signals.Stop

		g.SaveToFile(SaveFileName)
		return nil
	default:
		return event
	}
	return nil

}

func (g *Game) resetInteraction() {
	g.isPlayerReadyToInteract = false
	g.view.ResetInteraction()
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
			if g.level.Number <= maxLevelNumber {
				if err := g.level.GenerateWithExistingPlayer(g.level.Player); err != nil {
					panic("an error occurred when generating next level: " + err.Error())
				}
			} else {
				// @todo сделать победное окошко. Пока что будет как будто esc
				g.signal = signals.Stop
			}
		}

		g.resetInteraction()
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

func (g *Game) updateItemInfo(pos primitives.Point2D[int]) {
	g.view.UpdateItemInfo(dto.ConvertPositionalItemToDto(g.level.GetItemAtPosition(pos)))
}

func (g *Game) Update(float64) signals.Type {
	g.updateGameView()

	sig := g.signal
	g.signal = signals.NoSignal
	return sig
}

func (g *Game) Primitive() tview.Primitive {
	return g.view.GetRootPrimitive()
}

func (g *Game) SaveToFile(filename string) error {
	fogOfWarList := make([]primitives.Point2D[int], 0, len(g.level.FogOfWar.ExploredTiles))
	for point := range g.level.FogOfWar.ExploredTiles {
		fogOfWarList = append(fogOfWarList, point)
	}

	gameDto := dto.GameSaveDto{
		LevelNumber:  g.level.Number,
		Rooms:        g.level.Rooms,
		Passages:     g.level.Passages,
		Player:       *g.level.Player,
		FinishPortal: g.level.FinishPortal,
		FogOfWar: dto.FogOfWarDto{
			ExploredTiles: fogOfWarList,
			Width:         g.level.FogOfWar.Width,
			Height:        g.level.FogOfWar.Height,
		},
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(gameDto)
}
