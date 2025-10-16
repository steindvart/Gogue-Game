package entity

import "math/rand"

const (
	TreasureBaseValue uint32 = 5
	TreasureMaxValue  uint32 = 100
)

type Treasure struct {
	Shape Box
	Value uint
	Name  string
}

func NewTreasure(box Box) *Treasure {
	return &Treasure{
		Shape: box,
		Value: getTreasureRandomValue(),
		Name:  "Gold",
	}
}

func getTreasureRandomValue() uint {
	return uint(TreasureBaseValue + rand.Uint32()%TreasureMaxValue)
}
