package items

// import (
// 	"gogue/internal/model/entity"
// 	"gogue/internal/model/primitive"
// 	"slices"
// 	"testing"
// )

// func TestNewScroll(t *testing.T) {
// 	tests := []struct {
// 		name                       string
// 		box                        primitives.Box
// 		wantShape                  primitives.Box
// 		wantNames                  []string
// 		wantIncrementLessThan      uint
// 		wantAffectedAttributesList []primitives.Attributes
// 	}{
// 		{
// 			name:      "scroll constructor",
// 			box:       primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 1}},
// 			wantShape: primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 1}},
// 			wantNames: []string{
// 				"Scroll of Shadowstep",
// 				"Parchment of Eternal Flame",
// 				"Manuscript of Forgotten Truths",
// 				"Scroll of Iron Will",
// 				"Vellum of the Void",
// 				"Scroll of Whispers",
// 				"Tome of the Lost King",
// 				"Scroll of Unseen Paths",
// 				"Parchment of Thunderous Roar",
// 			},
// 			wantIncrementLessThan: uint(IncreaseAttributeBaseParcentage + IncreaseAttributeMaxPercentage + 1),
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
// 			got := NewScroll(tt.box)
// 			if !(got.Item.Shape == tt.wantShape) {
// 				t.Errorf("NewScroll(): got Shape %v, want %v", got.Item.Shape, tt.wantShape)
// 			}
// 			if !(slices.Contains(tt.wantNames, got.Item.Name)) {
// 				t.Errorf("NewScroll(): got Name %v, want in %v", got.Item.Name, tt.wantNames)
// 			}
// 			if !(got.Increment < tt.wantIncrementLessThan) {
// 				t.Errorf("NewScroll(): got Increment %v, want less than %v", got.Increment, tt.wantIncrementLessThan)
// 			}
// 			if !(slices.Contains(tt.wantAffectedAttributesList, got.AffectedAttribute)) {
// 				t.Errorf("NewScroll(): got AffectedAttribute %v, want in %v", got.AffectedAttribute, tt.wantAffectedAttributesList)
// 			}
// 		})
// 	}
// }

// func TestScroll_Take(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want primitives.Box
// 	}{
// 		{
// 			name: "scroll is taken",
// 			want: primitives.Box{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			scroll := &Scroll{
// 				Item: Item{
// 					Shape: primitives.Box{
// 						Point: primitives.Point2D[int]{X: 1, Y: 2},
// 						Size:  primitives.Size2D[uint]{Height: 1, Width: 1},
// 					},
// 					Name: "Awkward Scroll",
// 				},
// 				AffectedAttribute: primitives.Attributes{MaxHealth: 1, Agility: 0, Strength: 0},
// 				Increment:         1,
// 			}
// 			scroll.Take()
// 			if scroll.Item.Shape != tt.want {
// 				t.Errorf("Taken() = (%v), want (%v)", scroll.Item.Shape, tt.want)
// 			}
// 		})
// 	}
// }

// func TestScroll_Drop(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		box  primitives.Box
// 		want primitives.Box
// 	}{
// 		{
// 			name: "scroll is dropped",
// 			box:  primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 1}},
// 			want: primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 1}},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			scroll := &Scroll{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Scroll",
// 				},
// 				AffectedAttribute: primitives.Attributes{MaxHealth: 1, Agility: 0, Strength: 0},
// 				Increment:         1,
// 			}

// 			scroll.Drop(tt.box)

// 			if scroll.Item.Shape != tt.want {
// 				t.Errorf("Dropped(): Shape %v, want %v", scroll.Item.Shape, tt.want)
// 			}
// 		})
// 	}
// }

// func TestScroll_Use(t *testing.T) {
// 	tests := []struct {
// 		name   string
// 		scroll *Scroll
// 	}{
// 		{
// 			name: "use scroll",
// 			scroll: &Scroll{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Scroll",
// 				},
// 				AffectedAttribute: primitives.Attributes{MaxHealth: 1, Agility: 0, Strength: 0},
// 				Increment:         5,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			player := entities.NewPlayer(primitives.Box{})
// 			tt.scroll.Use(player)

// 			if player.Character.MaxHealth == float64(entities.AttributeRateAverage) {
// 				t.Errorf("Use(): got %#v, want %#v", player.Character.MaxHealth, float64(entities.AttributeRateAverage)+float64(tt.scroll.Increment))
// 			}
// 		})
// 	}
// }
