package entity

import (
	"fmt"
)

type Weapon struct {
	Item   Item
	Damage uint
}

func NewWeapon(box Box) *Weapon {
	weaponNames := []string{
		"Blade of the Forgotten Dawn",
		"Obsidian Reaver",
		"Fang of the Shadow Wolf",
		"Ironclad Cleaver",
		"Crimson Talon",
		"Thunderstrike Maul",
		"Serpent's Kiss Dagger",
		"Voidrend Sword",
		"Ebonheart Spear",
	}

	return &Weapon{
		Item: Item{
			Shape: box,
			Name:  getAttributeRandomName(weaponNames),
		},
		Damage: getAttributeRandomPercentIncrease(),
	}
}

func (w *Weapon) Taken() {
	w.Item.Taken()
}

func (w *Weapon) Dropped(box Box) {
	w.Item.Dropped(box)
}

func (w *Weapon) Use(p *Player) string {
	if p.Weapon != nil {
		currentWeapon := p.Weapon
		currentWeapon.Dropped(p.Character.Shape)
	}
	p.Weapon = w

	return fmt.Sprintf(
		"You picked up the %v, now all your attacks have %v extra damage",
		w.Item.Name,
		w.Damage,
	)
}
