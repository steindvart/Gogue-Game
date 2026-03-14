package world

import (
	"gogue/internal/common"
	"gogue/internal/model/primitives"
	"testing"
)

// TestIsVisible проверяет алгоритм ray casting для определения видимости
func TestIsVisible(t *testing.T) {
	tests := []struct {
		name    string
		field   [][]common.GameEntityType
		origin  primitives.Point2D[int]
		target  primitives.Point2D[int]
		want    bool
		comment string
	}{
		{
			name: "direct line of sight",
			field: [][]common.GameEntityType{
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
			},
			origin: primitives.Point2D[int]{X: 0, Y: 1},
			target: primitives.Point2D[int]{X: 4, Y: 1},
			want:   true,
		},
		{
			name: "wall blocking view",
			field: [][]common.GameEntityType{
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeWall, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
			},
			origin: primitives.Point2D[int]{X: 0, Y: 1},
			target: primitives.Point2D[int]{X: 4, Y: 1},
			want:   false,
		},
		{
			name: "none type blocking view",
			field: [][]common.GameEntityType{
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.EntityTypeNone, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
			},
			origin: primitives.Point2D[int]{X: 0, Y: 1},
			target: primitives.Point2D[int]{X: 4, Y: 1},
			want:   false,
		},
		{
			name: "diagonal view through passage",
			field: [][]common.GameEntityType{
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
				{common.WorldTypeRoomFloor, common.WorldTypePassage, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypePassage, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
			},
			origin: primitives.Point2D[int]{X: 0, Y: 0},
			target: primitives.Point2D[int]{X: 3, Y: 3},
			want:   true,
		},
		{
			name: "same position",
			field: [][]common.GameEntityType{
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
				{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
			},
			origin: primitives.Point2D[int]{X: 1, Y: 1},
			target: primitives.Point2D[int]{X: 1, Y: 1},
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			height := len(tt.field)
			width := 0
			if height > 0 {
				width = len(tt.field[0])
			}

			fog := NewFogOfWar(width, height)
			got := fog.isVisible(tt.field, tt.origin, tt.target)

			if got != tt.want {
				t.Errorf("isVisible() = %v, want %v (%s)", got, tt.want, tt.comment)
			}
		})
	}
}

// TestApplyFogOfWar проверяет применение тумана войны
func TestApplyFogOfWar(t *testing.T) {
	t.Run("player sees items in view radius", func(t *testing.T) {
		// Создаём двухслойное поле 5x5
		rf := common.NewRenderedField(5, 5)

		// Окружение - всё пол
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				rf.EnvironmentLayer[y][x] = common.WorldTypeRoomFloor
			}
		}

		// Объекты
		rf.ObjectLayer[1][1] = common.Food
		rf.ObjectLayer[2][2] = common.EntityTypePlayer
		rf.ObjectLayer[3][3] = common.EntityTypeZombie

		fog := NewFogOfWar(5, 5)
		playerPos := primitives.Point2D[int]{X: 2, Y: 2}
		viewRadius := 2

		result := fog.ApplyFogOfWar(rf, playerPos, viewRadius)

		// Игрок должен видеть себя
		if result.ObjectLayer[2][2] != common.EntityTypePlayer {
			t.Errorf("Player not visible at own position")
		}

		// Еда в радиусе видимости должна быть видна
		if result.ObjectLayer[1][1] != common.Food {
			t.Errorf("Food within view radius not visible")
		}

		// Враг в радиусе видимости должен быть виден
		if result.ObjectLayer[3][3] != common.EntityTypeZombie {
			t.Errorf("Enemy within view radius not visible")
		}
	})

	t.Run("walls are remembered after viewing", func(t *testing.T) {
		// Поле 5x3 со стенами
		rf := common.NewRenderedField(5, 3)

		// Окружение
		rf.EnvironmentLayer[0][0] = common.WorldTypeWall
		rf.EnvironmentLayer[0][1] = common.WorldTypeWall
		rf.EnvironmentLayer[0][2] = common.WorldTypeWall
		rf.EnvironmentLayer[0][3] = common.WorldTypeRoomFloor
		rf.EnvironmentLayer[0][4] = common.WorldTypeRoomFloor
		for x := 0; x < 5; x++ {
			rf.EnvironmentLayer[1][x] = common.WorldTypeRoomFloor
			rf.EnvironmentLayer[2][x] = common.WorldTypeRoomFloor
		}

		fog := NewFogOfWar(5, 3)
		playerPos := primitives.Point2D[int]{X: 1, Y: 1}
		viewRadius := 2

		// Первый просмотр - игрок видит стену
		result1 := fog.ApplyFogOfWar(rf, playerPos, viewRadius)

		if result1.EnvironmentLayer[0][0] != common.WorldTypeWall {
			t.Errorf("Wall not visible on first view")
		}

		// Игрок перемещается далеко от стены
		playerPos = primitives.Point2D[int]{X: 4, Y: 2}

		// Второй просмотр - игрок далеко от стены
		result2 := fog.ApplyFogOfWar(rf, playerPos, viewRadius)

		// Стена всё ещё должна быть видна (запомнена в EnvironmentLayer)
		if result2.EnvironmentLayer[0][0] != common.WorldTypeWall {
			t.Errorf("Wall not remembered after moving away")
		}
	})

	t.Run("enemies disappear when out of view", func(t *testing.T) {
		// Поле 5x3
		rf := common.NewRenderedField(5, 3)

		// Окружение - проход (passage) в (0,0), остальное - пол
		for y := 0; y < 3; y++ {
			for x := 0; x < 5; x++ {
				rf.EnvironmentLayer[y][x] = common.WorldTypeRoomFloor
			}
		}
		rf.EnvironmentLayer[0][0] = common.WorldTypePassage

		// Враг стоит на проходе
		rf.ObjectLayer[0][0] = common.EntityTypeZombie

		fog := NewFogOfWar(5, 3)
		playerPos := primitives.Point2D[int]{X: 1, Y: 0}
		viewRadius := 1

		// Первый просмотр - игрок видит врага
		result1 := fog.ApplyFogOfWar(rf, playerPos, viewRadius)

		if result1.ObjectLayer[0][0] != common.EntityTypeZombie {
			t.Errorf("Enemy not visible when in view radius")
		}

		// Игрок перемещается далеко от врага
		playerPos = primitives.Point2D[int]{X: 4, Y: 2}

		// Второй просмотр - игрок далеко от врага
		result2 := fog.ApplyFogOfWar(rf, playerPos, viewRadius)

		// Враг не должен быть виден в ObjectLayer (динамический объект)
		if result2.ObjectLayer[0][0] != common.EntityTypeNone {
			t.Errorf("Enemy should not be visible when out of view radius, got %v", result2.ObjectLayer[0][0])
		}

		// Проход должен быть запомнен в EnvironmentLayer (passage - статический тайл)
		if result2.EnvironmentLayer[0][0] != common.WorldTypePassage {
			t.Errorf("Passage should be remembered after exploration, got %v", result2.EnvironmentLayer[0][0])
		}
	})

	t.Run("passages and doors are remembered", func(t *testing.T) {
		// Поле 5x3 с проходом и дверью
		rf := common.NewRenderedField(5, 3)

		rf.EnvironmentLayer[0][0] = common.WorldTypePassage
		rf.EnvironmentLayer[0][1] = common.WorldTypeDoor
		rf.EnvironmentLayer[0][2] = common.WorldTypeRoomFloor
		rf.EnvironmentLayer[0][3] = common.WorldTypeRoomFloor
		rf.EnvironmentLayer[0][4] = common.WorldTypeRoomFloor
		for x := 0; x < 5; x++ {
			rf.EnvironmentLayer[1][x] = common.WorldTypeRoomFloor
			rf.EnvironmentLayer[2][x] = common.WorldTypeRoomFloor
		}

		fog := NewFogOfWar(5, 3)
		playerPos := primitives.Point2D[int]{X: 1, Y: 0}
		viewRadius := 1

		// Первый просмотр
		result1 := fog.ApplyFogOfWar(rf, playerPos, viewRadius)

		if result1.EnvironmentLayer[0][0] != common.WorldTypePassage {
			t.Errorf("Passage not visible")
		}
		if result1.EnvironmentLayer[0][1] != common.WorldTypeDoor {
			t.Errorf("Door not visible")
		}

		// Игрок перемещается далеко
		playerPos = primitives.Point2D[int]{X: 4, Y: 2}

		// Второй просмотр
		result2 := fog.ApplyFogOfWar(rf, playerPos, viewRadius)

		// Проход и дверь должны быть запомнены
		if result2.EnvironmentLayer[0][0] != common.WorldTypePassage {
			t.Errorf("Passage not remembered")
		}
		if result2.EnvironmentLayer[0][1] != common.WorldTypeDoor {
			t.Errorf("Door not remembered")
		}
	})
}

// TestIsStaticTile проверяет определение статичных тайлов
func TestIsStaticTile(t *testing.T) {
	fog := NewFogOfWar(10, 10)

	tests := []struct {
		name       string
		entityType common.GameEntityType
		want       bool
	}{
		{"wall is static", common.WorldTypeWall, true},
		{"passage is static", common.WorldTypePassage, true},
		{"door is static", common.WorldTypeDoor, true},
		{"portal is static", common.WorldTypePortal, true},
		{"player is not static", common.EntityTypePlayer, false},
		{"enemy is not static", common.EntityTypeZombie, false},
		{"food is not static", common.Food, false},
		{"weapon is not static", common.Weapon, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fog.isStaticTile(tt.entityType)
			if got != tt.want {
				t.Errorf("isStaticTile(%v) = %v, want %v", tt.entityType, got, tt.want)
			}
		})
	}
}

// TestReset проверяет сброс памяти тумана войны
func TestReset(t *testing.T) {
	fog := NewFogOfWar(5, 5)

	// Добавляем несколько исследованных точек
	fog.ExploredTiles[primitives.Point2D[int]{X: 0, Y: 0}] = true
	fog.ExploredTiles[primitives.Point2D[int]{X: 1, Y: 1}] = true
	fog.ExploredTiles[primitives.Point2D[int]{X: 2, Y: 2}] = true

	if len(fog.ExploredTiles) != 3 {
		t.Errorf("Expected 3 explored tiles before reset, got %d", len(fog.ExploredTiles))
	}

	// Сбрасываем
	fog.Reset()

	if len(fog.ExploredTiles) != 0 {
		t.Errorf("Expected 0 explored tiles after reset, got %d", len(fog.ExploredTiles))
	}
}

// BenchmarkApplyFogOfWar проверяет производительность применения тумана войны
func BenchmarkApplyFogOfWar(b *testing.B) {
	// Создаём большое поле 100x100
	width, height := 100, 100
	rf := common.NewRenderedField(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x+y)%10 == 0 {
				rf.EnvironmentLayer[y][x] = common.WorldTypeWall
			} else {
				rf.EnvironmentLayer[y][x] = common.WorldTypeRoomFloor
			}
		}
	}

	fog := NewFogOfWar(width, height)
	playerPos := primitives.Point2D[int]{X: 50, Y: 50}
	viewRadius := 10

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fog.ApplyFogOfWar(rf, playerPos, viewRadius)
	}
}
