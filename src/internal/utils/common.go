package utils

func CreateEmpty2DSlice[T any](rows, cols int) [][]T {
	field := make([][]T, rows)
	for y := range field {
		field[y] = make([]T, cols)
	}

	return field
}
