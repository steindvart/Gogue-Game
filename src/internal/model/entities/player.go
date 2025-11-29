package entities

import (
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
)

type Player struct {
	*Character
	*items.Backpack
	*items.Weapon
	Experience     uint
	CharacterLevel uint
}

func NewPlayer(box *primitives.Box) *Player {
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

func (p *Player) EquipWeapon(w *items.Weapon) {
	if w == nil {
		return
	}

	// @todo - сделать обработку случая когда уже есть экипированный предмет
	p.Weapon = w
	p.Character.Use(w)
}

func (p *Player) UnequipWeapon() *items.Weapon {
	w := p.Weapon
	if w != nil {
		p.ApplyEffect(&primitives.Effect{
			Attributes: primitives.Inverse(w.Effect.Attributes),
		})
		p.Weapon = nil
	}
	return w
}
