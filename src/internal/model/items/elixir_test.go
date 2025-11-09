package items

// import (
// 	"gogue/internal/model/entity"
// 	"gogue/internal/model/primitive"
// 	"slices"
// 	"testing"
// 	"time"
// )

// func TestNewElixir(t *testing.T) {
// 	tests := []struct {
// 		name                       string
// 		box                        primitives.Box
// 		wantShape                  primitives.Box
// 		wantNames                  []string
// 		wantEffectDurationLessThan time.Duration
// 		wantIncrementLessThan      uint
// 		wantAffectedAttributesList []primitives.Attributes
// 	}{
// 		{
// 			name:      "elixir constructor",
// 			box:       primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 1}},
// 			wantShape: primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 1}},
// 			wantNames: []string{
// 				"Elixir of the Jade Serpent",
// 				"Potion of the Phantom's Breath",
// 				"Vial of Crimson Vitality",
// 				"Draught of the Frozen Star",
// 				"Elixir of the Shattered Mind",
// 				"Potion of the Wandering Soul",
// 				"Vial of Ember Essence",
// 				"Elixir of the Obsidian Veil",
// 				"Potion of the Howling Wind",
// 			},
// 			wantEffectDurationLessThan: time.Duration(ElixirDurationBase+ElixirMaxDurationFactor) * time.Minute,
// 			wantIncrementLessThan:      uint(IncreaseAttributeBaseParcentage + IncreaseAttributeMaxPercentage + 1),
// 			wantAffectedAttributesList: []primitives.Attributes{
// 				{
// 					MaxHealth: 1,
// 					Agility:   0,
// 					Strength:  0,
// 				},
// 				{
// 					MaxHealth: 0,
// 					Agility:   1,
// 					Strength:  0,
// 				},
// 				{
// 					MaxHealth: 0,
// 					Agility:   0,
// 					Strength:  1,
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := NewElixir(tt.box)
// 			if !(got.Item.Shape == tt.wantShape) {
// 				t.Errorf("NewElixir(): got Shape %v, want %v", got.Item.Shape, tt.wantShape)
// 			}
// 			if !(slices.Contains(tt.wantNames, got.Item.Name)) {
// 				t.Errorf("NewElixir(): got Name %v, want in %v", got.Item.Name, tt.wantNames)
// 			}
// 			if !(got.EffectDuration < tt.wantEffectDurationLessThan) {
// 				t.Errorf("NewElixir(): got EffectDuration %v, want less than %v", got.EffectDuration, tt.wantEffectDurationLessThan)
// 			}
// 			if !(got.Increment < tt.wantIncrementLessThan) {
// 				t.Errorf("NewElixir(): got Increment %v, want less than %v", got.Increment, tt.wantIncrementLessThan)
// 			}
// 			if !(slices.Contains(tt.wantAffectedAttributesList, got.AffectedAttribute)) {
// 				t.Errorf("NewElixir(): got AffectedAttribute %#v, want in %#v", got.AffectedAttribute, tt.wantAffectedAttributesList)
// 			}
// 		})
// 	}
// }

// func TestElixir_Take(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want primitives.Box
// 	}{
// 		{
// 			name: "elixir is taken",
// 			want: primitives.Box{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			elixir := &Elixir{
// 				Item: Item{
// 					Shape: primitives.Box{
// 						Point: primitives.Point2D[int]{X: 1, Y: 2},
// 						Size:  primitives.Size2D[uint]{Height: 1, Width: 1},
// 					},
// 					Name: "Awkward Elixir",
// 				},
// 				AffectedAttribute: primitives.Attributes{MaxHealth: 1, Agility: 0, Strength: 0},
// 				Increment:         1,
// 				EffectDuration:    time.Minute,
// 			}
// 			elixir.Take()
// 			if elixir.Item.Shape != tt.want {
// 				t.Errorf("Taken() = (%v), want (%v)", elixir.Item.Shape, tt.want)
// 			}
// 		})
// 	}
// }

// func TestElixir_Drop(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		box  primitives.Box
// 		want primitives.Box
// 	}{
// 		{
// 			name: "elixir is dropped",
// 			box:  primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 1}},
// 			want: primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 1}},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			elixir := &Elixir{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Elixir",
// 				},
// 				AffectedAttribute: primitives.Attributes{MaxHealth: 1, Agility: 0, Strength: 0},
// 				Increment:         1,
// 				EffectDuration:    time.Minute,
// 			}

// 			elixir.Drop(tt.box)

// 			if elixir.Item.Shape != tt.want {
// 				t.Errorf("Dropped(): Shape %v, want %v", elixir.Item.Shape, tt.want)
// 			}
// 		})
// 	}
// }

// func TestElixir_Use(t *testing.T) {
// 	tests := []struct {
// 		name   string
// 		elixir *Elixir
// 	}{
// 		{
// 			name: "use elixir",
// 			elixir: &Elixir{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Elixir",
// 				},
// 				AffectedAttribute: primitives.Attributes{MaxHealth: 1, Agility: 0, Strength: 0},
// 				Increment:         5,
// 				EffectDuration:    20 * time.Millisecond,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			player := entities.NewPlayer(primitives.Box{})
// 			tt.elixir.Use(player)

// 			time.Sleep(2 * time.Millisecond)

// 			if player.Character.MaxHealth == float64(entities.AttributeRateAverage) {
// 				t.Errorf("Use(): got %#v, want %#v", player.Character.MaxHealth, float64(entities.AttributeRateAverage)+float64(tt.elixir.Increment))
// 			}

// 			time.Sleep(30 * time.Millisecond)

// 			if player.Character.MaxHealth != float64(entities.AttributeRateAverage) {
// 				t.Errorf("Use(): got %#v, want %#v", player.Character.MaxHealth, float64(entities.AttributeRateAverage))
// 			}
// 		})
// 	}
// }
