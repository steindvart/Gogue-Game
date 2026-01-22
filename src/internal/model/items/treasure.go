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

func newTreasure(box primitives.Box, t TreasureType, value int32) *Treasure {
	return &Treasure{
		Item:  &Item{Box: &box, Name: string(t)},
		Type:  t,
		Value: value,
	}
}

func NewTreasureBuiltin(rnd utils.Randomizer, box primitives.Box, t TreasureType) *Treasure {
	treasure, _ := newTreasureByConfig(rnd, box, GetTreasureConfig(t))
	return treasure
}

func newTreasureByConfig(rnd utils.Randomizer, box primitives.Box, cfg TreasureConfig) (*Treasure, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return newTreasure(box, cfg.Type, cfg.GenerateValue(rnd)), nil
}

func GetTreasureType(rnd utils.Randomizer, goldPercent, gemPercent, artifactPercent, mysteryPercent uint) (TreasureType, error) {
	return getRandomTreasureType(rnd, goldPercent, gemPercent, artifactPercent, mysteryPercent)
}

func (tr *Treasure) Take() int32 {
	return tr.Value
}
