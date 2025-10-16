package entity

import (
	"slices"
	"testing"
)

func TestNewtFood(t *testing.T) {
	tests := []struct {
		name                           string
		box                            Box
		wantShape                      Box
		wantNames                      []string
		wantHealthRegenerationLessThan uint
	}{
		{
			name:      "food constructor",
			box:       Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
			wantShape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
			wantNames: []string{
				"Ration of the Ironclad",
				"Crimson Berry Cluster",
				"Loaf of the Forgotten Baker",
				"Smoked Wyrm Jerky",
				"Golden Apple of Vitality",
				"Hardtack of the Endless March",
				"Spiced Venison Strips",
				"Honeyed Nectar Bread",
				"Dried Mushrooms of the Deep",
			},
			wantHealthRegenerationLessThan: uint(IncreaseAttributeBaseParcentage + IncreaseAttributeMaxPercentage + 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFood(tt.box)
			if !(got.Shape == tt.wantShape) {
				t.Errorf("NewFood(): got Shape %v, want %v", got.Shape, tt.wantShape)
			}
			if !(got.HealthRegeneration < tt.wantHealthRegenerationLessThan) {
				t.Errorf("NewFood(): got HealthRegeneration %v, want less than %v", got.HealthRegeneration, tt.wantHealthRegenerationLessThan)
			}
			if !(slices.Contains(tt.wantNames, got.Name)) {
				t.Errorf("NewFood(): got Name %v, want in %v", got.Name, tt.wantNames)
			}
		})
	}
}
