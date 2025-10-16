package entity

type Food struct {
	Shape              Box
	HealthRegeneration uint
	Name               string
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
		Shape:              box,
		HealthRegeneration: getAttributeRandomPercentIncrease(),
		Name:               getAttributeRandomName(foodNames),
	}
}
