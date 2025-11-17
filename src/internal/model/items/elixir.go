package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Elixir struct {
	*Item
	Type               ElixirType
	EffectDuration     uint32 // in steps
	AffectedAttributes primitives.Attributes
}

func NewElixir(rnd *utils.RandomGenerator, box primitives.Box, t ElixirType) *Elixir {
	cfg := GetElixirConfig(t)

	return &Elixir{
		Item: &Item{
			Shape: box,
			Name:  string(cfg.Type),
		},
		Type:               t,
		AffectedAttributes: cfg.GenerateAttributes(rnd),
		EffectDuration:     cfg.GenerateDuration(rnd),
	}
}

func NewElixirByConfig(rnd *utils.RandomGenerator, box primitives.Box, cfg ElixirConfig) (*Elixir, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Elixir{
		Item: &Item{
			Shape: box,
			Name:  string(cfg.Type),
		},
		Type:               ElixirTypeCustom,
		AffectedAttributes: cfg.GenerateAttributes(rnd),
		EffectDuration:     cfg.GenerateDuration(rnd),
	}, nil
}

func (e *Elixir) Use() primitives.Attributes {
	return e.AffectedAttributes
}

func AsElixir(item any) *Elixir {
	e, ok := item.(*Elixir)
	if ok {
		return e
	}
	return nil
}
