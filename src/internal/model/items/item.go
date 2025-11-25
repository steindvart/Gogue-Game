package items

import (
	"gogue/internal/model/primitives"
)

type NotItemError struct{}

func (NotItemError) Error() string {
	return "the provided value is not an item"
}

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
