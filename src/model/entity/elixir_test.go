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
				EffectDuration:    time.Minute,
				AffectedAttribute: Attributes{MaxHealth: 1, Agility: 0, Strength: 0},
				Increment: 1,
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
			name: "consumable is dropped",
			box:  Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
			want: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
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
				EffectDuration:    time.Minute,
				AffectedAttribute: Attributes{MaxHealth: 1, Agility: 0, Strength: 0},
				Increment: 1,
			}
			elixir.Dropped(tt.box)
			if elixir.Consumable.Shape != tt.want {
				t.Errorf("Dropped() = (%v), want (%v)", elixir.Consumable.Shape, tt.want)
			}
		})
	}
}

func TestElixir_Use(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "use elixir",
			want: "Awkward Elixir",
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
				EffectDuration:    time.Minute,
				AffectedAttribute: Attributes{MaxHealth: 1, Agility: 0, Strength: 0},
				Increment: 1,
			}
			name := elixir.Use()
			if name != tt.want {
				t.Errorf("Use() = (%v), want (%v)", name, tt.want)
			}
		})
	}
}
