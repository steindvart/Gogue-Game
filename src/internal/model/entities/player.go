package entities

import (
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
)

type Player struct {
	*Character
	Experience     uint
	CharacterLevel uint
	*items.Backpack
	*items.Weapon
}

func (p *Player) Attack() float64 {
	damage := 0.0
	if p.Weapon != nil {
		damage += p.Weapon.Effect.Attributes.Strength
	}

	damage += p.Character.Attack()
	return damage
}

func NewPlayer(box primitives.Box) *Player {
	return &Player{
		Character: &Character{
			Shape: box,
			Attributes: primitives.Attributes{
				Health:    100,
				MaxHealth: 100,
				Strength:  10,
				Agility:   5,
			},
		},
		Experience:     0,
		CharacterLevel: 1,
		Backpack:       items.NewBackpack(),
		Weapon:         nil,
	}
}
