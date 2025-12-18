package items

import (
	"gogue/internal/model/primitives"
)

type NotItemError struct{}

func (NotItemError) Error() string {
	return "the provided value is not an item"
}

type Dropable interface {
	Drop(position primitives.Point2D[int]) primitives.Box
}

type ItemLike interface {
	primitives.Positional2D[int]
	Dropable
}

type Item struct {
	primitives.Box
	Name string
}

func (i *Item) Drop(position primitives.Point2D[int]) primitives.Box {
	i.Box.Point = position
	return i.Box
}
