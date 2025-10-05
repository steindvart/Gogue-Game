package entity

type Point2D[T any] struct {
	X, Y T
}

type Size2D[T any] struct {
	Height, Width T
}

type Box struct {
	Point Point2D[int]
	Size  Size2D[uint]
}
