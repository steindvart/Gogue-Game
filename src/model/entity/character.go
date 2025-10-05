package entity

type Attributes struct {
	MaxHealth, Agility, Strength uint
}

type Character struct {
	Shape    Box
	Health   float64
	Strength uint
	Agility  uint
}
