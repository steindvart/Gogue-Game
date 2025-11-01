package entity

import (
	"slices"
	"testing"
	"time"
)

func TestNewElixir(t *testing.T) {
	tests := []struct {
		name                       string
		box                        Box
		wantShape                  Box
		wantNames                  []string
		wantEffectDurationLessThan time.Duration
		wantIncrementLessThan      uint
		wantAffectedAttributesList []Attributes
	}{
		{
			name:      "elixir constructor",
			box:       Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
			wantShape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
			wantNames: []string{
				"Elixir of the Jade Serpent",
				"Potion of the Phantom's Breath",
				"Vial of Crimson Vitality",
				"Draught of the Frozen Star",
				"Elixir of the Shattered Mind",
				"Potion of the Wandering Soul",
				"Vial of Ember Essence",
				"Elixir of the Obsidian Veil",
				"Potion of the Howling Wind",
			},
			wantEffectDurationLessThan: time.Duration(ElixirDurationBase+ElixirMaxDurationFactor) * time.Minute,
			wantIncrementLessThan:      uint(IncreaseAttributeBaseParcentage + IncreaseAttributeMaxPercentage + 1),
			wantAffectedAttributesList: []Attributes{
				{
					MaxHealth: 1,
					Agility:   0,
					Strength:  0,
				},
				{
					MaxHealth: 0,
					Agility:   1,
					Strength:  0,
				},
				{
					MaxHealth: 0,
					Agility:   0,
					Strength:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewElixir(tt.box)
			if !(got.Consumable.Shape == tt.wantShape) {
				t.Errorf("NewElixir(): got Shape %v, want %v", got.Consumable.Shape, tt.wantShape)
			}
			if !(slices.Contains(tt.wantNames, got.Consumable.Name)) {
				t.Errorf("NewElixir(): got Name %v, want in %v", got.Consumable.Name, tt.wantNames)
			}
			if !(got.EffectDuration < tt.wantEffectDurationLessThan) {
				t.Errorf("NewElixir(): got EffectDuration %v, want less than %v", got.EffectDuration, tt.wantEffectDurationLessThan)
			}
			if !(got.Increment < tt.wantIncrementLessThan) {
				t.Errorf("NewElixir(): got Increment %v, want less than %v", got.Increment, tt.wantIncrementLessThan)
			}
			if !(slices.Contains(tt.wantAffectedAttributesList, got.AffectedAttribute)) {
				t.Errorf("NewElixir(): got AffectedAttribute %#v, want in %#v", got.AffectedAttribute, tt.wantAffectedAttributesList)
			}
		})
	}
}

func TestElixir_Taken(t *testing.T) {
	tests := []struct {
		name string
		want Box
	}{
		{
			name: "elixir is taken",
			want: Box{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elixir := &Elixir{
				Consumable: Consumable{
					Shape: Box{
						Point: Point2D[int]{X: 1, Y: 2},
						Size:  Size2D[uint]{Height: 1, Width: 1},
					},
					Name: "Awkward Elixir",
				},
				AffectedAttribute: Attributes{MaxHealth: 1, Agility: 0, Strength: 0},
				Increment:         1,
				EffectDuration:    time.Minute,
			}
			elixir.Taken()
			if elixir.Consumable.Shape != tt.want {
				t.Errorf("Taken() = (%v), want (%v)", elixir.Consumable.Shape, tt.want)
			}
		})
	}
}

func TestElixir_Dropped(t *testing.T) {
	tests := []struct {
		name string
		box  Box
		want Box
	}{
		{
			name: "elixir is dropped",
			box:  Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
			want: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elixir := &Elixir{
				Consumable: Consumable{
					Shape: Box{},
					Name:  "Awkward Elixir",
				},
				AffectedAttribute: Attributes{MaxHealth: 1, Agility: 0, Strength: 0},
				Increment:         1,
				EffectDuration:    time.Minute,
			}
			item := elixir.Dropped(tt.box)

			if item != elixir {
				t.Errorf("Dropped() = %#v, want %#v", item, elixir)
			}

			if elixir.Consumable.Shape != tt.want {
				t.Errorf("Dropped(): Shape %v, want %v", elixir.Consumable.Shape, tt.want)
			}
		})
	}
}

func TestElixir_Use(t *testing.T) {
	tests := []struct {
		name   string
		elixir *Elixir
		want   string
	}{
		{
			name: "use elixir",
			elixir: &Elixir{
				Consumable: Consumable{
					Shape: Box{},
					Name:  "Awkward Elixir",
				},
				AffectedAttribute: Attributes{MaxHealth: 1, Agility: 0, Strength: 0},
				Increment:         5,
				EffectDuration:    20 * time.Millisecond,
			},
			want: "You drank the Awkward Elixir, your MaxHealth has increased by 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(Box{})
			line, ptr := tt.elixir.Use(player)

			if ptr != nil {
				t.Errorf("Use(): ConsumableLike pointer got %#v, want %#v", ptr, nil)
			}

			time.Sleep(2 * time.Millisecond)

			if player.Character.MaxHealth == float64(AttributeRateAverage) {
				t.Errorf("Use(): got %#v, want %#v", player.Character.MaxHealth, float64(AttributeRateAverage)+float64(tt.elixir.Increment))
			}

			time.Sleep(30 * time.Millisecond)

			if player.Character.MaxHealth != float64(AttributeRateAverage) {
				t.Errorf("Use(): got %#v, want %#v", player.Character.MaxHealth, float64(AttributeRateAverage))
			}

			if !(line == tt.want) {
				t.Errorf("Use() = (%v), want (%v)", line, tt.want)
			}
		})
	}
}
