package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
	"testing"
)

const defaultWeaponTestSeed int64 = 42

func TestWeapon_NewWeapon_BuiltinConfig(t *testing.T) {
	tests := []struct {
		name         string
		weaponType   WeaponType
		expectedName string
	}{
		{
			name:         "Dagger type uses Dagger weapon config",
			weaponType:   WeaponTypeDagger,
			expectedName: string(WeaponTypeDagger),
		},
		{
			name:         "Spear type uses Spear weapon config",
			weaponType:   WeaponTypeSpear,
			expectedName: string(WeaponTypeSpear),
		},
		{
			name:         "Sword type uses Sword weapon config",
			weaponType:   WeaponTypeSword,
			expectedName: string(WeaponTypeSword),
		},
		{
			name:         "Axe type uses Axe weapon config",
			weaponType:   WeaponTypeAxe,
			expectedName: string(WeaponTypeAxe),
		},
		{
			name:         "Maul type uses Maul weapon config",
			weaponType:   WeaponTypeMaul,
			expectedName: string(WeaponTypeMaul),
		},
		{
			name:         "Mystery type uses Mystery weapon config",
			weaponType:   WeaponTypeMystery,
			expectedName: string(WeaponTypeMystery),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomWithSeed(defaultWeaponTestSeed)
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

			weapon := NewWeaponBuiltin(rng, box, tt.weaponType)

			if weapon == nil {
				t.Fatal("Expected valid weapon, got nil")
			}

			if weapon.Type != tt.weaponType {
				t.Errorf("Expected type %q, got %q", tt.weaponType, weapon.Type)
			}

			if weapon.Name != tt.expectedName {
				t.Errorf("Expected name %q, got %q", tt.expectedName, weapon.Name)
			}

			// Verify attributes match the config's ranges
			config := GetWeaponConfig(tt.weaponType)
			if weapon.Effect.Attributes.Strength < config.StrengthRange.Min ||
				weapon.Effect.Attributes.Strength > config.StrengthRange.Max {
				t.Errorf("Strength out of config range: got %.2f, want [%.2f, %.2f]",
					weapon.Effect.Attributes.Strength, config.StrengthRange.Min, config.StrengthRange.Max)
			}

			if weapon.Effect.Attributes.Agility < config.AgilityRange.Min ||
				weapon.Effect.Attributes.Agility > config.AgilityRange.Max {
				t.Errorf("Agility out of config range: got %.2f, want [%.2f, %.2f]",
					weapon.Effect.Attributes.Agility, config.AgilityRange.Min, config.AgilityRange.Max)
			}

			if weapon.Name != string(config.Type) {
				t.Errorf("Expected name %q, got %q", string(config.Type), weapon.Name)
			}
		})
	}
}

func TestWeapon_NewWeapon_CustomConfig(t *testing.T) {
	tests := []struct {
		name      string
		seed      int64
		config    WeaponConfig
		box       primitives.Box
		wantError bool
	}{
		{
			name: "Valid custom config",
			seed: defaultWeaponTestSeed,
			config: WeaponConfig{
				Type:          "Custom Weapon",
				StrengthRange: primitives.AttributeRange{Min: 10, Max: 50},
				AgilityRange:  primitives.AttributeRange{Min: -5, Max: 15},
				Description:   "A custom test weapon",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 3, Y: 3}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: false,
		},
		{
			name: "Valid config with negative agility",
			seed: 100,
			config: WeaponConfig{
				Type:          "Heavy Custom Weapon",
				StrengthRange: primitives.AttributeRange{Min: 30, Max: 60},
				AgilityRange:  primitives.AttributeRange{Min: -20, Max: -5},
				Description:   "A heavy custom weapon",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 1}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: false,
		},
		{
			name: "Invalid strength range returns error",
			seed: defaultWeaponTestSeed,
			config: WeaponConfig{
				Type:          "Invalid Weapon",
				StrengthRange: primitives.AttributeRange{Min: 50, Max: 10}, // Min > Max
				AgilityRange:  primitives.AttributeRange{Min: 0, Max: 10},
				Description:   "Invalid config",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: true,
		},
		{
			name: "Invalid agility range returns error",
			seed: defaultWeaponTestSeed,
			config: WeaponConfig{
				Type:          "Invalid Weapon",
				StrengthRange: primitives.AttributeRange{Min: 0, Max: 10},
				AgilityRange:  primitives.AttributeRange{Min: 20, Max: 5}, // Min > Max
				Description:   "Invalid config",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomWithSeed(tt.seed)
			weapon, err := NewWeaponByConfig(rng, tt.box, tt.config)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if weapon != nil {
					t.Errorf("Expected nil weapon on error, got %v", weapon)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if weapon == nil {
				t.Fatal("Expected valid weapon, got nil")
			}

			if weapon.Type != tt.config.Type {
				t.Errorf("Expected type %q, got %q", tt.config.Type, weapon.Type)
			}

			// Verify attributes are within config ranges
			if weapon.Effect.Attributes.Strength < tt.config.StrengthRange.Min ||
				weapon.Effect.Attributes.Strength > tt.config.StrengthRange.Max {
				t.Errorf("Strength out of range: got %.2f, want [%.2f, %.2f]",
					weapon.Effect.Attributes.Strength, tt.config.StrengthRange.Min, tt.config.StrengthRange.Max)
			}

			if weapon.Effect.Attributes.Agility < tt.config.AgilityRange.Min ||
				weapon.Effect.Attributes.Agility > tt.config.AgilityRange.Max {
				t.Errorf("Agility out of range: got %.2f, want [%.2f, %.2f]",
					weapon.Effect.Attributes.Agility, tt.config.AgilityRange.Min, tt.config.AgilityRange.Max)
			}

			// Verify name matches config type
			if weapon.Name != string(tt.config.Type) {
				t.Errorf("Expected name %q, got %q", tt.config.Type, weapon.Name)
			}
		})
	}
}

func TestWeapon_NewWeapon_Determinism(t *testing.T) {
	tests := []struct {
		name       string
		seed       int64
		weaponType WeaponType
	}{
		{
			name:       "Same seed produces same Sword weapon",
			seed:       defaultWeaponTestSeed,
			weaponType: WeaponTypeSword,
		},
		{
			name:       "Same seed produces same Dagger weapon",
			seed:       defaultWeaponTestSeed,
			weaponType: WeaponTypeDagger,
		},
		{
			name:       "Same seed produces same Mystery weapon",
			seed:       defaultWeaponTestSeed,
			weaponType: WeaponTypeMystery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

			// Create first weapon
			rng1 := utils.NewRandomWithSeed(tt.seed)
			weapon1 := NewWeaponBuiltin(rng1, box, tt.weaponType)

			// Create second weapon with same seed
			rng2 := utils.NewRandomWithSeed(tt.seed)
			weapon2 := NewWeaponBuiltin(rng2, box, tt.weaponType)

			// Verify they are identical
			if weapon1.Effect.Attributes.Strength != weapon2.Effect.Attributes.Strength {
				t.Errorf("Strength mismatch: %f != %f", weapon1.Effect.Attributes.Strength, weapon2.Effect.Attributes.Strength)
			}

			if weapon1.Effect.Attributes.Agility != weapon2.Effect.Attributes.Agility {
				t.Errorf("Agility mismatch: %f != %f", weapon1.Effect.Attributes.Agility, weapon2.Effect.Attributes.Agility)
			}
		})
	}
}

func TestWeapon_NewWeapon_Randomness(t *testing.T) {
	tests := []struct {
		name       string
		weaponType WeaponType
		iterations int
	}{
		{
			name:       "Sword weapon produces varied values",
			weaponType: WeaponTypeSword,
			iterations: 50,
		},
		{
			name:       "Mystery weapon produces varied values",
			weaponType: WeaponTypeMystery,
			iterations: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
			strengthValues := make(map[float64]bool)
			agilityValues := make(map[float64]bool)

			// Generate multiple weapons with different seeds
			for i := 0; i < tt.iterations; i++ {
				rng := utils.NewRandomWithSeed(int64(i))
				weapon := NewWeaponBuiltin(rng, box, tt.weaponType)

				strengthValues[weapon.Effect.Attributes.Strength] = true
				agilityValues[weapon.Effect.Attributes.Agility] = true
			}

			// Check that we got varied values (at least 5 different values)
			minUniqueValues := 5

			// For attributes, check based on type
			config := GetWeaponConfig(tt.weaponType)
			if config.StrengthRange.Min != config.StrengthRange.Max {
				if len(strengthValues) < minUniqueValues {
					t.Errorf("Expected at least %d unique strength values, got %d", minUniqueValues, len(strengthValues))
				}
			}

			if config.AgilityRange.Min != config.AgilityRange.Max {
				if len(agilityValues) < minUniqueValues {
					t.Errorf("Expected at least %d unique agility values, got %d", minUniqueValues, len(agilityValues))
				}
			}
		})
	}
}

func TestWeapon_NewWeapon_ZeroSizedBoxIsOk(t *testing.T) {
	rng := utils.NewRandomWithSeed(defaultWeaponTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 0, Height: 0}}
	weapon := NewWeaponBuiltin(rng, box, WeaponTypeSword)

	if weapon == nil {
		t.Fatal("Expected valid weapon with zero-sized box")
	}

	if *weapon.Box != box {
		t.Errorf("Expected box %v, got %v", box, weapon.Box)
	}
}

func TestWeapon_Drop(t *testing.T) {
	rng := utils.NewRandomWithSeed(defaultWeaponTestSeed)
	initialBox := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	weapon := NewWeaponBuiltin(rng, initialBox, WeaponTypeSword)

	if *weapon.Box != initialBox {
		t.Errorf("Expected initial box %v, got %v", initialBox, weapon.Box)
	}

	newPosition := primitives.Point2D[int]{X: 10, Y: 10}
	resultBox := weapon.Drop(newPosition)

	if weapon.Box.Point != newPosition {
		t.Errorf("Expected position to be updated to %v, got %v", newPosition, weapon.Box.Point)
	}

	if resultBox.Point != newPosition {
		t.Errorf("Expected returned box to have position %v, got %v", newPosition, resultBox.Point)
	}
}

func TestWeapon_Use(t *testing.T) {
	tests := []struct {
		name       string
		seed       int64
		weaponType WeaponType
	}{
		{
			name:       "Use Sword weapon returns correct attributes",
			seed:       defaultWeaponTestSeed,
			weaponType: WeaponTypeSword,
		},
		{
			name:       "Use Dagger weapon returns correct attributes",
			seed:       1,
			weaponType: WeaponTypeDagger,
		},
		{
			name:       "Use Maul weapon returns correct attributes",
			seed:       2,
			weaponType: WeaponTypeMaul,
		},
		{
			name:       "Use Mystery weapon returns correct attributes",
			seed:       3,
			weaponType: WeaponTypeMystery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomWithSeed(tt.seed)
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
			weapon := NewWeaponBuiltin(rng, box, tt.weaponType)

			effect := weapon.Use()

			// Verify returned attributes match stored attributes
			if effect.Attributes.Strength != weapon.Effect.Attributes.Strength {
				t.Errorf("Expected strength %.2f, got %.2f", weapon.Effect.Attributes.Strength, effect.Attributes.Strength)
			}

			if effect.Attributes.Agility != weapon.Effect.Attributes.Agility {
				t.Errorf("Expected agility %.2f, got %.2f", weapon.Effect.Attributes.Agility, effect.Attributes.Agility)
			}
		})
	}
}

func TestWeapon_UseMultipleTimesIsOk(t *testing.T) {
	rng := utils.NewRandomWithSeed(defaultWeaponTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	weapon := NewWeaponBuiltin(rng, box, WeaponTypeSword)

	effect1 := weapon.Use()
	effect2 := weapon.Use()

	if effect1.Attributes.Strength != effect2.Attributes.Strength {
		t.Error("Use() should return consistent attributes on multiple calls")
	}
	if effect1.Attributes.Agility != effect2.Attributes.Agility {
		t.Error("Use() should return consistent attributes on multiple calls")
	}
}

// Benchmarks
func BenchmarkWeapon_NewWeapon(b *testing.B) {
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

	b.Run("Sword", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomWithSeed(int64(i))
			_ = NewWeaponBuiltin(rng, box, WeaponTypeSword)
		}
	})

	b.Run("Dagger", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomWithSeed(int64(i))
			_ = NewWeaponBuiltin(rng, box, WeaponTypeDagger)
		}
	})

	b.Run("Maul", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomWithSeed(int64(i))
			_ = NewWeaponBuiltin(rng, box, WeaponTypeMaul)
		}
	})

	b.Run("Mystery", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomWithSeed(int64(i))
			_ = NewWeaponBuiltin(rng, box, WeaponTypeMystery)
		}
	})
}

func BenchmarkWeapon_NewWeaponByConfig(b *testing.B) {
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	config := GetWeaponConfig(WeaponTypeSword)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng := utils.NewRandomWithSeed(int64(i))
		_, _ = NewWeaponByConfig(rng, box, config)
	}
}

func BenchmarkWeapon_Use(b *testing.B) {
	rng := utils.NewRandomWithSeed(defaultWeaponTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	weapon := NewWeaponBuiltin(rng, box, WeaponTypeSword)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = weapon.Use()
	}
}
