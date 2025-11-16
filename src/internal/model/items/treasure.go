package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Treasure struct {
	*Item
	Type  TreasureType
	Value int32
}

func NewTreasure(rnd *utils.RandomGenerator, box primitives.Box, t TreasureType) *Treasure {
	cfg := GetTreasureConfig(t)

	return &Treasure{
		Item: &Item{
			Shape: box,
			Name:  string(cfg.Type),
		},
		Type:  t,
		Value: cfg.GenerateValue(rnd),
	}
}

func NewTreasureByConfig(rnd *utils.RandomGenerator, box primitives.Box, cfg TreasureConfig) (*Treasure, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Treasure{
		Item: &Item{
			Shape: box,
			Name:  string(cfg.Type),
		},
		Type:  TreasureTypeCustom,
		Value: cfg.GenerateValue(rnd),
	}, nil
}

func (tr *Treasure) Take() int32 {
	return tr.Value
}

func AsTreasure(item any) *Treasure {
	s, ok := item.(*Treasure)
	if ok {
		return s
	}
	return nil
}
