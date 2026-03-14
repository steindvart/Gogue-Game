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
		Enemies:  []primitives.Positional2D[int]{},
		Player:   nil,
	}

	rf := renderer.RenderField(10, 10, level)

	if rf.Height != 10 {
		t.Errorf("Expected height 10, got %d", rf.Height)
	}
	if rf.Width != 10 {
		t.Errorf("Expected width 10, got %d", rf.Width)
	}

	// Проверяем, что оба слоя пустые
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			if rf.EnvironmentLayer[y][x] != common.EntityTypeNone {
				t.Errorf("Expected empty environment at (%d, %d), got %v", x, y, rf.EnvironmentLayer[y][x])
			}
			if rf.ObjectLayer[y][x] != common.EntityTypeNone {
				t.Errorf("Expected empty object at (%d, %d), got %v", x, y, rf.ObjectLayer[y][x])
			}
		}
	}
}

func TestDefaultFieldRenderer_RenderSingleRoom(t *testing.T) {
	renderer := NewDefaultFieldRenderer()

	room := Room{
		Box: &primitives.Box{
			Point: primitives.Point2D[int]{X: 2, Y: 2},
			Size:  primitives.Size2D[uint]{Width: 5, Height: 4},
		},
		Type: RoomTypeOrdinary,
	}

	level := &Level{
		Rooms:    []Room{room},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{},
		Enemies:  []primitives.Positional2D[int]{},
		Player:   nil,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: -1, Y: -1},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	rf := renderer.RenderField(10, 10, level)
	env := rf.EnvironmentLayer

	// Проверяем стены комнаты на слое окружения
	for x := 2; x <= 6; x++ {
		if env[2][x] != common.WorldTypeWall {
			t.Errorf("Expected wall at (%d, 2), got %v", x, env[2][x])
		}
	}

	for x := 2; x <= 6; x++ {
		if env[5][x] != common.WorldTypeWall {
			t.Errorf("Expected wall at (%d, 5), got %v", x, env[5][x])
		}
	}

	for y := 2; y <= 5; y++ {
		if env[y][2] != common.WorldTypeWall {
			t.Errorf("Expected wall at (2, %d), got %v", y, env[y][2])
		}
	}

	for y := 2; y <= 5; y++ {
		if env[y][6] != common.WorldTypeWall {
			t.Errorf("Expected wall at (6, %d), got %v", y, env[y][6])
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
		Enemies:  []primitives.Positional2D[int]{},
		Player:   nil,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: -1, Y: -1},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	rf := renderer.RenderField(10, 10, level)
	env := rf.EnvironmentLayer

	// Проверяем двери
	if env[3][3] != common.WorldTypeDoor {
		t.Errorf("Expected door at (3, 3), got %v", env[3][3])
	}
	if env[3][7] != common.WorldTypeDoor {
		t.Errorf("Expected door at (7, 3), got %v", env[3][7])
	}

	// Проверяем путь коридора
	for _, pt := range passage.Way {
		if env[pt.Y][pt.X] != common.WorldTypePassage {
			t.Errorf("Expected passage at (%d, %d), got %v", pt.X, pt.Y, env[pt.Y][pt.X])
		}
	}
}

func TestDefaultFieldRenderer_RenderPortal(t *testing.T) {
	renderer := NewDefaultFieldRenderer()

	room := Room{
		Box: &primitives.Box{
			Point: primitives.Point2D[int]{X: 1, Y: 1},
			Size:  primitives.Size2D[uint]{Width: 5, Height: 5},
		},
		Type: RoomTypeFinish,
	}

	level := &Level{
		Rooms:    []Room{room},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{},
		Enemies:  []primitives.Positional2D[int]{},
		Player:   nil,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: 3, Y: 3},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	rf := renderer.RenderField(10, 10, level)

	// Проверяем портал (портал - часть окружения)
	if rf.EnvironmentLayer[3][3] != common.WorldTypePortal {
		t.Errorf("Expected portal at (3, 3), got %v", rf.EnvironmentLayer[3][3])
	}
}

func TestDefaultFieldRenderer_RenderItems(t *testing.T) {
	renderer := NewDefaultFieldRenderer()

	food := items.Food{
		Type: items.FoodTypeBread,
		Item: &items.Item{
			Box: &primitives.Box{
				Point: primitives.Point2D[int]{X: 2, Y: 2},
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			},
		},
	}

	elixir := items.Elixir{
		Type: items.ElixirTypeStrength,
		Item: &items.Item{
			Box: &primitives.Box{
				Point: primitives.Point2D[int]{X: 3, Y: 3},
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			},
		},
	}

	scroll := items.Scroll{
		Type: items.ScrollTypeAgility,
		Item: &items.Item{
			Box: &primitives.Box{
				Point: primitives.Point2D[int]{X: 4, Y: 4},
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			},
		},
	}

	weapon := items.Weapon{
		Type: items.WeaponTypeSword,
		Item: &items.Item{
			Box: &primitives.Box{
				Point: primitives.Point2D[int]{X: 5, Y: 5},
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			},
		},
	}

	level := &Level{
		Rooms:    []Room{},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{&food, &elixir, &scroll, &weapon},
		Enemies:  []primitives.Positional2D[int]{},
		Player:   nil,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: -1, Y: -1},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	rf := renderer.RenderField(10, 10, level)
	obj := rf.ObjectLayer

	// Проверяем предметы (предметы - динамический слой)
	if obj[2][2] != common.Food {
		t.Errorf("Expected bread at (2, 2), got %v", obj[2][2])
	}
	if obj[3][3] != common.Elixir {
		t.Errorf("Expected strength elixir at (3, 3), got %v", obj[3][3])
	}
	if obj[4][4] != common.Scroll {
		t.Errorf("Expected agility scroll at (4, 4), got %v", obj[4][4])
	}
	if obj[5][5] != common.Weapon {
		t.Errorf("Expected weapon at (5, 5), got %v", obj[5][5])
	}
}

func TestDefaultFieldRenderer_RenderEnemies(t *testing.T) {
	renderer := NewDefaultFieldRenderer()

	zombie := entities.NewZombie(&primitives.Box{
		Point: primitives.Point2D[int]{X: 3, Y: 3},
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	})

	vampire := entities.NewVampire(&primitives.Box{
		Point: primitives.Point2D[int]{X: 5, Y: 5},
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	})

	level := &Level{
		Rooms:    []Room{},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{},
		Enemies:  []primitives.Positional2D[int]{zombie, vampire},
		Player:   nil,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: -1, Y: -1},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	rf := renderer.RenderField(10, 10, level)
	obj := rf.ObjectLayer

	// Проверяем врагов (враги - динамический слой)
	if obj[3][3] != common.EntityTypeZombie {
		t.Errorf("Expected zombie at (3, 3), got %v", obj[3][3])
	}
	if obj[5][5] != common.EntityTypeVampire {
		t.Errorf("Expected vampire at (5, 5), got %v", obj[5][5])
	}
}

func TestDefaultFieldRenderer_RenderPlayer(t *testing.T) {
	renderer := NewDefaultFieldRenderer()

	player := entities.NewPlayer(primitives.Box{
		Point: primitives.Point2D[int]{X: 5, Y: 5},
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	})

	level := &Level{
		Rooms:    []Room{},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{},
		Enemies:  []primitives.Positional2D[int]{},
		Player:   player,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: -1, Y: -1},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	rf := renderer.RenderField(10, 10, level)

	// Проверяем игрока (игрок - динамический слой)
	if rf.ObjectLayer[5][5] != common.EntityTypePlayer {
		t.Errorf("Expected player at (5, 5), got %v", rf.ObjectLayer[5][5])
	}
}

func TestDefaultFieldRenderer_PlayerOverlapsItem(t *testing.T) {
	// Игрок должен рисоваться поверх предметов
	renderer := NewDefaultFieldRenderer()

	food := items.Food{
		Type: items.FoodTypeBread,
		Item: &items.Item{
			Box: &primitives.Box{
				Point: primitives.Point2D[int]{X: 5, Y: 5},
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			},
		},
	}

	player := entities.NewPlayer(primitives.Box{
		Point: primitives.Point2D[int]{X: 5, Y: 5},
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	})

	level := &Level{
		Rooms:    []Room{},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{&food},
		Enemies:  []primitives.Positional2D[int]{},
		Player:   player,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: -1, Y: -1},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	rf := renderer.RenderField(10, 10, level)

	// Игрок должен быть поверх еды (оба в ObjectLayer, игрок перезаписывает)
	if rf.ObjectLayer[5][5] != common.EntityTypePlayer {
		t.Errorf("Expected player at (5, 5) overlapping food, got %v", rf.ObjectLayer[5][5])
	}
}

func TestDefaultFieldRenderer_BoundsChecking(t *testing.T) {
	// Тест проверяет, что объекты за границами карты не вызывают panic
	renderer := NewDefaultFieldRenderer()

	food := items.Food{
		Type: items.FoodTypeBread,
		Item: &items.Item{
			Box: &primitives.Box{
				Point: primitives.Point2D[int]{X: 15, Y: 15}, // За пределами карты 10x10
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			},
		},
	}

	level := &Level{
		Rooms:    []Room{},
		Passages: []Passage{},
		Items:    []primitives.Positional2D[int]{&food},
		Enemies:  []primitives.Positional2D[int]{},
		Player:   nil,
		FinishPortal: primitives.Box{
			Point: primitives.Point2D[int]{X: -1, Y: -1},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
	}

	// Не должно быть паники
	rf := renderer.RenderField(10, 10, level)

	if rf.Height != 10 || rf.Width != 10 {
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
			name:     "EnemyTypeZombie converts correctly",
			testFunc: func() common.GameEntityType { return convertEnemyToEntityType(&entities.Zombie{}) },
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
