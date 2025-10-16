package entity

type Scroll struct {
	Consumable        Consumable
	AffectedAttribute Attributes
	Increment         uint
}

func NewScroll(box Box) *Scroll {
	scrollNames := []string{
		"Scroll of Shadowstep",
		"Parchment of Eternal Flame",
		"Manuscript of Forgotten Truths",
		"Scroll of Iron Will",
		"Vellum of the Void",
		"Scroll of Whispers",
		"Tome of the Lost King",
		"Scroll of Unseen Paths",
		"Parchment of Thunderous Roar",
	}

	return &Scroll{
		Consumable: Consumable{
			Shape:             box,
			Name:              getAttributeRandomName(scrollNames),
		},
		AffectedAttribute: getRandomAttribute(),
		Increment:         getAttributeRandomPercentIncrease(),
	}
}
