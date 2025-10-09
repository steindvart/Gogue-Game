package entity

import "testing"

func TestNewRoom(t *testing.T) {
	tests := []struct {
		name     string
		roomType RoomType
		shape    Box
		portal   *Box
		want     Room
	}{
		{
			name:     "New room",
			roomType: RoomTypeFinish,
			shape:    Box{Point2D[int]{X: 6, Y: 3}, Size2D[uint]{Height: 10, Width: 10}},
			portal:   &Box{Point2D[int]{X: 8, Y: 7}, Size2D[uint]{Height: 1, Width: 1}},
			want: Room{
				Shape:   Box{Point2D[int]{X: 6, Y: 3}, Size2D[uint]{Height: 10, Width: 10}},
				Type:    RoomTypeFinish,
				Foods:   []Food{},
				Elixirs: []Elixir{},
				Scrolls: []Scroll{},
				Weapons: []Weapon{},
				Enemies: []Enemy{},
				Portal:  &Box{Point2D[int]{X: 8, Y: 7}, Size2D[uint]{Height: 1, Width: 1}},
			},
		},
		{
			name:     "New room",
			roomType: RoomTypeOrdinary,
			shape:    Box{Point2D[int]{X: 6, Y: 3}, Size2D[uint]{Height: 10, Width: 10}},
			portal:   nil,
			want: Room{
				Shape:   Box{Point2D[int]{X: 6, Y: 3}, Size2D[uint]{Height: 10, Width: 10}},
				Type:    RoomTypeOrdinary,
				Foods:   []Food{},
				Elixirs: []Elixir{},
				Scrolls: []Scroll{},
				Weapons: []Weapon{},
				Enemies: []Enemy{},
				Portal:  nil,
			},
		},
		{
			name:     "New room",
			roomType: RoomTypeStart,
			shape:    Box{Point2D[int]{X: 6, Y: 3}, Size2D[uint]{Height: 10, Width: 10}},
			portal:   nil,
			want: Room{
				Shape:   Box{Point2D[int]{X: 6, Y: 3}, Size2D[uint]{Height: 10, Width: 10}},
				Type:    RoomTypeStart,
				Foods:   []Food{},
				Elixirs: []Elixir{},
				Scrolls: []Scroll{},
				Weapons: []Weapon{},
				Enemies: []Enemy{},
				Portal:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := NewRoom(tt.roomType, tt.shape, tt.portal)
			if room.Shape != tt.want.Shape {
				t.Errorf("Expected shape to be %v, got %v", tt.want.Shape, room.Shape)
			}
			if room.Type != tt.want.Type {
				t.Errorf("Expected type to be %v, got %v", tt.want.Type, room.Type)
			}
			if len(room.Foods) != len(tt.want.Foods) {
				t.Errorf("Expected Foods to be empty, got %v", room.Foods)
			}
			if len(room.Elixirs) != len(tt.want.Elixirs) {
				t.Errorf("Expected Elixirs to be empty, got %v", room.Elixirs)
			}
			if len(room.Scrolls) != len(tt.want.Scrolls) {
				t.Errorf("Expected Scrolls to be empty, got %v", room.Scrolls)
			}
			if len(room.Weapons) != len(tt.want.Weapons) {
				t.Errorf("Expected Weapons to be empty, got %v", room.Weapons)
			}
			if len(room.Enemies) != len(tt.want.Enemies) {
				t.Errorf("Expected Enemies to be empty, got %v", room.Enemies)
			}
			if tt.want.Portal != nil {
				if room.Portal.Point != tt.want.Portal.Point {
					t.Errorf("Expected portal point to be %v, got %v", tt.want.Portal.Point, room.Portal.Point)
				}
				if room.Portal.Size != tt.want.Portal.Size {
					t.Errorf("Expected portal size to be %v, got %v", tt.want.Portal.Size, room.Portal.Size)
				}
			}
		})
	}
}
