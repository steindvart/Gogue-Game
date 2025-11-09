package items

import (
	"gogue/internal/model/primitives"
	"math/rand"
	"time"
)

const (
	ElixirDurationBase      uint32 = 1
	ElixirMaxDurationFactor uint32 = 3
)

type Elixir struct {
	Item              Item
	EffectDuration    time.Duration
	AffectedAttribute primitives.Attributes
}

func NewElixir(box primitives.Box) *Elixir {
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
		Item: Item{
			Shape: box,
			Name:  getAttributeRandomName(elixirNames),
		},
		AffectedAttribute: getRandomAttribute(),
		EffectDuration:    getRandomElixirDuration(),
	}
}

func (e *Elixir) Take() {
	e.Item.Take()
}

func (e *Elixir) Drop(box primitives.Box) {
	e.Item.Drop(box)
}

func (e *Elixir) Use() primitives.Attributes {
	return e.AffectedAttribute
}

func AsElixir(item any) *Elixir {
	e, ok := item.(*Elixir)
	if ok {
		return e
	}
	return nil
}

func getRandomElixirDuration() time.Duration {
	return time.Duration(time.Duration(ElixirDurationBase+rand.Uint32()%ElixirMaxDurationFactor) * time.Minute)
}
