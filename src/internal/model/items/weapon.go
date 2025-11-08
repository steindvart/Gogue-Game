package items

import (
	"gogue/internal/model/primitive"
)

type Weapon struct {
	Item   Item
	Damage float64
}

func NewWeapon(box primitive.Box) *Weapon {
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

func (w *Weapon) Take() {
	w.Item.Take()
}

func (w *Weapon) Drop(box primitive.Box) {
	w.Item.Drop(box)
}

func (w *Weapon) Use() primitive.Attributes {
	return primitive.Attributes{
		Strength: float64(w.Damage),
	}
}

func AsWeapon(item any) *Weapon {
	w, ok := item.(*Weapon)
	if ok {
		return w
	}
	return nil
}
