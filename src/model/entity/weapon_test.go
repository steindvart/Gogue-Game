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
			if !(got.Item.Shape == tt.wantShape) {
				t.Errorf("NewWeapon(): got Shape %v, want %v", got.Item.Shape, tt.wantShape)
			}
			if !(slices.Contains(tt.wantNames, got.Item.Name)) {
				t.Errorf("NewWeapon(): got Name %v, want in %v", got.Item.Name, tt.wantNames)
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
				Item: Item{
					Shape: Box{
						Point: Point2D[int]{X: 1, Y: 2},
						Size:  Size2D[uint]{Height: 1, Width: 1},
					},
					Name: "Awkward Weapon",
				},
				Damage: 5,
			}
			weapon.Taken()
			if weapon.Item.Shape != tt.want {
				t.Errorf("Taken() = (%v), want (%v)", weapon.Item.Shape, tt.want)
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
				Item: Item{
					Shape: Box{},
					Name:  "Awkward Weapon",
				},
				Damage: 5,
			}

			weapon.Dropped(tt.box)

			if weapon.Item.Shape != tt.want {
				t.Errorf("Dropped(): Shape %v, want %v", weapon.Item.Shape, tt.want)
			}
		})
	}
}

func TestWeapon_Use(t *testing.T) {
	tests := []struct {
		name    string
		weapons []*Weapon
		want    string
	}{
		{
			name: "use weapon",
			weapons: []*Weapon{
				{
					Item: Item{
						Shape: Box{},
						Name:  "Awkward Weapon 1",
					},
					Damage: 5,
				},
				{
					Item: Item{
						Shape: Box{},
						Name:  "Awkward Weapon 2",
					},
					Damage: 5,
				},
			},
			want: "You picked up the Awkward Weapon 2, now all your attacks have 5 extra damage",
		},
		{
			name: "use first weapon",
			weapons: []*Weapon{
				{
					Item: Item{
						Shape: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
						Name:  "Awkward Weapon",
					},
					Damage: 5,
				},
			},
			want: "You picked up the Awkward Weapon, now all your attacks have 5 extra damage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}})

			for idx := range tt.weapons {
				_ = player.TakeItem(tt.weapons[idx])
			}

			if len(tt.weapons) > 1 {
				player.Weapon = tt.weapons[0]
				currentWeapon := player.Weapon
				if currentWeapon.Item.Shape == player.Character.Shape {
					t.Errorf("Use(): Shape got %#v, got %#v", currentWeapon.Item.Shape, player.Character.Shape)
				}

				line := tt.weapons[1].Use(player)

				if !(currentWeapon.Item.Shape == player.Character.Shape) {
					t.Errorf("Use(): Shape got %#v, got %#v", currentWeapon.Item.Shape, player.Character.Shape)
				}

				if !(player.Weapon == tt.weapons[1]) {
					t.Errorf("Use(): Weapon got %#v, got %#v", player.Weapon, tt.weapons[1])
				}

				if !(line == tt.want) {
					t.Errorf("Use() = (%v), want (%v)", line, tt.want)
				}
			} else {
				line := tt.weapons[0].Use(player)

				if !(player.Weapon == tt.weapons[0]) {
					t.Errorf("Use(): Weapon got %#v, want %#v", player.Weapon, tt.weapons[0])
				}

				if !(line == tt.want) {
					t.Errorf("Use() = (%v), want (%v)", line, tt.want)
				}
			}
		})
	}
}
