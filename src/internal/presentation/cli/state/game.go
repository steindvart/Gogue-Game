package state

import (
	"gogue/internal/common"
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/model/signals"
	"gogue/internal/model/world"
	"gogue/internal/presentation/action"
	"gogue/internal/presentation/dto"
	"gogue/internal/presentation/save"
	viewcli "gogue/internal/view/cli"
	"math/rand"
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

const (
	SaveFileName  = "save.json"
	ScoreFileName = "score.json"
)

type Game struct {
	level  *world.Level
	view   *viewcli.Game
	signal signals.Type

	isPlayerReadyToInteract bool
	backpackDropMode        bool

	GameOverStats *dto.GameOverStats
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
	level, err := save.LoadGame(SaveFileName)
	if err != nil {
		return nil, err
	}

	game := &Game{
		level: level,
		view:  viewcli.NewGame(),
	}

	game.updateGameView()
	game.initInput()
	return game, nil
}

func (g *Game) updateGameView() {
	g.view.UpdateGameField(g.level.MakeCurrentField(MapWidth, MapHeight))

	player := g.level.Player
	if player != nil {
		g.view.UpdatePlayerInfo(dto.ConvertPlayerToDto(g.level.Player, g.level.Number))
		g.updateBackpackInfo()
	}
}

func (g *Game) initInput() {
	g.view.SetInputCapture(g.handleEvent)
}

func (g *Game) handleEvent(event *tcell.EventKey) *tcell.EventKey {
	switch g.eventToAction(event) {
	case action.MoveUp,
		action.MoveDown,
		action.MoveLeft,
		action.MoveRight,
		action.MoveLeftUpperCorner,
		action.MoveRightUpperCorner,
		action.MoveLefLowerCorner,
		action.MoveRightLowerCorner:
		g.view.ClearAttackInfos()

		// Если игрок оглушён - ход пропускается, но враги всё равно действуют.
		if g.level.Player.ProcessStun() {
			g.view.SetStunnedMessage(true)
			enemyAttacks := g.level.ProcessTurns(1)
			if len(enemyAttacks) > 0 {
				g.view.SetEnemyAttackInfos(dto.ConvertAttackResultsToDto(enemyAttacks))
			}
			g.checkPlayerDeath()
			break
		}
		g.view.SetStunnedMessage(false)

		g.handleMoveAction(g.eventToAction(event))
		switch g.level.CheckEntityCollision(g.level.Player.GetPosition()) {
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

		enemyAttacks := g.level.ProcessTurns(1)
		if len(enemyAttacks) > 0 {
			g.view.SetEnemyAttackInfos(dto.ConvertAttackResultsToDto(enemyAttacks))
		}
		g.checkPlayerDeath()
	case action.Select:
		g.handleSelectAction()
	case action.Take:
		g.handleTakeAction()
	case action.ToggleBackpack:
		g.view.ToggleSecondInfoViewMode()
		// При переключении в режим рюкзака сбрасываем drop mode (по умолчанию use)
		if g.view.SecondInfoViewMode == viewcli.SecondInfoViewModeBackpack {
			g.backpackDropMode = false
			g.view.SetBackpackDropMode(false)
		}
	case action.ToggleBackpackMode:
		if g.view.SecondInfoViewMode == viewcli.SecondInfoViewModeBackpack {
			g.backpackDropMode = !g.backpackDropMode
			g.view.SetBackpackDropMode(g.backpackDropMode)
		}
	case action.OpenWeaponsTab:
		g.view.SetBackpackTab(viewcli.BackpackTabWeapons)
		if g.view.SecondInfoViewMode == viewcli.SecondInfoViewModeEffects {
			g.view.ToggleSecondInfoViewMode()
			// @todo - подумать как лучше сделать
			g.backpackDropMode = false
			g.view.SetBackpackDropMode(false)
		}
	case action.OpenFoodTab:
		g.view.SetBackpackTab(viewcli.BackpackTabFood)
		if g.view.SecondInfoViewMode == viewcli.SecondInfoViewModeEffects {
			g.view.ToggleSecondInfoViewMode()
			g.backpackDropMode = false
			g.view.SetBackpackDropMode(false)
		}
	case action.OpenElixirsTab:
		g.view.SetBackpackTab(viewcli.BackpackTabElixirs)
		if g.view.SecondInfoViewMode == viewcli.SecondInfoViewModeEffects {
			g.view.ToggleSecondInfoViewMode()
			g.backpackDropMode = false
			g.view.SetBackpackDropMode(false)
		}
	case action.OpenScrollsTab:
		g.view.SetBackpackTab(viewcli.BackpackTabScrolls)
		if g.view.SecondInfoViewMode == viewcli.SecondInfoViewModeEffects {
			g.view.ToggleSecondInfoViewMode()
			g.backpackDropMode = false
			g.view.SetBackpackDropMode(false)
		}
	case action.Num0, action.Num1, action.Num2, action.Num3,
		action.Num4, action.Num5, action.Num6, action.Num7,
		action.Num8, action.Num9:
		if g.view.SecondInfoViewMode == viewcli.SecondInfoViewModeBackpack {
			if g.backpackDropMode {
				g.handleItemDrop(g.eventToAction(event))
			} else {
				g.handleItemSelection(g.eventToAction(event))
			}
		}
	case action.Exit:
		err := save.SaveGame(g.level, primitives.Size2D[uint]{Height: MapHeight, Width: MapWidth}, SaveFileName)
		if err != nil {
			panic("an error occurred while saving the game:" + err.Error())
		}

		g.signal = signals.Stop

		return nil
	default:
		return event
	}
	return nil

}

func (g *Game) handleMoveAction(a action.Type) {
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

	oldPos := g.level.Player.GetPosition()
	g.level.Player.Move(movementRegistry[a])

	if g.level.IsCollisionWithBorders(g.level.Player.GetPosition()) {
		g.level.Player.SetPosition(oldPos)
		return
	}

	if enemy := g.level.GetEnemyAtPosition(g.level.Player.GetPosition()); enemy != nil {
		attackResult := g.level.Attack(g.level.Player, enemy)

		var playerAttackInfo *dto.AttackInfo
		if attackResult != nil {
			attackResult.AttackerName = "Player"
			attackResult.DefenderName = g.getEnemyDisplayName(enemy)
			playerAttackInfo = dto.ConvertAttackResultToDto(attackResult)

			if attackResult.Evaded {
				g.level.Player.HitsMissed++
			} else {
				g.level.Player.HitsDealt++
			}
		}

		if provider, ok := enemy.(entities.CharacterProvider); ok {
			character := provider.GetCharacter()
			if !character.IsAlive() {
				item, err := g.level.GenerateTreasureForEnemy(g.getEnemyType(enemy))
				if err != nil {
					panic(err)
				}
				if playerAttackInfo != nil {
					playerAttackInfo.LootName = item.Name
					playerAttackInfo.LootValue = item.Value
				}
				g.level.Player.AddTreasure(item)
				g.level.RemoveEnemy(enemy)
				g.level.Player.EnemiesKilled++
			}
			g.level.Player.SetPosition(oldPos)
		}

		g.view.SetPlayerAttackInfo(playerAttackInfo)
		return
	}

	g.level.Player.CellsMoved++
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
			g.trackConsumableUseAtPosition(pos)
			g.level.PlayerUseItemAtPosition(pos)
		case world.CollisionTypeTeleport:
			if g.level.Number < maxLevelNumber {
				if err := g.level.GenerateWithExistingPlayer(g.level.Player); err != nil {
					panic(err.Error())
				}
			} else {
				if err := save.DeleteGame(SaveFileName); err != nil {
					panic(err.Error())
				}
				if err := save.SaveScore(g.buildScoreEntry(), ScoreFileName); err != nil {
					panic(err.Error())
				}
				g.GameOverStats = g.buildGameOverStats()
				g.signal = signals.GameWon
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
				g.view.SetInfoErrorMessage(err.Error())
				return
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
	case 'n', 'т':
		return action.ToggleBackpackMode
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

// getEnemyDisplayName возвращает читаемое имя типа врага для отображения в UI.
func (g *Game) getEnemyDisplayName(enemy primitives.Positional2D[int]) string {
	switch enemy.(type) {
	case *entities.Zombie:
		return "Zombie"
	case *entities.Vampire:
		return "Vampire"
	case *entities.Ghost:
		return "Ghost"
	case *entities.Ogre:
		return "Ogre"
	case *entities.SnakeMage:
		return "Snake Mage"
	case *entities.Mimic:
		return "Mimic"
	default:
		return "Unknown"
	}
}

func (g *Game) getEnemyType(enemy primitives.Positional2D[int]) entities.EnemyType {
	switch enemy.(type) {
	case *entities.Zombie:
		return entities.EnemyTypeZombie
	case *entities.Vampire:
		return entities.EnemyTypeVampire
	case *entities.Ghost:
		return entities.EnemyTypeGhost
	case *entities.Ogre:
		return entities.EnemyTypeOgre
	case *entities.SnakeMage:
		return entities.EnemyTypeSnakeMage
	case *entities.Mimic:
		return entities.EnemyTypeMimic
	default:
		return entities.EnemyTypeZombie
	}
}

func (g *Game) updateBackpackInfo() {
	if g.level.Player != nil && g.level.Player.Backpack != nil {
		g.view.UpdateBackpackInfo(dto.ConvertBackpackToDto(g.level.Player.Backpack))
	}
}

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
				g.view.SetInfoErrorMessage(err.Error())
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
		g.view.SetInfoErrorMessage(err.Error())
	} else {
		switch itemType {
		case entities.BackpackItemTypeFood:
			g.level.Player.FoodEaten++
		case entities.BackpackItemTypeElixir:
			g.level.Player.ElixirsDrunk++
		case entities.BackpackItemTypeScroll:
			g.level.Player.ScrollsRead++
		}
	}

	g.updateBackpackInfo()
}

func (g *Game) handleItemDrop(actionType action.Type) {
	if g.level.Player == nil || g.level.Player.Backpack == nil {
		return
	}

	itemIndex := int(actionType - action.Num0)

	var itemType entities.BackpackItemType
	var actualIndex int

	switch g.view.GetCurrentBackpackTab() {
	case viewcli.BackpackTabWeapons:
		itemType = entities.BackpackItemTypeWeapon
		// Для оружия индекс 0 - текущее экипированное оружие
		if itemIndex == 0 {
			if g.level.Player.Weapon != nil {
				g.dropItemToMap(g.level.Player.Weapon, true)
			}
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

	g.dropItemToMap(selectedItem, false)
}

// dropItemToMap выбрасывает предмет на карту в соседнюю свободную клетку
func (g *Game) dropItemToMap(item any, isEquipped bool) {
	if g.level.Player == nil || item == nil {
		return
	}

	positional, ok := item.(primitives.Positional2D[int])
	if !ok {
		return
	}

	playerPos := g.level.Player.GetPosition()

	// Ищем свободную соседнюю клетку по часовой стрелке
	// Порядок: N, NE, E, SE, S, SW, W, NW
	adjacentOffsets := []primitives.Point2D[int]{
		{X: 0, Y: -1},  // North
		{X: 1, Y: -1},  // North-East
		{X: 1, Y: 0},   // East
		{X: 1, Y: 1},   // South-East
		{X: 0, Y: 1},   // South
		{X: -1, Y: 1},  // South-West
		{X: -1, Y: 0},  // West
		{X: -1, Y: -1}, // North-West
	}

	var dropPos *primitives.Point2D[int]
	for _, offset := range adjacentOffsets {
		candidatePos := primitives.Point2D[int]{
			X: playerPos.X + offset.X,
			Y: playerPos.Y + offset.Y,
		}

		// Определение коллизий через получение полного поля карты
		// а не прохода по позициям логического представления сущностей.
		// В качестве примера и для учебного разнообразия.
		// Ну и так значительно проще код.
		field := g.level.GetFullField(MapWidth, MapHeight)

		if field.EnvironmentLayer[candidatePos.Y][candidatePos.X] == common.WorldTypeRoomFloor &&
			field.ObjectLayer[candidatePos.Y][candidatePos.X] == common.EntityTypeNone {
			dropPos = &candidatePos
			break
		}
	}

	if dropPos == nil {
		g.view.SetInfoErrorMessage("No free adjacent cells to drop item")
		return
	}

	positional.SetPosition(*dropPos)
	g.level.AddItem(positional)

	if isEquipped {
		_ = g.level.Player.UnequipWeapon()
	} else {
		_ = g.level.Player.Backpack.RemoveItem(item)
	}

	g.updateBackpackInfo()
}

// trackConsumableUseAtPosition проверяет, является ли предмет на указанной позиции расходным,
// и увеличивает соответствующий счётчик статистики игрока.
func (g *Game) trackConsumableUseAtPosition(pos primitives.Point2D[int]) {
	item := g.level.GetItemAtPosition(pos)
	if item == nil {
		return
	}

	switch item.(type) {
	case *items.Food:
		g.level.Player.FoodEaten++
	case *items.Elixir:
		g.level.Player.ElixirsDrunk++
	case *items.Scroll:
		g.level.Player.ScrollsRead++
	}
}

// checkPlayerDeath проверяет, жив ли игрок, и если нет - сигнализирует об окончании игры.
func (g *Game) checkPlayerDeath() {
	if g.level.Player == nil || g.level.Player.Character.IsAlive() {
		return
	}

	// Удаляем сохранение при гибели игрока
	if err := save.DeleteGame(SaveFileName); err != nil {
		panic(err.Error())
	}

	// Сохраняем рекорд (сокровища) в таблицу рекордов и при гибели
	if err := save.SaveScore(g.buildScoreEntry(), ScoreFileName); err != nil {
		panic(err.Error())
	}

	g.GameOverStats = g.buildGameOverStats()
	g.signal = signals.GameOver
}

// buildGameOverStats формирует DTO со статистикой для экрана завершения игры.
func (g *Game) buildGameOverStats() *dto.GameOverStats {
	if g.level.Player == nil {
		return &dto.GameOverStats{}
	}

	return &dto.GameOverStats{
		Treasures:     g.level.Player.Backpack.Treasures,
		EnemiesKilled: g.level.Player.EnemiesKilled,
		LevelReached:  g.level.Number,
		FoodEaten:     g.level.Player.FoodEaten,
		ElixirsDrunk:  g.level.Player.ElixirsDrunk,
		ScrollsRead:   g.level.Player.ScrollsRead,
		HitsDealt:     g.level.Player.HitsDealt,
		HitsMissed:    g.level.Player.HitsMissed,
		CellsMoved:    g.level.Player.CellsMoved,
	}
}

func (g *Game) buildScoreEntry() dto.ScoreEntry {
	if g.level.Player == nil {
		return dto.ScoreEntry{}
	}

	return dto.ScoreEntry{
		Treasures:     g.level.Player.Backpack.Treasures,
		LevelReached:  g.level.Number,
		EnemiesKilled: g.level.Player.EnemiesKilled,
		FoodEaten:     g.level.Player.FoodEaten,
		ElixirsDrunk:  g.level.Player.ElixirsDrunk,
		ScrollsRead:   g.level.Player.ScrollsRead,
		HitsDealt:     g.level.Player.HitsDealt,
		HitsMissed:    g.level.Player.HitsMissed,
		CellsMoved:    g.level.Player.CellsMoved,
	}
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
