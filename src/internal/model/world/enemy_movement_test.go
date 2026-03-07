package world

import (
	"gogue/internal/common"
	"gogue/internal/model/entities"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
	"testing"
)

// --- Вспомогательные функции для тестов ---

// makeBox создаёт Box с указанной позицией и размером 1x1.
func makeBox(x, y int) *primitives.Box {
	return &primitives.Box{
		Point: primitives.Point2D[int]{X: x, Y: y},
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}
}

// makeTestRoom создаёт Room с указанными параметрами.
func makeTestRoom(x, y int, w, h uint) Room {
	room, _ := NewRoom(RoomTypeOrdinary, primitives.Box{
		Point: primitives.Point2D[int]{X: x, Y: y},
		Size:  primitives.Size2D[uint]{Width: w, Height: h},
	})
	return *room
}

// makeFieldFromTypes создаёт двумерное поле из заданного среза типов.
func makeFieldFromTypes(rows [][]common.GameEntityType) [][]common.GameEntityType {
	return rows
}

// buildSimpleField создаёт простое поле: комната 7x7 с полом и стенами.
//
//	#######
//	#.....#
//	#.....#
//	#.....#
//	#.....#
//	#.....#
//	#######
func buildSimpleField() [][]common.GameEntityType {
	h, w := 7, 7
	field := make([][]common.GameEntityType, h)
	for y := 0; y < h; y++ {
		field[y] = make([]common.GameEntityType, w)
		for x := 0; x < w; x++ {
			if y == 0 || y == h-1 || x == 0 || x == w-1 {
				field[y][x] = common.WorldTypeWall
			} else {
				field[y][x] = common.WorldTypeRoomFloor
			}
		}
	}
	return field
}

func TestEnemyIdleMover_ZombieMovesCardinal(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	zombie := entities.NewZombie(makeBox(3, 3))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{zombie}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	startPos := zombie.GetPosition()
	moved := mover.MoveIdle(zombie, field, rooms, enemies)

	if !moved {
		t.Fatal("Zombie should be able to move in open room")
	}

	newPos := zombie.GetPosition()
	if newPos == startPos {
		t.Fatal("Zombie should have moved to a new position")
	}

	// Проверяем, что движение было кардинальным (N/S/W/E, не диагональ)
	dx := newPos.X - startPos.X
	dy := newPos.Y - startPos.Y

	isCardinal := (dx == 0 && (dy == 1 || dy == -1)) || (dy == 0 && (dx == 1 || dx == -1))
	if !isCardinal {
		t.Errorf("Zombie moved non-cardinally: from %+v to %+v", startPos, newPos)
	}
}

func TestEnemyIdleMover_VampireMovesCardinal(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	vampire := entities.NewVampire(makeBox(3, 3))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{vampire}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	startPos := vampire.GetPosition()
	moved := mover.MoveIdle(vampire, field, rooms, enemies)

	if !moved {
		t.Fatal("Vampire should be able to move in open room")
	}

	newPos := vampire.GetPosition()
	dx := newPos.X - startPos.X
	dy := newPos.Y - startPos.Y

	isCardinal := (dx == 0 && (dy == 1 || dy == -1)) || (dy == 0 && (dx == 1 || dx == -1))
	if !isCardinal {
		t.Errorf("Vampire moved non-cardinally: from %+v to %+v", startPos, newPos)
	}
}

func TestEnemyIdleMover_OgreMovesUpToTwoSteps(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	ogre := entities.NewOgre(makeBox(3, 3))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{ogre}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	startPos := ogre.GetPosition()
	moved := mover.MoveIdle(ogre, field, rooms, enemies)

	if !moved {
		t.Fatal("Ogre should be able to move in open room")
	}

	newPos := ogre.GetPosition()
	dx := newPos.X - startPos.X
	dy := newPos.Y - startPos.Y

	// Огр должен сделать ровно 2 шага (если есть место) или 1 шаг
	manhattanDist := abs(dx) + abs(dy)
	if manhattanDist != 2 && manhattanDist != 1 {
		t.Errorf("Ogre moved %d steps, expected 1 or 2: from %+v to %+v", manhattanDist, startPos, newPos)
	}

	// Движение должно быть по одной оси (кардинальное)
	isCardinal := (dx == 0 || dy == 0)
	if !isCardinal {
		t.Errorf("Ogre moved diagonally: from %+v to %+v", startPos, newPos)
	}
}

func TestEnemyIdleMover_GhostTeleportsInRoom(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	ghost := entities.NewGhost(makeBox(3, 3))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{ghost}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	startPos := ghost.GetPosition()

	// Ghost может телепортироваться куда угодно в комнате — проверяем за несколько итераций
	movedFar := false
	for i := 0; i < 20; i++ {
		ghost.SetPosition(startPos)
		mover.MoveIdle(ghost, field, rooms, enemies)

		newPos := ghost.GetPosition()
		dx := abs(newPos.X - startPos.X)
		dy := abs(newPos.Y - startPos.Y)

		// Телепортация — может переместиться более чем на 1 клетку
		if dx > 1 || dy > 1 {
			movedFar = true
			break
		}
	}

	if !movedFar {
		t.Error("Ghost should teleport to distant positions in the room (not just adjacent)")
	}
}

func TestEnemyIdleMover_GhostTogglesVisibility(t *testing.T) {
	rng := utils.NewRandomWithSeed(0)
	mover := NewEnemyIdleMover(rng)

	ghost := entities.NewGhost(makeBox(3, 3))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{ghost}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	initialVisibility := ghost.IsVisible
	changedVisibility := false

	// За 50 итераций невидимость должна хотя бы раз переключиться (~30% шанс)
	for i := 0; i < 50; i++ {
		ghost.SetPosition(primitives.Point2D[int]{X: 3, Y: 3})
		mover.MoveIdle(ghost, field, rooms, enemies)

		if ghost.IsVisible != initialVisibility {
			changedVisibility = true
			break
		}
	}

	if !changedVisibility {
		t.Error("Ghost visibility should toggle at some point during idle movement")
	}
}

func TestEnemyIdleMover_SnakeMageMovesDiagonally(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	snake := entities.NewSnakeMage(makeBox(3, 3))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{snake}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	startPos := snake.GetPosition()
	moved := mover.MoveIdle(snake, field, rooms, enemies)

	if !moved {
		t.Fatal("SnakeMage should be able to move in open room")
	}

	newPos := snake.GetPosition()
	dx := abs(newPos.X - startPos.X)
	dy := abs(newPos.Y - startPos.Y)

	// Движение по диагонали: |dx| == 1 && |dy| == 1
	if dx != 1 || dy != 1 {
		t.Errorf("SnakeMage should move diagonally: dx=%d, dy=%d (from %+v to %+v)", dx, dy, startPos, newPos)
	}
}

func TestEnemyIdleMover_SnakeMageChangesDirection(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	snake := entities.NewSnakeMage(makeBox(3, 3))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{snake}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	// Первый ход
	mover.MoveIdle(snake, field, rooms, enemies)
	dirAfterFirst := snake.Direction

	// Второй ход (возвращаем на центр)
	snake.SetPosition(primitives.Point2D[int]{X: 3, Y: 3})
	mover.MoveIdle(snake, field, rooms, enemies)
	dirAfterSecond := snake.Direction

	// Направление должно меняться после каждого хода
	if dirAfterFirst == dirAfterSecond {
		t.Errorf("SnakeMage direction should change: first=%d, second=%d", dirAfterFirst, dirAfterSecond)
	}
}

func TestEnemyIdleMover_EnemyDoesNotMoveIntoWall(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	// Зомби зажат в углу — может двигаться только вправо или вниз
	zombie := entities.NewZombie(makeBox(1, 1))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{zombie}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	for i := 0; i < 50; i++ {
		zombie.SetPosition(primitives.Point2D[int]{X: 1, Y: 1})
		mover.MoveIdle(zombie, field, rooms, enemies)

		pos := zombie.GetPosition()
		// Не должен оказаться на стене
		if field[pos.Y][pos.X] == common.WorldTypeWall {
			t.Fatalf("Zombie moved into a wall at %+v", pos)
		}
	}
}

func TestEnemyIdleMover_EnemyDoesNotMoveIntoOtherEnemy(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	zombie1 := entities.NewZombie(makeBox(2, 3))
	zombie2 := entities.NewZombie(makeBox(3, 3))

	// Создаем маленькое поле: zombie1 зажат, zombie2 рядом
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{zombie1, zombie2}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	for i := 0; i < 50; i++ {
		zombie1.SetPosition(primitives.Point2D[int]{X: 2, Y: 3})
		zombie2.SetPosition(primitives.Point2D[int]{X: 3, Y: 3})

		mover.MoveIdle(zombie1, field, rooms, enemies)

		// zombie1 не должен оказаться на позиции zombie2
		if zombie1.GetPosition() == zombie2.GetPosition() {
			t.Fatalf("Zombie1 moved into Zombie2's position at %+v", zombie1.GetPosition())
		}
	}
}

func TestEnemyIdleMover_NoMovementWhenBlocked(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	// Создаём крохотную комнату 3x3 (внутри 1 клетка пола)
	field := [][]common.GameEntityType{
		{common.WorldTypeWall, common.WorldTypeWall, common.WorldTypeWall},
		{common.WorldTypeWall, common.WorldTypeRoomFloor, common.WorldTypeWall},
		{common.WorldTypeWall, common.WorldTypeWall, common.WorldTypeWall},
	}

	zombie := entities.NewZombie(makeBox(1, 1))
	enemies := []primitives.Positional2D[int]{zombie}
	rooms := []Room{makeTestRoom(0, 0, 3, 3)}

	moved := mover.MoveIdle(zombie, field, rooms, enemies)
	if moved {
		t.Error("Zombie should not be able to move when fully surrounded by walls")
	}
}

func TestNextSnakeMageDirection_CyclesCorrectly(t *testing.T) {
	tests := []struct {
		input entities.Direction
		want  entities.Direction
	}{
		{entities.DirectionDiagonallyForwardRight, entities.DirectionDiagonallyBackRight},   // NE → SE
		{entities.DirectionDiagonallyBackRight, entities.DirectionDiagonallyBackLeft},       // SE → SW
		{entities.DirectionDiagonallyBackLeft, entities.DirectionDiagonallyForwardLeft},     // SW → NW
		{entities.DirectionDiagonallyForwardLeft, entities.DirectionDiagonallyForwardRight}, // NW → NE
	}

	for _, tt := range tests {
		got := nextSnakeMageDirection(tt.input)
		if got != tt.want {
			t.Errorf("nextSnakeMageDirection(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestSnakeMageDirectionToOffset_Roundtrip(t *testing.T) {
	directions := []entities.Direction{
		entities.DirectionDiagonallyForwardRight,
		entities.DirectionDiagonallyBackRight,
		entities.DirectionDiagonallyBackLeft,
		entities.DirectionDiagonallyForwardLeft,
	}

	for _, dir := range directions {
		offset := snakeMageDirectionToOffset(dir)
		restored := offsetToSnakeMageDirection(offset)
		if restored != dir {
			t.Errorf("Roundtrip failed: dir=%d → offset=%+v → dir=%d", dir, offset, restored)
		}
	}
}

func TestIsWalkableTile(t *testing.T) {
	field := [][]common.GameEntityType{
		{common.WorldTypeWall, common.WorldTypeRoomFloor, common.EntityTypeNone},
		{common.WorldTypePassage, common.WorldTypeDoor, common.WorldTypePortal},
		{common.Food, common.Elixir, common.Scroll},
	}

	tests := []struct {
		pos  primitives.Point2D[int]
		want bool
	}{
		{primitives.Point2D[int]{X: 0, Y: 0}, false},  // Wall
		{primitives.Point2D[int]{X: 1, Y: 0}, true},   // Floor
		{primitives.Point2D[int]{X: 2, Y: 0}, false},  // None
		{primitives.Point2D[int]{X: 0, Y: 1}, true},   // Passage
		{primitives.Point2D[int]{X: 1, Y: 1}, true},   // Door
		{primitives.Point2D[int]{X: 2, Y: 1}, true},   // Portal
		{primitives.Point2D[int]{X: 0, Y: 2}, true},   // Food
		{primitives.Point2D[int]{X: 1, Y: 2}, true},   // Elixir
		{primitives.Point2D[int]{X: 2, Y: 2}, true},   // Scroll
		{primitives.Point2D[int]{X: -1, Y: 0}, false}, // Out of bounds
		{primitives.Point2D[int]{X: 0, Y: 5}, false},  // Out of bounds
	}

	for _, tt := range tests {
		got := isWalkableTile(tt.pos, field)
		if got != tt.want {
			t.Errorf("isWalkableTile(%+v) = %v, want %v", tt.pos, got, tt.want)
		}
	}
}

// abs возвращает абсолютное значение int
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
