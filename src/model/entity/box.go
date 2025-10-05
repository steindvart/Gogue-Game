package entity

type Point2D[T any] struct {
	X, Y T
}

type Size2D[T any] struct {
	Height, Width T
}

type Box struct {
	Coordinate Point2D[int]
	Measure    Size2D[uint]
}
