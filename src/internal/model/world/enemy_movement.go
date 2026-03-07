package world

import (
	"gogue/internal/common"
	"gogue/internal/model/entities"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

// cardinalDirections - 4 кардинальных направления: N, S, W, E.
var cardinalDirections = []primitives.Point2D[int]{
	{X: 0, Y: -1}, // North
	{X: 0, Y: 1},  // South
	{X: -1, Y: 0}, // West
	{X: 1, Y: 0},  // East
}

// diagonalDirections - 4 диагональных направления.
var diagonalDirections = []primitives.Point2D[int]{
	{X: 1, Y: -1},  // NE
	{X: 1, Y: 1},   // SE
	{X: -1, Y: 1},  // SW
	{X: -1, Y: -1}, // NW
}

// EnemyIdleMover обрабатывает "бродячее" поведение врагов,
// когда они не преследуют игрока. Каждый тип врага имеет свой паттерн.
type EnemyIdleMover struct {
	random utils.Randomizer
}

func NewEnemyIdleMover(random utils.Randomizer) *EnemyIdleMover {
	return &EnemyIdleMover{random: random}
}

// MoveIdle выполняет один шаг "бродячего" движения для врага.
// Возвращает true, если враг был перемещён.
//
// Dispatch по конкретному типу:
//   - Zombie: случайный шаг на 1 клетку (N/S/W/E)
//   - Vampire: случайный шаг на 1 клетку (N/S/W/E)
//   - Ghost: телепортация в случайную свободную клетку комнаты
//   - переключение невидимости
//   - Ogre: случайный шаг на 2 клетки (N/S/W/E)
//   - SnakeMage: движение по диагонали с постоянной сменой стороны
func (m *EnemyIdleMover) MoveIdle(
	positionalEnemy primitives.Positional2D[int],
	fullField [][]common.GameEntityType,
	rooms []Room,
	otherEnemies []primitives.Positional2D[int],
) bool {
	switch enemy := positionalEnemy.(type) {
	case *entities.Zombie:
		return m.moveRandomCardinal(enemy.Enemy, fullField, otherEnemies, 1)
	case *entities.Vampire:
		return m.moveRandomCardinal(enemy.Enemy, fullField, otherEnemies, 1)
	case *entities.Ghost:
		return m.moveGhost(enemy, fullField, rooms, otherEnemies)
	case *entities.Ogre:
		return m.moveOgre(enemy, fullField, otherEnemies)
	case *entities.SnakeMage:
		return m.moveSnakeMage(enemy, fullField, otherEnemies)
	default:
		return false
	}
}

// moveRandomCardinal перемещает врага на случайную проходимую клетку
// из 4 кардинальных направлений. steps — количество клеток за один ход.
func (m *EnemyIdleMover) moveRandomCardinal(
	enemy *entities.Enemy,
	field [][]common.GameEntityType,
	otherEnemies []primitives.Positional2D[int],
	steps int,
) bool {
	pos := enemy.GetPosition()

	// Собираем все доступные конечные позиции
	candidates := m.getWalkableCandidates(pos, cardinalDirections, steps, field, otherEnemies, enemy)
	if len(candidates) == 0 {
		return false
	}

	target := candidates[m.random.Intn(len(candidates))]
	enemy.SetPosition(target)
	return true
}

// moveGhost телепортирует привидение в случайную свободную клетку в его комнате.
// Переключает невидимость с вероятностью ~30%.
func (m *EnemyIdleMover) moveGhost(
	ghost *entities.Ghost,
	field [][]common.GameEntityType,
	rooms []Room,
	otherEnemies []primitives.Positional2D[int],
) bool {
	ghostPos := ghost.GetPosition()

	// Переключение невидимости (~30% шанс стать невидимым/видимым каждый тик)
	const invisibilityChance = 0.3
	if m.random.Float64() < invisibilityChance {
		ghost.IsVisible = !ghost.IsVisible
	}

	// Ищем комнату, в которой находится привидение
	var ghostRoom *Room
	for i := range rooms {
		if isInRoom(ghostPos, rooms[i]) {
			ghostRoom = &rooms[i]
			break
		}
	}

	// Если привидение не в комнате (в проходе), двигаемся как обычный враг
	if ghostRoom == nil {
		return m.moveRandomCardinal(ghost.Enemy, field, otherEnemies, 1)
	}

	// Собираем все свободные клетки внутри комнаты
	freeTiles := m.getFreeTilesInRoom(*ghostRoom, field, otherEnemies, ghost.Enemy)
	if len(freeTiles) == 0 {
		return false
	}

	target := freeTiles[m.random.Intn(len(freeTiles))]
	ghost.SetPosition(target)
	return true
}

// moveOgre перемещает огра на 2 клетки в случайном кардинальном направлении.
// Если 2 шага недоступны — пытается сделать 1 шаг.
func (m *EnemyIdleMover) moveOgre(
	ogre *entities.Ogre,
	field [][]common.GameEntityType,
	otherEnemies []primitives.Positional2D[int],
) bool {
	// Сначала пробуем 2 шага
	if m.moveRandomCardinal(ogre.Enemy, field, otherEnemies, 2) {
		return true
	}
	// Фолбэк: 1 шаг
	return m.moveRandomCardinal(ogre.Enemy, field, otherEnemies, 1)
}

// moveSnakeMage перемещает Змея-мага по диагонали.
// Использует Direction врага для запоминания текущего диагонального направления.
// При каждом ходе пытается двигаться в текущем направлении;
// если заблокирован — меняет направление на случайное доступное.
func (m *EnemyIdleMover) moveSnakeMage(
	snake *entities.SnakeMage,
	field [][]common.GameEntityType,
	otherEnemies []primitives.Positional2D[int],
) bool {
	pos := snake.GetPosition()
	currentDir := snakeMageDirectionToOffset(snake.Direction)

	// Пытаемся продолжить движение в текущем направлении
	next := primitives.Point2D[int]{X: pos.X + currentDir.X, Y: pos.Y + currentDir.Y}
	if isWalkableTile(next, field) && !isOccupiedByOtherEnemy(next, otherEnemies, snake.Enemy) {
		snake.SetPosition(next)
		// Меняем направление для следующего хода
		snake.Direction = nextSnakeMageDirection(snake.Direction)
		return true
	}

	// Текущее направление заблокировано — выбираем случайное из доступных диагоналей
	candidates := m.getWalkableCandidates(pos, diagonalDirections, 1, field, otherEnemies, snake.Enemy)
	if len(candidates) == 0 {
		return false
	}

	target := candidates[m.random.Intn(len(candidates))]
	snake.SetPosition(target)

	// Устанавливаем новое направление на основе фактического движения
	snake.Direction = offsetToSnakeMageDirection(primitives.Point2D[int]{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	})
	// Меняем для следующего хода
	snake.Direction = nextSnakeMageDirection(snake.Direction)

	return true
}

// getWalkableCandidates возвращает все доступные конечные позиции
// при движении из pos в заданных направлениях на steps шагов.
// Проверяет, что ВСЕ промежуточные клетки проходимы (для steps > 1).
func (m *EnemyIdleMover) getWalkableCandidates(
	pos primitives.Point2D[int],
	directions []primitives.Point2D[int],
	steps int,
	field [][]common.GameEntityType,
	otherEnemies []primitives.Positional2D[int],
	self *entities.Enemy,
) []primitives.Point2D[int] {
	var candidates []primitives.Point2D[int]

	for _, dir := range directions {
		valid := true
		current := pos

		for s := 1; s <= steps; s++ {
			next := primitives.Point2D[int]{X: current.X + dir.X, Y: current.Y + dir.Y}

			if !isWalkableTile(next, field) || isOccupiedByOtherEnemy(next, otherEnemies, self) {
				valid = false
				break
			}
			current = next
		}

		if valid {
			candidates = append(candidates, current)
		}
	}

	return candidates
}

// getFreeTilesInRoom собирает все свободные клетки внутри комнаты (пол),
// исключая позиции других врагов и свою текущую позицию.
func (m *EnemyIdleMover) getFreeTilesInRoom(
	room Room,
	field [][]common.GameEntityType,
	otherEnemies []primitives.Positional2D[int],
	self *entities.Enemy,
) []primitives.Point2D[int] {
	var tiles []primitives.Point2D[int]

	selfPos := self.GetPosition()
	minX := room.Box.Point.X + 1
	minY := room.Box.Point.Y + 1
	maxX := minX + int(room.Box.Size.Width) - 3
	maxY := minY + int(room.Box.Size.Height) - 3

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			pos := primitives.Point2D[int]{X: x, Y: y}

			// Не телепортироваться на своё же место
			if pos == selfPos {
				continue
			}

			if !isWalkableTile(pos, field) {
				continue
			}

			if isOccupiedByOtherEnemy(pos, otherEnemies, self) {
				continue
			}

			tiles = append(tiles, pos)
		}
	}

	return tiles
}

// isWalkableTile проверяет, является ли клетка проходимой для idle-движения.
func isWalkableTile(pos primitives.Point2D[int], field [][]common.GameEntityType) bool {
	h := len(field)
	if h == 0 {
		return false
	}
	w := len(field[0])

	if pos.X < 0 || pos.X >= w || pos.Y < 0 || pos.Y >= h {
		return false
	}

	tile := field[pos.Y][pos.X]
	switch tile {
	case common.WorldTypeRoomFloor,
		common.WorldTypePassage,
		common.WorldTypeDoor,
		common.WorldTypePortal,
		common.Food,
		common.Elixir,
		common.Scroll,
		common.Weapon:
		return true
	default:
		return false
	}
}

// isOccupiedByOtherEnemy проверяет, занята ли позиция другим врагом (кроме self).
func isOccupiedByOtherEnemy(
	pos primitives.Point2D[int],
	enemies []primitives.Positional2D[int],
	self *entities.Enemy,
) bool {
	for _, other := range enemies {
		if provider, ok := other.(entities.EnemyProvider); ok {
			if provider.GetEnemy() == self {
				continue
			}
		}
		if other.GetPosition() == pos {
			return true
		}
	}
	return false
}

// --- Маппинг Direction - диагональные смещения для SnakeMage ---

// snakeMageDirectionToOffset конвертирует Direction врага в смещение по сетке.
// Используем только диагональные направления для Змея-мага.
func snakeMageDirectionToOffset(dir entities.Direction) primitives.Point2D[int] {
	switch dir {
	case entities.DirectionDiagonallyForwardRight:
		return primitives.Point2D[int]{X: 1, Y: -1} // NE
	case entities.DirectionDiagonallyBackRight:
		return primitives.Point2D[int]{X: 1, Y: 1} // SE
	case entities.DirectionDiagonallyBackLeft:
		return primitives.Point2D[int]{X: -1, Y: 1} // SW
	case entities.DirectionDiagonallyForwardLeft:
		return primitives.Point2D[int]{X: -1, Y: -1} // NW
	default:
		// По умолчанию — NE
		return primitives.Point2D[int]{X: 1, Y: -1}
	}
}

// offsetToSnakeMageDirection конвертирует смещение в Direction.
func offsetToSnakeMageDirection(offset primitives.Point2D[int]) entities.Direction {
	switch {
	case offset.X > 0 && offset.Y < 0:
		return entities.DirectionDiagonallyForwardRight // NE
	case offset.X > 0 && offset.Y > 0:
		return entities.DirectionDiagonallyBackRight // SE
	case offset.X < 0 && offset.Y > 0:
		return entities.DirectionDiagonallyBackLeft // SW
	case offset.X < 0 && offset.Y < 0:
		return entities.DirectionDiagonallyForwardLeft // NW
	default:
		return entities.DirectionDiagonallyForwardRight
	}
}

// nextSnakeMageDirection меняет диагональное направление на следующее (по часовой стрелке).
// NE -> SE -> SW -> NW -> NE ...
func nextSnakeMageDirection(dir entities.Direction) entities.Direction {
	switch dir {
	case entities.DirectionDiagonallyForwardRight:
		return entities.DirectionDiagonallyBackRight // NE -> SE
	case entities.DirectionDiagonallyBackRight:
		return entities.DirectionDiagonallyBackLeft // SE -> SW
	case entities.DirectionDiagonallyBackLeft:
		return entities.DirectionDiagonallyForwardLeft // SW -> NW
	case entities.DirectionDiagonallyForwardLeft:
		return entities.DirectionDiagonallyForwardRight // NW -> NE
	default:
		return entities.DirectionDiagonallyBackRight
	}
}
