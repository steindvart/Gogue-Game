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

// ApplyFogOfWar применяет туман войны к двухслойному полю.
// Для текущей зоны видимости - показываются оба слоя (окружение + объекты).
// Для ранее исследованных, но не видимых сейчас клеток - показывается только окружение,
// объекты (враги, предметы) скрываются.
func (f *FogOfWar) ApplyFogOfWar(
	rendered *common.RenderedField,
	playerPos primitives.Point2D[int],
	viewRadius int,
) *common.RenderedField {
	if rendered == nil || rendered.Height == 0 || rendered.Width == 0 {
		return rendered
	}

	result := common.NewRenderedField(rendered.Width, rendered.Height)

	// Вычисляем видимые клетки на основе слоя окружения (стены блокируют луч)
	visibleTiles := f.computeVisibleTiles(rendered.EnvironmentLayer, playerPos, viewRadius)

	for y := 0; y < rendered.Height; y++ {
		for x := 0; x < rendered.Width; x++ {
			pos := primitives.Point2D[int]{X: x, Y: y}
			envType := rendered.EnvironmentLayer[y][x]

			isCurrentlyVisible := visibleTiles[pos]

			if isCurrentlyVisible {
				// Игрок видит клетку прямо сейчас - показываем оба слоя
				result.EnvironmentLayer[y][x] = envType
				result.ObjectLayer[y][x] = rendered.ObjectLayer[y][x]

				// Запоминаем статическое окружение
				if f.isStaticTile(envType) {
					f.ExploredTiles[pos] = true
				}
			} else if f.ExploredTiles[pos] {
				// Клетка была исследована ранее - показываем ТОЛЬКО окружение
				// Объекты (враги, предметы) НЕ показываются
				result.EnvironmentLayer[y][x] = envType
				// result.ObjectLayer[y][x] остаётся EntityTypeNone
			}
			// Иначе: неисследованная клетка - оба слоя пустые (туман)
		}
	}

	return result
}

// computeVisibleTiles вычисляет все видимые клетки из позиции игрока используя ray casting.
// Принимает поле окружения (environment layer) для проверки блокировки лучей стенами.
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
	return HasLineOfSight(field, origin, target)
}

// isStaticTile проверяет, является ли тип клетки статичным (должен запоминаться)
func (f *FogOfWar) isStaticTile(entityType common.GameEntityType) bool {
	return entityType == common.WorldTypeWall ||
		entityType == common.WorldTypePassage ||
		entityType == common.WorldTypeDoor ||
		entityType == common.WorldTypePortal
}

// Reset сбрасывает память о посещённых областях (например, при переходе на новый уровень)
func (f *FogOfWar) Reset() {
	f.ExploredTiles = make(map[primitives.Point2D[int]]bool)
}
