package entities

// import (
// 	"gogue/internal/model/primitives"
// 	"reflect"
// 	"testing"
// )

// func TestNewZombie(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		box  primitives.Box
// 		want *Zombie
// 	}{
// 		{
// 			name: "zombie constructor",
// 			box:  primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 2}},
// 			want: &Zombie{
// 				Enemy: Enemy{
// 					Character: Character{
// 						Shape:     primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 2}},
// 						Agility:   uint(AttributeRateLow),
// 						Strength:  uint(AttributeRateAverage),
// 						Health:    float64(AttributeRateHigh),
// 						MaxHealth: float64(AttributeRateHigh),
// 					},
// 					Type:            EnemyTypeZombie,
// 					HostilityRadius: HostilityRadiusAverage,
// 					IsChasing:       false,
// 					Direction:       DirectionStop,
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := NewZombie(tt.box)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewZombie() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestNewVampire(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		box  primitives.Box
// 		want *Vampire
// 	}{
// 		{
// 			name: "vampire constructor",
// 			box:  primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 2}},
// 			want: &Vampire{
// 				Enemy: Enemy{
// 					Character: Character{
// 						Shape:     primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 2}},
// 						Agility:   uint(AttributeRateHigh),
// 						Strength:  uint(AttributeRateAverage),
// 						Health:    float64(AttributeRateHigh),
// 						MaxHealth: float64(AttributeRateHigh),
// 					},
// 					Type:            EnemyTypeVampire,
// 					HostilityRadius: HostilityRadiusHigh,
// 					IsChasing:       false,
// 					Direction:       DirectionStop,
// 				},
// 				AbsoluteEvasions: 0,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := NewVampire(tt.box)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewVampire() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestNewGhost(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		box  primitives.Box
// 		want *Ghost
// 	}{
// 		{
// 			name: "ghost constructor",
// 			box:  primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 2}},
// 			want: &Ghost{
// 				Enemy: Enemy{
// 					Character: Character{
// 						Shape:     primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 2}},
// 						Agility:   uint(AttributeRateHigh),
// 						Strength:  uint(AttributeRateLow),
// 						Health:    float64(AttributeRateLow),
// 						MaxHealth: float64(AttributeRateLow),
// 					},
// 					Type:            EnemyTypeGhost,
// 					HostilityRadius: HostilityRadiusLow,
// 					IsChasing:       false,
// 					Direction:       DirectionStop,
// 				},
// 				IsVisible: true,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := NewGhost(tt.box)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewGhost() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestNewOgre(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		box  primitives.Box
// 		want *Ogre
// 	}{
// 		{
// 			name: "ogre constructor",
// 			box:  primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 2}},
// 			want: &Ogre{
// 				Enemy: Enemy{
// 					Character: Character{
// 						Shape:     primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 2}},
// 						Agility:   uint(AttributeRateLow),
// 						Strength:  uint(AttributeRateVeryHigh),
// 						Health:    float64(AttributeRateVeryHigh),
// 						MaxHealth: float64(AttributeRateVeryHigh),
// 					},
// 					Type:            EnemyTypeOgre,
// 					HostilityRadius: HostilityRadiusAverage,
// 					IsChasing:       false,
// 					Direction:       DirectionStop,
// 				},
// 				IsResting: false,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := NewOgre(tt.box)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewOgre() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestNewSnakeMage(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		box  primitives.Box
// 		want *SnakeMage
// 	}{
// 		{
// 			name: "snake mage constructor",
// 			box:  primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 2}},
// 			want: &SnakeMage{
// 				Enemy: Enemy{
// 					Character: Character{
// 						Shape:     primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Height: 1, Width: 2}},
// 						Agility:   uint(AttributeRateVeryHigh),
// 						Strength:  uint(AttributeRateAverage),
// 						Health:    float64(AttributeRateHigh),
// 						MaxHealth: float64(AttributeRateHigh),
// 					},
// 					Type:            EnemyTypeSnakeMage,
// 					HostilityRadius: HostilityRadiusHigh,
// 					IsChasing:       false,
// 					Direction:       DirectionStop,
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := NewSnakeMage(tt.box)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewSnakeMage() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }
