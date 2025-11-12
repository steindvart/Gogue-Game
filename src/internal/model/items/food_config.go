package items

import (
	"fmt"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type FoodType string

const (
	// Random health food
	FoodTypePotatoes FoodType = "Potatoes" // +Health(1-10)
	FoodTypeBread    FoodType = "Bread"    // +Health(10-20)
	FoodTypeMeat     FoodType = "Meat"     // +Health(20-30)
	FoodTypeMistery  FoodType = "Mistery"  // +Health(-10-40)

	// No random health food
	FoodTypeBeer FoodType = "Beer" // +Health(20-20)
)

type FoodConfig struct {
	Type        FoodType
	HealthRange primitives.AttributeRange
	Description string
}

var FoodRegistry = map[FoodType]FoodConfig{
	// Random health food
	FoodTypePotatoes: {
		Type:        FoodTypePotatoes,
		HealthRange: primitives.AttributeRange{Min: 1, Max: 10},
		Description: "Increases health by a small amount",
	},
	FoodTypeBread: {
		Type:        FoodTypeBread,
		HealthRange: primitives.AttributeRange{Min: 10, Max: 20},
		Description: "Increases health by a medium amount",
	},
	FoodTypeMeat: {
		Type:        FoodTypeMeat,
		HealthRange: primitives.AttributeRange{Min: 20, Max: 30},
		Description: "Increases health by a high amount",
	},
	FoodTypeMistery: {
		Type:        FoodTypeMistery,
		HealthRange: primitives.AttributeRange{Min: -10, Max: 40},
		Description: "Increases health by a random amount (or maybe decreases it...)",
	},

	// No random health food
	FoodTypeBeer: {
		Type:        FoodTypeBeer,
		HealthRange: primitives.AttributeRange{Min: 20, Max: 20},
		Description: "Increases health by a medium amount without random and makes you feel good",
	},
}

func GetFoodConfig(t FoodType) FoodConfig {
	if config, exists := FoodRegistry[t]; exists {
		return config
	}

	return FoodRegistry[FoodTypeMistery]
}

func (cfg *FoodConfig) GenerateAttributes(rng *utils.RandomGenerator) primitives.Attributes {
	return primitives.Attributes{
		Health: float64(utils.RandomRoundedFloatInRange(rng, cfg.HealthRange.Min, cfg.HealthRange.Max)),
	}
}

func (cfg *FoodConfig) Validate() error {
	if cfg.HealthRange.Min > cfg.HealthRange.Max {
		return fmt.Errorf("invalid health range: min=%.2f > max=%.2f", cfg.HealthRange.Min, cfg.HealthRange.Max)
	}

	return nil
}
