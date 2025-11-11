package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

const (
	ElixirDurationMin    int = 5
	ElixirMaxDurationMax int = 30
)

type ElixirType string

const (
	ElixirTypeStrength      = "Elixir of Strength" // +Strength
	ElixirTypeAgility       = "Elixir of Agility"  // +Agility
	ElixirTypePhantomBreath = "Phantom's Breath"   // +Agility -Strength
	ElixirTypeFrozenStar    = "Frozen Star"        // +Strength -Agility
	ElixirTypeMystery       = "Elixir of Mystery"  // All random
)

type Elixir struct {
	*Item
	EffectDuration     uint32 // in steps
	AffectedAttributes primitives.Attributes
}

func NewElixir(rnd utils.RandomGenerator, box primitives.Box, t ElixirType) *Elixir {
	switch t {
	case ElixirTypeStrength:
		return createElixir(box, t, primitives.Attributes{
			Strength: float64(utils.RandomRoundedFloatInRange(&rnd, 5, 20)),
		}, uint32(utils.RandomIntInRange(&rnd, ElixirDurationMin, ElixirMaxDurationMax)))
	case ElixirTypeAgility:
		return createElixir(box, t, primitives.Attributes{
			Agility: float64(utils.RandomRoundedFloatInRange(&rnd, 5, 20)),
		}, uint32(utils.RandomIntInRange(&rnd, ElixirDurationMin, ElixirMaxDurationMax)))
	case ElixirTypePhantomBreath:
		return createElixir(box, t, primitives.Attributes{
			Agility:  float64(utils.RandomRoundedFloatInRange(&rnd, 10, 30)),
			Strength: float64(utils.RandomRoundedFloatInRange(&rnd, -10, -2)),
		}, uint32(utils.RandomIntInRange(&rnd, ElixirDurationMin, ElixirMaxDurationMax)))
	case ElixirTypeFrozenStar:
		return createElixir(box, t, primitives.Attributes{
			Agility:  float64(utils.RandomRoundedFloatInRange(&rnd, -10, -2)),
			Strength: float64(utils.RandomRoundedFloatInRange(&rnd, 10, 30)),
		}, uint32(utils.RandomIntInRange(&rnd, ElixirDurationMin, ElixirMaxDurationMax)))
	case ElixirTypeMystery:
		fallthrough
	default:
		return createElixir(box, ElixirTypeMystery, primitives.Attributes{
			Agility:  float64(utils.RandomRoundedFloatInRange(&rnd, -20, 30)),
			Strength: float64(utils.RandomRoundedFloatInRange(&rnd, -20, 30)),
		}, uint32(utils.RandomIntInRange(&rnd, ElixirDurationMin, ElixirMaxDurationMax)))
	}
}

func createElixir(box primitives.Box, t ElixirType, affected primitives.Attributes, duration uint32) *Elixir {
	return &Elixir{
		Item: &Item{
			Shape: box,
			Name:  string(t),
		},
		AffectedAttributes: affected,
		EffectDuration:     duration,
	}
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
