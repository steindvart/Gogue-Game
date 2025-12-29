package items

import (
	"gogue/internal/model/primitives"
)

type NotItemError struct{}

func (NotItemError) Error() string {
	return "the provided value is not an item"
}

type Usable interface {
	Use() *primitives.Effect
}

type Takeable interface {
	Take() int32
}

type Dropable interface {
	Drop(position primitives.Point2D[int]) primitives.Positional2D[int]
}

type Item struct {
	*primitives.Box
	Name string
}

func (i *Item) Drop(position primitives.Point2D[int]) primitives.Positional2D[int] {
	i.Box.Point = position
	return i
}
