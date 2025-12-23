package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Scroll struct {
	*Item
	*primitives.Effect
	Type ScrollType
}

func NewScroll(box primitives.Box, t ScrollType, e *primitives.Effect) *Scroll {
	return &Scroll{
		Item:   &Item{Box: &box, Name: string(t)},
		Effect: e,
		Type:   t,
	}
}

func NewScrollBuiltin(rnd utils.Randomizer, box primitives.Box, t ScrollType) *Scroll {
	s, _ := NewScrollByConfig(rnd, box, GetScrollConfig(t))
	return s
}

func NewScrollByConfig(rnd utils.Randomizer, box primitives.Box, cfg ScrollConfig) (*Scroll, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return NewScroll(box, cfg.Type, &primitives.Effect{
		Attributes: cfg.GenerateAttributes(rnd),
	}), nil
}

func (s *Scroll) Use() *primitives.Effect {
	return s.Effect
}

// Ничего не возвращаем, но добавляем функцию, чтобы отметить Scroll как Takeable.
func (s *Scroll) Take() int32 {
	return 0
}
