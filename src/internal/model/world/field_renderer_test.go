package world

import (
	"gogue/internal/common"
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"testing"
)

func TestDefaultFieldRenderer_RenderMap_EmptyLevel(t *testing.T) {
	renderer := NewDefaultFieldRenderer()
	level := &Level{
		Rooms:    []Room{},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{},
		Enemies:  []entities.Enemy{},
		Player:   nil,
	}

	field := renderer.RenderField(10, 10, level)

	if len(field) != 10 {
		t.Errorf("Expected height 10, got %d", len(field))
	}
	if len(field[0]) != 10 {
		t.Errorf("Expected width 10, got %d", len(field[0]))
	}

	// Проверяем, что все ячейки пустые
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			if field[y][x] != common.GameEntityType(0) {
				t.Errorf("Expected empty cell at (%d, %d), got %v", x, y, field[y][x])
			}
		}
	}
}

func TestDefaultFieldRenderer_RenderSingleRoom(t *testing.T) {
	renderer := NewDefaultFieldRenderer()

	room := Room{
		Box: primitives.Box{
			Point: primitives.Point2D[int]{X: 2, Y: 2},
			Size:  primitives.Size2D[uint]{Width: 5, Height: 4},
		},
		Type: RoomTypeOrdinary,
	}

	level := &Level{
		Rooms:    []Room{room},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{},
		Enemies:  []entities.Enemy{},
		Player:   nil,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: -1, Y: -1},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	field := renderer.RenderField(10, 10, level)

	// Проверяем стены комнаты
	// Верхняя стена (y=2, x=2..6)
	for x := 2; x <= 6; x++ {
		if field[2][x] != common.WorldTypeWall {
			t.Errorf("Expected wall at (%d, 2), got %v", x, field[2][x])
		}
	}

	// Нижняя стена (y=5, x=2..6)
	for x := 2; x <= 6; x++ {
		if field[5][x] != common.WorldTypeWall {
			t.Errorf("Expected wall at (%d, 5), got %v", x, field[5][x])
		}
	}

	// Левая стена (x=2, y=2..5)
	for y := 2; y <= 5; y++ {
		if field[y][2] != common.WorldTypeWall {
			t.Errorf("Expected wall at (2, %d), got %v", y, field[y][2])
		}
	}

	// Правая стена (x=6, y=2..5)
	for y := 2; y <= 5; y++ {
		if field[y][6] != common.WorldTypeWall {
			t.Errorf("Expected wall at (6, %d), got %v", y, field[y][6])
		}
	}
}

func TestDefaultFieldRenderer_RenderPassage(t *testing.T) {
	renderer := NewDefaultFieldRenderer()

	passage := Passage{
		DoorOne: primitives.Point2D[int]{X: 3, Y: 3},
		DoorTwo: primitives.Point2D[int]{X: 7, Y: 3},
		Way: []primitives.Point2D[int]{
			{X: 4, Y: 3},
			{X: 5, Y: 3},
			{X: 6, Y: 3},
		},
	}

	level := &Level{
		Rooms:    []Room{},
		Passages: []Passage{passage},
		Items:    []primitives.Positional2D[int]{},
		Enemies:  []entities.Enemy{},
		Player:   nil,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: -1, Y: -1},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	field := renderer.RenderField(10, 10, level)

	// Проверяем двери
	if field[3][3] != common.WorldTypeDoor {
		t.Errorf("Expected door at (3, 3), got %v", field[3][3])
	}
	if field[3][7] != common.WorldTypeDoor {
		t.Errorf("Expected door at (7, 3), got %v", field[3][7])
	}

	// Проверяем путь коридора
	for _, pt := range passage.Way {
		if field[pt.Y][pt.X] != common.WorldTypePassage {
			t.Errorf("Expected passage at (%d, %d), got %v", pt.X, pt.Y, field[pt.Y][pt.X])
		}
	}
}

func TestDefaultFieldRenderer_RenderPortal(t *testing.T) {
	renderer := NewDefaultFieldRenderer()

	room := Room{
		Box: primitives.Box{
			Point: primitives.Point2D[int]{X: 1, Y: 1},
			Size:  primitives.Size2D[uint]{Width: 5, Height: 5},
		},
		Type: RoomTypeFinish,
	}

	level := &Level{
		Rooms:    []Room{room},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{},
		Enemies:  []entities.Enemy{},
		Player:   nil,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: 3, Y: 3},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	field := renderer.RenderField(10, 10, level)

	// Проверяем портал
	if field[3][3] != common.WorldTypePortal {
		t.Errorf("Expected portal at (3, 3), got %v", field[3][3])
	}
}

func TestDefaultFieldRenderer_RenderItems(t *testing.T) {
	renderer := NewDefaultFieldRenderer()

	food := items.Food{
		Type: items.FoodTypeBread,
		Item: &items.Item{
			Box: primitives.Box{
				Point: primitives.Point2D[int]{X: 2, Y: 2},
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			},
		},
	}

	elixir := items.Elixir{
		Type: items.ElixirTypeStrength,
		Item: &items.Item{
			Box: primitives.Box{
				Point: primitives.Point2D[int]{X: 3, Y: 3},
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			},
		},
	}

	scroll := items.Scroll{
		Type: items.ScrollTypeAgility,
		Item: &items.Item{
			Box: primitives.Box{
				Point: primitives.Point2D[int]{X: 4, Y: 4},
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			},
		},
	}

	weapon := items.Weapon{
		Type: items.WeaponTypeSword,
		Item: &items.Item{
			Box: primitives.Box{
				Point: primitives.Point2D[int]{X: 5, Y: 5},
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			},
		},
	}

	level := &Level{
		Rooms:    []Room{},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{&food, &elixir, &scroll, &weapon},
		Enemies:  []entities.Enemy{},
		Player:   nil,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: -1, Y: -1},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	field := renderer.RenderField(10, 10, level)

	// Проверяем предметы
	if field[2][2] != common.FoodTypeBread {
		t.Errorf("Expected bread at (2, 2), got %v", field[2][2])
	}
	if field[3][3] != common.ElixirTypeStrength {
		t.Errorf("Expected strength elixir at (3, 3), got %v", field[3][3])
	}
	if field[4][4] != common.ScrollTypeAgility {
		t.Errorf("Expected agility scroll at (4, 4), got %v", field[4][4])
	}
	if field[5][5] != common.Weapon {
		t.Errorf("Expected weapon at (5, 5), got %v", field[5][5])
	}
}

func TestDefaultFieldRenderer_RenderEnemies(t *testing.T) {
	renderer := NewDefaultFieldRenderer()

	zombie := entities.Enemy{
		Type: entities.EnemyTypeZombie,
		Character: &entities.Character{
			Box: &primitives.Box{
				Point: primitives.Point2D[int]{X: 3, Y: 3},
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			},
		},
	}

	vampire := entities.Enemy{
		Type: entities.EnemyTypeVampire,
		Character: &entities.Character{
			Box: &primitives.Box{
				Point: primitives.Point2D[int]{X: 5, Y: 5},
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			},
		},
	}

	level := &Level{
		Rooms:    []Room{},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{},
		Enemies:  []entities.Enemy{zombie, vampire},
		Player:   nil,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: -1, Y: -1},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	field := renderer.RenderField(10, 10, level)

	// Проверяем врагов
	if field[3][3] != common.EntityTypeZombie {
		t.Errorf("Expected zombie at (3, 3), got %v", field[3][3])
	}
	if field[5][5] != common.EntityTypeVampire {
		t.Errorf("Expected vampire at (5, 5), got %v", field[5][5])
	}
}

func TestDefaultFieldRenderer_RenderPlayer(t *testing.T) {
	renderer := NewDefaultFieldRenderer()

	player := entities.NewPlayer(&primitives.Box{
		Point: primitives.Point2D[int]{X: 5, Y: 5},
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	})

	level := &Level{
		Rooms:    []Room{},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{},
		Enemies:  []entities.Enemy{},
		Player:   player,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: -1, Y: -1},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	field := renderer.RenderField(10, 10, level)

	// Проверяем игрока
	if field[5][5] != common.EntityTypePlayer {
		t.Errorf("Expected player at (5, 5), got %v", field[5][5])
	}
}

func TestDefaultFieldRenderer_PlayerOverlapsItem(t *testing.T) {
	// Игрок должен рисоваться поверх предметов
	renderer := NewDefaultFieldRenderer()

	food := items.Food{
		Type: items.FoodTypeBread,
		Item: &items.Item{
			Box: primitives.Box{
				Point: primitives.Point2D[int]{X: 5, Y: 5},
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			},
		},
	}

	player := entities.NewPlayer(&primitives.Box{
		Point: primitives.Point2D[int]{X: 5, Y: 5},
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	})

	level := &Level{
		Rooms:    []Room{},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{&food},
		Enemies:  []entities.Enemy{},
		Player:   player,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: -1, Y: -1},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	field := renderer.RenderField(10, 10, level)

	// Игрок должен быть поверх еды
	if field[5][5] != common.EntityTypePlayer {
		t.Errorf("Expected player at (5, 5) overlapping food, got %v", field[5][5])
	}
}

func TestDefaultFieldRenderer_BoundsChecking(t *testing.T) {
	// Тест проверяет, что объекты за границами карты не вызывают panic
	renderer := NewDefaultFieldRenderer()

	food := items.Food{
		Type: items.FoodTypeBread,
		Item: &items.Item{
			Box: primitives.Box{
				Point: primitives.Point2D[int]{X: 15, Y: 15}, // За пределами карты 10x10
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			},
		},
	}

	level := &Level{
		Rooms:    []Room{},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{&food},
		Enemies:  []entities.Enemy{},
		Player:   nil,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: -1, Y: -1},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	// Не должно быть паники
	field := renderer.RenderField(10, 10, level)

	if len(field) != 10 || len(field[0]) != 10 {
		t.Error("Field size changed after rendering out-of-bounds object")
	}
}

func TestConversionFunctions(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() common.GameEntityType
		expected common.GameEntityType
	}{
		{
			name:     "FoodTypeBread converts correctly",
			testFunc: func() common.GameEntityType { return convertFoodToEntityType(items.FoodTypeBread) },
			expected: common.FoodTypeBread,
		},
		{
			name:     "ElixirTypeStrength converts correctly",
			testFunc: func() common.GameEntityType { return convertElixirToEntityType(items.ElixirTypeStrength) },
			expected: common.ElixirTypeStrength,
		},
		{
			name:     "ScrollTypeAgility converts correctly",
			testFunc: func() common.GameEntityType { return convertScrollToEntityType(items.ScrollTypeAgility) },
			expected: common.ScrollTypeAgility,
		},
		{
			name:     "EnemyTypeZombie converts correctly",
			testFunc: func() common.GameEntityType { return convertEnemyToEntityType(entities.EnemyTypeZombie) },
			expected: common.EntityTypeZombie,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
