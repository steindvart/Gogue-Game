package state

import (
	"gogue/internal/model/entities"
	"gogue/internal/model/primitives"
	"gogue/internal/model/signals"
	"gogue/internal/model/world"
	"gogue/internal/presentation/action"
	"gogue/internal/presentation/dto"
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
		g.view.UpdatePlayerInfo(dto.ConvertPlayerToDto(g.level.Player))
		g.updateBackpackInfo()
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
	case action.Take:
		g.handleTakeAction()
	case action.ToggleBackpack:
		g.view.ToggleSecondInfoViewMode()
	case action.OpenWeaponsTab:
		g.view.SetBackpackTab(viewcli.BackpackTabWeapons)
		if g.view.SecondInfoViewMode == viewcli.SecondInfoViewModeEffects {
			g.view.ToggleSecondInfoViewMode()
		}
	case action.OpenFoodTab:
		g.view.SetBackpackTab(viewcli.BackpackTabFood)
		if g.view.SecondInfoViewMode == viewcli.SecondInfoViewModeEffects {
			g.view.ToggleSecondInfoViewMode()
		}
	case action.OpenElixirsTab:
		g.view.SetBackpackTab(viewcli.BackpackTabElixirs)
		if g.view.SecondInfoViewMode == viewcli.SecondInfoViewModeEffects {
			g.view.ToggleSecondInfoViewMode()
		}
	case action.OpenScrollsTab:
		g.view.SetBackpackTab(viewcli.BackpackTabScrolls)
		if g.view.SecondInfoViewMode == viewcli.SecondInfoViewModeEffects {
			g.view.ToggleSecondInfoViewMode()
		}
	case action.Num0, action.Num1, action.Num2, action.Num3,
		action.Num4, action.Num5, action.Num6, action.Num7,
		action.Num8, action.Num9:
		if g.view.SecondInfoViewMode == viewcli.SecondInfoViewModeBackpack {
			g.handleItemSelection(g.eventToAction(event))
		}
	case action.Exit:
		g.signal = signals.Stop
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
			if err := g.level.GenerateWithExistingPlayer(g.level.Player); err != nil {
				panic("an error occurred when generating next level: " + err.Error())
			}
		}

		g.resetInteraction()
	}
}

func (g *Game) handleTakeAction() {
	if g.level.Player == nil {
		return
	}

	if g.isPlayerReadyToInteract {
		pos := g.level.Player.GetPosition()

		switch g.level.CheckEntityCollision(pos) {
		case world.CollisionTypeItem:
			if err := g.level.PlayerTakeItemAtPosition(pos); err != nil {
				// @todo - вывод сообщения об ошибке в view
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
	case 'r', 'к':
		return action.Take
	case 'b', 'и':
		return action.ToggleBackpack
	case 'e', 'у':
		return action.Select
	case 'z', 'я':
		return action.OpenWeaponsTab
	case 'x', 'ч':
		return action.OpenFoodTab
	case 'c', 'с':
		return action.OpenElixirsTab
	case 'v', 'м':
		return action.OpenScrollsTab
	case '0':
		return action.Num0
	case '1':
		return action.Num1
	case '2':
		return action.Num2
	case '3':
		return action.Num3
	case '4':
		return action.Num4
	case '5':
		return action.Num5
	case '6':
		return action.Num6
	case '7':
		return action.Num7
	case '8':
		return action.Num8
	case '9':
		return action.Num9
	}

	return action.NoAction
}

func (g *Game) updateItemInfo(pos primitives.Point2D[int]) {
	g.view.UpdateItemInfo(dto.ConvertPositionalItemToDto(g.level.GetItemAtPosition(pos)))
}

func (g *Game) updateBackpackInfo() {
	if g.level.Player != nil && g.level.Player.Backpack != nil {
		g.view.UpdateBackpackInfo(dto.ConvertBackpackToDto(g.level.Player.Backpack))
	}
}

// handleItemSelection обрабатывает выбор предмета из рюкзака
func (g *Game) handleItemSelection(actionType action.Type) {
	if g.level.Player == nil || g.level.Player.Backpack == nil {
		return
	}

	itemIndex := int(actionType - action.Num0)

	var itemType entities.BackpackItemType
	var actualIndex int

	switch g.view.GetCurrentBackpackTab() {
	case viewcli.BackpackTabWeapons:
		itemType = entities.BackpackItemTypeWeapon
		// Для оружия индекс 0 - снятие текущего оружия
		if itemIndex == 0 {
			if err := g.level.Player.DropEquipWeaponToBackpack(); err != nil {
				// @todo - вывод сообщения об ошибке в view
			}
			g.updateBackpackInfo()
			return
		}
		actualIndex = itemIndex - 1

	case viewcli.BackpackTabFood:
		itemType = entities.BackpackItemTypeFood
		actualIndex = itemIndex - 1

	case viewcli.BackpackTabElixirs:
		itemType = entities.BackpackItemTypeElixir
		actualIndex = itemIndex - 1

	case viewcli.BackpackTabScrolls:
		itemType = entities.BackpackItemTypeScroll
		actualIndex = itemIndex - 1
	}

	selectedItem := g.level.Player.GetItemFromBackpackByIndex(itemType, actualIndex)
	if selectedItem == nil {
		return
	}

	if err := g.level.Player.UseItemFromBackpack(selectedItem); err != nil {
		// @todo - вывод сообщения об ошибке в view
	}

	g.updateBackpackInfo()
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
