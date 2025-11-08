package entity

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

func (w *Weapon) Take() {
	w.Item.Take()
}

func (w *Weapon) Drop(box Box) {
	w.Item.Drop(box)
}

func (w *Weapon) Use(p *Player) {
	if p.Weapon != nil {
		currentWeapon := p.Weapon
		currentWeapon.Drop(p.Character.Shape)
	}
	p.Weapon = w
}
