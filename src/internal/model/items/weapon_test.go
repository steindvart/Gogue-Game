package items

// import (
// 	"gogue/internal/model/entity"
// 	"gogue/internal/model/primitive"
// 	"slices"
// 	"testing"
// )

// func TestNewWeapon(t *testing.T) {
// 	tests := []struct {
// 		name               string
// 		box                primitive.Box
// 		wantShape          primitive.Box
// 		wantNames          []string
// 		wantDamageLessThan uint
// 	}{
// 		{
// 			name:      "weapon constructor",
// 			box:       primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 1}},
// 			wantShape: primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 1}},
// 			wantNames: []string{
// 				"Blade of the Forgotten Dawn",
// 				"Obsidian Reaver",
// 				"Fang of the Shadow Wolf",
// 				"Ironclad Cleaver",
// 				"Crimson Talon",
// 				"Thunderstrike Maul",
// 				"Serpent's Kiss Dagger",
// 				"Voidrend Sword",
// 				"Ebonheart Spear",
// 			},
// 			wantDamageLessThan: uint(IncreaseAttributeBaseParcentage + IncreaseAttributeMaxPercentage + 1),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := NewWeapon(tt.box)
// 			if !(got.Item.Shape == tt.wantShape) {
// 				t.Errorf("NewWeapon(): got Shape %v, want %v", got.Item.Shape, tt.wantShape)
// 			}
// 			if !(slices.Contains(tt.wantNames, got.Item.Name)) {
// 				t.Errorf("NewWeapon(): got Name %v, want in %v", got.Item.Name, tt.wantNames)
// 			}
// 			if !(got.Damage < tt.wantDamageLessThan) {
// 				t.Errorf("NewWeapon(): got Damage %v, want less than %v", got.Damage, tt.wantDamageLessThan)
// 			}
// 		})
// 	}
// }

// func TestWeapon_Take(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want primitive.Box
// 	}{
// 		{
// 			name: "weapon is taken",
// 			want: primitive.Box{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			weapon := &Weapon{
// 				Item: Item{
// 					Shape: primitive.Box{
// 						Point: primitive.Point2D[int]{X: 1, Y: 2},
// 						Size:  primitive.Size2D[uint]{Height: 1, Width: 1},
// 					},
// 					Name: "Awkward Weapon",
// 				},
// 				Damage: 5,
// 			}
// 			weapon.Take()
// 			if weapon.Item.Shape != tt.want {
// 				t.Errorf("Taken() = (%v), want (%v)", weapon.Item.Shape, tt.want)
// 			}
// 		})
// 	}
// }

// func TestWeapon_Drop(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		box  primitive.Box
// 		want primitive.Box
// 	}{
// 		{
// 			name: "weapon is dropped",
// 			box:  primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 1}},
// 			want: primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 1}},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			weapon := &Weapon{
// 				Item: Item{
// 					Shape: primitive.Box{},
// 					Name:  "Awkward Weapon",
// 				},
// 				Damage: 5,
// 			}

// 			weapon.Drop(tt.box)

// 			if weapon.Item.Shape != tt.want {
// 				t.Errorf("Dropped(): Shape %v, want %v", weapon.Item.Shape, tt.want)
// 			}
// 		})
// 	}
// }

// func TestWeapon_Use(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		weapons []*Weapon
// 	}{
// 		{
// 			name: "use weapon",
// 			weapons: []*Weapon{
// 				{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 1",
// 					},
// 					Damage: 5,
// 				},
// 				{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 2",
// 					},
// 					Damage: 5,
// 				},
// 			},
// 		},
// 		{
// 			name: "use first weapon",
// 			weapons: []*Weapon{
// 				{
// 					Item: Item{
// 						Shape: primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 1}},
// 						Name:  "Awkward Weapon",
// 					},
// 					Damage: 5,
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			player := entities.NewPlayer(primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 1}})

// 			for idx := range tt.weapons {
// 				_ = player.TakeItem(tt.weapons[idx])
// 			}

// 			if len(tt.weapons) > 1 {
// 				player.Weapon = tt.weapons[0]
// 				currentWeapon := player.Weapon
// 				if currentWeapon.Item.Shape == player.Character.Shape {
// 					t.Errorf("Use(): Shape got %#v, got %#v", currentWeapon.Item.Shape, player.Character.Shape)
// 				}

// 				tt.weapons[1].Use(player)

// 				if !(currentWeapon.Item.Shape == player.Character.Shape) {
// 					t.Errorf("Use(): Shape got %#v, got %#v", currentWeapon.Item.Shape, player.Character.Shape)
// 				}

// 				if !(player.Weapon == tt.weapons[1]) {
// 					t.Errorf("Use(): Weapon got %#v, got %#v", player.Weapon, tt.weapons[1])
// 				}
// 			} else {
// 				tt.weapons[0].Use(player)

// 				if !(player.Weapon == tt.weapons[0]) {
// 					t.Errorf("Use(): Weapon got %#v, want %#v", player.Weapon, tt.weapons[0])
// 				}
// 			}
// 		})
// 	}
// }
