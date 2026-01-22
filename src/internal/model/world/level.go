package world

import (
	"gogue/internal/common"
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Level struct {
	// Геометрия
	Rooms        []Room
	Passages     []Passage
	FinishPortal primitives.Box

	// Сущности, хранящиеся непосредственно в уровне
	Player  *entities.Player
	Enemies []primitives.Positional2D[int]
	Items   []primitives.Positional2D[int]

	// Метаданные
	Number uint

	// Конфигурация и зависимости
	config        LevelConfig
	random        utils.Randomizer
	roomGen       RoomGenerator
	passageGen    PassageGenerator
	entitySpawner EntitySpawner
	playerSpawner PlayerSpawner
	fieldRenderer FieldRenderer
	FogOfWar      *FogOfWar
}

type CollisionType int

const (
	CollisionTypeNone CollisionType = iota
	CollisionTypeEnemy
	CollisionTypeItem
	CollisionTypeTeleport
)

func NewLevelWithDefaults(random utils.Randomizer, mapSize primitives.Size2D[uint]) *Level {
	cfg := DefaultLevelConfig(mapSize)
	return newLevelWithComponents(
		random,
		cfg,
		1,
		NewGridRoomGenerator(),
		NewConnectionTreePassageGenerator(),
		NewRoomBasedEntitySpawner(),
		NewStartRoomPlayerSpawner(),
		NewDefaultFieldRenderer(),
	)
}

func newLevelWithComponents(
	random utils.Randomizer,
	cfg LevelConfig,
	levelNumber uint,
	roomGen RoomGenerator,
	passageGen PassageGenerator,
	entitySpawner EntitySpawner,
	playerSpawner PlayerSpawner,
	fieldRenderer FieldRenderer,
) *Level {
	return &Level{
		config:        cfg,
		random:        random,
		roomGen:       roomGen,
		passageGen:    passageGen,
		entitySpawner: entitySpawner,
		playerSpawner: playerSpawner,
		Number:        levelNumber,
		fieldRenderer: fieldRenderer,
		FogOfWar:      NewFogOfWar(int(cfg.MapSize.Width), int(cfg.MapSize.Height)),
		Enemies:       []primitives.Positional2D[int]{},
		Items:         []primitives.Positional2D[int]{},
	}
}

func (l *Level) ProcessTurns(turns uint32) {
	l.Player.ProcessTurns(turns)
	// @todo - движение врагов и их взаимодействие с миром (хождение по миру, агрессия и нападение на игрока)
}

// Generate генерирует геометрию, сущности и игрока
func (l *Level) Generate() error {
	err := l.generateEnvironment()
	if err != nil {
		return err
	}

	player, err := l.playerSpawner.SpawnPlayer(l.Rooms, l.random)
	if err != nil {
		return err
	}
	l.Player = player

	return nil
}

// GenerateWithExistingPlayer генерирует мир, но использует переданного игрока
func (l *Level) GenerateWithExistingPlayer(player *entities.Player) error {
	err := l.generateEnvironment()
	if err != nil {
		return err
	}

	l.Number++
	l.Player = player
	startPos, err := l.playerSpawner.GetStartPosition(l.Rooms, l.random)
	if err != nil {
		return err
	}
	l.Player.SetPosition(*startPos)
	return nil
}

func (l *Level) generateEnvironment() error {
	// Сбрасываем туман войны при генерации нового уровня
	l.FogOfWar.Reset()

	// Геометрия
	rooms, finishPortal, err := l.roomGen.GenerateRooms(l.config, l.random)
	if err != nil {
		return err
	}
	l.Rooms = rooms
	l.FinishPortal = finishPortal

	passages, err := l.passageGen.GeneratePassages(l.Rooms, l.config, l.random)
	if err != nil {
		return err
	}
	l.Passages = passages

	// Сущности
	spawned, err := l.entitySpawner.SpawnEntities(l.Rooms, l.config.ItemsSpawnConfig, l.config.EnemiesSpawnConfig, l.random)
	if err != nil {
		return err
	}

	l.Items = spawned.Items
	l.Enemies = spawned.Enemies

	return nil
}

func (l *Level) CheckEntityCollision(pos primitives.Point2D[int]) CollisionType {
	if l.isCollisionWithEnemy(pos) {
		return CollisionTypeEnemy
	}
	if l.isCollisionWithItem(pos) {
		return CollisionTypeItem
	}
	if l.isCollisionWithTeleport(pos) {
		return CollisionTypeTeleport
	}

	return CollisionTypeNone
}

func (l *Level) IsCollisionWithBorders(pos primitives.Point2D[int]) bool {
	if l.isCollisionWithMapBorders(pos) {
		return true
	}
	if l.isCollisionWithRoomsWall(pos) {
		return true
	}
	if !isInSomeRoom(pos, l.Rooms) && !l.isCollisionWithPassages(pos) {
		return true
	}

	return false
}

func (l *Level) isCollisionWithMapBorders(pos primitives.Point2D[int]) bool {
	if pos.X < 0 || pos.Y < 0 {
		return true
	}
	if pos.X >= int(l.config.MapSize.Width) || pos.Y >= int(l.config.MapSize.Height) {
		return true
	}
	return false
}

func (l *Level) isCollisionWithRoomsWall(pos primitives.Point2D[int]) bool {
	for _, room := range l.Rooms {
		if isInRoom(pos, room) && isCollisionWithRoomWall(pos, room) {
			return true
		}
	}
	return false
}

func isInSomeRoom(pos primitives.Point2D[int], rooms []Room) bool {
	for _, room := range rooms {
		if isInRoom(pos, room) {
			return true
		}
	}
	return false
}

func isInRoom(pos primitives.Point2D[int], room Room) bool {
	leftEndX := room.Box.Point.X + 1
	rightEndX := leftEndX + int(room.Box.Size.Width) - 3
	topEndY := room.Box.Point.Y + 1
	downEndY := topEndY + int(room.Box.Size.Height) - 3

	return (pos.X >= leftEndX && pos.X <= rightEndX) && (pos.Y >= topEndY && pos.Y <= downEndY)
}

func isCollisionWithRoomWall(pos primitives.Point2D[int], room Room) bool {
	if isCollisionWithDoors(pos, room.Doors) {
		return false
	}

	leftEndX := room.Box.Point.X
	rightEndX := leftEndX + int(room.Box.Size.Width)
	topEndY := room.Box.Point.Y
	downEndY := topEndY + int(room.Box.Size.Height)

	if (pos.X == leftEndX || pos.X == rightEndX) || (pos.Y == topEndY || pos.Y == downEndY) {
		return true
	}
	return false
}

func isCollisionWithDoors(pos primitives.Point2D[int], doors []primitives.Point2D[int]) bool {
	for _, door := range doors {
		if pos == door {
			return true
		}
	}
	return false
}

func (l *Level) isCollisionWithPassages(newPos primitives.Point2D[int]) bool {
	for _, passage := range l.Passages {
		if isInPassage(newPos, passage) {
			return true
		}
	}
	return false
}

func isInPassage(pos primitives.Point2D[int], passage Passage) bool {
	if pos == passage.DoorOne || pos == passage.DoorTwo {
		return true
	}
	for _, wayPoint := range passage.Way {
		if pos == wayPoint {
			return true
		}
	}
	return false
}

func (l *Level) isCollisionWithTeleport(delta primitives.Point2D[int]) bool {
	return delta == l.FinishPortal.Point
}

func (l *Level) GetEnemyAtPosition(pos primitives.Point2D[int]) primitives.Positional2D[int] {
	for _, enemy := range l.Enemies {
		if pos == enemy.GetPosition() {
			return enemy
		}
	}
	return nil
}

func (l *Level) isCollisionWithEnemy(pos primitives.Point2D[int]) bool {
	return l.GetEnemyAtPosition(pos) != nil
}

func (l *Level) isCollisionWithItem(pos primitives.Point2D[int]) bool {
	for _, item := range l.Items {
		if pos == item.GetPosition() {
			return true
		}
	}
	return false
}

// MakeCurrentField создаёт двумерное представление карты уровня с учётом тумана войны
func (l *Level) MakeCurrentField(w, h int) [][]common.GameEntityType {
	// Сначала рендерим полное поле
	fullField := l.GetFullField(w, h)

	// Если игрока нет, возвращаем пустое поле (всё скрыто туманом войны)
	if l.Player == nil {
		return utils.CreateEmpty2DSlice[common.GameEntityType](h, w)
	}

	// Фильтруем поле с учётом тумана войны и радиуса обзора игрока
	return l.FogOfWar.ApplyFogOfWar(fullField, l.Player.GetPosition(), l.Player.ViewRadius)
}

// GetFullField возвращает двумерное представление карты уровня бещ учёта тумана войны
func (l *Level) GetFullField(w, h int) [][]common.GameEntityType {
	return l.fieldRenderer.RenderField(w, h, l)
}

func (l *Level) PlayerUseItemAtPosition(pos primitives.Point2D[int]) {
	item := l.GetItemAtPosition(pos)
	if item == nil {
		return
	}

	// @todo - кривая логика, нужно декомпозировать и перенести в Player
	if weapon, ok := item.(*items.Weapon); ok {
		_ = l.Player.EquipWeapon(weapon)
	} else if usableItem, ok := item.(items.Usable); ok {
		l.Player.Character.Use(usableItem)
	}

	l.RemoveItem(item)
}

func (l *Level) PlayerTakeItemAtPosition(pos primitives.Point2D[int]) error {
	item := l.GetItemAtPosition(pos)
	if item == nil {
		return nil
	}

	if takeableItem, ok := item.(items.Takeable); ok {
		if err := l.Player.Backpack.AddItem(takeableItem); err != nil {
			return err
		}
		l.RemoveItem(item)
	}

	return nil
}

func (l *Level) GetItemAtPosition(pos primitives.Point2D[int]) primitives.Positional2D[int] {
	for _, item := range l.Items {
		if item.GetPosition() == pos {
			return item
		}
	}
	return nil
}

func (l *Level) RemoveItem(item primitives.Positional2D[int]) {
	for i, it := range l.Items {
		if it == item {
			l.Items = append(l.Items[:i], l.Items[i+1:]...)
			return
		}
	}
}

func (l *Level) AddItem(item primitives.Positional2D[int]) {
	l.Items = append(l.Items, item)
}

func (l *Level) RemoveEnemy(enemy primitives.Positional2D[int]) {
	for i, it := range l.Enemies {
		if it == enemy {
			l.Enemies = append(l.Enemies[:i], l.Enemies[i+1:]...)
			return
		}
	}
}

type Attacker interface {
	Attack(defender *entities.Character, rnd utils.Randomizer)
}

// @todo - возвращать информацию о совершенной атаке - кто был атакован, какой был нанесён урон, сколько здоровья осталось и т.д.
func (l *Level) Attack(attacker Attacker, defender primitives.Positional2D[int]) {
	if attacker == nil || defender == nil {
		return
	}

	if enemy, ok := defender.(entities.CharacterProvider); ok {
		attacker.Attack(enemy.GetCharacter(), l.random)
		return
	}
}

// @todo - с повышением lvl повысится кол-во созданных сокровищ
func (l *Level) GenerateTreasure(goldPercent, gemPercent, artifactPercent, mysteryPercent uint) (*items.Treasure, error) {
	item, err := items.GetTreasureType(l.random, goldPercent, gemPercent, artifactPercent, mysteryPercent)
	if err != nil {
		return nil, err
	}

	return items.NewTreasureBuiltin(
		l.random,
		primitives.Box{
			Point: primitives.Point2D[int]{X: 0, Y: 0},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
		item), nil
}
