package primitives

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 | uint | uint8 | uint16 | uint32 | uint64
}

type Point2D[T Number] struct {
	X, Y T
}

func (p *Point2D[T]) Move(delta Point2D[T]) {
	p.X += delta.X
	p.Y += delta.Y
}

type Size2D[T Number] struct {
	Height, Width T
}

type Box struct {
	Point Point2D[int]
	Size  Size2D[uint]
}

func (b *Box) Move(delta Point2D[int]) {
	b.Point.Move(delta)
}
