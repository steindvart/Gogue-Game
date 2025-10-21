package entity

type Weapon struct {
	Consumable   Consumable
	StrengthBuff uint
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
		StrengthBuff: getAttributeRandomPercentIncrease(),
	}
}

func (w *Weapon) Taken() {
	w.Consumable.Shape = Box{}
}

func (w *Weapon) Dropped(box Box) {
	w.Consumable.Shape = box
}

func (w *Weapon) Use() string {
	return w.Consumable.Name
}
