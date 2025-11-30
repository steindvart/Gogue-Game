package primitives

type Range[T any] struct {
	Min T
	Max T
}

type AttributeRange = Range[float64]

type Attributes struct {
	MaxHealth, Health, Agility, Strength float64
}

func (a *Attributes) Affect(delta Attributes) {
	a.MaxHealth += delta.MaxHealth
	a.Health += delta.Health
	a.Agility += delta.Agility
	a.Strength += delta.Strength
}

func Inverse(a Attributes) Attributes {
	return Attributes{
		MaxHealth: -a.MaxHealth,
		Health:    -a.Health,
		Agility:   -a.Agility,
		Strength:  -a.Strength,
	}
}
