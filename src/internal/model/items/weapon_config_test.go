package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
	"testing"
)

const defaultWeaponConfigTestSeed int64 = 42

func TestGetWeaponConfig(t *testing.T) {
	tests := []struct {
		name              string
		weaponType        WeaponType
		wantType          WeaponType
		wantStrengthRange primitives.AttributeRange
		wantAgilityRange  primitives.AttributeRange
	}{
		{
			name:              "Get Dagger config",
			weaponType:        WeaponTypeDagger,
			wantType:          WeaponTypeDagger,
			wantStrengthRange: primitives.AttributeRange{Min: 1, Max: 5},
			wantAgilityRange:  primitives.AttributeRange{Min: 3, Max: 8},
		},
		{
			name:              "Get Spear config",
			weaponType:        WeaponTypeSpear,
			wantType:          WeaponTypeSpear,
			wantStrengthRange: primitives.AttributeRange{Min: 3, Max: 10},
			wantAgilityRange:  primitives.AttributeRange{Min: 1, Max: 5},
		},
		{
			name:              "Get Sword config",
			weaponType:        WeaponTypeSword,
			wantType:          WeaponTypeSword,
			wantStrengthRange: primitives.AttributeRange{Min: 7, Max: 13},
			wantAgilityRange:  primitives.AttributeRange{Min: -2, Max: 3},
		},
		{
			name:              "Get Axe config",
			weaponType:        WeaponTypeAxe,
			wantType:          WeaponTypeAxe,
			wantStrengthRange: primitives.AttributeRange{Min: 9, Max: 19},
			wantAgilityRange:  primitives.AttributeRange{Min: -7, Max: -3},
		},
		{
			name:              "Get Maul config",
			weaponType:        WeaponTypeMaul,
			wantType:          WeaponTypeMaul,
			wantStrengthRange: primitives.AttributeRange{Min: 11, Max: 23},
			wantAgilityRange:  primitives.AttributeRange{Min: -12, Max: -4},
		},
		{
			name:              "Get Mystery config",
			weaponType:        WeaponTypeMystery,
			wantType:          WeaponTypeMystery,
			wantStrengthRange: primitives.AttributeRange{Min: 1, Max: 26},
			wantAgilityRange:  primitives.AttributeRange{Min: -15, Max: 15},
		},
		{
			name:              "Unknown type returns Mystery",
			weaponType:        "Unknown Weapon",
			wantType:          WeaponTypeMystery,
			wantStrengthRange: primitives.AttributeRange{Min: 1, Max: 26},
			wantAgilityRange:  primitives.AttributeRange{Min: -15, Max: 15},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := GetWeaponConfig(tt.weaponType)

			if config.Type != tt.wantType {
				t.Errorf("Expected type %q, got %q", tt.wantType, config.Type)
			}

			if config.StrengthRange != tt.wantStrengthRange {
				t.Errorf("Expected StrengthRange %v, got %v", tt.wantStrengthRange, config.StrengthRange)
			}

			if config.AgilityRange != tt.wantAgilityRange {
				t.Errorf("Expected AgilityRange %v, got %v", tt.wantAgilityRange, config.AgilityRange)
			}

			// Validate the config
			if err := config.Validate(); err != nil {
				t.Errorf("Config validation failed: %v", err)
			}

			// Check description is not empty
			if config.Description == "" {
				t.Errorf("Config for %q has empty description", tt.weaponType)
			}
		})
	}
}

func TestWeaponConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    WeaponConfig
		wantError bool
	}{
		{
			name: "Valid config",
			config: WeaponConfig{
				Type:          WeaponTypeSword,
				StrengthRange: primitives.AttributeRange{Min: 10, Max: 25},
				AgilityRange:  primitives.AttributeRange{Min: -2, Max: 5},
			},
			wantError: false,
		},
		{
			name: "Valid config with negative agility",
			config: WeaponConfig{
				Type:          WeaponTypeMaul,
				StrengthRange: primitives.AttributeRange{Min: 20, Max: 40},
				AgilityRange:  primitives.AttributeRange{Min: -15, Max: -5},
			},
			wantError: false,
		},
		{
			name: "Invalid strength range is error",
			config: WeaponConfig{
				Type:          WeaponTypeSword,
				StrengthRange: primitives.AttributeRange{Min: 25, Max: 10}, // Min > Max
				AgilityRange:  primitives.AttributeRange{Min: 0, Max: 5},
			},
			wantError: true,
		},
		{
			name: "Invalid agility range is error",
			config: WeaponConfig{
				Type:          WeaponTypeDagger,
				StrengthRange: primitives.AttributeRange{Min: 5, Max: 15},
				AgilityRange:  primitives.AttributeRange{Min: 18, Max: 8}, // Min > Max
			},
			wantError: true,
		},
		{
			name: "Zero range is valid",
			config: WeaponConfig{
				Type:          "Test Weapon",
				StrengthRange: primitives.AttributeRange{Min: 10, Max: 10},
				AgilityRange:  primitives.AttributeRange{Min: 0, Max: 0},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantError && err == nil {
				t.Error("Expected error, got nil")
			}

			if !tt.wantError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}

func TestWeaponRegistry_AllConfigsValid(t *testing.T) {
	for weaponType, config := range WeaponRegistry {
		t.Run(string(weaponType), func(t *testing.T) {
			if err := config.Validate(); err != nil {
				t.Errorf("Registry config for %q is invalid: %v", weaponType, err)
			}

			if config.Type != weaponType {
				t.Errorf("Config type mismatch: registry key=%q, config.Type=%q", weaponType, config.Type)
			}

			if config.Description == "" {
				t.Errorf("Config for %q has empty description", weaponType)
			}
		})
	}
}

func TestWeaponConfig_GenerateAttributes(t *testing.T) {
	tests := []struct {
		name   string
		seed   int64
		config WeaponConfig
	}{
		{
			name: "Sword config generates deterministic attributes",
			seed: defaultWeaponConfigTestSeed,
			config: WeaponConfig{
				Type:          WeaponTypeSword,
				StrengthRange: primitives.AttributeRange{Min: 10, Max: 25},
				AgilityRange:  primitives.AttributeRange{Min: -2, Max: 5},
			},
		},
		{
			name: "Dagger config generates light weapon attributes",
			seed: 100,
			config: WeaponConfig{
				Type:          WeaponTypeDagger,
				StrengthRange: primitives.AttributeRange{Min: 5, Max: 15},
				AgilityRange:  primitives.AttributeRange{Min: 8, Max: 18},
			},
		},
		{
			name: "Maul config generates heavy weapon attributes",
			seed: 200,
			config: WeaponConfig{
				Type:          WeaponTypeMaul,
				StrengthRange: primitives.AttributeRange{Min: 20, Max: 40},
				AgilityRange:  primitives.AttributeRange{Min: -15, Max: -5},
			},
		},
		{
			name: "Mystery config generates varied attributes",
			seed: 300,
			config: WeaponConfig{
				Type:          WeaponTypeMystery,
				StrengthRange: primitives.AttributeRange{Min: 1, Max: 50},
				AgilityRange:  primitives.AttributeRange{Min: -20, Max: 20},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate twice with same seed
			rng1 := utils.NewRandomWithSeed(tt.seed)
			attrs1 := tt.config.GenerateAttributes(rng1)

			rng2 := utils.NewRandomWithSeed(tt.seed)
			attrs2 := tt.config.GenerateAttributes(rng2)

			// Should be deterministic
			if attrs1.Strength != attrs2.Strength {
				t.Errorf("Strength not deterministic: %f != %f", attrs1.Strength, attrs2.Strength)
			}
			if attrs1.Agility != attrs2.Agility {
				t.Errorf("Agility not deterministic: %f != %f", attrs1.Agility, attrs2.Agility)
			}

			// Should be within range
			if attrs1.Strength < tt.config.StrengthRange.Min || attrs1.Strength > tt.config.StrengthRange.Max {
				t.Errorf("Strength out of range: got %.2f, want [%.2f, %.2f]",
					attrs1.Strength, tt.config.StrengthRange.Min, tt.config.StrengthRange.Max)
			}
			if attrs1.Agility < tt.config.AgilityRange.Min || attrs1.Agility > tt.config.AgilityRange.Max {
				t.Errorf("Agility out of range: got %.2f, want [%.2f, %.2f]",
					attrs1.Agility, tt.config.AgilityRange.Min, tt.config.AgilityRange.Max)
			}
		})
	}
}

func TestWeaponConfig_GenerateAttributes_Randomness(t *testing.T) {
	config := WeaponConfig{
		Type:          WeaponTypeMystery,
		StrengthRange: primitives.AttributeRange{Min: 1, Max: 50},
		AgilityRange:  primitives.AttributeRange{Min: -20, Max: 20},
	}

	strengthValues := make(map[float64]bool)
	agilityValues := make(map[float64]bool)
	iterations := 50

	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomWithSeed(int64(i))
		attrs := config.GenerateAttributes(rng)
		strengthValues[attrs.Strength] = true
		agilityValues[attrs.Agility] = true
	}

	minUniqueValues := 10
	if len(strengthValues) < minUniqueValues {
		t.Errorf("Expected at least %d unique strength values, got %d", minUniqueValues, len(strengthValues))
	}
	if len(agilityValues) < minUniqueValues {
		t.Errorf("Expected at least %d unique agility values, got %d", minUniqueValues, len(agilityValues))
	}
}

func TestWeaponConfig_GenerateAttributes_LightWeaponBalance(t *testing.T) {
	config := GetWeaponConfig(WeaponTypeDagger)

	iterations := 30
	strengthSum := 0.0
	agilitySum := 0.0

	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomWithSeed(int64(i))
		attrs := config.GenerateAttributes(rng)
		strengthSum += attrs.Strength
		agilitySum += attrs.Agility
	}

	avgStrength := strengthSum / float64(iterations)
	avgAgility := agilitySum / float64(iterations)

	// Light weapons should have lower average strength than agility
	if avgStrength >= avgAgility {
		t.Errorf("Light weapon (Dagger) should have higher agility than strength on average: strength=%.2f, agility=%.2f",
			avgStrength, avgAgility)
	}

	// Both should be positive for light weapons
	if avgStrength < 0 || avgAgility < 0 {
		t.Errorf("Light weapon should have positive average attributes: strength=%.2f, agility=%.2f",
			avgStrength, avgAgility)
	}
}

func TestWeaponConfig_GenerateAttributes_HeavyWeaponBalance(t *testing.T) {
	config := GetWeaponConfig(WeaponTypeMaul)

	iterations := 30
	strengthSum := 0.0
	agilitySum := 0.0

	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomWithSeed(int64(i))
		attrs := config.GenerateAttributes(rng)
		strengthSum += attrs.Strength
		agilitySum += attrs.Agility
	}

	avgStrength := strengthSum / float64(iterations)
	avgAgility := agilitySum / float64(iterations)

	// Heavy weapons should have much higher strength than agility
	if avgStrength <= avgAgility {
		t.Errorf("Heavy weapon (Maul) should have much higher strength than agility on average: strength=%.2f, agility=%.2f",
			avgStrength, avgAgility)
	}

	// Strength should be positive, agility should be negative for heavy weapons
	if avgStrength < 0 {
		t.Errorf("Heavy weapon should have positive average strength: %.2f", avgStrength)
	}
	if avgAgility > 0 {
		t.Errorf("Heavy weapon should have negative average agility: %.2f", avgAgility)
	}
}

func TestWeaponConfig_GenerateAttributes_FixedRanges(t *testing.T) {
	tests := []struct {
		name   string
		config WeaponConfig
	}{
		{
			name: "Fixed strength, varied agility",
			config: WeaponConfig{
				Type:          "Test Weapon 1",
				StrengthRange: primitives.AttributeRange{Min: 10, Max: 10},
				AgilityRange:  primitives.AttributeRange{Min: -5, Max: 5},
			},
		},
		{
			name: "Varied strength, fixed agility",
			config: WeaponConfig{
				Type:          "Test Weapon 2",
				StrengthRange: primitives.AttributeRange{Min: 5, Max: 20},
				AgilityRange:  primitives.AttributeRange{Min: 0, Max: 0},
			},
		},
		{
			name: "Both fixed",
			config: WeaponConfig{
				Type:          "Test Weapon 3",
				StrengthRange: primitives.AttributeRange{Min: 15, Max: 15},
				AgilityRange:  primitives.AttributeRange{Min: -3, Max: -3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomWithSeed(defaultWeaponConfigTestSeed)
			attrs := tt.config.GenerateAttributes(rng)

			// Check if ranges are fixed
			if tt.config.StrengthRange.Min == tt.config.StrengthRange.Max {
				if attrs.Strength != tt.config.StrengthRange.Min {
					t.Errorf("Expected fixed strength %.2f, got %.2f",
						tt.config.StrengthRange.Min, attrs.Strength)
				}
			}

			if tt.config.AgilityRange.Min == tt.config.AgilityRange.Max {
				if attrs.Agility != tt.config.AgilityRange.Min {
					t.Errorf("Expected fixed agility %.2f, got %.2f",
						tt.config.AgilityRange.Min, attrs.Agility)
				}
			}
		})
	}
}

func TestWeaponConfig_GenerateAttributes_LargeRanges(t *testing.T) {
	config := WeaponConfig{
		Type:          "Wide Range Weapon",
		StrengthRange: primitives.AttributeRange{Min: -100, Max: 100},
		AgilityRange:  primitives.AttributeRange{Min: -100, Max: 100},
	}

	iterations := 100
	strengthValues := make(map[float64]bool)
	agilityValues := make(map[float64]bool)

	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomWithSeed(int64(i))
		attrs := config.GenerateAttributes(rng)

		if attrs.Strength < config.StrengthRange.Min || attrs.Strength > config.StrengthRange.Max {
			t.Errorf("Strength out of range: got %.2f, want [%.2f, %.2f]",
				attrs.Strength, config.StrengthRange.Min, config.StrengthRange.Max)
		}

		if attrs.Agility < config.AgilityRange.Min || attrs.Agility > config.AgilityRange.Max {
			t.Errorf("Agility out of range: got %.2f, want [%.2f, %.2f]",
				attrs.Agility, config.AgilityRange.Min, config.AgilityRange.Max)
		}

		strengthValues[attrs.Strength] = true
		agilityValues[attrs.Agility] = true
	}

	// Should have good distribution
	minUniqueValues := 30
	if len(strengthValues) < minUniqueValues {
		t.Errorf("Expected at least %d unique strength values in large range, got %d",
			minUniqueValues, len(strengthValues))
	}
	if len(agilityValues) < minUniqueValues {
		t.Errorf("Expected at least %d unique agility values in large range, got %d",
			minUniqueValues, len(agilityValues))
	}
}

// Benchmarks
func BenchmarkWeaponConfig_GenerateAttributes(b *testing.B) {
	config := GetWeaponConfig(WeaponTypeSword)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng := utils.NewRandomWithSeed(int64(i))
		_ = config.GenerateAttributes(rng)
	}
}

func BenchmarkWeaponConfig_Validate(b *testing.B) {
	config := GetWeaponConfig(WeaponTypeMystery)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}

func BenchmarkGetWeaponConfig(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetWeaponConfig(WeaponTypeSword)
	}
}
