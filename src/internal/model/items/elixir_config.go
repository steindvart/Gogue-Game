package items

import (
	"fmt"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

const (
	ElixirDurationMin    uint32 = 5
	ElixirMaxDurationMax uint32 = 30
)

type Range[T any] struct {
	Min T
	Max T
}

type ElixirAttributeRange = Range[float64]
type ElixirDurationStepsRange = Range[uint32]

var defaultDurationRange = ElixirDurationStepsRange{Min: ElixirDurationMin, Max: ElixirMaxDurationMax}

type ElixirType string

const (
	ElixirTypeStrength      ElixirType = "Elixir of Strength" // +Strength
	ElixirTypeAgility       ElixirType = "Elixir of Agility"  // +Agility
	ElixirTypePhantomBreath ElixirType = "Phantom's Breath"   // +Agility -Strength
	ElixirTypeFrozenStar    ElixirType = "Frozen Star"        // +Strength -Agility
	ElixirTypeMystery       ElixirType = "Elixir of Mystery"  // All random
)

type ElixirConfig struct {
	Type               ElixirType
	StrengthRange      ElixirAttributeRange
	AgilityRange       ElixirAttributeRange
	DurationStepsRange ElixirDurationStepsRange
	Description        string
}

var ElixirRegistry = map[ElixirType]ElixirConfig{
	ElixirTypeStrength: {
		Type:               ElixirTypeStrength,
		StrengthRange:      ElixirAttributeRange{Min: 5, Max: 20},
		AgilityRange:       ElixirAttributeRange{Min: 0, Max: 0},
		DurationStepsRange: defaultDurationRange,
		Description:        "Increases strength for a short time",
	},
	ElixirTypeAgility: {
		Type:               ElixirTypeAgility,
		StrengthRange:      ElixirAttributeRange{Min: 0, Max: 0},
		AgilityRange:       ElixirAttributeRange{Min: 5, Max: 20},
		DurationStepsRange: defaultDurationRange,
		Description:        "Increases agility for a short time",
	},
	ElixirTypePhantomBreath: {
		Type:               ElixirTypePhantomBreath,
		StrengthRange:      ElixirAttributeRange{Min: -10, Max: -2},
		AgilityRange:       ElixirAttributeRange{Min: 10, Max: 30},
		DurationStepsRange: defaultDurationRange,
		Description:        "Greatly increases agility but weakens strength",
	},
	ElixirTypeFrozenStar: {
		Type:               ElixirTypeFrozenStar,
		StrengthRange:      ElixirAttributeRange{Min: 10, Max: 30},
		AgilityRange:       ElixirAttributeRange{Min: -10, Max: -2},
		DurationStepsRange: defaultDurationRange,
		Description:        "Greatly increases strength but reduces agility",
	},
	ElixirTypeMystery: {
		Type:               ElixirTypeMystery,
		StrengthRange:      ElixirAttributeRange{Min: -20, Max: 30},
		AgilityRange:       ElixirAttributeRange{Min: -20, Max: 30},
		DurationStepsRange: defaultDurationRange,
		Description:        "Random effects on both strength and agility",
	},
}

func GetElixirConfig(t ElixirType) ElixirConfig {
	if config, exists := ElixirRegistry[t]; exists {
		return config
	}

	return ElixirRegistry[ElixirTypeMystery]
}

func (cfg *ElixirConfig) GenerateAttributes(rng *utils.RandomGenerator) primitives.Attributes {
	return primitives.Attributes{
		Strength: float64(utils.RandomRoundedFloatInRange(rng, cfg.StrengthRange.Min, cfg.StrengthRange.Max)),
		Agility:  float64(utils.RandomRoundedFloatInRange(rng, cfg.AgilityRange.Min, cfg.AgilityRange.Max)),
	}
}

func (cfg *ElixirConfig) GenerateDuration(rng *utils.RandomGenerator) uint32 {
	return uint32(utils.RandomIntInRange(rng, int(cfg.DurationStepsRange.Min), int(cfg.DurationStepsRange.Max)))
}

func (cfg *ElixirConfig) Validate() error {
	if cfg.StrengthRange.Min > cfg.StrengthRange.Max {
		return fmt.Errorf("invalid strength range: min=%.2f > max=%.2f", cfg.StrengthRange.Min, cfg.StrengthRange.Max)
	}
	if cfg.AgilityRange.Min > cfg.AgilityRange.Max {
		return fmt.Errorf("invalid agility range: min=%.2f > max=%.2f", cfg.AgilityRange.Min, cfg.AgilityRange.Max)
	}
	if cfg.DurationStepsRange.Min > cfg.DurationStepsRange.Max {
		return fmt.Errorf("invalid duration range: min=%d > max=%d", cfg.DurationStepsRange.Min, cfg.DurationStepsRange.Max)
	}

	return nil
}
