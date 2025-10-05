package entity

type Attributes[T any] struct {
	MaxHealth, Agility, Strength T
}

type Character struct {
	Geometry Object
	Health   float64
	Strength uint
	Agility  uint
}
