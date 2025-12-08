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
		Item:   &Item{Shape: box, Name: string(t)},
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

func (e *Scroll) Use() *primitives.Effect {
	return e.Effect
}

func AsScroll(item any) *Scroll {
	s, ok := item.(*Scroll)
	if ok {
		return s
	}
	return nil
}
