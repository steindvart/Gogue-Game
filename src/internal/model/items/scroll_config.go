package items

import (
	"fmt"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type ScrollType string

const (
	ScrollTypeStrength  ScrollType = "Strength"   // +Strength
	ScrollTypeAgility   ScrollType = "Agility"    // +Agility
	ScrollTypeUltimate  ScrollType = "Ultimate"   // +Agility +Strength
	ScrollTypeMaxHealth ScrollType = "Max Health" // +MaxHealth
	ScrollTypeMystery   ScrollType = "Mystery"    // All random
	ScrollTypeCustom    ScrollType = "Custom"     // For custom scroll - dynamicly or from external data created
)

type ScrollConfig struct {
	Type           ScrollType
	StrengthRange  primitives.AttributeRange
	AgilityRange   primitives.AttributeRange
	MaxHealthRange primitives.AttributeRange
	Description    string
}

var ScrollRegistry = map[ScrollType]ScrollConfig{
	ScrollTypeStrength: {
		Type:           ScrollTypeStrength,
		StrengthRange:  primitives.AttributeRange{Min: 1, Max: 3},
		AgilityRange:   primitives.AttributeRange{Min: 0, Max: 0},
		MaxHealthRange: primitives.AttributeRange{Min: 0, Max: 0},
		Description:    "Increases strength",
	},
	ScrollTypeAgility: {
		Type:           ScrollTypeAgility,
		StrengthRange:  primitives.AttributeRange{Min: 0, Max: 0},
		AgilityRange:   primitives.AttributeRange{Min: 1, Max: 3},
		MaxHealthRange: primitives.AttributeRange{Min: 0, Max: 0},
		Description:    "Increases agility",
	},
	ScrollTypeMaxHealth: {
		Type:           ScrollTypeMaxHealth,
		StrengthRange:  primitives.AttributeRange{Min: 0, Max: 0},
		AgilityRange:   primitives.AttributeRange{Min: 0, Max: 0},
		MaxHealthRange: primitives.AttributeRange{Min: 2, Max: 10},
		Description:    "Increases maximus health",
	},
	ScrollTypeMystery: {
		Type:           ScrollTypeMystery,
		StrengthRange:  primitives.AttributeRange{Min: -3, Max: 8},
		AgilityRange:   primitives.AttributeRange{Min: -3, Max: 8},
		MaxHealthRange: primitives.AttributeRange{Min: -3, Max: 8},
		Description:    "Increases or decreases random attributes",
	},
}

func GetScrollConfig(t ScrollType) ScrollConfig {
	if cfg, exists := ScrollRegistry[t]; exists {
		return cfg
	}

	return ScrollRegistry[ScrollTypeMystery]
}

func (cfg *ScrollConfig) GenerateAttributes(rng utils.Randomizer) primitives.Attributes {
	return primitives.Attributes{
		Strength:  float64(utils.RandomRoundedFloatInRange(rng, cfg.StrengthRange.Min, cfg.StrengthRange.Max)),
		Agility:   float64(utils.RandomRoundedFloatInRange(rng, cfg.AgilityRange.Min, cfg.AgilityRange.Max)),
		MaxHealth: float64(utils.RandomRoundedFloatInRange(rng, cfg.MaxHealthRange.Min, cfg.MaxHealthRange.Max)),
	}
}

func (cfg *ScrollConfig) Validate() error {
	if cfg.StrengthRange.Min > cfg.StrengthRange.Max {
		return fmt.Errorf("invalid strength range: min=%.2f > max=%.2f", cfg.StrengthRange.Min, cfg.StrengthRange.Max)
	}
	if cfg.AgilityRange.Min > cfg.AgilityRange.Max {
		return fmt.Errorf("invalid agility range: min=%.2f > max=%.2f", cfg.AgilityRange.Min, cfg.AgilityRange.Max)
	}
	if cfg.MaxHealthRange.Min > cfg.MaxHealthRange.Max {
		return fmt.Errorf("invalid duration range: min=%.2f > max=%.2f", cfg.MaxHealthRange.Min, cfg.MaxHealthRange.Max)
	}

	return nil
}
