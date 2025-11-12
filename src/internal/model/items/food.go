package items

import "gogue/internal/model/primitives"

type Food struct {
	Item               Item
	HealthRegeneration float64
}

func NewFood(box primitives.Box) *Food {
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
		Item: Item{
			Shape: box,
			Name:  getAttributeRandomName(foodNames),
		},
		HealthRegeneration: getAttributeRandomPercentIncrease(),
	}
}

func (e *Food) Use() primitives.Attributes {
	return primitives.Attributes{
		Health: float64(e.HealthRegeneration),
	}
}

func AsFood(item any) *Food {
	f, ok := item.(*Food)
	if ok {
		return f
	}
	return nil
}
