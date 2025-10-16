package entity

type Weapon struct {
	Shape        Box
	StrengthBuff uint
	Name         string
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
		Shape:        box,
		StrengthBuff: getAttributeRandomPercentIncrease(),
		Name:         getAttributeRandomName(weaponNames),
	}
}
