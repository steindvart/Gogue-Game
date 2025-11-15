package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
	"testing"
)

func TestGetScrollConfig(t *testing.T) {
	tests := []struct {
		name               string
		scrollType         ScrollType
		wantType           ScrollType
		wantStrengthRange  primitives.AttributeRange
		wantAgilityRange   primitives.AttributeRange
		wantMaxHealthRange primitives.AttributeRange
	}{
		{
			name:               "Get Strength config",
			scrollType:         ScrollTypeStrength,
			wantType:           ScrollTypeStrength,
			wantStrengthRange:  primitives.AttributeRange{Min: 1, Max: 10},
			wantAgilityRange:   primitives.AttributeRange{Min: 0, Max: 0},
			wantMaxHealthRange: primitives.AttributeRange{Min: 0, Max: 0},
		},
		{
			name:               "Get Agility config",
			scrollType:         ScrollTypeAgility,
			wantType:           ScrollTypeAgility,
			wantStrengthRange:  primitives.AttributeRange{Min: 0, Max: 0},
			wantAgilityRange:   primitives.AttributeRange{Min: 1, Max: 10},
			wantMaxHealthRange: primitives.AttributeRange{Min: 0, Max: 0},
		},
		{
			name:               "Get MaxHealth config",
			scrollType:         ScrollTypeMaxHealth,
			wantType:           ScrollTypeMaxHealth,
			wantStrengthRange:  primitives.AttributeRange{Min: 0, Max: 0},
			wantAgilityRange:   primitives.AttributeRange{Min: 0, Max: 0},
			wantMaxHealthRange: primitives.AttributeRange{Min: 5, Max: 20},
		},
		{
			name:               "Get Mystery config - all attributes",
			scrollType:         ScrollTypeMystery,
			wantType:           ScrollTypeMystery,
			wantStrengthRange:  primitives.AttributeRange{Min: -5, Max: 20},
			wantAgilityRange:   primitives.AttributeRange{Min: -5, Max: 20},
			wantMaxHealthRange: primitives.AttributeRange{Min: -5, Max: 20},
		},
		{
			name:               "Unknown type returns Mystery",
			scrollType:         "Unknown Scroll",
			wantType:           ScrollTypeMystery,
			wantStrengthRange:  primitives.AttributeRange{Min: -5, Max: 20},
			wantAgilityRange:   primitives.AttributeRange{Min: -5, Max: 20},
			wantMaxHealthRange: primitives.AttributeRange{Min: -5, Max: 20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := GetScrollConfig(tt.scrollType)

			if config.Type != tt.wantType {
				t.Errorf("Expected type %q, got %q", tt.wantType, config.Type)
			}

			if config.StrengthRange != tt.wantStrengthRange {
				t.Errorf("Expected StrengthRange %v, got %v", tt.wantStrengthRange, config.StrengthRange)
			}

			if config.AgilityRange != tt.wantAgilityRange {
				t.Errorf("Expected AgilityRange %v, got %v", tt.wantAgilityRange, config.AgilityRange)
			}

			if config.MaxHealthRange != tt.wantMaxHealthRange {
				t.Errorf("Expected MaxHealthRange %v, got %v", tt.wantMaxHealthRange, config.MaxHealthRange)
			}

			// Validate the config
			if err := config.Validate(); err != nil {
				t.Errorf("Config validation failed: %v", err)
			}

			// Check description is not empty
			if config.Description == "" {
				t.Errorf("Config for %q has empty description", tt.scrollType)
			}
		})
	}
}

func TestScrollConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    ScrollConfig
		wantError bool
	}{
		{
			name: "Valid config with all attributes",
			config: ScrollConfig{
				Type:           ScrollTypeMystery,
				StrengthRange:  primitives.AttributeRange{Min: -5, Max: 20},
				AgilityRange:   primitives.AttributeRange{Min: -5, Max: 20},
				MaxHealthRange: primitives.AttributeRange{Min: -5, Max: 20},
				Description:    "Valid mystery scroll",
			},
			wantError: false,
		},
		{
			name: "Valid config with single attribute",
			config: ScrollConfig{
				Type:           ScrollTypeStrength,
				StrengthRange:  primitives.AttributeRange{Min: 1, Max: 10},
				AgilityRange:   primitives.AttributeRange{Min: 0, Max: 0},
				MaxHealthRange: primitives.AttributeRange{Min: 0, Max: 0},
				Description:    "Valid strength scroll",
			},
			wantError: false,
		},
		{
			name: "Valid config with zero ranges",
			config: ScrollConfig{
				Type:           "Zero Scroll",
				StrengthRange:  primitives.AttributeRange{Min: 0, Max: 0},
				AgilityRange:   primitives.AttributeRange{Min: 0, Max: 0},
				MaxHealthRange: primitives.AttributeRange{Min: 0, Max: 0},
				Description:    "No attribute changes",
			},
			wantError: false,
		},
		{
			name: "Invalid strength range - Min > Max",
			config: ScrollConfig{
				Type:           "Invalid Scroll",
				StrengthRange:  primitives.AttributeRange{Min: 20, Max: 5}, // Min > Max
				AgilityRange:   primitives.AttributeRange{Min: 0, Max: 0},
				MaxHealthRange: primitives.AttributeRange{Min: 0, Max: 0},
				Description:    "Invalid range",
			},
			wantError: true,
		},
		{
			name: "Invalid agility range - Min > Max",
			config: ScrollConfig{
				Type:           "Invalid Scroll",
				StrengthRange:  primitives.AttributeRange{Min: 0, Max: 0},
				AgilityRange:   primitives.AttributeRange{Min: 15, Max: 5}, // Min > Max
				MaxHealthRange: primitives.AttributeRange{Min: 0, Max: 0},
				Description:    "Invalid range",
			},
			wantError: true,
		},
		{
			name: "Invalid MaxHealth range - Min > Max",
			config: ScrollConfig{
				Type:           "Invalid Scroll",
				StrengthRange:  primitives.AttributeRange{Min: 0, Max: 0},
				AgilityRange:   primitives.AttributeRange{Min: 0, Max: 0},
				MaxHealthRange: primitives.AttributeRange{Min: 30, Max: 10}, // Min > Max
				Description:    "Invalid range",
			},
			wantError: true,
		},
		{
			name: "Invalid multiple ranges",
			config: ScrollConfig{
				Type:           "Invalid Scroll",
				StrengthRange:  primitives.AttributeRange{Min: 20, Max: 5},  // Min > Max
				AgilityRange:   primitives.AttributeRange{Min: 15, Max: 3},  // Min > Max
				MaxHealthRange: primitives.AttributeRange{Min: 30, Max: 10}, // Min > Max
				Description:    "All ranges invalid",
			},
			wantError: true,
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

func TestScrollRegistry_AllConfigsValid(t *testing.T) {
	if len(ScrollRegistry) == 0 {
		t.Fatal("ScrollRegistry is empty")
	}

	expectedCount := 4
	if len(ScrollRegistry) != expectedCount {
		t.Errorf("Expected %d scroll types in registry, got %d", expectedCount, len(ScrollRegistry))
	}

	for scrollType, config := range ScrollRegistry {
		t.Run(string(scrollType), func(t *testing.T) {
			// Validate config
			if err := config.Validate(); err != nil {
				t.Errorf("Registry config for %q is invalid: %v", scrollType, err)
			}

			// Check type matches registry key
			if config.Type != scrollType {
				t.Errorf("Config type mismatch: registry key=%q, config.Type=%q", scrollType, config.Type)
			}

			// Check description is not empty
			if config.Description == "" {
				t.Errorf("Config for %q has empty description", scrollType)
			}

			// Check ranges are reasonable
			if config.StrengthRange.Min > config.StrengthRange.Max {
				t.Errorf("Config for %q has invalid strength range: Min=%f > Max=%f",
					scrollType, config.StrengthRange.Min, config.StrengthRange.Max)
			}
			if config.AgilityRange.Min > config.AgilityRange.Max {
				t.Errorf("Config for %q has invalid agility range: Min=%f > Max=%f",
					scrollType, config.AgilityRange.Min, config.AgilityRange.Max)
			}
			if config.MaxHealthRange.Min > config.MaxHealthRange.Max {
				t.Errorf("Config for %q has invalid MaxHealth range: Min=%f > Max=%f",
					scrollType, config.MaxHealthRange.Min, config.MaxHealthRange.Max)
			}
		})
	}
}

func TestScrollRegistry_ExpectedTypes(t *testing.T) {
	expectedTypes := []ScrollType{
		ScrollTypeStrength,
		ScrollTypeAgility,
		ScrollTypeMaxHealth,
		ScrollTypeMystery,
	}

	for _, scrollType := range expectedTypes {
		t.Run(string(scrollType), func(t *testing.T) {
			config, exists := ScrollRegistry[scrollType]
			if !exists {
				t.Errorf("Expected scroll type %q not found in registry", scrollType)
				return
			}

			if config.Type != scrollType {
				t.Errorf("Expected Type=%q, got %q", scrollType, config.Type)
			}
		})
	}
}

func TestScrollConfig_GenerateAttributes(t *testing.T) {
	tests := []struct {
		name   string
		seed   int64
		config ScrollConfig
	}{
		{
			name: "Strength config generates deterministic attributes",
			seed: defaultItemsTestSeed,
			config: ScrollConfig{
				Type:           ScrollTypeStrength,
				StrengthRange:  primitives.AttributeRange{Min: 1, Max: 10},
				AgilityRange:   primitives.AttributeRange{Min: 0, Max: 0},
				MaxHealthRange: primitives.AttributeRange{Min: 0, Max: 0},
			},
		},
		{
			name: "Agility config generates deterministic attributes",
			seed: 100,
			config: ScrollConfig{
				Type:           ScrollTypeAgility,
				StrengthRange:  primitives.AttributeRange{Min: 0, Max: 0},
				AgilityRange:   primitives.AttributeRange{Min: 1, Max: 10},
				MaxHealthRange: primitives.AttributeRange{Min: 0, Max: 0},
			},
		},
		{
			name: "MaxHealth config generates deterministic attributes",
			seed: 200,
			config: ScrollConfig{
				Type:           ScrollTypeMaxHealth,
				StrengthRange:  primitives.AttributeRange{Min: 0, Max: 0},
				AgilityRange:   primitives.AttributeRange{Min: 0, Max: 0},
				MaxHealthRange: primitives.AttributeRange{Min: 5, Max: 20},
			},
		},
		{
			name: "Mystery config with all attributes",
			seed: 300,
			config: ScrollConfig{
				Type:           ScrollTypeMystery,
				StrengthRange:  primitives.AttributeRange{Min: -5, Max: 20},
				AgilityRange:   primitives.AttributeRange{Min: -5, Max: 20},
				MaxHealthRange: primitives.AttributeRange{Min: -5, Max: 20},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate twice with same seed
			rng1 := utils.NewRandomGeneratorWithSeed(tt.seed)
			attrs1 := tt.config.GenerateAttributes(rng1)

			rng2 := utils.NewRandomGeneratorWithSeed(tt.seed)
			attrs2 := tt.config.GenerateAttributes(rng2)

			// Should be deterministic
			if attrs1.Strength != attrs2.Strength {
				t.Errorf("Strength not deterministic: %f != %f", attrs1.Strength, attrs2.Strength)
			}
			if attrs1.Agility != attrs2.Agility {
				t.Errorf("Agility not deterministic: %f != %f", attrs1.Agility, attrs2.Agility)
			}
			if attrs1.MaxHealth != attrs2.MaxHealth {
				t.Errorf("MaxHealth not deterministic: %f != %f", attrs1.MaxHealth, attrs2.MaxHealth)
			}

			// Should be within primitives.AttributeRange
			if attrs1.Strength < tt.config.StrengthRange.Min || attrs1.Strength > tt.config.StrengthRange.Max {
				t.Errorf("Strength out of AttributeRange: got %.2f, want [%.2f, %.2f]",
					attrs1.Strength, tt.config.StrengthRange.Min, tt.config.StrengthRange.Max)
			}
			if attrs1.Agility < tt.config.AgilityRange.Min || attrs1.Agility > tt.config.AgilityRange.Max {
				t.Errorf("Agility out of AttributeRange: got %.2f, want [%.2f, %.2f]",
					attrs1.Agility, tt.config.AgilityRange.Min, tt.config.AgilityRange.Max)
			}
			if attrs1.MaxHealth < tt.config.MaxHealthRange.Min || attrs1.MaxHealth > tt.config.MaxHealthRange.Max {
				t.Errorf("MaxHealth out of AttributeRange: got %.2f, want [%.2f, %.2f]",
					attrs1.MaxHealth, tt.config.MaxHealthRange.Min, tt.config.MaxHealthRange.Max)
			}

			// Health should be zero (scrolls don't affect current health)
			if attrs1.Health != 0 {
				t.Errorf("Expected Health=0, got %f", attrs1.Health)
			}
		})
	}
}

func TestScrollConfig_GenerateAttributes_Randomness(t *testing.T) {
	config := ScrollConfig{
		Type:           ScrollTypeMystery,
		StrengthRange:  primitives.AttributeRange{Min: -5, Max: 20},
		AgilityRange:   primitives.AttributeRange{Min: -5, Max: 20},
		MaxHealthRange: primitives.AttributeRange{Min: -5, Max: 20},
	}

	strengthValues := make(map[float64]bool)
	agilityValues := make(map[float64]bool)
	maxHealthValues := make(map[float64]bool)
	iterations := 50

	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomGeneratorWithSeed(int64(i))
		attrs := config.GenerateAttributes(rng)
		strengthValues[attrs.Strength] = true
		agilityValues[attrs.Agility] = true
		maxHealthValues[attrs.MaxHealth] = true
	}

	minUniqueValues := 10
	if len(strengthValues) < minUniqueValues {
		t.Errorf("Expected at least %d unique strength values, got %d", minUniqueValues, len(strengthValues))
	}
	if len(agilityValues) < minUniqueValues {
		t.Errorf("Expected at least %d unique agility values, got %d", minUniqueValues, len(agilityValues))
	}
	if len(maxHealthValues) < minUniqueValues {
		t.Errorf("Expected at least %d unique MaxHealth values, got %d", minUniqueValues, len(maxHealthValues))
	}

	// Verify all values are within range
	for strength := range strengthValues {
		if strength < config.StrengthRange.Min || strength > config.StrengthRange.Max {
			t.Errorf("Strength value %f is out of range [%f, %f]",
				strength, config.StrengthRange.Min, config.StrengthRange.Max)
		}
	}
}

func TestScrollConfig_GenerateAttributes_ZeroRanges(t *testing.T) {
	config := ScrollConfig{
		Type:           ScrollTypeStrength,
		StrengthRange:  primitives.AttributeRange{Min: 5, Max: 10},
		AgilityRange:   primitives.AttributeRange{Min: 0, Max: 0},
		MaxHealthRange: primitives.AttributeRange{Min: 0, Max: 0},
	}

	iterations := 20
	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomGeneratorWithSeed(int64(i))
		attrs := config.GenerateAttributes(rng)

		// Zero ranges should always produce zero
		if attrs.Agility != 0 {
			t.Errorf("Expected Agility=0 for zero range, got %f", attrs.Agility)
		}
		if attrs.MaxHealth != 0 {
			t.Errorf("Expected MaxHealth=0 for zero range, got %f", attrs.MaxHealth)
		}
		if attrs.Health != 0 {
			t.Errorf("Expected Health=0, got %f", attrs.Health)
		}
	}
}

func TestScrollConfig_GenerateAttributes_NegativeValues(t *testing.T) {
	config := ScrollConfig{
		Type:           "Curse Scroll",
		StrengthRange:  primitives.AttributeRange{Min: -10, Max: -1},
		AgilityRange:   primitives.AttributeRange{Min: -10, Max: -1},
		MaxHealthRange: primitives.AttributeRange{Min: -20, Max: -5},
	}

	iterations := 30
	foundNegativeStr := false
	foundNegativeAgi := false
	foundNegativeMaxHP := false

	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomGeneratorWithSeed(int64(i))
		attrs := config.GenerateAttributes(rng)

		if attrs.Strength < 0 {
			foundNegativeStr = true
		}
		if attrs.Agility < 0 {
			foundNegativeAgi = true
		}
		if attrs.MaxHealth < 0 {
			foundNegativeMaxHP = true
		}

		// Verify all are within ranges
		if attrs.Strength < config.StrengthRange.Min || attrs.Strength > config.StrengthRange.Max {
			t.Errorf("Strength out of range: got %.2f, want [%.2f, %.2f]",
				attrs.Strength, config.StrengthRange.Min, config.StrengthRange.Max)
		}
		if attrs.Agility < config.AgilityRange.Min || attrs.Agility > config.AgilityRange.Max {
			t.Errorf("Agility out of range: got %.2f, want [%.2f, %.2f]",
				attrs.Agility, config.AgilityRange.Min, config.AgilityRange.Max)
		}
		if attrs.MaxHealth < config.MaxHealthRange.Min || attrs.MaxHealth > config.MaxHealthRange.Max {
			t.Errorf("MaxHealth out of range: got %.2f, want [%.2f, %.2f]",
				attrs.MaxHealth, config.MaxHealthRange.Min, config.MaxHealthRange.Max)
		}
	}

	if !foundNegativeStr {
		t.Error("Expected to find at least one negative Strength value")
	}
	if !foundNegativeAgi {
		t.Error("Expected to find at least one negative Agility value")
	}
	if !foundNegativeMaxHP {
		t.Error("Expected to find at least one negative MaxHealth value")
	}
}

// Benchmarks
func BenchmarkGetScrollConfig(b *testing.B) {
	b.Run("Known type", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = GetScrollConfig(ScrollTypeStrength)
		}
	})

	b.Run("Unknown type", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = GetScrollConfig("Unknown")
		}
	})
}

func BenchmarkScrollConfig_Validate(b *testing.B) {
	config := ScrollConfig{
		Type:           ScrollTypeMystery,
		StrengthRange:  primitives.AttributeRange{Min: -5, Max: 20},
		AgilityRange:   primitives.AttributeRange{Min: -5, Max: 20},
		MaxHealthRange: primitives.AttributeRange{Min: -5, Max: 20},
		Description:    "Test scroll",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}

func BenchmarkScrollConfig_GenerateAttributes(b *testing.B) {
	config := GetScrollConfig(ScrollTypeMystery)
	rng := utils.NewRandomGeneratorWithSeed(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.GenerateAttributes(rng)
	}
}
