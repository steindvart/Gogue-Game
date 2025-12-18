package utils

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 | uint | uint8 | uint16 | uint32 | uint64
}

func CreateEmpty2DSlice[T any](rows, cols int) [][]T {
	field := make([][]T, rows)
	for y := range field {
		field[y] = make([]T, cols)
	}

	return field
}
