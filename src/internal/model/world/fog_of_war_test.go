package world

import (
	"gogue/internal/common"
	"gogue/internal/model/primitives"
	"testing"
)

func TestBresenhamLine(t *testing.T) {
	fog := NewFogOfWar(10, 10)

	tests := []struct {
		name     string
		from     primitives.Point2D[int]
		to       primitives.Point2D[int]
		wantLen  int
		wantLast primitives.Point2D[int]
	}{
		{
			name:     "horizontal line right",
			from:     primitives.Point2D[int]{X: 0, Y: 0},
			to:       primitives.Point2D[int]{X: 5, Y: 0},
			wantLen:  6,
			wantLast: primitives.Point2D[int]{X: 5, Y: 0},
		},
		{
			name:     "vertical line down",
			from:     primitives.Point2D[int]{X: 0, Y: 0},
			to:       primitives.Point2D[int]{X: 0, Y: 5},
			wantLen:  6,
			wantLast: primitives.Point2D[int]{X: 0, Y: 5},
		},
		{
			name:     "diagonal line",
			from:     primitives.Point2D[int]{X: 0, Y: 0},
			to:       primitives.Point2D[int]{X: 3, Y: 3},
			wantLen:  4,
			wantLast: primitives.Point2D[int]{X: 3, Y: 3},
		},
		{
			name:     "single point",
			from:     primitives.Point2D[int]{X: 5, Y: 5},
			to:       primitives.Point2D[int]{X: 5, Y: 5},
			wantLen:  1,
			wantLast: primitives.Point2D[int]{X: 5, Y: 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := fog.bresenhamLine(tt.from, tt.to)

			if len(line) != tt.wantLen {
				t.Errorf("bresenhamLine() length = %d, want %d", len(line), tt.wantLen)
			}

			if len(line) > 0 {
				lastPoint := line[len(line)-1]
				if lastPoint != tt.wantLast {
					t.Errorf("bresenhamLine() last point = %+v, want %+v", lastPoint, tt.wantLast)
				}
			}

			// Проверяем, что первая точка - это начало
			if len(line) > 0 && line[0] != tt.from {
				t.Errorf("bresenhamLine() first point = %+v, want %+v", line[0], tt.from)
			}
		})
	}
}

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
		// Создаём простое поле 5x5
		fullField := [][]common.GameEntityType{
			{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
			{common.WorldTypeRoomFloor, common.FoodTypeBread, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
			{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.EntityTypePlayer, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
			{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.EntityTypeZombie, common.WorldTypeRoomFloor},
			{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
		}

		fog := NewFogOfWar(5, 5)
		playerPos := primitives.Point2D[int]{X: 2, Y: 2}
		viewRadius := 2

		result := fog.ApplyFogOfWar(fullField, playerPos, viewRadius)

		// Игрок должен видеть себя
		if result[2][2] != common.EntityTypePlayer {
			t.Errorf("Player not visible at own position")
		}

		// Еда в радиусе видимости должна быть видна
		if result[1][1] != common.FoodTypeBread {
			t.Errorf("Food within view radius not visible")
		}

		// Враг в радиусе видимости должен быть виден
		if result[3][3] != common.EntityTypeZombie {
			t.Errorf("Enemy within view radius not visible")
		}
	})

	t.Run("walls are remembered after viewing", func(t *testing.T) {
		// Поле со стеной
		fullField := [][]common.GameEntityType{
			{common.WorldTypeWall, common.WorldTypeWall, common.WorldTypeWall, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
			{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
			{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
		}

		fog := NewFogOfWar(5, 3)
		playerPos := primitives.Point2D[int]{X: 1, Y: 1}
		viewRadius := 2

		// Первый просмотр - игрок видит стену
		result1 := fog.ApplyFogOfWar(fullField, playerPos, viewRadius)

		// Стена должна быть видна
		if result1[0][0] != common.WorldTypeWall {
			t.Errorf("Wall not visible on first view")
		}

		// Игрок перемещается далеко от стены
		playerPos = primitives.Point2D[int]{X: 4, Y: 2}

		// Второй просмотр - игрок далеко от стены
		result2 := fog.ApplyFogOfWar(fullField, playerPos, viewRadius)

		// Стена всё ещё должна быть видна (запомнена)
		if result2[0][0] != common.WorldTypeWall {
			t.Errorf("Wall not remembered after moving away")
		}
	})

	t.Run("enemies disappear when out of view", func(t *testing.T) {
		// Поле с врагом
		fullField := [][]common.GameEntityType{
			{common.EntityTypeZombie, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
			{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
			{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
		}

		fog := NewFogOfWar(5, 3)
		playerPos := primitives.Point2D[int]{X: 1, Y: 0}
		viewRadius := 1

		// Первый просмотр - игрок видит врага
		result1 := fog.ApplyFogOfWar(fullField, playerPos, viewRadius)

		// Враг должен быть виден
		if result1[0][0] != common.EntityTypeZombie {
			t.Errorf("Enemy not visible when in view radius")
		}

		// Игрок перемещается далеко от врага
		playerPos = primitives.Point2D[int]{X: 4, Y: 2}

		// Второй просмотр - игрок далеко от врага
		result2 := fog.ApplyFogOfWar(fullField, playerPos, viewRadius)

		// Враг не должен быть виден (динамический объект)
		if result2[0][0] != 0 {
			t.Errorf("Enemy should not be visible when out of view radius")
		}
	})

	t.Run("passages and doors are remembered", func(t *testing.T) {
		// Поле с проходом и дверью
		fullField := [][]common.GameEntityType{
			{common.WorldTypePassage, common.WorldTypeDoor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
			{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
			{common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor, common.WorldTypeRoomFloor},
		}

		fog := NewFogOfWar(5, 3)
		playerPos := primitives.Point2D[int]{X: 1, Y: 0}
		viewRadius := 1

		// Первый просмотр
		result1 := fog.ApplyFogOfWar(fullField, playerPos, viewRadius)

		// Проход и дверь должны быть видны
		if result1[0][0] != common.WorldTypePassage {
			t.Errorf("Passage not visible")
		}
		if result1[0][1] != common.WorldTypeDoor {
			t.Errorf("Door not visible")
		}

		// Игрок перемещается далеко
		playerPos = primitives.Point2D[int]{X: 4, Y: 2}

		// Второй просмотр
		result2 := fog.ApplyFogOfWar(fullField, playerPos, viewRadius)

		// Проход и дверь должны быть запомнены
		if result2[0][0] != common.WorldTypePassage {
			t.Errorf("Passage not remembered")
		}
		if result2[0][1] != common.WorldTypeDoor {
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
		{"food is not static", common.FoodTypeBread, false},
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
	fog.exploredTiles[primitives.Point2D[int]{X: 0, Y: 0}] = true
	fog.exploredTiles[primitives.Point2D[int]{X: 1, Y: 1}] = true
	fog.exploredTiles[primitives.Point2D[int]{X: 2, Y: 2}] = true

	if len(fog.exploredTiles) != 3 {
		t.Errorf("Expected 3 explored tiles before reset, got %d", len(fog.exploredTiles))
	}

	// Сбрасываем
	fog.Reset()

	if len(fog.exploredTiles) != 0 {
		t.Errorf("Expected 0 explored tiles after reset, got %d", len(fog.exploredTiles))
	}
}

// BenchmarkApplyFogOfWar проверяет производительность применения тумана войны
func BenchmarkApplyFogOfWar(b *testing.B) {
	// Создаём большое поле 100x100
	width, height := 100, 100
	fullField := make([][]common.GameEntityType, height)
	for y := range fullField {
		fullField[y] = make([]common.GameEntityType, width)
		for x := range fullField[y] {
			// Заполняем поле случайными статичными элементами
			if (x+y)%10 == 0 {
				fullField[y][x] = common.WorldTypeWall
			}
		}
	}

	fog := NewFogOfWar(width, height)
	playerPos := primitives.Point2D[int]{X: 50, Y: 50}
	viewRadius := 10

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fog.ApplyFogOfWar(fullField, playerPos, viewRadius)
	}
}
