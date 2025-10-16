package entity

import "math/rand"

const (
	IncreaseAttributeBaseParcentage int = 5
	IncreaseAttributeMaxPercentage  int = 20
)

type ConsumableLike interface {
	Taken()
	Dropped(box Box)
	Use() string
}

type Consumable struct {
	Shape Box
	Name  string
}

type Attributes struct {
	MaxHealth, Agility, Strength uint
}

func (consumable *Consumable) Taken() {
	consumable.Shape = Box{}
}

func (consumable *Consumable) Dropped(box Box) {
	consumable.Shape = box
}

func (consumable *Consumable) Use() string {
	return consumable.Name
}

func getAttributeRandomPercentIncrease() uint {
	return uint(IncreaseAttributeBaseParcentage + rand.Intn(IncreaseAttributeMaxPercentage+1))
}

func getAttributeRandomName(names []string) string {
	return names[rand.Intn(len(names))]
}

func getRandomAttribute() Attributes {
	attributesNum := 3
	attribute := rand.Intn(attributesNum)

	attributes := Attributes{}

	switch attribute {
	case 0:
		attributes.MaxHealth = 1
	case 1:
		attributes.Agility = 1
	case 2:
		attributes.Strength = 1
	}

	return attributes
}
