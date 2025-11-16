package items

import (
	"gogue/internal/model/primitives"
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
