package items

import "gogue/internal/model/primitives"

type Weapon struct {
	Item   Item
	Damage float64
}

func NewWeapon(box primitives.Box) *Weapon {
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

func (w *Weapon) Drop(box primitives.Box) {
	w.Item.Drop(box)
}

func (w *Weapon) Use() primitives.Attributes {
	return primitives.Attributes{
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
