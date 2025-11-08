package primitives

type Attributes struct {
	MaxHealth, Health, Agility, Strength float64
}

func (a *Attributes) GetAffectedAttributeName() string {
	if a.MaxHealth == 1 && a.Agility == 0 && a.Strength == 0 {
		return "MaxHealth"
	} else if a.Agility == 1 && a.MaxHealth == 0 && a.Strength == 0 {
		return "Agility"
	} else if a.Strength == 1 && a.MaxHealth == 0 && a.Agility == 0 {
		return "Strength"
	} else {
		return "None"
	}
}
