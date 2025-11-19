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

func NewScroll(rnd *utils.RandomGenerator, box primitives.Box, t ScrollType) *Scroll {
	cfg := GetScrollConfig(t)

	return &Scroll{
		Item: &Item{
			Shape: box,
			Name:  string(cfg.Type),
		},
		Effect: &primitives.Effect{
			Attributes: cfg.GenerateAttributes(rnd),
			Duration:   0,
		},
		Type: t,
	}
}

func NewScrollByConfig(rnd *utils.RandomGenerator, box primitives.Box, cfg ScrollConfig) (*Scroll, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Scroll{
		Item: &Item{
			Shape: box,
			Name:  string(cfg.Type),
		},
		Effect: &primitives.Effect{
			Attributes: cfg.GenerateAttributes(rnd),
			Duration:   0,
		},
		Type: ScrollTypeCustom,
	}, nil
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
