package world

import (
	"gogue/internal/model/primitives"
	"math/rand"
	"testing"
)

const testSeed = 42

func TestGridRoomGenerator_GenerateRooms(t *testing.T) {
	tests := []struct {
		name      string
		mapSize   primitives.Size2D[uint]
		wantErr   bool
		errorText string
	}{
		{
			name:    "Valid map size 30x90",
			mapSize: primitives.Size2D[uint]{Height: 30, Width: 90},
			wantErr: false,
		},
		{
			name:    "Valid map size 45x150",
			mapSize: primitives.Size2D[uint]{Height: 45, Width: 150},
			wantErr: false,
		},
		{
			name:      "Map too small - below minimum",
			mapSize:   primitives.Size2D[uint]{Height: 6, Width: 6},
			wantErr:   true,
			errorText: "game map size is too small",
		},
		{
			name:      "Map too big - above maximum",
			mapSize:   primitives.Size2D[uint]{Height: 300, Width: 300},
			wantErr:   true,
			errorText: "game map size is too big",
		},
		{
			name:    "Minimum valid size",
			mapSize: primitives.Size2D[uint]{Height: 12, Width: 12},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := rand.New(rand.NewSource(testSeed))
			config := DefaultLevelConfig(tt.mapSize)
			generator := NewGridRoomGenerator()

			rooms, portal, err := generator.GenerateRooms(config, source)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("Expected error containing %q, but got nil", tt.errorText)
				}
				if err.Error() != tt.errorText {
					t.Errorf("Expected error %q, got %q", tt.errorText, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				// Проверяем количество комнат
				if len(rooms) != config.RoomsCount {
					t.Errorf("Expected %d rooms, got %d", config.RoomsCount, len(rooms))
				}

				// Проверяем что портал создан
				if portal.Size.Width == 0 || portal.Size.Height == 0 {
					t.Error("Finish portal not generated")
				}

				// Проверяем типы комнат
				var hasStart, hasFinish bool
				for _, room := range rooms {
					if room.Type == RoomTypeStart {
						hasStart = true
					}
					if room.Type == RoomTypeFinish {
						hasFinish = true
					}
				}

				if !hasStart {
					t.Error("No start room found")
				}
				if !hasFinish {
					t.Error("No finish room found")
				}

				// Проверяем что комнаты находятся в пределах карты
				for i, room := range rooms {
					if room.Shape.Point.X < 0 || room.Shape.Point.Y < 0 {
						t.Errorf("Room %d has negative coordinates: %+v", i, room.Shape.Point)
					}
					maxX := room.Shape.Point.X + int(room.Shape.Size.Width)
					maxY := room.Shape.Point.Y + int(room.Shape.Size.Height)
					if maxX > int(tt.mapSize.Width) || maxY > int(tt.mapSize.Height) {
						t.Errorf("Room %d exceeds map boundaries: maxX=%d, maxY=%d, mapSize=%+v",
							i, maxX, maxY, tt.mapSize)
					}
				}
			}
		})
	}
}

func TestGridRoomGenerator_GenerateRooms_Deterministic(t *testing.T) {
	source := rand.New(rand.NewSource(testSeed))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}
	config := DefaultLevelConfig(mapSize)
	generator := NewGridRoomGenerator()

	rooms, portal, err := generator.GenerateRooms(config, source)
	if err != nil {
		t.Fatalf("GenerateRooms failed: %v", err)
	}

	// Проверяем детерминированность
	wantRoomsCount := 9
	if len(rooms) != wantRoomsCount {
		t.Errorf("Expected %d rooms, got %d", wantRoomsCount, len(rooms))
	}

	// Проверяем портал
	if portal.Size.Width != 1 || portal.Size.Height != 1 {
		t.Errorf("Expected portal size 1x1, got %dx%d", portal.Size.Width, portal.Size.Height)
	}

	// Проверяем что room[7] - стартовая комната
	if rooms[7].Type != RoomTypeStart {
		t.Errorf("Expected room[7] to be Start, got %v", rooms[7].Type)
	}

	// Проверяем что room[5] - финишная комната
	if rooms[5].Type != RoomTypeFinish {
		t.Errorf("Expected room[5] to be Finish, got %v", rooms[5].Type)
	}
}

func TestGridRoomGenerator_RoomDistribution(t *testing.T) {
	source := rand.New(rand.NewSource(testSeed))
	mapSize := primitives.Size2D[uint]{Height: 45, Width: 135}
	config := DefaultLevelConfig(mapSize)
	generator := NewGridRoomGenerator()

	rooms, _, err := generator.GenerateRooms(config, source)
	if err != nil {
		t.Fatalf("GenerateRooms failed: %v", err)
	}

	// Проверяем что комнаты распределены по сетке 3x3
	// Разбиваем карту на 9 секций и проверяем что в каждой есть комната
	sectionWidth := int(mapSize.Width) / 3
	sectionHeight := int(mapSize.Height) / 3

	sectionsOccupied := make(map[int]bool)
	for _, room := range rooms {
		sectionX := room.Shape.Point.X / sectionWidth
		sectionY := room.Shape.Point.Y / sectionHeight
		sectionIndex := sectionY*3 + sectionX
		sectionsOccupied[sectionIndex] = true
	}

	if len(sectionsOccupied) != 9 {
		t.Errorf("Expected rooms in all 9 sections, got %d occupied sections", len(sectionsOccupied))
	}
}
