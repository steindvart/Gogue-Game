package items

import (
	"gogue/internal/model/primitive"
)

type Food struct {
	Item               Item
	HealthRegeneration float64
}

func NewFood(box primitive.Box) *Food {
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

func (f *Food) Take() {
	f.Item.Take()
}

func (f *Food) Drop(box primitive.Box) {
	f.Item.Drop(box)
}

func (e *Food) Use() primitive.Attributes {
	return primitive.Attributes{
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
