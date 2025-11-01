package entity

import "fmt"

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
			Shape: box,
			Name:  getAttributeRandomName(scrollNames),
		},
		AffectedAttribute: getRandomAttribute(),
		Increment:         getAttributeRandomPercentIncrease(),
	}
}

func (s *Scroll) Taken() {
	s.Consumable.Taken()
}

func (s *Scroll) Dropped(box Box) ConsumableLike {
	s.Consumable.Dropped(box)
	return s
}

func (s *Scroll) Use(p *Player) (string, ConsumableLike) {
	if s.AffectedAttribute.MaxHealth == 1 {
		p.Character.MaxHealth += float64(s.Increment)
	}
	if s.AffectedAttribute.Agility == 1 {
		p.Character.Agility += s.Increment
	}
	if s.AffectedAttribute.Strength == 1 {
		p.Character.Strength += s.Increment
	}

	return fmt.Sprintf(
		"You read the %v, your %v has increased by %v",
		s.Consumable.Name,
		s.AffectedAttribute.GetAffectedAttributeName(),
		s.Increment,
	), nil
}
