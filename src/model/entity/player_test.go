package entity

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestPlayer_IsAlive(t *testing.T) {
	tests := []struct {
		name   string
		health float64
		want   bool
	}{
		{
			name:   "alive",
			health: 10,
			want:   true,
		},
		{
			name:   "dead",
			health: 0,
			want:   false,
		},
		{
			name:   "negative health",
			health: -5,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Character: Character{Health: tt.health}}
			got := p.IsAlive()
			if got != tt.want {
				t.Errorf("IsAlive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_Move(t *testing.T) {
	tests := []struct {
		name  string
		start Point2D[int]
		delta Point2D[int]
		want  Point2D[int]
	}{
		{
			name:  "move positive",
			start: Point2D[int]{X: 0, Y: 0},
			delta: Point2D[int]{X: 2, Y: 3},
			want:  Point2D[int]{X: 2, Y: 3},
		},
		{
			name:  "move negative",
			start: Point2D[int]{X: 5, Y: 5},
			delta: Point2D[int]{X: -2, Y: -3},
			want:  Point2D[int]{X: 3, Y: 2},
		},
		{
			name:  "move zero",
			start: Point2D[int]{X: 1, Y: 1},
			delta: Point2D[int]{X: 0, Y: 0},
			want:  Point2D[int]{X: 1, Y: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Character: Character{Shape: Box{Point: tt.start}}}
			p.Move(tt.delta)
			if p.Character.Shape.Point != tt.want {
				t.Errorf("Move() = (%v), want (%v)", p.Character.Shape.Point, tt.want)
			}
		})
	}
}

func TestPlayer_TakeDamage(t *testing.T) {
	tests := []struct {
		name   string
		health float64
		damage float64
		want   float64
	}{
		{
			name:   "normal damage",
			health: 10,
			damage: 4,
			want:   6,
		},
		{
			name:   "overkill",
			health: 5,
			damage: 10,
			want:   0,
		},
		{
			name:   "zero damage",
			health: 7,
			damage: 0,
			want:   7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Character: Character{Health: tt.health}}
			p.TakeDamage(tt.damage)
			if p.Character.Health != tt.want {
				t.Errorf("TakeDamage() = %v, want %v", p.Character.Health, tt.want)
			}
		})
	}
}

func TestPlayer_Heal(t *testing.T) {
	tests := []struct {
		name      string
		health    float64
		maxHealth float64
		heal      float64
		want      float64
	}{
		{
			name:      "normal heal",
			health:    5,
			maxHealth: 10,
			heal:      3,
			want:      8,
		},
		{
			name:      "overheal",
			health:    8,
			maxHealth: 10,
			heal:      5,
			want:      10,
		},
		{
			name:      "zero heal",
			health:    7,
			maxHealth: 10,
			heal:      0,
			want:      7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Character: Character{Health: tt.health, MaxHealth: tt.maxHealth}}
			p.Heal(tt.heal)
			if p.Character.Health != tt.want {
				t.Errorf("Heal() = %v, want %v", p.Character.Health, tt.want)
			}
		})
	}
}

func TestPlayer_Attack(t *testing.T) {
	tests := []struct {
		name     string
		strength uint
		want     uint
	}{
		{
			name:     "normal attack",
			strength: 7,
			want:     7,
		},
		{
			name:     "zero strength",
			strength: 0,
			want:     0,
		},
		{
			name:     "high strength",
			strength: 100,
			want:     100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Character: Character{Strength: tt.strength}}
			got := p.Attack()
			if got != tt.want {
				t.Errorf("Attack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_CheckEvasion(t *testing.T) {
	tests := []struct {
		name         string
		agility      uint
		wantAllFalse bool
		wantHighRate bool
	}{
		{
			name:         "zero agility",
			agility:      0,
			wantAllFalse: true,
			wantHighRate: false,
		},
		{
			name:         "high agility",
			agility:      1000,
			wantAllFalse: false,
			wantHighRate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Character: Character{Agility: tt.agility}}
			tries := 1000
			count := 0
			for i := 0; i < tries; i++ {
				if p.CheckEvasion() {
					count++
				}
			}
			if tt.wantAllFalse && count != 0 {
				t.Errorf("CheckEvasion() with agility 0 should always be false, got %d/%d", count, tries)
			}
			if tt.wantHighRate && count < tries/2 {
				t.Errorf("CheckEvasion() with high agility should be high rate, got %d/%d", count, tries)
			}
			if count == tries {
				t.Errorf("CheckEvasion() should never be 100%%")
			}
		})
	}
}

func TestPlayer_TakeConsumableLike(t *testing.T) {
	tests := []struct {
		name        string
		consumables []ConsumableLike
		want        Backpack
		wantErrors  []error
	}{
		{
			name: "take elixir",
			consumables: []ConsumableLike{
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 1,
						Agility:   0,
						Strength:  0,
					},
					Increment: 10,
				},
			},
			want: Backpack{
				Capacity:    BackpackDefaultCapacity,
				ItemsNum:    1,
				Consumables: []ConsumableLike{},
				Treasures:   0,
			},
			wantErrors: []error{nil},
		},
		{
			name: "take two elixirs",
			consumables: []ConsumableLike{
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir One",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 1,
						Agility:   0,
						Strength:  0,
					},
					Increment: 10,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 3}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Two",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
			},
			want: Backpack{
				Capacity:    BackpackDefaultCapacity,
				ItemsNum:    2,
				Consumables: []ConsumableLike{},
				Treasures:   0,
			},
			wantErrors: []error{nil, nil},
		},
		{
			name: "take elixir and two foods",
			consumables: []ConsumableLike{
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 1,
						Agility:   0,
						Strength:  0,
					},
					Increment: 10,
				},
				&Food{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 3}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Food One",
					},
					HealthRegeneration: 15,
				},
				&Food{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 4}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Food Two",
					},
					HealthRegeneration: 15,
				},
			},
			want: Backpack{
				Capacity:    BackpackDefaultCapacity,
				ItemsNum:    3,
				Consumables: []ConsumableLike{},
				Treasures:   0,
			},
			wantErrors: []error{nil, nil, nil},
		},
		{
			name: "take elixir and food and scroll and weapon",
			consumables: []ConsumableLike{
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 1,
						Agility:   0,
						Strength:  0,
					},
					Increment: 10,
				},
				&Food{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 3}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Food",
					},
					HealthRegeneration: 15,
				},
				&Scroll{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 4}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Scroll",
					},
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 10,
				},
				&Weapon{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 5}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Weapon",
					},
					StrengthBuff: 10,
				},
			},
			want: Backpack{
				Capacity:    BackpackDefaultCapacity,
				ItemsNum:    4,
				Consumables: []ConsumableLike{},
				Treasures:   0,
			},
			wantErrors: []error{nil, nil, nil, nil},
		},
		{
			name: "take nine elixirs",
			consumables: []ConsumableLike{
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir One",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 1,
						Agility:   0,
						Strength:  0,
					},
					Increment: 10,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 3}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Two",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 4}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Three",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   0,
						Strength:  1,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 5}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Four",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 6}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Five",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 7}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Six",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 8}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Seven",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 9}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Eight",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 10}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Nine",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
			},
			want: Backpack{
				Capacity:    BackpackDefaultCapacity,
				ItemsNum:    9,
				Consumables: []ConsumableLike{},
				Treasures:   0,
			},
			wantErrors: []error{nil, nil, nil, nil, nil, nil, nil, nil, nil},
		},
		{
			name: "take ten elixirs",
			consumables: []ConsumableLike{
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir One",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 1,
						Agility:   0,
						Strength:  0,
					},
					Increment: 10,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 3}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Two",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 4}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Three",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   0,
						Strength:  1,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 5}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Four",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 6}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Five",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 7}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Six",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 8}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Seven",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 9}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Eight",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 10}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Nine",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 10}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir Ten",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 15,
				},
			},
			want: Backpack{
				Capacity:    BackpackDefaultCapacity,
				ItemsNum:    9,
				Consumables: []ConsumableLike{},
				Treasures:   0,
			},
			wantErrors: []error{nil, nil, nil, nil, nil, nil, nil, nil, nil, BackpackIsFullError{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}})

			for idx := range tt.consumables {
				err := p.TakeConsumableLike(tt.consumables[idx])
				errWant := tt.wantErrors[idx]

				if !errors.Is(err, errWant) {
					t.Errorf("TakeConsumableLike(): error %#v, want %#v", err, errWant)
				} else if err == nil {
					tt.want.Consumables = append(tt.want.Consumables, tt.consumables[idx])
				}
			}

			got := *p.Backpack

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TakeConsumableLike(): Backpack %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestPlayer_DropConsumableLike(t *testing.T) {
	tests := []struct {
		name             string
		consumables      []ConsumableLike
		consumableToDropIdx int
		wantConsumables []ConsumableLike
		want             Backpack
	}{
		{
			name: "drop elixir from elixir and food and scroll and weapon",
			consumables: []ConsumableLike{
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 1,
						Agility:   0,
						Strength:  0,
					},
					Increment: 10,
				},
				&Food{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 3}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Food",
					},
					HealthRegeneration: 15,
				},
				&Scroll{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 4}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Scroll",
					},
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 10,
				},
				&Weapon{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 5}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Weapon",
					},
					StrengthBuff: 10,
				},
			},
			consumableToDropIdx: 0,
			want: Backpack{
				Capacity:    BackpackDefaultCapacity,
				ItemsNum:    3,
				Consumables: []ConsumableLike{},
				Treasures:   0,
			},
		},
		{
			name: "drop food from elixir and food and scroll and weapon",
			consumables: []ConsumableLike{
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 1,
						Agility:   0,
						Strength:  0,
					},
					Increment: 10,
				},
				&Food{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 3}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Food",
					},
					HealthRegeneration: 15,
				},
				&Scroll{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 4}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Scroll",
					},
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 10,
				},
				&Weapon{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 5}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Weapon",
					},
					StrengthBuff: 10,
				},
			},
			consumableToDropIdx: 1,
			want: Backpack{
				Capacity:    BackpackDefaultCapacity,
				ItemsNum:    3,
				Consumables: []ConsumableLike{},
				Treasures:   0,
			},
		},
		{
			name: "drop scroll from elixir and food and scroll and weapon",
			consumables: []ConsumableLike{
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 1,
						Agility:   0,
						Strength:  0,
					},
					Increment: 10,
				},
				&Food{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 3}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Food",
					},
					HealthRegeneration: 15,
				},
				&Scroll{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 4}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Scroll",
					},
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 10,
				},
				&Weapon{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 5}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Weapon",
					},
					StrengthBuff: 10,
				},
			},
			consumableToDropIdx: 2,
			want: Backpack{
				Capacity:    BackpackDefaultCapacity,
				ItemsNum:    3,
				Consumables: []ConsumableLike{},
				Treasures:   0,
			},
		},
		{
			name: "drop weapon from elixir and food and scroll and weapon",
			consumables: []ConsumableLike{
				&Elixir{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Elixir",
					},
					EffectDuration: time.Minute,
					AffectedAttribute: Attributes{
						MaxHealth: 1,
						Agility:   0,
						Strength:  0,
					},
					Increment: 10,
				},
				&Food{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 3}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Food",
					},
					HealthRegeneration: 15,
				},
				&Scroll{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 4}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Scroll",
					},
					AffectedAttribute: Attributes{
						MaxHealth: 0,
						Agility:   1,
						Strength:  0,
					},
					Increment: 10,
				},
				&Weapon{
					Consumable: Consumable{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 5}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Weapon",
					},
					StrengthBuff: 10,
				},
			},
			consumableToDropIdx: 3,
			want: Backpack{
				Capacity:    BackpackDefaultCapacity,
				ItemsNum:    3,
				Consumables: []ConsumableLike{},
				Treasures:   0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curShape := Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}}

			tt.want.Consumables = make([]ConsumableLike, len(tt.consumables))
			copy(tt.want.Consumables, tt.consumables)
			tt.want.Consumables = append(tt.want.Consumables[:tt.consumableToDropIdx], tt.want.Consumables[tt.consumableToDropIdx+1:]...)

			var wantDroppedItem = tt.consumables[tt.consumableToDropIdx]
			wantDroppedItem.Dropped(curShape)

			p := NewPlayer(curShape)
			for _, item := range tt.consumables {
				p.TakeConsumableLike(item)
			}

			gotItem := p.DropConsumableLike(&tt.consumables[tt.consumableToDropIdx], curShape)
			if gotItem != wantDroppedItem {
				t.Errorf("DropConsumableLike() = %#v, want %#v", gotItem, wantDroppedItem)
			}

			gotBackpack := *p.Backpack
			if !reflect.DeepEqual(gotBackpack, tt.want) {
				t.Errorf("DropConsumableLike(): Backpack %#v, want %#v", gotBackpack, tt.want)
			}
		})
	}
}

func TestNewPlayer(t *testing.T) {
	tests := []struct {
		name string
		box  Box
		want *Player
	}{
		{
			name: "player constructor",
			box:  Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}},
			want: &Player{
				Character: Character{
					Shape:     Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}},
					Health:    float64(AttributeRateAverage),
					MaxHealth: float64(AttributeRateAverage),
					Strength:  uint(AttributeRateAverage),
					Agility:   uint(AttributeRateAverage),
				},
				Experience:     0,
				CharacterLevel: 1,
				Backpack:       NewBackpack(),
				Weapon:         nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPlayer(tt.box)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPlayer() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
