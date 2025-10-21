package entity

import "math/rand"

const (
	TreasureBaseValue uint32 = 5
	TreasureMaxValue  uint32 = 100
)

type Treasure struct {
	Shape Box
	Name  string
	Value      uint
}

func NewTreasure(box Box) *Treasure {
	return &Treasure{
		Shape: box,
		Name:  "Gold",
		Value: getTreasureRandomValue(),
	}
}

func getTreasureRandomValue() uint {
	return uint(TreasureBaseValue + rand.Uint32()%TreasureMaxValue)
}
