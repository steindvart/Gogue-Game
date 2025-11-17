package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Food struct {
	*Item
	Type               FoodType
	AffectedAttributes primitives.Attributes
}

func NewFood(rnd *utils.RandomGenerator, box primitives.Box, t FoodType) *Food {
	cfg := GetFoodConfig(t)

	return &Food{
		Item: &Item{
			Shape: box,
			Name:  string(cfg.Type),
		},
		Type:               t,
		AffectedAttributes: cfg.GenerateAttributes(rnd),
	}
}

func NewFoodByConfig(rnd *utils.RandomGenerator, box primitives.Box, cfg FoodConfig) (*Food, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Food{
		Item: &Item{
			Shape: box,
			Name:  string(cfg.Type),
		},
		Type:               FoodTypeCustom,
		AffectedAttributes: cfg.GenerateAttributes(rnd),
	}, nil
}

func AsFood(item any) *Food {
	f, ok := item.(*Food)
	if ok {
		return f
	}
	return nil
}
