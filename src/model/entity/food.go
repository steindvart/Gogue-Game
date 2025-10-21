package entity

type Food struct {
	Consumable         Consumable
	HealthRegeneration uint
}

func NewFood(box Box) *Food {
	foodNames := []string{
		"Ration of the Ironclad",
		"Crimson Berry Cluster",
		"Loaf of the Forgotten Baker",
		"Smoked Wyrm Jerky",
		"Golden Apple of Vitality",
		"Hardtack of the Endless March",
		"Spiced Venison Strips",
		"Honeyed Nectar Bread",
		"Dried Mushrooms of the Deep",
	}

	return &Food{
		Consumable: Consumable{
			Shape: box,
			Name:  getAttributeRandomName(foodNames),
		},
		HealthRegeneration: getAttributeRandomPercentIncrease(),
	}
}

func (f *Food) Taken() {
	f.Consumable.Shape = Box{}
}

func (f *Food) Dropped(box Box) {
	f.Consumable.Shape = box
}

func (f *Food) Use() string {
	return f.Consumable.Name
}
