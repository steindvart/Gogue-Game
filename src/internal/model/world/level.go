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
	Enemies []entities.Enemy
	Foods   []items.Food
	Elixirs []items.Elixir
	Scrolls []items.Scroll
	Weapons []items.Weapon

	// Метаданные
	Number uint

	// Конфигурация и зависимости
	config        LevelConfig
	random        utils.Randomizer
	roomGen       RoomGenerator
	passageGen    PassageGenerator
	entitySpawner EntitySpawner
	playerSpawner PlayerSpawner
}

func NewLevelWithDefaults(random utils.Randomizer, mapSize primitives.Size2D[uint]) *Level {
	cfg := DefaultLevelConfig(mapSize)
	return NewLevelWithComponents(random, cfg, NewGridRoomGenerator(), NewConnectionTreePassageGenerator(), NewRoomBasedEntitySpawner(), NewStartRoomPlayerSpawner())
}

func NewLevelWithComponents(random utils.Randomizer, cfg LevelConfig, roomGen RoomGenerator, passageGen PassageGenerator, entitySpawner EntitySpawner, playerSpawner PlayerSpawner) *Level {
	return &Level{
		config:        cfg,
		random:        random,
		roomGen:       roomGen,
		passageGen:    passageGen,
		entitySpawner: entitySpawner,
		playerSpawner: playerSpawner,
		Enemies:       []entities.Enemy{},
		Foods:         []items.Food{},
		Elixirs:       []items.Elixir{},
		Scrolls:       []items.Scroll{},
		Weapons:       []items.Weapon{},
	}
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

	l.Player = player
	startPos, err := l.playerSpawner.GetStartPosition(l.Rooms, l.random)
	if err != nil {
		return err
	}
	l.Player.SetPosition(*startPos)
	return nil
}

func (l *Level) generateEnvironment() error {
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
	spawned, err := l.entitySpawner.SpawnEntities(l.Rooms, l.config.ItemCounts, l.random)
	if err != nil {
		return err
	}

	l.Foods = spawned.Foods
	l.Elixirs = spawned.Elixirs
	l.Scrolls = spawned.Scrolls
	l.Weapons = spawned.Weapons
	l.Enemies = spawned.Enemies

	return nil
}

func (l *Level) MovePlayerWithCheckCollision(delta primitives.Point2D[int]) {
	oldPlayerPos := l.Player.GetPosition()
	l.Player.Move(delta)

	if l.checkCollision(l.Player.GetPosition()) {
		l.Player.SetPosition(oldPlayerPos)
	}
}

func (l *Level) checkCollision(pos primitives.Point2D[int]) bool {
	if l.checkCollisionWithMapBorders(pos) {
		return true
	}
	if l.checkCollisionWithRoomsWall(pos) {
		return true
	}
	if l.checkCollisionWithEnemy(pos) {
		return true
	}

	if !isInSomeRoom(pos, l.Rooms) && !l.checkCollisionWithPassages(pos) {
		return true
	}

	return false
}

func (l *Level) checkCollisionWithMapBorders(pos primitives.Point2D[int]) bool {
	if pos.X < 0 || pos.Y < 0 {
		return true
	}
	if pos.X >= int(l.config.MapSize.Width) || pos.Y >= int(l.config.MapSize.Height) {
		return true
	}
	return false
}

func (l *Level) checkCollisionWithRoomsWall(pos primitives.Point2D[int]) bool {
	for _, room := range l.Rooms {
		if isInRoom(pos, room) && checkCollisionWithRoomWall(pos, room) {
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
	leftEndX := room.Shape.Point.X + 1
	rightEndX := leftEndX + int(room.Shape.Size.Width) - 3
	topEndY := room.Shape.Point.Y + 1
	downEndY := topEndY + int(room.Shape.Size.Height) - 3

	return (pos.X >= leftEndX && pos.X <= rightEndX) && (pos.Y >= topEndY && pos.Y <= downEndY)
}

func checkCollisionWithRoomWall(pos primitives.Point2D[int], room Room) bool {
	if checkCollisionWithDoors(pos, room.Doors) {
		return false
	}

	leftEndX := room.Shape.Point.X
	rightEndX := leftEndX + int(room.Shape.Size.Width)
	topEndY := room.Shape.Point.Y
	downEndY := topEndY + int(room.Shape.Size.Height)

	if (pos.X == leftEndX || pos.X == rightEndX) || (pos.Y == topEndY || pos.Y == downEndY) {
		return true
	}
	return false
}

func checkCollisionWithDoors(pos primitives.Point2D[int], doors []primitives.Point2D[int]) bool {
	for _, door := range doors {
		if pos == door {
			return true
		}
	}
	return false
}

func (l *Level) checkCollisionWithPassages(newPos primitives.Point2D[int]) bool {
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

func (l *Level) checkCollisionWithEnemy(pos primitives.Point2D[int]) bool {
	for _, enemy := range l.Enemies {
		if pos == enemy.GetPosition() {
			return true
		}
	}
	return false
}

// MakeCurrentMap формирует двумерный массив примитивов для отрисовки (модельная ответственность)
func (l *Level) MakeCurrentMap(w, h int) [][]common.GameEntityType {
	field := make([][]common.GameEntityType, h)
	for y := range field {
		field[y] = make([]common.GameEntityType, w)
	}

	// Рисуем комнаты
	for _, room := range l.Rooms {
		l.putRoom(room, l.FinishPortal, field)
	}

	// Рисуем коридоры
	for _, passage := range l.Passages {
		l.putPassage(passage, field)
	}

	// Рисуем предметы и врагов (берём их напрямую из Level)
	for _, food := range l.Foods {
		pt := food.Item.Shape.Point
		if pt.Y >= 0 && pt.Y < h && pt.X >= 0 && pt.X < w {
			field[pt.Y][pt.X] = foodToEntityType(food.Type)
		}
	}
	for _, elixir := range l.Elixirs {
		pt := elixir.Item.Shape.Point
		if pt.Y >= 0 && pt.Y < h && pt.X >= 0 && pt.X < w {
			field[pt.Y][pt.X] = elixirToEntityType(elixir.Type)
		}
	}
	for _, scroll := range l.Scrolls {
		pt := scroll.Item.Shape.Point
		if pt.Y >= 0 && pt.Y < h && pt.X >= 0 && pt.X < w {
			field[pt.Y][pt.X] = scrollToEntityType(scroll.Type)
		}
	}
	for _, weapon := range l.Weapons {
		pt := weapon.Item.Shape.Point
		if pt.Y >= 0 && pt.Y < h && pt.X >= 0 && pt.X < w {
			field[pt.Y][pt.X] = common.Weapon
		}
	}

	for _, e := range l.Enemies {
		pt := e.GetPosition()
		if pt.Y >= 0 && pt.Y < h && pt.X >= 0 && pt.X < w {
			field[pt.Y][pt.X] = enemyToEntityType(e.Type)
		}
	}

	// Рисуем игрока поверх остальных
	if l.Player != nil {
		px := l.Player.Character.Shape.Point.X
		py := l.Player.Character.Shape.Point.Y
		if py >= 0 && py < h && px >= 0 && px < w {
			field[py][px] = common.EntityTypePlayer
		}
	}

	return field
}

func (l *Level) putRoom(room Room, finishPortal primitives.Box, field [][]common.GameEntityType) {
	width := int(room.Shape.Size.Width)
	height := int(room.Shape.Size.Height)

	startX := room.Shape.Point.X
	endX := startX + width - 1
	startY := room.Shape.Point.Y
	endY := startY + height - 1

	for col := startX; col <= endX; col++ {
		field[startY][col] = common.WorldTypeWall
		field[endY][col] = common.WorldTypeWall
	}
	for row := startY; row <= endY; row++ {
		field[row][startX] = common.WorldTypeWall
		field[row][endX] = common.WorldTypeWall
	}
	// Портал
	if finishPortal.Point.Y >= 0 && finishPortal.Point.Y < len(field) && finishPortal.Point.X >= 0 && finishPortal.Point.X < len(field[0]) {
		field[finishPortal.Point.Y][finishPortal.Point.X] = common.WorldTypePortal
	}
}

func (l *Level) putPassage(passage Passage, field [][]common.GameEntityType) {
	for _, p := range passage.Way {
		field[p.Y][p.X] = common.WorldTypePassage
	}
	field[passage.DoorOne.Y][passage.DoorOne.X] = common.WorldTypeDoor
	field[passage.DoorTwo.Y][passage.DoorTwo.X] = common.WorldTypeDoor
}

// Вспомогательные функции для конвертации типов
func foodToEntityType(foodType items.FoodType) common.GameEntityType {
	switch foodType {
	case items.FoodTypePotatoes:
		return common.FoodTypePotatoes
	case items.FoodTypeBread:
		return common.FoodTypeBread
	case items.FoodTypeMeat:
		return common.FoodTypeMeat
	case items.FoodTypeMistery:
		return common.FoodTypeMistery
	case items.FoodTypeBeer:
		return common.FoodTypeBeer
	default:
		return common.FoodTypeMistery
	}
}

func elixirToEntityType(elixirType items.ElixirType) common.GameEntityType {
	switch elixirType {
	case items.ElixirTypeStrength:
		return common.ElixirTypeStrength
	case items.ElixirTypeAgility:
		return common.ElixirTypeAgility
	case items.ElixirTypeDwarfism:
		return common.ElixirTypeDwarfism
	case items.ElixirTypeGiantism:
		return common.ElixirTypeGiantism
	case items.ElixirTypeMystery:
		return common.ElixirTypeMystery
	default:
		return common.ElixirTypeMystery
	}
}

func scrollToEntityType(scrollType items.ScrollType) common.GameEntityType {
	switch scrollType {
	case items.ScrollTypeStrength:
		return common.ScrollTypeStrength
	case items.ScrollTypeAgility:
		return common.ScrollTypeAgility
	case items.ScrollTypeUltimate:
		return common.ScrollTypeUltimate
	case items.ScrollTypeMaxHealth:
		return common.ScrollTypeMaxHealth
	case items.ScrollTypeMystery:
		return common.ScrollTypeMystery
	default:
		return common.ScrollTypeMystery
	}
}

func enemyToEntityType(t entities.EnemyType) common.GameEntityType {
	switch t {
	case entities.EnemyTypeZombie:
		return common.EntityTypeZombie
	case entities.EnemyTypeVampire:
		return common.EntityTypeVampire
	case entities.EnemyTypeGhost:
		return common.EntityTypeGhost
	case entities.EnemyTypeOgre:
		return common.EntityTypeOgre
	case entities.EnemyTypeSnakeMage:
		return common.EntityTypeSnakeMage
	default:
		return common.EntityTypeZombie
	}
}

// GetItemAtPosition возвращает предмет на позиции (если есть)
func (l *Level) GetItemAtPosition(pos primitives.Point2D[int]) interface{} {
	for i := range l.Foods {
		if l.Foods[i].Item.Shape.Point == pos {
			return &l.Foods[i]
		}
	}
	for i := range l.Elixirs {
		if l.Elixirs[i].Item.Shape.Point == pos {
			return &l.Elixirs[i]
		}
	}
	for i := range l.Scrolls {
		if l.Scrolls[i].Item.Shape.Point == pos {
			return &l.Scrolls[i]
		}
	}
	for i := range l.Weapons {
		if l.Weapons[i].Item.Shape.Point == pos {
			return &l.Weapons[i]
		}
	}
	return nil
}

// Удаление предметов по индексу
func (l *Level) RemoveFood(index int) {
	l.Foods = append(l.Foods[:index], l.Foods[index+1:]...)
}
func (l *Level) RemoveElixir(index int) {
	l.Elixirs = append(l.Elixirs[:index], l.Elixirs[index+1:]...)
}
func (l *Level) RemoveScroll(index int) {
	l.Scrolls = append(l.Scrolls[:index], l.Scrolls[index+1:]...)
}
func (l *Level) RemoveWeapon(index int) {
	l.Weapons = append(l.Weapons[:index], l.Weapons[index+1:]...)
}

// AddWeapon добавляет оружие на уровень (например, при выбрасывании)
func (l *Level) AddWeapon(weapon items.Weapon) {
	l.Weapons = append(l.Weapons, weapon)
}
