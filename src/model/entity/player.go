package entity

type Player struct {
	BaseCharacter Character
	MaxHealth     uint
	Backpack      *Backpack
	Weapon        *Weapon
}
