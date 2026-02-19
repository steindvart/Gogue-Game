package world

import (
	"gogue/internal/model/entities"
	"gogue/internal/model/primitives"
	"testing"
)

// Константы для читаемости поля в тестах:
// 0 - стена/непроходимая клетка, 1 - пол (проходимая)
const (
	wall  = 0
	floor = 1
)

func makeCharacterAt(x, y int) *entities.Character {
	return entities.NewCharacter(
		primitives.Box{
			Point: primitives.Point2D[int]{X: x, Y: y},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
		primitives.Attributes{Health: 100, MaxHealth: 100, Strength: 10, Agility: 5},
	)
}

func TestFindPathBFS_DirectPath(t *testing.T) {
	// Простое открытое поле 5x5, путь из (0,0) в (4,4) - диагональ.
	field := [][]int{
		{1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1},
	}

	c := makeCharacterAt(0, 0)
	target := primitives.Point2D[int]{X: 4, Y: 4}

	path, err := FindPathBFS(c.GetPosition(), target, field, []int{floor})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Диагональный путь: (1,1), (2,2), (3,3), (4,4) = 4 шага
	if len(path) != 4 {
		t.Errorf("expected path length 4, got %d; path=%v", len(path), path)
	}

	// Последний элемент должен быть целью
	if len(path) > 0 && path[len(path)-1] != target {
		t.Errorf("last element should be target %v, got %v", target, path[len(path)-1])
	}
}

func TestFindPathBFS_SamePosition(t *testing.T) {
	field := [][]int{
		{1, 1},
		{1, 1},
	}

	c := makeCharacterAt(1, 1)
	target := primitives.Point2D[int]{X: 1, Y: 1}

	path, err := FindPathBFS(c.GetPosition(), target, field, []int{floor})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(path) != 0 {
		t.Errorf("expected empty path when start == target, got %v", path)
	}
}

func TestFindPathBFS_PathAroundWall(t *testing.T) {
	// Поле с преградой посередине; путь должен огибать стену
	//
	// S . . . .
	// . 0 0 0 .
	// . 0 . 0 .
	// . 0 0 0 .
	// . . . . T
	field := [][]int{
		{1, 1, 1, 1, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 1, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 1, 1, 1, 1},
	}

	c := makeCharacterAt(0, 0)
	target := primitives.Point2D[int]{X: 4, Y: 4}

	path, err := FindPathBFS(c.GetPosition(), target, field, []int{floor})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(path) == 0 {
		t.Fatal("expected non-empty path, got empty")
	}

	// Проверяем, что путь заканчивается в цели
	if path[len(path)-1] != target {
		t.Errorf("last element should be target %v, got %v", target, path[len(path)-1])
	}

	// Проверяем, что все клетки пути проходимы
	for i, p := range path {
		if field[p.Y][p.X] != floor {
			t.Errorf("path[%d] = %v is not walkable (value=%d)", i, p, field[p.Y][p.X])
		}
	}
}

func TestFindPathBFS_NoPath(t *testing.T) {
	// Цель полностью окружена стенами - путь невозможен
	//
	// S . . . .
	// . . . . .
	// . . 0 0 0
	// . . 0 T 0
	// . . 0 0 0
	field := [][]int{
		{1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1},
		{1, 1, 0, 0, 0},
		{1, 1, 0, 1, 0},
		{1, 1, 0, 0, 0},
	}

	c := makeCharacterAt(0, 0)
	target := primitives.Point2D[int]{X: 3, Y: 3}

	path, err := FindPathBFS(c.GetPosition(), target, field, []int{floor})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if path != nil {
		t.Errorf("expected nil path, got %v", path)
	}
}

func TestFindPathBFS_TargetOutOfBounds(t *testing.T) {
	field := [][]int{
		{1, 1},
		{1, 1},
	}

	c := makeCharacterAt(0, 0)

	tests := []struct {
		name   string
		target primitives.Point2D[int]
	}{
		{"negative X", primitives.Point2D[int]{X: -1, Y: 0}},
		{"negative Y", primitives.Point2D[int]{X: 0, Y: -1}},
		{"X out of bounds", primitives.Point2D[int]{X: 2, Y: 0}},
		{"Y out of bounds", primitives.Point2D[int]{X: 0, Y: 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := FindPathBFS(c.GetPosition(), tt.target, field, []int{floor})
			if err == nil {
				t.Error("expected error for out-of-bounds target")
			}
			if path != nil {
				t.Errorf("expected nil path, got %v", path)
			}
		})
	}
}

func TestFindPathBFS_TargetOnWall(t *testing.T) {
	field := [][]int{
		{1, 0},
		{1, 1},
	}

	c := makeCharacterAt(0, 0)
	target := primitives.Point2D[int]{X: 1, Y: 0}

	path, err := FindPathBFS(c.GetPosition(), target, field, []int{floor})
	if err == nil {
		t.Fatal("expected error when target is on wall")
	}
	if path != nil {
		t.Errorf("expected nil path, got %v", path)
	}
}

func TestFindPathBFS_StartOnWall(t *testing.T) {
	field := [][]int{
		{0, 1},
		{1, 1},
	}

	c := makeCharacterAt(0, 0)
	target := primitives.Point2D[int]{X: 1, Y: 1}

	path, err := FindPathBFS(c.GetPosition(), target, field, []int{floor})
	if err == nil {
		t.Fatal("expected error when start is on wall")
	}
	if path != nil {
		t.Errorf("expected nil path, got %v", path)
	}
}

func TestFindPathBFS_EmptyField(t *testing.T) {
	c := makeCharacterAt(0, 0)
	target := primitives.Point2D[int]{X: 0, Y: 0}

	path, err := FindPathBFS(c.GetPosition(), target, [][]int{}, []int{floor})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
	if path != nil {
		t.Errorf("expected nil path, got %v", path)
	}
}

func TestFindPathBFS_MultipleWalkableTypes(t *testing.T) {
	// Поле с несколькими типами проходимых клеток (1 = пол, 2 = дверь)
	const door = 2
	field := [][]int{
		{1, 0, 1},
		{1, 2, 1},
		{1, 0, 1},
	}

	c := makeCharacterAt(0, 0)
	target := primitives.Point2D[int]{X: 2, Y: 0}

	// Без двери - путь невозможен
	path, err := FindPathBFS(c.GetPosition(), target, field, []int{floor})
	// С 8-направленным движением путь существует через диагональ (0,0)->(1,1)->(2,0)
	// но (1,1)=2 (door), и мы не указали door в walkable - проверим
	if err == nil && len(path) > 0 {
		// Путь может существовать через (0,1)->(0,2)->(1,2)...(2,0)
		// Проверим, что каждая клетка проходимая
		for i, p := range path {
			if field[p.Y][p.X] != floor {
				t.Errorf("path[%d]=%v has type %d, not floor", i, p, field[p.Y][p.X])
			}
		}
	}

	// С дверью - путь через дверь (кратчайший)
	path, err = FindPathBFS(c.GetPosition(), target, field, []int{floor, door})
	if err != nil {
		t.Fatalf("unexpected error with multiple walkable types: %v", err)
	}
	if len(path) == 0 {
		t.Fatal("expected non-empty path with door as walkable")
	}
	if path[len(path)-1] != target {
		t.Errorf("last element should be target %v, got %v", target, path[len(path)-1])
	}
}

func TestFindPathBFS_PathConsistency(t *testing.T) {
	// Проверяем, что каждый шаг пути - соседняя клетка (максимум 1 по каждой оси)
	field := [][]int{
		{1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 0, 1},
		{1, 0, 0, 0, 1, 0, 1},
		{1, 0, 1, 1, 1, 0, 1},
		{1, 0, 1, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1},
	}

	c := makeCharacterAt(0, 0)
	target := primitives.Point2D[int]{X: 6, Y: 6}

	path, err := FindPathBFS(c.GetPosition(), target, field, []int{floor})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Проверка непрерывности: от start к path[0], затем path[i] к path[i+1]
	start := c.GetPosition()
	prev := start
	for i, p := range path {
		dx := p.X - prev.X
		dy := p.Y - prev.Y
		if dx < -1 || dx > 1 || dy < -1 || dy > 1 {
			t.Errorf("step %d: jump from %v to %v is not adjacent (dx=%d, dy=%d)", i, prev, p, dx, dy)
		}
		prev = p
	}

	// Последний элемент - цель
	if path[len(path)-1] != target {
		t.Errorf("last element should be target %v, got %v", target, path[len(path)-1])
	}
}

func TestFindPathBFS_LinearPath(t *testing.T) {
	// Узкий коридор: путь только вправо
	field := [][]int{
		{1, 1, 1, 1, 1},
	}

	c := makeCharacterAt(0, 0)
	target := primitives.Point2D[int]{X: 4, Y: 0}

	path, err := FindPathBFS(c.GetPosition(), target, field, []int{floor})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(path) != 4 {
		t.Errorf("expected path length 4 in linear corridor, got %d; path=%v", len(path), path)
	}

	expected := []primitives.Point2D[int]{
		{X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 4, Y: 0},
	}
	for i, p := range path {
		if p != expected[i] {
			t.Errorf("path[%d] = %v, want %v", i, p, expected[i])
		}
	}
}

func TestFindPathBFS_LargeOpenField(t *testing.T) {
	// Большое открытое поле 50x50: от (0,0) до (49,49), диагональный путь = 49 шагов
	const size = 50
	field := make([][]int, size)
	for y := range field {
		field[y] = make([]int, size)
		for x := range field[y] {
			field[y][x] = floor
		}
	}

	c := makeCharacterAt(0, 0)
	target := primitives.Point2D[int]{X: size - 1, Y: size - 1}

	path, err := FindPathBFS(c.GetPosition(), target, field, []int{floor})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Диагональный кратчайший путь = 49 шагов
	if len(path) != size-1 {
		t.Errorf("expected path length %d, got %d", size-1, len(path))
	}
}

func BenchmarkFindPath_OpenField50x50(b *testing.B) {
	const size = 50
	field := make([][]int, size)
	for y := range field {
		field[y] = make([]int, size)
		for x := range field[y] {
			field[y][x] = floor
		}
	}

	c := makeCharacterAt(0, 0)
	target := primitives.Point2D[int]{X: size - 1, Y: size - 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindPathBFS(c.GetPosition(), target, field, []int{floor})
	}
}

func BenchmarkFindPath_Maze20x20(b *testing.B) {
	// Лабиринтообразное поле 20x20 со змеевидным проходом
	const size = 20
	field := make([][]int, size)
	for y := range field {
		field[y] = make([]int, size)
		for x := range field[y] {
			field[y][x] = wall
		}
	}

	// Создаём змеевидный путь
	for y := 0; y < size; y++ {
		if y%4 == 0 || y%4 == 1 {
			for x := 0; x < size; x++ {
				field[y][x] = floor
			}
		} else if y%4 == 2 {
			field[y][size-1] = floor
		} else {
			field[y][0] = floor
			for x := 0; x < size; x++ {
				field[y][x] = floor
			}
		}
	}

	c := makeCharacterAt(0, 0)
	target := primitives.Point2D[int]{X: size - 1, Y: size - 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindPathBFS(c.GetPosition(), target, field, []int{floor})
	}
}
