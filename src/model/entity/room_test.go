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
		portal   *Box
		want     Room
	}{
		{
			name:     "New finish room",
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
			name:     "New ordinary room",
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
			name:     "New start room",
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
			reflect.DeepEqual(room, tt.want)
		})
	}
}
