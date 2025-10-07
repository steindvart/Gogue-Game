package entity

import (
	"math/rand"
	"testing"
)

func TestNewRoom(t *testing.T) {
	source := rand.NewSource(42)
	rng := rand.New(source)

	tests := []struct {
		name        string
		roomType    RoomType
		x, y        int
		mapHeight   int
		mapWidth    int
		expectError bool
	}{
		{
			name:        "successful creation ordinary room",
			roomType:    RoomTypeOrdinary,
			x:           5,
			y:           5,
			mapHeight:   20,
			mapWidth:    30,
			expectError: false,
		},
		{
			name:        "successful creation start room",
			roomType:    RoomTypeStart,
			x:           0,
			y:           0,
			mapHeight:   15,
			mapWidth:    15,
			expectError: false,
		},
		{
			name:        "successful creation finish room",
			roomType:    RoomTypeFinish,
			x:           10,
			y:           10,
			mapHeight:   12,
			mapWidth:    12,
			expectError: false,
		},
		{
			name:        "map too small",
			roomType:    RoomTypeStart,
			x:           0,
			y:           0,
			mapHeight:   2,
			mapWidth:    2,
			expectError: true,
		},
		{
			name:        "boundary case minimum size",
			roomType:    RoomTypeOrdinary,
			x:           0,
			y:           0,
			mapHeight:   12,
			mapWidth:    9,
			expectError: true,
		},
		{
			name:        "boundary case exact minimum",
			roomType:    RoomTypeOrdinary,
			x:           0,
			y:           0,
			mapHeight:   12,
			mapWidth:    12,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room, err := NewRoom(tt.roomType, tt.x, tt.y, tt.mapHeight, tt.mapWidth, rng)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				expectedMsg := "room cannot be built: map is too small to fit a room of minimum required size"
				if err.Error() != expectedMsg {
					t.Errorf("expected error message '%s', got '%s'", expectedMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if room == nil {
				t.Fatal("expected room to be created, got nil")
			}

			if room.Type != tt.roomType {
				t.Errorf("expected room type %d, got %d", tt.roomType, room.Type)
			}

			if room.Shape.Point.X != tt.x || room.Shape.Point.Y != tt.y {
				t.Errorf("expected position (%d, %d), got (%d, %d)", tt.x, tt.y, room.Shape.Point.X, room.Shape.Point.Y)
			}

			// Проверяем, что размеры в допустимых пределах
			if room.Shape.Size.Width < RoomMinWidth {
				t.Errorf("expected width >= %d, got %d", RoomMinWidth, room.Shape.Size.Width)
			}

			if room.Shape.Size.Height < RoomMinHeight {
				t.Errorf("expected height >= %d, got %d", RoomMinHeight, room.Shape.Size.Height)
			}

			// Проверяем, что слайсы инициализированы пустыми
			if len(room.Foods) != 0 {
				t.Errorf("expected empty Foods slice, got length %d", len(room.Foods))
			}

			if len(room.Elixirs) != 0 {
				t.Errorf("expected empty Elixirs slice, got length %d", len(room.Elixirs))
			}

			if len(room.Scrolls) != 0 {
				t.Errorf("expected empty Scrolls slice, got length %d", len(room.Scrolls))
			}

			if len(room.Weapons) != 0 {
				t.Errorf("expected empty Weapons slice, got length %d", len(room.Weapons))
			}

			if len(room.Enemies) != 0 {
				t.Errorf("expected empty Enemies slice, got length %d", len(room.Enemies))
			}
		})
	}
}

func TestCalculateRoomSize(t *testing.T) {
	source := rand.NewSource(42)
	rng := rand.New(source)

	tests := []struct {
		name        string
		mapHeight   int
		mapWidth    int
		expectError bool
	}{
		{
			name:        "normal calculation",
			mapHeight:   12,
			mapWidth:    18,
			expectError: false,
		},
		{
			name:        "map too small",
			mapHeight:   2,
			mapWidth:    2,
			expectError: true,
		},
		{
			name:        "boundary case exact minimum",
			mapHeight:   12,
			mapWidth:    12,
			expectError: false,
		},
		{
			name:        "width too small",
			mapHeight:   12,
			mapWidth:    6,
			expectError: true,
		},
		{
			name:        "height too small",
			mapHeight:   6,
			mapWidth:    12,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width, height, err := calculateRoomSize(tt.mapHeight, tt.mapWidth, rng)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				expectedMsg := "room cannot be built: map is too small to fit a room of minimum required size"
				if err.Error() != expectedMsg {
					t.Errorf("expected error message '%s', got '%s'", expectedMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			expectedMaxWidth := uint(tt.mapWidth / RoomsCountEntireLevelWidth)
			expectedMaxHeight := uint(tt.mapHeight / RoomsCountEntireLevelHeight)

			if width < RoomMinWidth || width > expectedMaxWidth {
				t.Errorf("expected width between %d and %d, got %d", RoomMinWidth, expectedMaxWidth, width)
			}

			if height < RoomMinHeight || height > expectedMaxHeight {
				t.Errorf("expected height between %d and %d, got %d", RoomMinHeight, expectedMaxHeight, height)
			}
		})
	}
}
