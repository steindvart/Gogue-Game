package items

import "gogue/internal/model/primitives"

type Scroll struct {
	Item              Item
	AffectedAttribute primitives.Attributes
}

func NewScroll(box primitives.Box) *Scroll {
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
		Item: Item{
			Shape: box,
			Name:  getAttributeRandomName(scrollNames),
		},
		AffectedAttribute: getRandomAttribute(),
	}
}

func (s *Scroll) Take() {
	s.Item.Take()
}

func (s *Scroll) Drop(box primitives.Box) {
	s.Item.Drop(box)
}

func (s *Scroll) Use() primitives.Attributes {
	return s.AffectedAttribute
}

func AsScroll(item any) *Scroll {
	s, ok := item.(*Scroll)
	if ok {
		return s
	}
	return nil
}
