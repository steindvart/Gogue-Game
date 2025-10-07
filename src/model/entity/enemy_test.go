package entity

import (
	"testing"
)

func TestZombieConstructor(t *testing.T) {
	tests := []struct {
		name                string
		box                 Box
		wantBox             Box
		wantAgility         uint
		wantStrength        uint
		wantHealth          float64
		wantMaxHealth       float64
		wantType            uint
		wantHostilityRadius uint
		wantIsChasing       bool
		wantDiretion        uint
	}{
		{
			name:                "zombie constructor",
			box:                 Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}},
			wantBox:             Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}},
			wantAgility:         25,
			wantStrength:        50,
			wantHealth:          75.0,
			wantMaxHealth:       0.0,
			wantType:            0,
			wantHostilityRadius: 4,
			wantIsChasing:       false,
			wantDiretion:        8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := NewZombie(tt.box)

			if z == nil {
				t.Fatalf("NewZombie(): want valid pointer, got %v", z)
			}

			if z.Enemy.Character.Shape != tt.wantBox {
				t.Errorf("NewZombie(): Shape got %v, want %v", z.Enemy.Character.Shape, tt.wantBox)
			}

			if z.Enemy.Character.Agility != tt.wantAgility {
				t.Errorf("NewZombie(): Agility got %v, want %v", z.Enemy.Character.Agility, tt.wantAgility)
			}

			if z.Enemy.Character.Strength != tt.wantStrength {
				t.Errorf("NewZombie(): Strength got %v, want %v", z.Enemy.Character.Strength, tt.wantStrength)
			}

			if z.Enemy.Character.Health != tt.wantHealth {
				t.Errorf("NewZombie(): Health got %v, want %v", z.Enemy.Character.Health, tt.wantHealth)
			}

			if z.Enemy.Character.MaxHealth != tt.wantMaxHealth {
				t.Errorf("NewZombie(): MaxHealth got %v, want %v", z.Enemy.Character.MaxHealth, tt.wantMaxHealth)
			}

			if uint(z.Enemy.Type) != tt.wantType {
				t.Errorf("NewZombie(): Type got %v, want %v", z.Enemy.Type, tt.wantType)
			}

			if uint(z.Enemy.HostilityRadius) != tt.wantHostilityRadius {
				t.Errorf("NewZombie(): HostilityRadius got %v, want %v", z.Enemy.HostilityRadius, tt.wantHostilityRadius)
			}

			if z.Enemy.IsChasing != tt.wantIsChasing {
				t.Errorf("NewZombie(): IsChasing got %v, want %v", z.Enemy.IsChasing, tt.wantDiretion)
			}

			if uint(z.Enemy.Direction) != tt.wantDiretion {
				t.Errorf("NewZombie(): Direction got %v, want %v", z.Enemy.Direction, tt.wantDiretion)
			}
		})
	}
}

func TestVampireConstructor(t *testing.T) {
	tests := []struct {
		name                string
		box                 Box
		wantBox             Box
		wantAgility         uint
		wantStrength        uint
		wantHealth          float64
		wantMaxHealth       float64
		wantType            uint
		wantHostilityRadius uint
		wantIsChasing       bool
		wantDiretion        uint
		wantHadFirstDamage  bool
	}{
		{
			name:                "vampire constructor",
			box:                 Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}},
			wantBox:             Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}},
			wantAgility:         75,
			wantStrength:        50,
			wantHealth:          75.0,
			wantMaxHealth:       0.0,
			wantType:            1,
			wantHostilityRadius: 6,
			wantIsChasing:       false,
			wantDiretion:        8,
			wantHadFirstDamage:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := NewVampire(tt.box)

			if z == nil {
				t.Fatalf("NewVampire(): want valid pointer, got %v", z)
			}

			if z.Enemy.Character.Shape != tt.wantBox {
				t.Errorf("NewVampire(): Shape got %v, want %v", z.Enemy.Character.Shape, tt.wantBox)
			}

			if z.Enemy.Character.Agility != tt.wantAgility {
				t.Errorf("NewVampire(): Agility got %v, want %v", z.Enemy.Character.Agility, tt.wantAgility)
			}

			if z.Enemy.Character.Strength != tt.wantStrength {
				t.Errorf("NewVampire(): Strength got %v, want %v", z.Enemy.Character.Strength, tt.wantStrength)
			}

			if z.Enemy.Character.Health != tt.wantHealth {
				t.Errorf("NewVampire(): Health got %v, want %v", z.Enemy.Character.Health, tt.wantHealth)
			}

			if z.Enemy.Character.MaxHealth != tt.wantMaxHealth {
				t.Errorf("NewVampire(): MaxHealth got %v, want %v", z.Enemy.Character.MaxHealth, tt.wantMaxHealth)
			}

			if uint(z.Enemy.Type) != tt.wantType {
				t.Errorf("NewVampire(): Type got %v, want %v", z.Enemy.Type, tt.wantType)
			}

			if uint(z.Enemy.HostilityRadius) != tt.wantHostilityRadius {
				t.Errorf("NewVampire(): HostilityRadius got %v, want %v", z.Enemy.HostilityRadius, tt.wantHostilityRadius)
			}

			if z.Enemy.IsChasing != tt.wantIsChasing {
				t.Errorf("NewVampire(): IsChasing got %v, want %v", z.Enemy.IsChasing, tt.wantDiretion)
			}

			if uint(z.Enemy.Direction) != tt.wantDiretion {
				t.Errorf("NewVampire(): Direction got %v, want %v", z.Enemy.Direction, tt.wantDiretion)
			}

			if z.HadFirstDamage != tt.wantHadFirstDamage {
				t.Errorf("NewVampire(): HadFirstDamage got %v, want %v", z.HadFirstDamage, tt.wantHadFirstDamage)
			}
		})
	}
}

func TestGhostConstructor(t *testing.T) {
	tests := []struct {
		name                string
		box                 Box
		wantBox             Box
		wantAgility         uint
		wantStrength        uint
		wantHealth          float64
		wantMaxHealth       float64
		wantType            uint
		wantHostilityRadius uint
		wantIsChasing       bool
		wantDiretion        uint
		wantIsVisible       bool
	}{
		{
			name:                "ghost constructor",
			box:                 Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}},
			wantBox:             Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}},
			wantAgility:         75,
			wantStrength:        25,
			wantHealth:          25.0,
			wantMaxHealth:       0.0,
			wantType:            2,
			wantHostilityRadius: 2,
			wantIsChasing:       false,
			wantDiretion:        8,
			wantIsVisible:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := NewGhost(tt.box)

			if z == nil {
				t.Fatalf("NewGhost(): want valid pointer, got %v", z)
			}

			if z.Enemy.Character.Shape != tt.wantBox {
				t.Errorf("NewGhost(): Shape got %v, want %v", z.Enemy.Character.Shape, tt.wantBox)
			}

			if z.Enemy.Character.Agility != tt.wantAgility {
				t.Errorf("NewGhost(): Agility got %v, want %v", z.Enemy.Character.Agility, tt.wantAgility)
			}

			if z.Enemy.Character.Strength != tt.wantStrength {
				t.Errorf("NewGhost(): Strength got %v, want %v", z.Enemy.Character.Strength, tt.wantStrength)
			}

			if z.Enemy.Character.Health != tt.wantHealth {
				t.Errorf("NewGhost(): Health got %v, want %v", z.Enemy.Character.Health, tt.wantHealth)
			}

			if z.Enemy.Character.MaxHealth != tt.wantMaxHealth {
				t.Errorf("NewGhost(): MaxHealth got %v, want %v", z.Enemy.Character.MaxHealth, tt.wantMaxHealth)
			}

			if uint(z.Enemy.Type) != tt.wantType {
				t.Errorf("NewGhost(): Type got %v, want %v", z.Enemy.Type, tt.wantType)
			}

			if uint(z.Enemy.HostilityRadius) != tt.wantHostilityRadius {
				t.Errorf("NewGhost(): HostilityRadius got %v, want %v", z.Enemy.HostilityRadius, tt.wantHostilityRadius)
			}

			if z.Enemy.IsChasing != tt.wantIsChasing {
				t.Errorf("NewGhost(): IsChasing got %v, want %v", z.Enemy.IsChasing, tt.wantDiretion)
			}

			if uint(z.Enemy.Direction) != tt.wantDiretion {
				t.Errorf("NewGhost(): Direction got %v, want %v", z.Enemy.Direction, tt.wantDiretion)
			}

			if z.IsVisible != tt.wantIsVisible {
				t.Errorf("NewGhost(): HadFirstDamage got %v, want %v", z.IsVisible, tt.wantIsVisible)
			}
		})
	}
}

func TestOgreConstructor(t *testing.T) {
	tests := []struct {
		name                string
		box                 Box
		wantBox             Box
		wantAgility         uint
		wantStrength        uint
		wantHealth          float64
		wantMaxHealth       float64
		wantType            uint
		wantHostilityRadius uint
		wantIsChasing       bool
		wantDiretion        uint
		wantIsResting       bool
	}{
		{
			name:                "ogre constructor",
			box:                 Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}},
			wantBox:             Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}},
			wantAgility:         25,
			wantStrength:        100,
			wantHealth:          100.0,
			wantMaxHealth:       0.0,
			wantType:            3,
			wantHostilityRadius: 4,
			wantIsChasing:       false,
			wantDiretion:        8,
			wantIsResting:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := NewOgre(tt.box)

			if z == nil {
				t.Fatalf("NewOgre(): want valid pointer, got %v", z)
			}

			if z.Enemy.Character.Shape != tt.wantBox {
				t.Errorf("NewOgre(): Shape got %v, want %v", z.Enemy.Character.Shape, tt.wantBox)
			}

			if z.Enemy.Character.Agility != tt.wantAgility {
				t.Errorf("NewOgre(): Agility got %v, want %v", z.Enemy.Character.Agility, tt.wantAgility)
			}

			if z.Enemy.Character.Strength != tt.wantStrength {
				t.Errorf("NewOgre(): Strength got %v, want %v", z.Enemy.Character.Strength, tt.wantStrength)
			}

			if z.Enemy.Character.Health != tt.wantHealth {
				t.Errorf("NewOgre(): Health got %v, want %v", z.Enemy.Character.Health, tt.wantHealth)
			}

			if z.Enemy.Character.MaxHealth != tt.wantMaxHealth {
				t.Errorf("NewOgre(): MaxHealth got %v, want %v", z.Enemy.Character.MaxHealth, tt.wantMaxHealth)
			}

			if uint(z.Enemy.Type) != tt.wantType {
				t.Errorf("NewOgre(): Type got %v, want %v", z.Enemy.Type, tt.wantType)
			}

			if uint(z.Enemy.HostilityRadius) != tt.wantHostilityRadius {
				t.Errorf("NewOgre(): HostilityRadius got %v, want %v", z.Enemy.HostilityRadius, tt.wantHostilityRadius)
			}

			if z.Enemy.IsChasing != tt.wantIsChasing {
				t.Errorf("NewOgre(): IsChasing got %v, want %v", z.Enemy.IsChasing, tt.wantDiretion)
			}

			if uint(z.Enemy.Direction) != tt.wantDiretion {
				t.Errorf("NewOgre(): Direction got %v, want %v", z.Enemy.Direction, tt.wantDiretion)
			}

			if z.IsResting != tt.wantIsResting {
				t.Errorf("NewOgre(): HadFirstDamage got %v, want %v", z.IsResting, tt.wantIsResting)
			}
		})
	}
}

func TestSnakeMageConstructor(t *testing.T) {
	tests := []struct {
		name                string
		box                 Box
		wantBox             Box
		wantAgility         uint
		wantStrength        uint
		wantHealth          float64
		wantMaxHealth       float64
		wantType            uint
		wantHostilityRadius uint
		wantIsChasing       bool
		wantDiretion        uint
	}{
		{
			name:                "snake mage constructor",
			box:                 Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}},
			wantBox:             Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}},
			wantAgility:         100,
			wantStrength:        50,
			wantHealth:          75.0,
			wantMaxHealth:       0.0,
			wantType:            4,
			wantHostilityRadius: 6,
			wantIsChasing:       false,
			wantDiretion:        8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := NewSnakeMage(tt.box)

			if z == nil {
				t.Fatalf("NewSnakeMage(): want valid pointer, got %v", z)
			}

			if z.Enemy.Character.Shape != tt.wantBox {
				t.Errorf("NewSnakeMage(): Shape got %v, want %v", z.Enemy.Character.Shape, tt.wantBox)
			}

			if z.Enemy.Character.Agility != tt.wantAgility {
				t.Errorf("NewSnakeMage(): Agility got %v, want %v", z.Enemy.Character.Agility, tt.wantAgility)
			}

			if z.Enemy.Character.Strength != tt.wantStrength {
				t.Errorf("NewSnakeMage(): Strength got %v, want %v", z.Enemy.Character.Strength, tt.wantStrength)
			}

			if z.Enemy.Character.Health != tt.wantHealth {
				t.Errorf("NewSnakeMage(): Health got %v, want %v", z.Enemy.Character.Health, tt.wantHealth)
			}

			if z.Enemy.Character.MaxHealth != tt.wantMaxHealth {
				t.Errorf("NewSnakeMage(): MaxHealth got %v, want %v", z.Enemy.Character.MaxHealth, tt.wantMaxHealth)
			}

			if uint(z.Enemy.Type) != tt.wantType {
				t.Errorf("NewSnakeMage(): Type got %v, want %v", z.Enemy.Type, tt.wantType)
			}

			if uint(z.Enemy.HostilityRadius) != tt.wantHostilityRadius {
				t.Errorf("NewSnakeMage(): HostilityRadius got %v, want %v", z.Enemy.HostilityRadius, tt.wantHostilityRadius)
			}

			if z.Enemy.IsChasing != tt.wantIsChasing {
				t.Errorf("NewSnakeMage(): IsChasing got %v, want %v", z.Enemy.IsChasing, tt.wantDiretion)
			}

			if uint(z.Enemy.Direction) != tt.wantDiretion {
				t.Errorf("NewSnakeMage(): Direction got %v, want %v", z.Enemy.Direction, tt.wantDiretion)
			}
		})
	}
}
