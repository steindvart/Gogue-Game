package entity

import (
	"testing"
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
		heal      uint
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
