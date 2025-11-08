package world

import (
	"gogue/internal/model/entity"
	"gogue/internal/model/items"
	"gogue/internal/model/primitive"
)

type RoomType uint

const (
	RoomTypeStart RoomType = iota
	RoomTypeOrdinary
	RoomTypeFinish
)

type Room struct {
	Shape   primitive.Box
	Type    RoomType
	Foods   []items.Food
	Elixirs []items.Elixir
	Scrolls []items.Scroll
	Weapons []items.Weapon
	Enemies []entity.Enemy
}

func NewRoom(roomType RoomType, shape primitive.Box) *Room {
	return &Room{
		Shape:   shape,
		Type:    roomType,
		Foods:   []items.Food{},
		Elixirs: []items.Elixir{},
		Scrolls: []items.Scroll{},
		Weapons: []items.Weapon{},
		Enemies: []entity.Enemy{},
	}
}
