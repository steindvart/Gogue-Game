package entity

import (
	"slices"
	"testing"
)

func TestNewFood(t *testing.T) {
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
			if !(got.Item.Shape == tt.wantShape) {
				t.Errorf("NewFood(): got Shape %v, want %v", got.Item.Shape, tt.wantShape)
			}
			if !(got.HealthRegeneration < tt.wantHealthRegenerationLessThan) {
				t.Errorf("NewFood(): got HealthRegeneration %v, want less than %v", got.HealthRegeneration, tt.wantHealthRegenerationLessThan)
			}
			if !(slices.Contains(tt.wantNames, got.Item.Name)) {
				t.Errorf("NewFood(): got Name %v, want in %v", got.Item.Name, tt.wantNames)
			}
		})
	}
}

func TestFood_Taken(t *testing.T) {
	tests := []struct {
		name string
		want Box
	}{
		{
			name: "food is taken",
			want: Box{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			food := &Food{
				Item: Item{
					Shape: Box{
						Point: Point2D[int]{X: 1, Y: 2},
						Size:  Size2D[uint]{Height: 1, Width: 1},
					},
					Name: "Awkward Food",
				},
				HealthRegeneration: 5,
			}
			food.Taken()
			if food.Item.Shape != tt.want {
				t.Errorf("Taken() = (%v), want (%v)", food.Item.Shape, tt.want)
			}
		})
	}
}

func TestFood_Dropped(t *testing.T) {
	tests := []struct {
		name string
		box  Box
		want Box
	}{
		{
			name: "food is dropped",
			box:  Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
			want: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			food := &Food{
				Item: Item{
					Shape: Box{},
					Name:  "Awkward Food",
				},
				HealthRegeneration: 5,
			}

			food.Dropped(tt.box)

			if food.Item.Shape != tt.want {
				t.Errorf("Dropped(): Shape %v, want %v", food.Item.Shape, tt.want)
			}
		})
	}
}

func TestFood_Use(t *testing.T) {
	tests := []struct {
		name string
		food *Food
	}{
		{
			name: "use food",
			food: &Food{
				Item: Item{
					Shape: Box{},
					Name:  "Awkward Food",
				},
				HealthRegeneration: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(Box{})
			tt.food.Use(player)

			if player.Character.Health == float64(AttributeRateAverage) {
				t.Errorf("Use(): got %#v, want %#v", player.Character.MaxHealth, float64(AttributeRateAverage)+float64(tt.food.HealthRegeneration))
			}
		})
	}
}
