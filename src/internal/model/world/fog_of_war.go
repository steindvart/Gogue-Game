package world

import (
	"gogue/internal/common"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

// FogOfWar управляет видимостью игрока и запоминает ранее посещённые области
type FogOfWar struct {
	// ExploredTiles хранит координаты клеток, которые были исследованы (стены, проходы, двери)
	// Эти клетки остаются видимыми даже когда игрок уходит
	ExploredTiles map[primitives.Point2D[int]]bool

	// Размеры поля для проверки границ
	Width  int
	Height int
}

func NewFogOfWar(width, height int) *FogOfWar {
	return &FogOfWar{
		ExploredTiles: make(map[primitives.Point2D[int]]bool),
		Width:         width,
		Height:        height,
	}
}

// ApplyFogOfWar применяет туман войны к полному полю, возвращая отфильтрованное поле
// с учётом радиуса видимости игрока и запомненных ранее областей
func (f *FogOfWar) ApplyFogOfWar(
	fullField [][]common.GameEntityType,
	playerPos primitives.Point2D[int],
	viewRadius int,
) [][]common.GameEntityType {
	if len(fullField) == 0 {
		return fullField
	}

	height := len(fullField)
	width := len(fullField[0])

	filteredField := utils.CreateEmpty2DSlice[common.GameEntityType](height, width)
	visibleTiles := f.computeVisibleTiles(fullField, playerPos, viewRadius)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pos := primitives.Point2D[int]{X: x, Y: y}
			entityType := fullField[y][x]

			isCurrentlyVisible := visibleTiles[pos]

			if isCurrentlyVisible {
				filteredField[y][x] = entityType

				if f.isStaticTile(entityType) {
					f.ExploredTiles[pos] = true
				}
			} else if f.ExploredTiles[pos] {
				filteredField[y][x] = entityType
			}
		}
	}

	return filteredField
}

// computeVisibleTiles вычисляет все видимые клетки из позиции игрока используя ray casting
func (f *FogOfWar) computeVisibleTiles(
	field [][]common.GameEntityType,
	origin primitives.Point2D[int],
	radius int,
) map[primitives.Point2D[int]]bool {
	visible := make(map[primitives.Point2D[int]]bool)

	// Смотрящий всегда видит свою позицию
	visible[origin] = true

	// Определяем границы области видимости (квадрат вокруг смотрящего)
	minX := utils.Max(0, origin.X-radius)
	maxX := utils.Min(f.Width-1, origin.X+radius)
	minY := utils.Max(0, origin.Y-radius)
	maxY := utils.Min(f.Height-1, origin.Y+radius)

	// Проверяем каждую клетку в области видимости
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			target := primitives.Point2D[int]{X: x, Y: y}

			if isInVisibleAreaCircle(origin, target, radius) {
				visible[target] = f.isVisible(field, origin, target)
			}
		}
	}

	return visible
}

// isInVisibleAreaCircle проверяет, находится ли целевая точка в области видимости
// Использует октагональную метрику (круг) для более естественной формы на grid-сетке
func isInVisibleAreaCircle(origin, target primitives.Point2D[int], viewRadius int) bool {
	dx := utils.Abs(target.X - origin.X)
	dy := utils.Abs(target.Y - origin.Y)

	minDist := utils.Min(dx, dy)
	maxDist := utils.Max(dx, dy)
	octagonalDistance := maxDist + minDist/2

	return octagonalDistance <= viewRadius
}

// isVisible проверяет, видна ли целевая точка из точки origin используя ray casting.
// Делегирует проверку пакетной функции HasLineOfSight (line_of_sight.go).
func (f *FogOfWar) isVisible(
	field [][]common.GameEntityType,
	origin, target primitives.Point2D[int],
) bool {
	// @todo - ввести в реализацию или удалить (пока больше склоняюсь ко второму варианту)
	// Краевой случай - когда смотрящий находится в проходе, который расположен вплотную к стене
	// Стену изнутри прохода не должно быть видно (или должно?)
	// if passage := f.findPassage(origin, passages); passage != nil {
	// 	if !f.isInBounds(target) {
	// 		return false
	// 	}

	// 	if field[target.Y][target.X] == common.WorldTypeWall {
	// 		return false
	// 	}
	// }

	return HasLineOfSight(field, origin, target)
}

// isStaticTile проверяет, является ли тип клетки статичным (должен запоминаться)
func (f *FogOfWar) isStaticTile(entityType common.GameEntityType) bool {
	return entityType == common.WorldTypeWall ||
		entityType == common.WorldTypePassage ||
		entityType == common.WorldTypeDoor ||
		entityType == common.WorldTypePortal
}

// isInBounds проверяет, находится ли точка в границах текущего поля
func (f *FogOfWar) isInBounds(pos primitives.Point2D[int]) bool {
	return pos.X >= 0 && pos.X < f.Width && pos.Y >= 0 && pos.Y < f.Height
}

// @todo - ввести в реализацию или удалить (пока больше склоняюсь ко второму варианту)
// findPlayerPassage определяет, находится ли игрок в каком-либо проходе
// func (f *FogOfWar) findPassage(pos primitives.Point2D[int], passages []Passage) *Passage {
// 	for i := range passages {
// 		if f.isInPassage(pos, passages[i]) {
// 			return &passages[i]
// 		}
// 	}
// 	return nil
// }

// // isInPassage проверяет, находится ли позиция внутри прохода
// func (f *FogOfWar) isInPassage(pos primitives.Point2D[int], passage Passage) bool {
// 	// Исключение - если игрок в двери, то он не в проходе
// 	if pos == passage.DoorOne || pos == passage.DoorTwo {
// 		return false
// 	}

// 	// Проверяем путь прохода
// 	for _, wayPoint := range passage.Way {
// 		if pos == wayPoint {
// 			return true
// 		}
// 	}

// 	return false
// }

// Reset сбрасывает память о посещённых областях (например, при переходе на новый уровень)
func (f *FogOfWar) Reset() {
	f.ExploredTiles = make(map[primitives.Point2D[int]]bool)
}
