package world

import (
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
)

type RoomType uint

const (
	RoomTypeStart RoomType = iota
	RoomTypeOrdinary
	RoomTypeFinish
)

type Room struct {
	Shape   primitives.Box
	Type    RoomType
	Doors   []primitives.Point2D[int]
	Foods   []items.Food
	Elixirs []items.Elixir
	Scrolls []items.Scroll
	Weapons []items.Weapon
	Enemies []entities.Enemy
}

func NewRoom(roomType RoomType, shape primitives.Box) *Room {
	return &Room{
		Shape:   shape,
		Type:    roomType,
		Foods:   []items.Food{},
		Elixirs: []items.Elixir{},
		Scrolls: []items.Scroll{},
		Weapons: []items.Weapon{},
		Enemies: []entities.Enemy{},
	}
}
