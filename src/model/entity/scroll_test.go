package entity

import (
	"slices"
	"testing"
)

func TestNewScroll(t *testing.T) {
	tests := []struct {
		name                        string
		box                         Box
		wantShape                   Box
		wantNames                   []string
		wantIncrementLessThan       uint
		wantAffectedAttributesList []Attributes
	}{
		{
			name:      "scroll constructor",
			box:       Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
			wantShape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
			wantNames: []string{
				"Scroll of Shadowstep",
				"Parchment of Eternal Flame",
				"Manuscript of Forgotten Truths",
				"Scroll of Iron Will",
				"Vellum of the Void",
				"Scroll of Whispers",
				"Tome of the Lost King",
				"Scroll of Unseen Paths",
				"Parchment of Thunderous Roar",
			},
			wantIncrementLessThan: uint(IncreaseAttributeBaseParcentage + IncreaseAttributeMaxPercentage + 1),
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
			got := NewScroll(tt.box)
			if !(got.Consumable.Shape == tt.wantShape) {
				t.Errorf("NewScroll(): got Shape %v, want %v", got.Consumable.Shape, tt.wantShape)
			}
			if !(slices.Contains(tt.wantNames, got.Consumable.Name)) {
				t.Errorf("NewScroll(): got Name %v, want in %v", got.Consumable.Name, tt.wantNames)
			}
			if !(got.Increment < tt.wantIncrementLessThan) {
				t.Errorf("NewScroll(): got Increment %v, want less than %v", got.Increment, tt.wantIncrementLessThan)
			}
			if !(slices.Contains(tt.wantAffectedAttributesList, got.AffectedAttribute)) {
				t.Errorf("NewScroll(): got AffectedAttribute %v, want in %v", got.AffectedAttribute, tt.wantAffectedAttributesList)
			}
		})
	}
}
