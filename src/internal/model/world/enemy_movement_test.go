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

// --- Тесты патрулирования (Zombie/Vampire/Ogre) ---

func TestEnemyIdleMover_ZombiePatrolKeepsDirection(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	zombie := entities.NewZombie(makeBox(3, 3))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{zombie}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	// Первый MoveIdle назначает направление и StepsRemaining
	mover.MoveIdle(zombie, field, rooms, enemies)
	dirAfterFirst := zombie.Direction
	stepsAfterFirst := zombie.StepsRemaining

	if stepsAfterFirst < minPatrolSteps-1 || stepsAfterFirst > maxPatrolSteps-1 {
		// Один шаг уже декрементирован, поэтому диапазон [min-1, max-1]
		t.Errorf("StepsRemaining after first move = %d, expected in [%d, %d]",
			stepsAfterFirst, minPatrolSteps-1, maxPatrolSteps-1)
	}

	// Ставим обратно в центр, делаем ещё один шаг - направление должно сохраниться
	zombie.SetPosition(primitives.Point2D[int]{X: 3, Y: 3})
	mover.MoveIdle(zombie, field, rooms, enemies)

	if zombie.Direction != dirAfterFirst {
		t.Errorf("Zombie should keep direction during patrol: first=%d, second=%d",
			dirAfterFirst, zombie.Direction)
	}
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

	dx := newPos.X - startPos.X
	dy := newPos.Y - startPos.Y

	isCardinal := (dx == 0 && (dy == 1 || dy == -1)) || (dy == 0 && (dx == 1 || dx == -1))
	if !isCardinal {
		t.Errorf("Zombie moved non-cardinally: from %+v to %+v", startPos, newPos)
	}
}

func TestEnemyIdleMover_VampirePatrolKeepsDirection(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	vampire := entities.NewVampire(makeBox(3, 3))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{vampire}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	mover.MoveIdle(vampire, field, rooms, enemies)
	dirFirst := vampire.Direction

	vampire.SetPosition(primitives.Point2D[int]{X: 3, Y: 3})
	mover.MoveIdle(vampire, field, rooms, enemies)

	if vampire.Direction != dirFirst {
		t.Errorf("Vampire should keep direction during patrol: first=%d, second=%d",
			dirFirst, vampire.Direction)
	}
}

func TestEnemyIdleMover_PatrolChangesDirectionAfterStepsExpire(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	zombie := entities.NewZombie(makeBox(3, 3))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{zombie}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	// Первый вызов - назначает направление
	mover.MoveIdle(zombie, field, rooms, enemies)
	firstDir := zombie.Direction

	// Искусственно обнуляем шаги - следующий вызов должен выбрать новое направление
	zombie.StepsRemaining = 0
	zombie.SetPosition(primitives.Point2D[int]{X: 3, Y: 3})

	// Делаем достаточно попыток чтобы хотя бы раз направление сменилось
	dirChanged := false
	for i := 0; i < 20; i++ {
		zombie.StepsRemaining = 0
		zombie.SetPosition(primitives.Point2D[int]{X: 3, Y: 3})
		mover.MoveIdle(zombie, field, rooms, enemies)
		if zombie.Direction != firstDir {
			dirChanged = true
			break
		}
	}

	if !dirChanged {
		t.Error("Zombie should eventually change direction after steps expire")
	}
}

func TestEnemyIdleMover_PatrolStandsWhenBlocked(t *testing.T) {
	rng := utils.NewRandomWithSeed(0)
	mover := NewEnemyIdleMover(rng)

	zombie := entities.NewZombie(makeBox(1, 1))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{zombie}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	// Принудительно ставим направление «на север» - (1,1) уже у стены
	zombie.Direction = entities.DirectionForward // North -> (1, 0) = wall
	zombie.StepsRemaining = 3

	startPos := zombie.GetPosition()
	moved := mover.MoveIdle(zombie, field, rooms, enemies)

	if moved {
		t.Error("Zombie facing a wall should not move")
	}

	if zombie.GetPosition() != startPos {
		t.Errorf("Zombie should stay at %+v, got %+v", startPos, zombie.GetPosition())
	}

	// StepsRemaining должен был декрементироваться, даже стоя на месте
	if zombie.StepsRemaining != 2 {
		t.Errorf("StepsRemaining should be 2, got %d", zombie.StepsRemaining)
	}
}

func TestEnemyIdleMover_OgreMovesUpToTwoStepsPerMove(t *testing.T) {
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

	manhattanDist := abs(dx) + abs(dy)
	if manhattanDist != 2 && manhattanDist != 1 {
		t.Errorf("Ogre moved %d steps, expected 1 or 2: from %+v to %+v", manhattanDist, startPos, newPos)
	}

	isCardinal := (dx == 0 || dy == 0)
	if !isCardinal {
		t.Errorf("Ogre moved diagonally: from %+v to %+v", startPos, newPos)
	}
}

func TestEnemyIdleMover_OgrePatrolKeepsDirection(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	ogre := entities.NewOgre(makeBox(3, 3))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{ogre}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	mover.MoveIdle(ogre, field, rooms, enemies)
	dirFirst := ogre.Direction

	ogre.SetPosition(primitives.Point2D[int]{X: 3, Y: 3})
	mover.MoveIdle(ogre, field, rooms, enemies)

	if ogre.Direction != dirFirst {
		t.Errorf("Ogre should keep direction during patrol: first=%d, second=%d",
			dirFirst, ogre.Direction)
	}
}

// --- Тесты Ghost ---

func TestEnemyIdleMover_GhostTeleportsInRoom(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	ghost := entities.NewGhost(makeBox(3, 3))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{ghost}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	startPos := ghost.GetPosition()

	movedFar := false
	for i := 0; i < 20; i++ {
		ghost.SetPosition(startPos)
		mover.MoveIdle(ghost, field, rooms, enemies)

		newPos := ghost.GetPosition()
		dx := abs(newPos.X - startPos.X)
		dy := abs(newPos.Y - startPos.Y)

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

// --- Тесты SnakeMage ---

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

	if dx != 1 || dy != 1 {
		t.Errorf("SnakeMage should move diagonally: dx=%d, dy=%d (from %+v to %+v)", dx, dy, startPos, newPos)
	}
}

func TestEnemyIdleMover_SnakeMageKeepsDirectionUntilBlocked(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	// Начальная позиция (3,3), комната 7x7 (пол 1..5)
	snake := entities.NewSnakeMage(makeBox(3, 3))
	snake.Direction = entities.DirectionDiagonallyForwardRight // NE: dx=+1, dy=-1
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{snake}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	// Первый ход: (3,3) -> NE -> (4,2) - свободно, направление НЕ должно измениться
	moved := mover.MoveIdle(snake, field, rooms, enemies)
	if !moved {
		t.Fatal("SnakeMage should move NE from (3,3)")
	}
	if snake.GetPosition().X != 4 || snake.GetPosition().Y != 2 {
		t.Fatalf("Expected (4,2), got %+v", snake.GetPosition())
	}
	if snake.Direction != entities.DirectionDiagonallyForwardRight {
		t.Errorf("SnakeMage should keep NE direction while path is clear, got %d", snake.Direction)
	}

	// Второй ход: (4,2) -> NE -> (5,1) - свободно
	moved = mover.MoveIdle(snake, field, rooms, enemies)
	if !moved {
		t.Fatal("SnakeMage should move NE from (4,2)")
	}
	if snake.GetPosition().X != 5 || snake.GetPosition().Y != 1 {
		t.Fatalf("Expected (5,1), got %+v", snake.GetPosition())
	}
	if snake.Direction != entities.DirectionDiagonallyForwardRight {
		t.Errorf("SnakeMage should keep NE direction, got %d", snake.Direction)
	}

	// Третий ход: (5,1) -> NE -> (6,0) = стена - должен сменить направление на случайную диагональ
	moved = mover.MoveIdle(snake, field, rooms, enemies)
	if !moved {
		t.Fatal("SnakeMage should find an alternative diagonal from (5,1)")
	}

	// Направление должно измениться (теперь не NE, потому что NE заблокирован)
	// Точное направление зависит от RNG, но оно должно быть одним из диагональных
	newDir := snake.Direction
	isDiagonal := newDir == entities.DirectionDiagonallyForwardRight ||
		newDir == entities.DirectionDiagonallyBackRight ||
		newDir == entities.DirectionDiagonallyBackLeft ||
		newDir == entities.DirectionDiagonallyForwardLeft
	if !isDiagonal {
		t.Errorf("SnakeMage should have a diagonal direction after bounce, got %d", newDir)
	}
}

func TestEnemyIdleMover_SnakeMageRandomDirectionOnBlock(t *testing.T) {
	// Проверяем, что при блокировке SnakeMage не всегда выбирает одно и то же
	// направление (т.е. выбор случайный, а не циклический)
	field := buildSimpleField()
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	directionsSeen := make(map[entities.Direction]bool)

	for seed := int64(0); seed < 50; seed++ {
		rng := utils.NewRandomWithSeed(seed)
		mover := NewEnemyIdleMover(rng)

		snake := entities.NewSnakeMage(makeBox(5, 1))
		snake.Direction = entities.DirectionDiagonallyForwardRight // NE: (6,0) = wall
		enemies := []primitives.Positional2D[int]{snake}

		mover.MoveIdle(snake, field, rooms, enemies)
		directionsSeen[snake.Direction] = true
	}

	// Должно быть несколько различных направлений (не один фиксированный цикл)
	if len(directionsSeen) < 2 {
		t.Errorf("SnakeMage should choose random directions on block, saw only %d unique", len(directionsSeen))
	}
}

// --- Общие тесты ---

func TestEnemyIdleMover_EnemyDoesNotMoveIntoWall(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

	zombie := entities.NewZombie(makeBox(1, 1))
	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{zombie}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	for i := 0; i < 50; i++ {
		zombie.SetPosition(primitives.Point2D[int]{X: 1, Y: 1})
		zombie.StepsRemaining = 0 // сбрасываем, чтобы каждый раз перевыбирало
		mover.MoveIdle(zombie, field, rooms, enemies)

		pos := zombie.GetPosition()
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

	field := buildSimpleField()
	enemies := []primitives.Positional2D[int]{zombie1, zombie2}
	rooms := []Room{makeTestRoom(0, 0, 7, 7)}

	for i := 0; i < 50; i++ {
		zombie1.SetPosition(primitives.Point2D[int]{X: 2, Y: 3})
		zombie2.SetPosition(primitives.Point2D[int]{X: 3, Y: 3})
		zombie1.StepsRemaining = 0

		mover.MoveIdle(zombie1, field, rooms, enemies)

		if zombie1.GetPosition() == zombie2.GetPosition() {
			t.Fatalf("Zombie1 moved into Zombie2's position at %+v", zombie1.GetPosition())
		}
	}
}

func TestEnemyIdleMover_NoMovementWhenBlocked(t *testing.T) {
	rng := utils.NewRandomWithSeed(42)
	mover := NewEnemyIdleMover(rng)

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

// --- Тесты вспомогательных функций ---

func TestCardinalDirectionToOffset(t *testing.T) {
	tests := []struct {
		dir  entities.Direction
		want primitives.Point2D[int]
	}{
		{entities.DirectionForward, primitives.Point2D[int]{X: 0, Y: -1}},
		{entities.DirectionBack, primitives.Point2D[int]{X: 0, Y: 1}},
		{entities.DirectionLeft, primitives.Point2D[int]{X: -1, Y: 0}},
		{entities.DirectionRight, primitives.Point2D[int]{X: 1, Y: 0}},
		{entities.DirectionStop, primitives.Point2D[int]{X: 0, Y: 0}},
	}

	for _, tt := range tests {
		got := cardinalDirectionToOffset(tt.dir)
		if got != tt.want {
			t.Errorf("cardinalDirectionToOffset(%d) = %+v, want %+v", tt.dir, got, tt.want)
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
			t.Errorf("Roundtrip failed: dir=%d -> offset=%+v -> dir=%d", dir, offset, restored)
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
