package entity

import (
	"slices"
	"testing"
)

func TestNewWeapon(t *testing.T) {
	tests := []struct {
		name                     string
		box                      Box
		wantShape                Box
		wantNames                []string
		wantStrengthBuffLessThan uint
	}{
		{
			name:      "weapon constructor",
			box:       Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
			wantShape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
			wantNames: []string{
				"Blade of the Forgotten Dawn",
				"Obsidian Reaver",
				"Fang of the Shadow Wolf",
				"Ironclad Cleaver",
				"Crimson Talon",
				"Thunderstrike Maul",
				"Serpent's Kiss Dagger",
				"Voidrend Sword",
				"Ebonheart Spear",
			},
			wantStrengthBuffLessThan: uint(IncreaseAttributeBaseParcentage + IncreaseAttributeMaxPercentage + 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWeapon(tt.box)
			if !(got.Consumable.Shape == tt.wantShape) {
				t.Errorf("NewWeapon(): got Shape %v, want %v", got.Consumable.Shape, tt.wantShape)
			}
			if !(slices.Contains(tt.wantNames, got.Consumable.Name)) {
				t.Errorf("NewWeapon(): got Name %v, want in %v", got.Consumable.Name, tt.wantNames)
			}
			if !(got.StrengthBuff < tt.wantStrengthBuffLessThan) {
				t.Errorf("NewWeapon(): got StrengthBuff %v, want less than %v", got.StrengthBuff, tt.wantStrengthBuffLessThan)
			}
		})
	}
}
