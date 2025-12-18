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

func NewTreasure(box primitives.Box, t TreasureType, value int32) *Treasure {
	return &Treasure{
		Item:  &Item{Box: box, Name: string(t)},
		Type:  t,
		Value: value,
	}
}

func NewTreasureBuiltin(rnd *utils.Random, box primitives.Box, t TreasureType) *Treasure {
	treasure, _ := NewTreasureByConfig(rnd, box, GetTreasureConfig(t))
	return treasure
}

func NewTreasureByConfig(rnd *utils.Random, box primitives.Box, cfg TreasureConfig) (*Treasure, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return NewTreasure(box, cfg.Type, cfg.GenerateValue(rnd)), nil
}

func (tr *Treasure) Take() int32 {
	return tr.Value
}
