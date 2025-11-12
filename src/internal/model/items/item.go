package items

import (
	"gogue/internal/model/primitives"
	"math/rand"
)

const (
	IncreaseAttributeBaseParcentage int = 5
	IncreaseAttributeMaxPercentage  int = 20
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
	Shape primitives.Box
	Name  string
}

func (i *Item) Drop(position primitives.Point2D[int]) primitives.Box {
	i.Shape.Point = position
	return i.Shape
}

func getAttributeRandomPercentIncrease() float64 {
	return float64(IncreaseAttributeBaseParcentage + rand.Intn(IncreaseAttributeMaxPercentage+1))
}

func getAttributeRandomName(names []string) string {
	return names[rand.Intn(len(names))]
}
