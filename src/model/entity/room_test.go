package entity

import (
	"reflect"
	"testing"
)

func TestNewRoom(t *testing.T) {
	tests := []struct {
		name     string
		roomType RoomType
		shape    Box
		want     Room
	}{
		{
			name:     "New finish room",
			roomType: RoomTypeFinish,
			shape:    Box{Point2D[int]{X: 6, Y: 3}, Size2D[uint]{Height: 10, Width: 10}},
			want: Room{
				Shape:   Box{Point2D[int]{X: 6, Y: 3}, Size2D[uint]{Height: 10, Width: 10}},
				Type:    RoomTypeFinish,
				Foods:   []Food{},
				Elixirs: []Elixir{},
				Scrolls: []Scroll{},
				Weapons: []Weapon{},
				Enemies: []Enemy{},
			},
		},
		{
			name:     "New ordinary room",
			roomType: RoomTypeOrdinary,
			shape:    Box{Point2D[int]{X: 6, Y: 3}, Size2D[uint]{Height: 10, Width: 10}},
			want: Room{
				Shape:   Box{Point2D[int]{X: 6, Y: 3}, Size2D[uint]{Height: 10, Width: 10}},
				Type:    RoomTypeOrdinary,
				Foods:   []Food{},
				Elixirs: []Elixir{},
				Scrolls: []Scroll{},
				Weapons: []Weapon{},
				Enemies: []Enemy{},
			},
		},
		{
			name:     "New start room",
			roomType: RoomTypeStart,
			shape:    Box{Point2D[int]{X: 6, Y: 3}, Size2D[uint]{Height: 10, Width: 10}},
			want: Room{
				Shape:   Box{Point2D[int]{X: 6, Y: 3}, Size2D[uint]{Height: 10, Width: 10}},
				Type:    RoomTypeStart,
				Foods:   []Food{},
				Elixirs: []Elixir{},
				Scrolls: []Scroll{},
				Weapons: []Weapon{},
				Enemies: []Enemy{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := NewRoom(tt.roomType, tt.shape)
			reflect.DeepEqual(room, tt.want)
		})
	}
}
