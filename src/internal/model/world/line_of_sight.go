package world

import (
	"gogue/internal/common"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

// HasLineOfSight проверяет прямую видимость между двумя точками на карте
// с помощью алгоритма Брезенхэма (ray casting).
//
// Луч строится от origin до target. Если на пути встречается стена
// (WorldTypeWall) или пустое пространство (EntityTypeNone), видимость считается
// заблокированной. Проверяются все промежуточные точки — стартовая и конечная
// точки не участвуют в проверке на препятствия.
func HasLineOfSight(
	field [][]common.GameEntityType,
	origin, target primitives.Point2D[int],
) bool {
	h := len(field)
	if h == 0 {
		return false
	}
	w := len(field[0])

	rayPoints := BresenhamLine(origin, target)

	// Проверяем все промежуточные точки (кроме origin и target)
	for i := 1; i < len(rayPoints)-1; i++ {
		p := rayPoints[i]

		// Проверяем границы
		if p.X < 0 || p.X >= w || p.Y < 0 || p.Y >= h {
			return false
		}

		entityType := field[p.Y][p.X]
		if entityType == common.WorldTypeWall || entityType == common.EntityTypeNone {
			return false
		}
	}

	return true
}

// BresenhamLine строит линию между двумя точками по алгоритму Брезенхэма.
// Возвращает упорядоченный набор точек от from до to включительно.
//
// Wiki: https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
func BresenhamLine(from, to primitives.Point2D[int]) []primitives.Point2D[int] {
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
