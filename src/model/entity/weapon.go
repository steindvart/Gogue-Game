package entity

import (
	"fmt"
)

type Weapon struct {
	Consumable Consumable
	Damage     uint
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
		Consumable: Consumable{
			Shape: box,
			Name:  getAttributeRandomName(weaponNames),
		},
		Damage: getAttributeRandomPercentIncrease(),
	}
}

func (w *Weapon) Taken() {
	w.Consumable.Taken()
}

func (w *Weapon) Dropped(box Box) ConsumableLike {
	w.Consumable.Dropped(box)
	return w
}

func (w *Weapon) Use(p *Player) (string, ConsumableLike) {
	currentWeapon := p.Weapon
	p.Weapon = w

	return fmt.Sprintf(
		"You picked up the %v, now all your attacks have %v extra damage",
		w.Consumable.Name,
		w.Damage,
	), currentWeapon
}
