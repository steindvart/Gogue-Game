package entity

type Coordinate2D[T any] struct {
	X, Y T
}

type Measure2D[T any] struct {
	Height, Width T
}

type Object struct {
	Coordinate Coordinate2D[int]
	Measure    Measure2D[uint]
}
