package world

import (
	"gogue/internal/common"
	"gogue/internal/model/primitives"
	"testing"
)

func TestBresenhamLine_PackageLevel(t *testing.T) {
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
		{
			name:     "reverse horizontal",
			from:     primitives.Point2D[int]{X: 5, Y: 0},
			to:       primitives.Point2D[int]{X: 0, Y: 0},
			wantLen:  6,
			wantLast: primitives.Point2D[int]{X: 0, Y: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := BresenhamLine(tt.from, tt.to)

			if len(line) != tt.wantLen {
				t.Errorf("BresenhamLine() length = %d, want %d", len(line), tt.wantLen)
			}

			if len(line) > 0 {
				first := line[0]
				if first != tt.from {
					t.Errorf("BresenhamLine() first = %+v, want %+v", first, tt.from)
				}
				last := line[len(line)-1]
				if last != tt.wantLast {
					t.Errorf("BresenhamLine() last = %+v, want %+v", last, tt.wantLast)
				}
			}
		})
	}
}

// Вспомогательная функция для создания карты из ASCII-представления.
// Обозначения:
//
//	'.' = WorldTypeRoomFloor (проходимо, прозрачно)
//	'#' = WorldTypeWall      (блокирует видимость)
//	' ' = EntityTypeNone     (пустота, блокирует видимость)
//	'P' = EntityTypePlayer   (как пол для LoS)
//	'E' = EntityTypeZombie   (как пол для LoS)
func makeFieldFromASCII(lines []string) [][]common.GameEntityType {
	h := len(lines)
	if h == 0 {
		return nil
	}
	w := len(lines[0])

	field := make([][]common.GameEntityType, h)
	for y := 0; y < h; y++ {
		field[y] = make([]common.GameEntityType, w)
		for x := 0; x < w; x++ {
			if x < len(lines[y]) {
				switch lines[y][x] {
				case '#':
					field[y][x] = common.WorldTypeWall
				case '.':
					field[y][x] = common.WorldTypeRoomFloor
				case 'P':
					field[y][x] = common.EntityTypePlayer
				case 'E':
					field[y][x] = common.EntityTypeZombie
				default:
					field[y][x] = common.EntityTypeNone
				}
			}
		}
	}
	return field
}

func TestHasLineOfSight(t *testing.T) {
	tests := []struct {
		name   string
		field  []string
		origin primitives.Point2D[int]
		target primitives.Point2D[int]
		want   bool
	}{
		{
			name: "clear line of sight in open room",
			field: []string{
				"#####",
				"#...#",
				"#...#",
				"#...#",
				"#####",
			},
			origin: primitives.Point2D[int]{X: 1, Y: 1},
			target: primitives.Point2D[int]{X: 3, Y: 3},
			want:   true,
		},
		{
			name: "wall blocks line of sight",
			field: []string{
				"#####",
				"#...#",
				"##..#",
				"#...#",
				"#####",
			},
			origin: primitives.Point2D[int]{X: 1, Y: 1},
			target: primitives.Point2D[int]{X: 1, Y: 3},
			want:   false,
		},
		{
			name: "adjacent tiles always visible",
			field: []string{
				"#####",
				"#...#",
				"#...#",
				"#####",
			},
			origin: primitives.Point2D[int]{X: 1, Y: 1},
			target: primitives.Point2D[int]{X: 2, Y: 1},
			want:   true,
		},
		{
			name: "same position",
			field: []string{
				"#####",
				"#...#",
				"#####",
			},
			origin: primitives.Point2D[int]{X: 2, Y: 1},
			target: primitives.Point2D[int]{X: 2, Y: 1},
			want:   true,
		},
		{
			name: "enemy in another room separated by wall",
			field: []string{
				"#########",
				"#...#...#",
				"#...#...#",
				"#...#...#",
				"#########",
			},
			origin: primitives.Point2D[int]{X: 2, Y: 2},
			target: primitives.Point2D[int]{X: 6, Y: 2},
			want:   false,
		},
		{
			name: "enemy visible through door",
			field: []string{
				"#########",
				"#...#...#",
				"#.......#",
				"#...#...#",
				"#########",
			},
			origin: primitives.Point2D[int]{X: 2, Y: 2},
			target: primitives.Point2D[int]{X: 6, Y: 2},
			want:   true,
		},
		{
			name: "void (EntityTypeNone) blocks line of sight",
			field: []string{
				"#####",
				"#...#",
				"#. .#",
				"#...#",
				"#####",
			},
			origin: primitives.Point2D[int]{X: 1, Y: 2},
			target: primitives.Point2D[int]{X: 3, Y: 2},
			want:   false,
		},
		{
			name: "long diagonal clear",
			field: []string{
				"#######",
				"#.....#",
				"#.....#",
				"#.....#",
				"#.....#",
				"#.....#",
				"#######",
			},
			origin: primitives.Point2D[int]{X: 1, Y: 1},
			target: primitives.Point2D[int]{X: 5, Y: 5},
			want:   true,
		},
		{
			name:   "empty field returns false",
			field:  []string{},
			origin: primitives.Point2D[int]{X: 0, Y: 0},
			target: primitives.Point2D[int]{X: 1, Y: 1},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := makeFieldFromASCII(tt.field)
			got := HasLineOfSight(field, tt.origin, tt.target)
			if got != tt.want {
				t.Errorf("HasLineOfSight() = %v, want %v", got, tt.want)
			}
		})
	}
}
