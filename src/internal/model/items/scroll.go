package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Scroll struct {
	*Item
	Type               ScrollType
	AffectedAttributes primitives.Attributes
}

func NewScroll(rnd *utils.RandomGenerator, box primitives.Box, t ScrollType) *Scroll {
	cfg := GetScrollConfig(t)

	return &Scroll{
		Item: &Item{
			Shape: box,
			Name:  string(cfg.Type),
		},
		Type:               t,
		AffectedAttributes: cfg.GenerateAttributes(rnd),
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
		Type:               ScrollTypeCustom,
		AffectedAttributes: cfg.GenerateAttributes(rnd),
	}, nil
}

func (e *Scroll) Use() primitives.Attributes {
	return e.AffectedAttributes
}

func AsScroll(item any) *Scroll {
	s, ok := item.(*Scroll)
	if ok {
		return s
	}
	return nil
}
