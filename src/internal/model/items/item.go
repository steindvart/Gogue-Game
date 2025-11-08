package items

import (
	"gogue/internal/model/primitive"
	"math/rand"
)

const (
	IncreaseAttributeBaseParcentage int = 5
	IncreaseAttributeMaxPercentage  int = 20
	AttributesNum                   int = 3
)

type AttributeType int

const (
	AttributeTypeMaxHealth AttributeType = iota
	AttributeTypeAgility
	AttributeTypeStrength
)

type Type interface {
	Food | Elixir | Scroll | Weapon
}

type Item struct {
	Shape primitive.Box
	Name  string
}

func (i *Item) Take() {
	i.Shape = primitive.Box{}
}

func (i *Item) Drop(box primitive.Box) {
	i.Shape = box
}

func getAttributeRandomPercentIncrease() float64 {
	return float64(IncreaseAttributeBaseParcentage + rand.Intn(IncreaseAttributeMaxPercentage+1))
}

func getAttributeRandomName(names []string) string {
	return names[rand.Intn(len(names))]
}

func getRandomAttribute() primitive.Attributes {
	attribute := rand.Intn(AttributesNum)
	attributes := primitive.Attributes{}

	switch attribute {
	case int(AttributeTypeMaxHealth):
		attributes.MaxHealth = 1
	case int(AttributeTypeAgility):
		attributes.Agility = 1
	case int(AttributeTypeStrength):
		attributes.Strength = 1
	}

	return attributes
}
