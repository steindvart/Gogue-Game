package world

import (
	"errors"
	"gogue/internal/model/primitives"
)

var ErrPathNotFound = errors.New("path not found: target is unreachable")
var ErrInvalidTarget = errors.New("path not found: target position is invalid or not walkable")
var ErrInvalidStart = errors.New("path not found: start position is invalid or not walkable")
var ErrEmptyField = errors.New("path not found: field is empty")

// bfsNode - узел для очереди BFS: хранит позицию и ссылку на предыдущий узел
// для восстановления пути.
type bfsNode struct {
	pos    primitives.Point2D[int]
	parent *bfsNode
}

// FindPathBFS находит кратчайший путь от стартовой позиции start до целевой точки target
// на двумерной карте field. Перемещение допускается только по клеткам, значение которых
// присутствует в списке walkable. Поддерживает движение в 8 направлениях (включая диагонали).
//
// Возвращает срез последовательных координат (от первого шага после start до target,
// включая target). Если путь не найден - возвращает nil и ошибку.
//
// Алгоритм: BFS (Breadth-First Search).
// Обоснование выбора: все клетки имеют одинаковую "стоимость" перемещения (1 шаг),
// поэтому BFS гарантирует кратчайший путь при O(V + E) сложности.
func FindPathBFS(
	start primitives.Point2D[int],
	target primitives.Point2D[int],
	field [][]int,
	walkable []int,
	directions []primitives.Point2D[int],
) ([]primitives.Point2D[int], error) {
	// Валидация поля
	if len(field) == 0 || len(field[0]) == 0 {
		return nil, ErrEmptyField
	}

	height := len(field)
	width := len(field[0])

	// Построим set из допустимых значений для быстрой проверки O(1)
	walkableSet := make(map[int]struct{}, len(walkable))
	for _, w := range walkable {
		walkableSet[w] = struct{}{}
	}

	// Валидация стартовой позиции
	if !isInBounds(start, width, height) || !isWalkable(field, start, walkableSet) {
		return nil, ErrInvalidStart
	}

	// Валидация целевой позиции
	if !isInBounds(target, width, height) || !isWalkable(field, target, walkableSet) {
		return nil, ErrInvalidTarget
	}

	// Тривиальный случай: старт == цель
	if start == target {
		return []primitives.Point2D[int]{}, nil
	}

	visited := make([][]bool, height)
	for i := range visited {
		visited[i] = make([]bool, width)
	}

	startNode := &bfsNode{pos: start, parent: nil}
	visited[start.Y][start.X] = true

	queue := []*bfsNode{startNode}

	for len(queue) > 0 {
		current := queue[0]
		// Сдвигаем очередь на один элемент (как бы "удаляя" текущий узел)
		queue = queue[1:]

		for _, dir := range directions {
			next := primitives.Point2D[int]{
				X: current.pos.X + dir.X,
				Y: current.pos.Y + dir.Y,
			}

			if !isInBounds(next, width, height) {
				continue
			}
			if visited[next.Y][next.X] {
				continue
			}
			if !isWalkable(field, next, walkableSet) {
				continue
			}

			nextNode := &bfsNode{pos: next, parent: current}

			if next == target {
				return reconstructPath(nextNode), nil
			}

			visited[next.Y][next.X] = true
			queue = append(queue, nextNode)
		}
	}

	return nil, ErrPathNotFound
}

// reconstructPath восстанавливает путь от конечного узла до стартового,
// исключая стартовую позицию. Возвращает срез в правильном порядке: от первого шага до цели.
func reconstructPath(node *bfsNode) []primitives.Point2D[int] {
	var path []primitives.Point2D[int]

	for n := node; n != nil && n.parent != nil; n = n.parent {
		path = append(path, n.pos)
	}

	// Разворачиваем путь: сейчас он от цели к старту, а нужно от старта к цели
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

// isInBounds проверяет, что позиция находится в пределах карты.
func isInBounds(pos primitives.Point2D[int], width, height int) bool {
	return pos.X >= 0 && pos.Y >= 0 && pos.X < width && pos.Y < height
}

// isWalkable проверяет, что клетка на данной позиции является проходимой.
func isWalkable(field [][]int, pos primitives.Point2D[int], walkableSet map[int]struct{}) bool {
	_, ok := walkableSet[field[pos.Y][pos.X]]
	return ok
}
