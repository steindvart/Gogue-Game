package entity

import "fmt"

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
	f.Consumable.Taken()
}

func (f *Food) Dropped(box Box) ConsumableLike {
	f.Consumable.Dropped(box)
	return f
}

func (f *Food) Use(p *Player) (string, ConsumableLike) {
	p.Character.Health += float64(f.HealthRegeneration)

	return fmt.Sprintf(
		"You ate the %v, your Health has increased by %v",
		f.Consumable.Name,
		f.HealthRegeneration,
	), nil
}
