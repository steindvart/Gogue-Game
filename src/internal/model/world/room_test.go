package world

import (
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"reflect"
	"testing"
)

func TestNewRoom(t *testing.T) {
	tests := []struct {
		name     string
		roomType RoomType
		shape    primitives.Box
		want     Room
	}{
		{
			name:     "New finish room",
			roomType: RoomTypeFinish,
			shape:    primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 10, Width: 10}},
			want: Room{
				Shape:   primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 10, Width: 10}},
				Type:    RoomTypeFinish,
				Foods:   []items.Food{},
				Elixirs: []items.Elixir{},
				Scrolls: []items.Scroll{},
				Weapons: []items.Weapon{},
				Enemies: []entities.Enemy{},
			},
		},
		{
			name:     "New ordinary room",
			roomType: RoomTypeOrdinary,
			shape:    primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 10, Width: 10}},
			want: Room{
				Shape:   primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 10, Width: 10}},
				Type:    RoomTypeOrdinary,
				Foods:   []items.Food{},
				Elixirs: []items.Elixir{},
				Scrolls: []items.Scroll{},
				Weapons: []items.Weapon{},
				Enemies: []entities.Enemy{},
			},
		},
		{
			name:     "New start room",
			roomType: RoomTypeStart,
			shape:    primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 10, Width: 10}},
			want: Room{
				Shape:   primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 10, Width: 10}},
				Type:    RoomTypeStart,
				Foods:   []items.Food{},
				Elixirs: []items.Elixir{},
				Scrolls: []items.Scroll{},
				Weapons: []items.Weapon{},
				Enemies: []entities.Enemy{},
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
