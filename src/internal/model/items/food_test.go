package items

// import (
// 	"gogue/internal/model/entity"
// 	"gogue/internal/model/primitive"
// 	"slices"
// 	"testing"
// )

// func TestNewFood(t *testing.T) {
// 	tests := []struct {
// 		name                           string
// 		box                            primitive.Box
// 		wantShape                      primitive.Box
// 		wantNames                      []string
// 		wantHealthRegenerationLessThan uint
// 	}{
// 		{
// 			name:      "food constructor",
// 			box:       primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 1}},
// 			wantShape: primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 1}},
// 			wantNames: []string{
// 				"Ration of the Ironclad",
// 				"Crimson Berry Cluster",
// 				"Loaf of the Forgotten Baker",
// 				"Smoked Wyrm Jerky",
// 				"Golden Apple of Vitality",
// 				"Hardtack of the Endless March",
// 				"Spiced Venison Strips",
// 				"Honeyed Nectar Bread",
// 				"Dried Mushrooms of the Deep",
// 			},
// 			wantHealthRegenerationLessThan: uint(IncreaseAttributeBaseParcentage + IncreaseAttributeMaxPercentage + 1),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := NewFood(tt.box)
// 			if !(got.Item.Shape == tt.wantShape) {
// 				t.Errorf("NewFood(): got Shape %v, want %v", got.Item.Shape, tt.wantShape)
// 			}
// 			if !(got.HealthRegeneration < tt.wantHealthRegenerationLessThan) {
// 				t.Errorf("NewFood(): got HealthRegeneration %v, want less than %v", got.HealthRegeneration, tt.wantHealthRegenerationLessThan)
// 			}
// 			if !(slices.Contains(tt.wantNames, got.Item.Name)) {
// 				t.Errorf("NewFood(): got Name %v, want in %v", got.Item.Name, tt.wantNames)
// 			}
// 		})
// 	}
// }

// func TestFood_Take(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want primitive.Box
// 	}{
// 		{
// 			name: "food is taken",
// 			want: primitive.Box{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			food := &Food{
// 				Item: Item{
// 					Shape: primitive.Box{
// 						Point: primitive.Point2D[int]{X: 1, Y: 2},
// 						Size:  primitive.Size2D[uint]{Height: 1, Width: 1},
// 					},
// 					Name: "Awkward Food",
// 				},
// 				HealthRegeneration: 5,
// 			}
// 			food.Take()
// 			if food.Item.Shape != tt.want {
// 				t.Errorf("Taken() = (%v), want (%v)", food.Item.Shape, tt.want)
// 			}
// 		})
// 	}
// }

// func TestFood_Drop(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		box  primitive.Box
// 		want primitive.Box
// 	}{
// 		{
// 			name: "food is dropped",
// 			box:  primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 1}},
// 			want: primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 1}},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			food := &Food{
// 				Item: Item{
// 					Shape: primitive.Box{},
// 					Name:  "Awkward Food",
// 				},
// 				HealthRegeneration: 5,
// 			}

// 			food.Drop(tt.box)

// 			if food.Item.Shape != tt.want {
// 				t.Errorf("Dropped(): Shape %v, want %v", food.Item.Shape, tt.want)
// 			}
// 		})
// 	}
// }

// func TestFood_Use(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		food *Food
// 	}{
// 		{
// 			name: "use food",
// 			food: &Food{
// 				Item: Item{
// 					Shape: primitive.Box{},
// 					Name:  "Awkward Food",
// 				},
// 				HealthRegeneration: 5,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			player := entities.NewPlayer(primitive.Box{})
// 			tt.food.Use(player)

// 			if player.Character.Health == float64(entities.AttributeRateAverage) {
// 				t.Errorf("Use(): got %#v, want %#v", player.Character.MaxHealth, float64(entities.AttributeRateAverage)+float64(tt.food.HealthRegeneration))
// 			}
// 		})
// 	}
// }
