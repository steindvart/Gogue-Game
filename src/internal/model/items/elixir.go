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

func NewElixir(box primitives.Box, t ElixirType, e *primitives.Effect) *Elixir {
	return &Elixir{
		Item: &Item{
			Box:  &box,
			Name: string(t),
		},
		Effect: e,
		Type:   t,
	}
}

func NewElixirBuiltin(rnd utils.Randomizer, box primitives.Box, t ElixirType) *Elixir {
	e, _ := NewElixirByConfig(rnd, box, GetElixirConfig(t))
	return e
}

func NewElixirByConfig(rnd utils.Randomizer, box primitives.Box, cfg ElixirConfig) (*Elixir, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return NewElixir(box, cfg.Type, &primitives.Effect{
		Duration: primitives.EffectDuration{
			Type:  primitives.EffectDurationTypeAllTemporaryHealPermanent,
			Steps: cfg.GenerateDuration(rnd),
		},
		Attributes: cfg.GenerateAttributes(rnd),
	}), nil
}

func (e *Elixir) Use() *primitives.Effect {
	return e.Effect
}
