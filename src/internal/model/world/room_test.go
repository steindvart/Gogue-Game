package world

import (
	"gogue/internal/model/entity"
	"gogue/internal/model/items"
	"gogue/internal/model/primitive"
	"reflect"
	"testing"
)

func TestNewRoom(t *testing.T) {
	tests := []struct {
		name     string
		roomType RoomType
		shape    primitive.Box
		portal   *primitive.Box
		want     Room
	}{
		{
			name:     "New finish room",
			roomType: RoomTypeFinish,
			shape:    primitive.Box{Point: primitive.Point2D[int]{X: 6, Y: 3}, Size: primitive.Size2D[uint]{Height: 10, Width: 10}},
			want: Room{
				Shape:   primitive.Box{Point: primitive.Point2D[int]{X: 6, Y: 3}, Size: primitive.Size2D[uint]{Height: 10, Width: 10}},
				Type:    RoomTypeFinish,
				Foods:   []items.Food{},
				Elixirs: []items.Elixir{},
				Scrolls: []items.Scroll{},
				Weapons: []items.Weapon{},
				Enemies: []entity.Enemy{},
			},
		},
		{
			name:     "New ordinary room",
			roomType: RoomTypeOrdinary,
			shape:    primitive.Box{Point: primitive.Point2D[int]{X: 6, Y: 3}, Size: primitive.Size2D[uint]{Height: 10, Width: 10}},
			want: Room{
				Shape:   primitive.Box{Point: primitive.Point2D[int]{X: 6, Y: 3}, Size: primitive.Size2D[uint]{Height: 10, Width: 10}},
				Type:    RoomTypeOrdinary,
				Foods:   []items.Food{},
				Elixirs: []items.Elixir{},
				Scrolls: []items.Scroll{},
				Weapons: []items.Weapon{},
				Enemies: []entity.Enemy{},
			},
		},
		{
			name:     "New start room",
			roomType: RoomTypeStart,
			shape:    primitive.Box{Point: primitive.Point2D[int]{X: 6, Y: 3}, Size: primitive.Size2D[uint]{Height: 10, Width: 10}},
			want: Room{
				Shape:   primitive.Box{Point: primitive.Point2D[int]{X: 6, Y: 3}, Size: primitive.Size2D[uint]{Height: 10, Width: 10}},
				Type:    RoomTypeStart,
				Foods:   []items.Food{},
				Elixirs: []items.Elixir{},
				Scrolls: []items.Scroll{},
				Weapons: []items.Weapon{},
				Enemies: []entity.Enemy{},
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
