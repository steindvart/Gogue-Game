package entity

// import (
// 	"errors"
// 	"gogue/internal/model/entity/primitive"
// 	"reflect"
// 	"testing"
// 	"time"
// )

// func TestPlayer_IsAlive(t *testing.T) {
// 	tests := []struct {
// 		name   string
// 		health float64
// 		want   bool
// 	}{
// 		{
// 			name:   "alive",
// 			health: 10,
// 			want:   true,
// 		},
// 		{
// 			name:   "dead",
// 			health: 0,
// 			want:   false,
// 		},
// 		{
// 			name:   "negative health",
// 			health: -5,
// 			want:   false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := &Player{Character: Character{Health: tt.health}}
// 			got := p.IsAlive()
// 			if got != tt.want {
// 				t.Errorf("IsAlive() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPlayer_Move(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		start primitive.Point2D[int]
// 		delta primitive.Point2D[int]
// 		want  primitive.Point2D[int]
// 	}{
// 		{
// 			name:  "move positive",
// 			start: primitive.Point2D[int]{X: 0, Y: 0},
// 			delta: primitive.Point2D[int]{X: 2, Y: 3},
// 			want:  primitive.Point2D[int]{X: 2, Y: 3},
// 		},
// 		{
// 			name:  "move negative",
// 			start: primitive.Point2D[int]{X: 5, Y: 5},
// 			delta: primitive.Point2D[int]{X: -2, Y: -3},
// 			want:  primitive.Point2D[int]{X: 3, Y: 2},
// 		},
// 		{
// 			name:  "move zero",
// 			start: primitive.Point2D[int]{X: 1, Y: 1},
// 			delta: primitive.Point2D[int]{X: 0, Y: 0},
// 			want:  primitive.Point2D[int]{X: 1, Y: 1},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := &Player{Character: Character{Shape: primitive.Box{Point: tt.start}}}
// 			p.Move(tt.delta)
// 			if p.Character.Shape.Point != tt.want {
// 				t.Errorf("Move() = (%v), want (%v)", p.Character.Shape.Point, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPlayer_TakeDamage(t *testing.T) {
// 	tests := []struct {
// 		name   string
// 		health float64
// 		damage float64
// 		want   float64
// 	}{
// 		{
// 			name:   "normal damage",
// 			health: 10,
// 			damage: 4,
// 			want:   6,
// 		},
// 		{
// 			name:   "overkill",
// 			health: 5,
// 			damage: 10,
// 			want:   0,
// 		},
// 		{
// 			name:   "zero damage",
// 			health: 7,
// 			damage: 0,
// 			want:   7,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := &Player{Character: Character{Health: tt.health}}
// 			p.TakeDamage(tt.damage)
// 			if p.Character.Health != tt.want {
// 				t.Errorf("TakeDamage() = %v, want %v", p.Character.Health, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPlayer_Heal(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		health    float64
// 		maxHealth float64
// 		heal      float64
// 		want      float64
// 	}{
// 		{
// 			name:      "normal heal",
// 			health:    5,
// 			maxHealth: 10,
// 			heal:      3,
// 			want:      8,
// 		},
// 		{
// 			name:      "overheal",
// 			health:    8,
// 			maxHealth: 10,
// 			heal:      5,
// 			want:      10,
// 		},
// 		{
// 			name:      "zero heal",
// 			health:    7,
// 			maxHealth: 10,
// 			heal:      0,
// 			want:      7,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := &Player{Character: Character{Health: tt.health, MaxHealth: tt.maxHealth}}
// 			p.Heal(tt.heal)
// 			if p.Character.Health != tt.want {
// 				t.Errorf("Heal() = %v, want %v", p.Character.Health, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPlayer_Attack(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		strength uint
// 		want     uint
// 	}{
// 		{
// 			name:     "normal attack",
// 			strength: 7,
// 			want:     7,
// 		},
// 		{
// 			name:     "zero strength",
// 			strength: 0,
// 			want:     0,
// 		},
// 		{
// 			name:     "high strength",
// 			strength: 100,
// 			want:     100,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := &Player{Character: Character{Strength: tt.strength}}
// 			got := p.Attack()
// 			if got != tt.want {
// 				t.Errorf("Attack() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPlayer_AttackWithWeapon(t *testing.T) {
// 	tests := []struct {
// 		name   string
// 		weapon *Weapon
// 		want   uint
// 	}{
// 		{
// 			name: "damage 0",
// 			weapon: &Weapon{
// 				Item: Item{
// 					Shape: primitive.Box{},
// 					Name:  "Awkward Weapon",
// 				},
// 				Damage: 0,
// 			},
// 			want: uint(AttributeRateAverage),
// 		},
// 		{
// 			name: "damage 10",
// 			weapon: &Weapon{
// 				Item: Item{
// 					Shape: primitive.Box{},
// 					Name:  "Awkward Weapon",
// 				},
// 				Damage: 10,
// 			},
// 			want: uint(AttributeRateAverage) + 10,
// 		},
// 		{
// 			name: "damage 50",
// 			weapon: &Weapon{
// 				Item: Item{
// 					Shape: primitive.Box{},
// 					Name:  "Awkward Weapon",
// 				},
// 				Damage: 50,
// 			},
// 			want: uint(AttributeRateAverage) + 50,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := NewPlayer(primitive.Box{})
// 			p.Weapon = tt.weapon

// 			got := p.Attack()
// 			if got != tt.want {
// 				t.Errorf("Attack() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPlayer_CheckEvasion(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		agility      uint
// 		wantAllFalse bool
// 		wantHighRate bool
// 	}{
// 		{
// 			name:         "zero agility",
// 			agility:      0,
// 			wantAllFalse: true,
// 			wantHighRate: false,
// 		},
// 		{
// 			name:         "high agility",
// 			agility:      1000,
// 			wantAllFalse: false,
// 			wantHighRate: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := &Player{Character: Character{Agility: tt.agility}}
// 			tries := 1000
// 			count := 0
// 			for i := 0; i < tries; i++ {
// 				if p.CheckEvasion() {
// 					count++
// 				}
// 			}
// 			if tt.wantAllFalse && count != 0 {
// 				t.Errorf("CheckEvasion() with agility 0 should always be false, got %d/%d", count, tries)
// 			}
// 			if tt.wantHighRate && count < tries/2 {
// 				t.Errorf("CheckEvasion() with high agility should be high rate, got %d/%d", count, tries)
// 			}
// 			if count == tries {
// 				t.Errorf("CheckEvasion() should never be 100%%")
// 			}
// 		})
// 	}
// }

// func TestPlayer_TakeTreasure(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		treasure Treasure
// 		want     uint
// 	}{
// 		{
// 			name: "take treasure cost 10",
// 			treasure: Treasure{
// 				Shape: primitive.Box{},
// 				Name:  "Gold",
// 				Value: 10,
// 			},
// 			want: 10,
// 		},
// 		{
// 			name: "take treasure cost 0",
// 			treasure: Treasure{
// 				Shape: primitive.Box{},
// 				Name:  "Gold",
// 				Value: 0,
// 			},
// 			want: 0,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := NewPlayer(primitive.Box{})

// 			p.TakeTreasure(&tt.treasure)
// 			got := p.Backpack.Treasures
// 			if got != tt.want {
// 				t.Errorf("TakeTreasure(): Backpack Treasures got %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPlayer_TakeItem(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		items        []ItemLike
// 		wantBackpack Backpack
// 		want         []error
// 	}{
// 		{
// 			name: "take elixir",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			wantBackpack: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 1,
// 				Elixirs: map[string][]Elixir{
// 					"Awkward Elixir": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 				},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: []error{nil},
// 		},
// 		{
// 			name: "take nine elixirs",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			wantBackpack: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 9,
// 				Elixirs: map[string][]Elixir{
// 					"Awkward Elixir": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 				},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: []error{nil, nil, nil, nil, nil, nil, nil, nil, nil},
// 		},
// 		{
// 			name: "take ten elixirs",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			wantBackpack: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 9,
// 				Elixirs: map[string][]Elixir{
// 					"Awkward Elixir": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 				},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: []error{nil, nil, nil, nil, nil, nil, nil, nil, nil, BackpackIsFullError{}},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := NewPlayer(primitive.Box{})
// 			for idx := range tt.items {
// 				err := p.TakeItem(tt.items[idx])
// 				errWant := tt.want[idx]

// 				if !errors.Is(err, errWant) {
// 					t.Errorf("TakeItem() = %#v, want %#v", err, errWant)
// 				}
// 			}

// 			if !reflect.DeepEqual(*p.Backpack, tt.wantBackpack) {
// 				t.Errorf("TakeItem(): Backpack got %#v, want %#v", *p.Backpack, tt.wantBackpack)
// 			}

// 		})
// 	}
// }

// func TestPlayer_DropItem(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		items        []ItemLike
// 		item         ItemLike
// 		wantBackpack Backpack
// 		want         error
// 	}{
// 		{
// 			name:  "drop elixir from empty backpack",
// 			items: []ItemLike{},
// 			item: &Elixir{
// 				Item: Item{
// 					Shape: primitive.Box{},
// 					Name:  "Awkward Elixir",
// 				},
// 				EffectDuration: time.Minute,
// 				AffectedAttribute: primitive.Attributes{
// 					MaxHealth: 1,
// 					Agility:   0,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			wantBackpack: Backpack{
// 				Capacity:  BackpackDefaultCapacity,
// 				ItemsNum:  0,
// 				Elixirs:   map[string][]Elixir{},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: ItemIsNotInBackpackError{},
// 		},
// 		{
// 			name: "drop elixir",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			item: &Elixir{
// 				Item: Item{
// 					Shape: primitive.Box{},
// 					Name:  "Awkward Elixir",
// 				},
// 				EffectDuration: time.Minute,
// 				AffectedAttribute: primitive.Attributes{
// 					MaxHealth: 1,
// 					Agility:   0,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			wantBackpack: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 2,
// 				Elixirs: map[string][]Elixir{
// 					"Awkward Elixir": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitive.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitive.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 				},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := NewPlayer(primitive.Box{})

// 			for idx := range tt.items {
// 				_ = p.Backpack.AddItem(tt.items[idx])
// 			}

// 			got := p.DropItem(tt.item)

// 			if !errors.Is(got, tt.want) {
// 				t.Errorf("DropItem() = %#v, want %#v", got, tt.want)
// 			}

// 			if !reflect.DeepEqual(*p.Backpack, tt.wantBackpack) {
// 				t.Errorf("DropItem(): Backpack got %#v, want %#v", *p.Backpack, tt.wantBackpack)
// 			}
// 		})
// 	}
// }

// func TestPlayer_UseItem(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		items      []ItemLike
// 		item       ItemLike
// 		wantPlayer Player
// 		wantError  error
// 	}{
// 		{
// 			name:  "use elixir, which is not in backpack",
// 			items: []ItemLike{},
// 			item: &Elixir{
// 				Item: Item{
// 					Shape: primitive.Box{},
// 					Name:  "Awkward Elixir",
// 				},
// 				EffectDuration: time.Minute,
// 				AffectedAttribute: primitive.Attributes{
// 					MaxHealth: 1,
// 					Agility:   0,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantError: ItemIsNotInBackpackError{},
// 		},
// 		{
// 			name: "use elixir",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			item: &Elixir{
// 				Item: Item{
// 					Shape: primitive.Box{},
// 					Name:  "Awkward Elixir",
// 				},
// 				EffectDuration: time.Minute,
// 				AffectedAttribute: primitive.Attributes{
// 					MaxHealth: 1,
// 					Agility:   0,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage) + float64(10),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantError: nil,
// 		},
// 		{
// 			name: "use scroll",
// 			items: []ItemLike{
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			item: &Scroll{
// 				Item: Item{
// 					Shape: primitive.Box{},
// 					Name:  "Awkward Scroll",
// 				},
// 				AffectedAttribute: primitive.Attributes{
// 					MaxHealth: 0,
// 					Agility:   1,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage) + 10,
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantError: nil,
// 		},
// 		{
// 			name: "use food",
// 			items: []ItemLike{
// 				&Food{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Food",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 			},
// 			item: &Food{
// 				Item: Item{
// 					Shape: primitive.Box{},
// 					Name:  "Awkward Food",
// 				},
// 				HealthRegeneration: 15,
// 			},
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage) + float64(15),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantError: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := NewPlayer(primitive.Box{})

// 			for idx := range tt.items {
// 				_ = p.Backpack.AddItem(tt.items[idx])
// 			}

// 			gotErr := p.UseItem(tt.item)

// 			time.Sleep(1 * time.Millisecond)

// 			if !errors.Is(gotErr, tt.wantError) {
// 				t.Errorf("UseItem(): error got %#v, want %#v", gotErr, tt.wantError)
// 			}

// 			if !reflect.DeepEqual(*p, tt.wantPlayer) {
// 				t.Errorf("UseItem(): Player got %#v, want %#v", *p, tt.wantPlayer)
// 			}
// 		})
// 	}
// }

// func TestPlayer_UseWeapon(t *testing.T) {
// 	tests := []struct {
// 		name             string
// 		weapons          []ItemLike
// 		currentWeaponIdx int
// 		useWeaponIdx     int
// 		wantPlayer       Player
// 		wantErr          error
// 	}{
// 		{
// 			name: "use first weapon with current weapon nil",
// 			weapons: []ItemLike{
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 1",
// 					},
// 					Damage: 1,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 2",
// 					},
// 					Damage: 10,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 3",
// 					},
// 					Damage: 50,
// 				},
// 			},
// 			currentWeaponIdx: -1,
// 			useWeaponIdx:     0,
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantErr: nil,
// 		},
// 		{
// 			name: "use second weapon with current weapon nil",
// 			weapons: []ItemLike{
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 1",
// 					},
// 					Damage: 1,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 2",
// 					},
// 					Damage: 10,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 3",
// 					},
// 					Damage: 50,
// 				},
// 			},
// 			currentWeaponIdx: -1,
// 			useWeaponIdx:     1,
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantErr: nil,
// 		},
// 		{
// 			name: "use first weapon with current weapon third",
// 			weapons: []ItemLike{
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 1",
// 					},
// 					Damage: 1,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 2",
// 					},
// 					Damage: 10,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 3",
// 					},
// 					Damage: 50,
// 				},
// 			},
// 			currentWeaponIdx: 2,
// 			useWeaponIdx:     0,
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantErr: nil,
// 		},
// 		{
// 			name: "use second weapon with current weapon third",
// 			weapons: []ItemLike{
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 1",
// 					},
// 					Damage: 1,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 2",
// 					},
// 					Damage: 10,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 3",
// 					},
// 					Damage: 50,
// 				},
// 			},
// 			currentWeaponIdx: 2,
// 			useWeaponIdx:     1,
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantErr: nil,
// 		},
// 		{
// 			name: "use third weapon with current weapon third",
// 			weapons: []ItemLike{
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 1",
// 					},
// 					Damage: 1,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 2",
// 					},
// 					Damage: 10,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon 3",
// 					},
// 					Damage: 50,
// 				},
// 			},
// 			currentWeaponIdx: 2,
// 			useWeaponIdx:     2,
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantErr: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := NewPlayer(primitive.Box{})

// 			for idx := range tt.weapons {
// 				_ = p.Backpack.AddItem(tt.weapons[idx])

// 				if tt.useWeaponIdx != tt.currentWeaponIdx && idx != tt.currentWeaponIdx || tt.useWeaponIdx == tt.currentWeaponIdx {
// 					_ = tt.wantPlayer.Backpack.AddItem(tt.weapons[idx])
// 				}
// 			}

// 			if tt.currentWeaponIdx != -1 {
// 				p.Weapon = &p.Backpack.Weapons[tt.weapons[tt.currentWeaponIdx].(*Weapon).Item.Name][0]
// 			}
// 			tt.wantPlayer.Weapon = &tt.wantPlayer.Backpack.Weapons[tt.weapons[tt.useWeaponIdx].(*Weapon).Item.Name][0]

// 			gotErr := p.UseItem(&p.Backpack.Weapons[tt.weapons[tt.useWeaponIdx].(*Weapon).Item.Name][0])

// 			if !errors.Is(gotErr, tt.wantErr) {
// 				t.Errorf("UseItem(): error got %#v, want %#v", gotErr, tt.wantErr)
// 			}

// 			if !reflect.DeepEqual(*p, tt.wantPlayer) {
// 				t.Errorf("UseItem(): Player got %#v, want %#v", *p, tt.wantPlayer)
// 			}
// 		})
// 	}
// }

// func TestPlayer_GetItemsListAndUseItem(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		items      []ItemLike
// 		wantPlayer Player
// 		wantError  error
// 	}{
// 		{
// 			name: "get and use elixir",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage) + float64(10),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantError: nil,
// 		},
// 		{
// 			name: "get and use scroll",
// 			items: []ItemLike{
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage) + 10,
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantError: nil,
// 		},
// 		{
// 			name: "get and use food",
// 			items: []ItemLike{
// 				&Food{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Food",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 			},
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage) + float64(15),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantError: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := NewPlayer(primitive.Box{})

// 			for idx := range tt.items {
// 				_ = p.Backpack.AddItem(tt.items[idx])
// 			}

// 			l := p.Backpack.GetItemsList()

// 			gotErr := p.UseItem(l[0].Item)

// 			time.Sleep(1 * time.Millisecond)

// 			if !errors.Is(gotErr, tt.wantError) {
// 				t.Errorf("UseItem(): error got %#v, want %#v", gotErr, tt.wantError)
// 			}

// 			if !reflect.DeepEqual(*p, tt.wantPlayer) {
// 				t.Errorf("UseItem(): Player got %#v, want %#v", *p, tt.wantPlayer)
// 			}
// 		})
// 	}
// }

// func TestPlayer_GetItemsListAndDropItem(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		items      []ItemLike
// 		wantPlayer Player
// 		wantError  error
// 	}{
// 		{
// 			name: "get and drop elixir",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantError: nil,
// 		},
// 		{
// 			name: "get and drop scroll",
// 			items: []ItemLike{
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitive.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantError: nil,
// 		},
// 		{
// 			name: "get and drop food",
// 			items: []ItemLike{
// 				&Food{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Food",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 			},
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantError: nil,
// 		},
// 		{
// 			name: "get and drop weapon",
// 			items: []ItemLike{
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitive.Box{},
// 						Name:  "Awkward Weapon",
// 					},
// 					Damage: 10,
// 				},
// 			},
// 			wantPlayer: Player{
// 				Character: Character{
// 					Shape:     primitive.Box{},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 			wantError: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := NewPlayer(primitive.Box{})

// 			for idx := range tt.items {
// 				_ = p.Backpack.AddItem(tt.items[idx])
// 			}

// 			l := p.Backpack.GetItemsList()

// 			gotErr := p.DropItem(l[0].Item)

// 			if !errors.Is(gotErr, tt.wantError) {
// 				t.Errorf("UseItem(): error got %#v, want %#v", gotErr, tt.wantError)
// 			}

// 			if !reflect.DeepEqual(*p, tt.wantPlayer) {
// 				t.Errorf("UseItem(): Player got %#v, want %#v", *p, tt.wantPlayer)
// 			}
// 		})
// 	}
// }

// func TestNewPlayer(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		box  primitive.Box
// 		want *Player
// 	}{
// 		{
// 			name: "player constructor",
// 			box:  primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 2}},
// 			want: &Player{
// 				Character: Character{
// 					Shape:     primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 2}},
// 					Health:    float64(AttributeRateAverage),
// 					MaxHealth: float64(AttributeRateAverage),
// 					Strength:  uint(AttributeRateAverage),
// 					Agility:   uint(AttributeRateAverage),
// 				},
// 				Experience:     0,
// 				CharacterLevel: 1,
// 				Backpack:       NewBackpack(),
// 				Weapon:         nil,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := NewPlayer(tt.box)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewPlayer() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }
