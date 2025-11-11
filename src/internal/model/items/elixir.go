package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Elixir struct {
	*Item
	EffectDuration     uint32 // in steps
	AffectedAttributes primitives.Attributes
}

func NewElixir(rnd utils.RandomGenerator, box primitives.Box, t ElixirType) *Elixir {
	config := GetElixirConfig(t)

	return &Elixir{
		Item: &Item{
			Shape: box,
			Name:  string(config.Type),
		},
		AffectedAttributes: config.GenerateAttributes(&rnd),
		EffectDuration:     config.GenerateDuration(&rnd),
	}
}

func NewElixirWithConfig(rnd utils.RandomGenerator, box primitives.Box, config ElixirConfig) (*Elixir, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Elixir{
		Item: &Item{
			Shape: box,
			Name:  string(config.Type),
		},
		AffectedAttributes: config.GenerateAttributes(&rnd),
		EffectDuration:     config.GenerateDuration(&rnd),
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
