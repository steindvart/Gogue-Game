package entity

type Scroll struct {
	Shape             Box
	AffectedAttribute Attributes
	Increment         uint
	Name              string
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
		Shape:             box,
		AffectedAttribute: getRandomAttribute(),
		Increment:         getAttributeRandomPercentIncrease(),
		Name:              getAttributeRandomName(scrollNames),
	}
}
