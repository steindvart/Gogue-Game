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
	Item              Item
	EffectDuration    time.Duration
	AffectedAttribute Attributes
	Increment         uint
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
		Item: Item{
			Shape: box,
			Name:  getAttributeRandomName(elixirNames),
		},
		AffectedAttribute: getRandomAttribute(),
		Increment:         getAttributeRandomPercentIncrease(),
		EffectDuration:    getRandomElixirDuration(),
	}
}

func (e *Elixir) Taken() {
	e.Item.Taken()
}

func (e *Elixir) Dropped(box Box) {
	e.Item.Dropped(box)
}

func (e *Elixir) Use(p *Player) {
	go func() {
		defer func() {
			p.Character.MaxHealth -= float64(e.Increment)
			p.Character.Agility -= e.Increment
			p.Character.Strength -= e.Increment
		}()

		if e.AffectedAttribute.MaxHealth == 1 {
			p.Character.MaxHealth += float64(e.Increment)
		}
		if e.AffectedAttribute.Agility == 1 {
			p.Character.Agility += e.Increment
		}
		if e.AffectedAttribute.Strength == 1 {
			p.Character.Strength += e.Increment
		}

		time.Sleep(e.EffectDuration)
	}()
}

func getRandomElixirDuration() time.Duration {
	return time.Duration(time.Duration(ElixirDurationBase+rand.Uint32()%ElixirMaxDurationFactor) * time.Minute)
}
