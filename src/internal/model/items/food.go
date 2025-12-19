package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Food struct {
	*Item
	*primitives.Effect
	Type FoodType
}

func NewFood(box primitives.Box, t FoodType, e *primitives.Effect) *Food {
	return &Food{
		Item:   &Item{Box: box, Name: string(t)},
		Effect: e,
		Type:   t,
	}
}

func NewFoodBuiltin(rnd utils.Randomizer, box primitives.Box, t FoodType) *Food {
	f, _ := NewFoodByConfig(rnd, box, GetFoodConfig(t))
	return f
}

func NewFoodByConfig(rnd utils.Randomizer, box primitives.Box, cfg FoodConfig) (*Food, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return NewFood(box, cfg.Type, &primitives.Effect{
		Attributes: cfg.GenerateAttributes(rnd),
	}), nil
}
