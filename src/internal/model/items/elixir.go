package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Elixir struct {
	*Item
	*primitives.Effect
	Type ElixirType
}

func NewElixir(rnd *utils.RandomGenerator, box primitives.Box, t ElixirType) *Elixir {
	cfg := GetElixirConfig(t)

	return &Elixir{
		Item: &Item{
			Shape: box,
			Name:  string(cfg.Type),
		},
		Effect: &primitives.Effect{
			Duration: primitives.EffectDuration{
				Type:  primitives.EffectDurationTypeAllTemporaryHealPermanent,
				Steps: cfg.GenerateDuration(rnd),
			},
			Attributes: cfg.GenerateAttributes(rnd),
		},
		Type: t,
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
		Effect: &primitives.Effect{
			Duration: primitives.EffectDuration{
				Type:  primitives.EffectDurationTypeAllTemporaryHealPermanent,
				Steps: cfg.GenerateDuration(rnd),
			},
			Attributes: cfg.GenerateAttributes(rnd),
		},
		Type: ElixirTypeCustom,
	}, nil
}

func (e *Elixir) Use() *primitives.Effect {
	return e.Effect
}

func AsElixir(item any) *Elixir {
	e, ok := item.(*Elixir)
	if ok {
		return e
	}
	return nil
}
