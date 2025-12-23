package entities

import (
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
)

type Player struct {
	*Character
	*items.Backpack
	*items.Weapon
	Experience uint
	Level      uint
	ViewRadius int
}

func NewPlayer(box primitives.Box) *Player {
	return &Player{
		Character: &Character{
			Box: &box,
			Attributes: &primitives.Attributes{
				Health:    100,
				MaxHealth: 100,
				Strength:  10,
				Agility:   5,
			},
		},
		Experience: 0,
		Level:      1,
		ViewRadius: 3,
		Backpack:   items.NewBackpack(),
		Weapon:     nil,
	}
}

func (p *Player) DropEquipWeaponToBackpack() error {
	if p.Weapon == nil {
		return nil
	}

	if err := p.Backpack.AddItem(p.Weapon); err != nil {
		return err
	}

	_ = p.UnequipWeapon()

	return nil
}

func (p *Player) EquipWeapon(w *items.Weapon) error {
	if w == nil {
		return nil
	}

	// Если уже есть экипированный предмет, пытаемся положить его в рюкзак.
	if err := p.DropEquipWeaponToBackpack(); err != nil {
		return err
	}

	p.Weapon = w
	p.Character.Use(w)

	return nil
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
