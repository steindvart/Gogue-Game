package entity

import (
	"slices"
	"testing"
)

func TestNewWeapon(t *testing.T) {
	tests := []struct {
		name               string
		box                Box
		wantShape          Box
		wantNames          []string
		wantDamageLessThan uint
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
			wantDamageLessThan: uint(IncreaseAttributeBaseParcentage + IncreaseAttributeMaxPercentage + 1),
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
			if !(got.Damage < tt.wantDamageLessThan) {
				t.Errorf("NewWeapon(): got Damage %v, want less than %v", got.Damage, tt.wantDamageLessThan)
			}
		})
	}
}

func TestWeapon_Taken(t *testing.T) {
	tests := []struct {
		name string
		want Box
	}{
		{
			name: "weapon is taken",
			want: Box{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weapon := &Weapon{
				Consumable: Consumable{
					Shape: Box{
						Point: Point2D[int]{X: 1, Y: 2},
						Size:  Size2D[uint]{Height: 1, Width: 1},
					},
					Name: "Awkward Weapon",
				},
				Damage: 5,
			}
			weapon.Taken()
			if weapon.Consumable.Shape != tt.want {
				t.Errorf("Taken() = (%v), want (%v)", weapon.Consumable.Shape, tt.want)
			}
		})
	}
}

func TestWeapon_Dropped(t *testing.T) {
	tests := []struct {
		name string
		box  Box
		want Box
	}{
		{
			name: "weapon is dropped",
			box:  Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
			want: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weapon := &Weapon{
				Consumable: Consumable{
					Shape: Box{},
					Name:  "Awkward Weapon",
				},
				Damage: 5,
			}
			item := weapon.Dropped(tt.box)

			if item != weapon {
				t.Errorf("Dropped() = %#v, want %#v", item, weapon)
			}

			if weapon.Consumable.Shape != tt.want {
				t.Errorf("Dropped(): Shape %v, want %v", weapon.Consumable.Shape, tt.want)
			}
		})
	}
}

func TestWeapon_Use(t *testing.T) {
	tests := []struct {
		name   string
		weapon *Weapon
		want   string
	}{
		{
			name: "use weapon",
			weapon: &Weapon{
				Consumable: Consumable{
					Shape: Box{},
					Name:  "Awkward Weapon",
				},
				Damage: 5,
			},
			want: "You picked up the Awkward Weapon, now all your attacks have 5 extra damage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(Box{})
			currentWeapon := player.Weapon

			line, ptr := tt.weapon.Use(player)

			if ptr != currentWeapon {
				t.Errorf("Use(): ConsumableLike pointer got %#v, want %#v", ptr, nil)
			}

			if !(line == tt.want) {
				t.Errorf("Use() = (%v), want (%v)", line, tt.want)
			}
		})
	}
}
