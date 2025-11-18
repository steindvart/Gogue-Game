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
	Shape      primitives.Box
	Type       RoomType
	Foods      []items.Food
	Elixirs    []items.Elixir
	Scrolls    []items.Scroll
	Weapons    []items.Weapon
	Zombies    []entities.Zombie
	Vampires   []entities.Vampire
	Ghosts     []entities.Ghost
	Ogres      []entities.Ogre
	SnakeMages []entities.SnakeMage
}

func NewRoom(roomType RoomType, shape primitives.Box) *Room {
	return &Room{
		Shape:      shape,
		Type:       roomType,
		Foods:      []items.Food{},
		Elixirs:    []items.Elixir{},
		Scrolls:    []items.Scroll{},
		Weapons:    []items.Weapon{},
		Zombies:    []entities.Zombie{},
		Vampires:   []entities.Vampire{},
		Ghosts:     []entities.Ghost{},
		Ogres:      []entities.Ogre{},
		SnakeMages: []entities.SnakeMage{},
	}
}

func (r *Room) checkBoxOccupied(b *primitives.Box) bool {
	var occupied bool

	for _, food := range r.Foods {
		if *b == food.Shape {
			occupied = true
		}
	}

	for _, elixir := range r.Elixirs {
		if *b == elixir.Shape {
			occupied = true
		}
	}

	for _, scroll := range r.Scrolls {
		if *b == scroll.Shape {
			occupied = true
		}
	}

	for _, weapon := range r.Weapons {
		if *b == weapon.Shape {
			occupied = true
		}
	}

	for _, zombie := range r.Zombies {
		if *b == zombie.Enemy.Character.Shape {
			occupied = true
		}
	}

	for _, vampire := range r.Vampires {
		if *b == vampire.Enemy.Character.Shape {
			occupied = true
		}
	}

	for _, ghost := range r.Ghosts {
		if *b == ghost.Enemy.Character.Shape {
			occupied = true
		}
	}

	for _, ogre := range r.Ogres {
		if *b == ogre.Enemy.Character.Shape {
			occupied = true
		}
	}

	for _, snakeMage := range r.SnakeMages {
		if *b == snakeMage.Enemy.Character.Shape {
			occupied = true
		}
	}

	return occupied
}