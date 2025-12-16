package world

import (
	"gogue/internal/common"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
	"math"
)

// FogOfWar управляет видимостью игрока и запоминает ранее посещённые области
type FogOfWar struct {
	// exploredTiles хранит координаты клеток, которые были исследованы (стены, проходы, двери)
	// Эти клетки остаются видимыми даже когда игрок уходит
	exploredTiles map[primitives.Point2D[int]]bool

	// Размеры поля для проверки границ
	width  int
	height int
}

func NewFogOfWar(width, height int) *FogOfWar {
	return &FogOfWar{
		exploredTiles: make(map[primitives.Point2D[int]]bool),
		width:         width,
		height:        height,
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
					f.exploredTiles[pos] = true
				}
			} else if f.exploredTiles[pos] {
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
	maxX := utils.Min(f.width-1, origin.X+radius)
	minY := utils.Max(0, origin.Y-radius)
	maxY := utils.Min(f.height-1, origin.Y+radius)

	// Проверяем каждую клетку в области видимости
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			target := primitives.Point2D[int]{X: x, Y: y}

			if isInVisibleArea(origin, target, radius) {
				visible[target] = f.isVisible(field, origin, target)
			}
		}
	}

	return visible
}

func isInVisibleArea(origin, target primitives.Point2D[int], viewRadius int) bool {
	dx := float64(target.X - origin.X)
	dy := float64(target.Y - origin.Y)
	distance := math.Sqrt(dx*dx + dy*dy)

	return distance <= float64(viewRadius)
}

// isVisible проверяет, видна ли целевая точка из точки origin используя ray casting
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

	// Получаем все точки на луче от origin до target
	rayPoints := f.bresenhamLine(origin, target)

	// Проходим по всем точкам луча (кроме последней - самой цели)
	for i := 0; i < len(rayPoints)-1; i++ {
		point := rayPoints[i]

		// Пропускаем стартовую точку
		if point == origin {
			continue
		}

		// Проверяем границы
		if !f.isInBounds(point) {
			return false
		}

		// Стена и пустое пространство блокирует видимость
		entityType := field[point.Y][point.X]
		if entityType == common.WorldTypeWall || entityType == common.EntityTypeNone {
			return false
		}
	}

	// Луч не встретил препятствий - цель видна
	return true
}

// bresenhamLine реализует алгоритм Брезенхэма для построения линии между двумя точками
// Wiki: https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
func (f *FogOfWar) bresenhamLine(from, to primitives.Point2D[int]) []primitives.Point2D[int] {
	var points []primitives.Point2D[int]

	x0, y0 := from.X, from.Y
	x1, y1 := to.X, to.Y

	dx := utils.Abs(x1 - x0)
	dy := utils.Abs(y1 - y0)

	sx := 1
	if x0 > x1 {
		sx = -1
	}

	sy := 1
	if y0 > y1 {
		sy = -1
	}

	e1 := dx - dy

	x, y := x0, y0

	for {
		points = append(points, primitives.Point2D[int]{X: x, Y: y})

		if x == x1 && y == y1 {
			break
		}

		e2 := 2 * e1

		if e2 > -dy {
			e1 -= dy
			x += sx
		}

		if e2 < dx {
			e1 += dx
			y += sy
		}
	}

	return points
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
	return pos.X >= 0 && pos.X < f.width && pos.Y >= 0 && pos.Y < f.height
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
	f.exploredTiles = make(map[primitives.Point2D[int]]bool)
}
