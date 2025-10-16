package entity

import (
	"math/rand"
	"time"
)

const (
	ElixirDurationBase      uint32 = 1
	ElixirMaxDurationFactor uint32 = 3
)

type Elixir struct {
	Shape             Box
	EffectDuration    time.Duration
	AffectedAttribute Attributes
	Increment         uint
	Name              string
}

func NewElixir(box Box) *Elixir {
	elixirNames := []string{
		"Elixir of the Jade Serpent",
		"Potion of the Phantom's Breath",
		"Vial of Crimson Vitality",
		"Draught of the Frozen Star",
		"Elixir of the Shattered Mind",
		"Potion of the Wandering Soul",
		"Vial of Ember Essence",
		"Elixir of the Obsidian Veil",
		"Potion of the Howling Wind",
	}

	return &Elixir{
		Shape:             box,
		EffectDuration:    getRandomElixirDuration(),
		AffectedAttribute: getRandomAttribute(),
		Increment:         getAttributeRandomPercentIncrease(),
		Name:              getAttributeRandomName(elixirNames),
	}
}

func getRandomElixirDuration() time.Duration {
	return time.Duration(time.Duration(ElixirDurationBase+rand.Uint32()%ElixirMaxDurationFactor) * time.Minute)
}
