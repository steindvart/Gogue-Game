package entity

import (
	"testing"
)

func TestNewRoom(t *testing.T) {
	t.Run("NewRoom - Ordinary Room", func(t *testing.T) {
		room, err := NewRoom(RoomTypeOrdinary, 0, 0, 10, 10)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if room.Type != RoomTypeOrdinary {
			t.Errorf("Expected room type %v, got %v", RoomTypeOrdinary, room.Type)
		}

		if room.Portal != nil {
			t.Error("Expected no portal for ordinary room")
		}

		if room.Shape.Size.Width < RoomMinWidth || room.Shape.Size.Height < RoomMinHeight {
			t.Errorf("Room size is smaller than minimum required: width=%d, height=%d", room.Shape.Size.Width, room.Shape.Size.Height)
		}

		if room.Shape.Point.X < 0 || room.Shape.Point.X >= 10 {
			t.Errorf("Room X position is out of bounds: %d", room.Shape.Point.X)
		}
		if room.Shape.Point.Y < 0 || room.Shape.Point.Y >= 10 {
			t.Errorf("Room Y position is out of bounds: %d", room.Shape.Point.Y)
		}
	})

	t.Run("NewRoom - Start Room with Portal", func(t *testing.T) {
		room, err := NewRoom(RoomTypeStart, 0, 0, 10, 10)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if room.Type != RoomTypeStart {
			t.Errorf("Expected room type %v, got %v", RoomTypeStart, room.Type)
		}

		if room.Portal == nil {
			t.Error("Expected portal for start room")
		}
	})

	t.Run("NewRoom - Finish Room with Portal", func(t *testing.T) {
		room, err := NewRoom(RoomTypeFinish, 0, 0, 10, 10)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if room.Type != RoomTypeFinish {
			t.Errorf("Expected room type %v, got %v", RoomTypeFinish, room.Type)
		}

		if room.Portal == nil {
			t.Error("Expected portal for finish room")
		}
	})

	t.Run("NewRoom - Error on Too Small Map", func(t *testing.T) {
		_, err := NewRoom(RoomTypeOrdinary, 0, 0, 3, 3)
		if err == nil {
			t.Fatal("Expected error for too small map, got nil")
		}

		expectedErr := "room cannot be built: map is too small to fit a room of minimum required size"
		if err.Error() != expectedErr {
			t.Errorf("Expected error message '%s', got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("NewRoom - Empty Slices Initialization", func(t *testing.T) {
		room, err := NewRoom(RoomTypeOrdinary, 0, 0, 10, 10)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(room.Foods) != 0 {
			t.Error("Expected empty Foods slice")
		}
		if len(room.Elixirs) != 0 {
			t.Error("Expected empty Elixirs slice")
		}
		if len(room.Scrolls) != 0 {
			t.Error("Expected empty Scrolls slice")
		}
		if len(room.Weapons) != 0 {
			t.Error("Expected empty Weapons slice")
		}
		if len(room.Enemies) != 0 {
			t.Error("Expected empty Enemies slice")
		}
	})
}

func TestCalculateRoomSize(t *testing.T) {
	t.Run("CalculateRoomSize - Valid Size", func(t *testing.T) {
		width, height, err := calculateRoomSize(10, 8)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if width < RoomMinWidth || width > 10 {
			t.Errorf("Width out of range: expected [%d, 10], got %d", RoomMinWidth, width)
		}

		if height < RoomMinHeight || height > 8 {
			t.Errorf("Height out of range: expected [%d, 8], got %d", RoomMinHeight, height)
		}
	})

	t.Run("CalculateRoomSize - Minimum Size", func(t *testing.T) {
		width, height, err := calculateRoomSize(RoomMinWidth, RoomMinHeight)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if width != RoomMinWidth {
			t.Errorf("Expected width %d, got %d", RoomMinWidth, width)
		}

		if height != RoomMinHeight {
			t.Errorf("Expected height %d, got %d", RoomMinHeight, height)
		}
	})
}

func TestGeneratePortal(t *testing.T) {
	t.Run("GeneratePortal - Valid Position", func(t *testing.T) {
		portal := generatePortal(0, 0, 10, 8)

		if portal.Shape.Size.Width != 1 || portal.Shape.Size.Height != 1 {
			t.Errorf("Expected portal size (1,1), got (%d,%d)", portal.Shape.Size.Width, portal.Shape.Size.Height)
		}

		if portal.Shape.Point.X < 1 || portal.Shape.Point.X >= 10-1 {
			t.Errorf("Portal X position out of valid range: %d", portal.Shape.Point.X)
		}

		if portal.Shape.Point.Y < 1 || portal.Shape.Point.Y >= 8-1 {
			t.Errorf("Portal Y position out of valid range: %d", portal.Shape.Point.Y)
		}
	})

	t.Run("GeneratePortal - Different Coordinates", func(t *testing.T) {
		portal := generatePortal(5, 5, 10, 8)

		if portal.Shape.Point.X < 6 || portal.Shape.Point.X >= 14 {
			t.Errorf("Portal X position out of valid range [6,14): %d", portal.Shape.Point.X)
		}

		if portal.Shape.Point.Y < 6 || portal.Shape.Point.Y >= 12 {
			t.Errorf("Portal Y position out of valid range [6,12): %d", portal.Shape.Point.Y)
		}
	})
}
