package primitives

import "gogue/internal/utils"

type Point2D[T utils.Number] struct {
	X, Y T
}

type Positional2D[T utils.Number] interface {
	GetPosition() Point2D[T]
}

func (p *Point2D[T]) Move(delta Point2D[T]) {
	p.X += delta.X
	p.Y += delta.Y
}

type Size2D[T utils.Number] struct {
	Height, Width T
}

type Box struct {
	Point Point2D[int]
	Size  Size2D[uint]
}

func (b *Box) Move(delta Point2D[int]) {
	b.Point.Move(delta)
}

func (b *Box) GetPosition() Point2D[int] {
	return b.Point
}
