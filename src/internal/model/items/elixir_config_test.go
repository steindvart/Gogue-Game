package items

import (
	"gogue/internal/utils"
	"testing"
)

func TestGetElixirConfig(t *testing.T) {
	tests := []struct {
		name          string
		elixirType    ElixirType
		wantType      ElixirType
		checkStrength bool
		checkAgility  bool
	}{
		{
			name:          "Get Strength config",
			elixirType:    ElixirTypeStrength,
			wantType:      ElixirTypeStrength,
			checkStrength: true,
			checkAgility:  false,
		},
		{
			name:          "Get Agility config",
			elixirType:    ElixirTypeAgility,
			wantType:      ElixirTypeAgility,
			checkStrength: false,
			checkAgility:  true,
		},
		{
			name:          "Unknown type returns Mystery",
			elixirType:    "Unknown Elixir",
			wantType:      ElixirTypeMystery,
			checkStrength: true,
			checkAgility:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := GetElixirConfig(tt.elixirType)

			if config.Type != tt.wantType {
				t.Errorf("Expected type %q, got %q", tt.wantType, config.Type)
			}

			// Check that config has proper ranges
			if tt.checkStrength && config.StrengthRange.Min == 0 && config.StrengthRange.Max == 0 {
				t.Error("Expected non-zero strength range")
			}

			if tt.checkAgility && config.AgilityRange.Min == 0 && config.AgilityRange.Max == 0 {
				t.Error("Expected non-zero agility range")
			}

			// Validate the config
			if err := config.Validate(); err != nil {
				t.Errorf("Config validation failed: %v", err)
			}
		})
	}
}

func TestElixirConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    ElixirConfig
		wantError bool
	}{
		{
			name: "Valid config",
			config: ElixirConfig{
				Type:               ElixirTypeStrength,
				StrengthRange:      ElixirAttributeRange{Min: 5, Max: 20},
				AgilityRange:       ElixirAttributeRange{Min: 0, Max: 0},
				DurationStepsRange: defaultDurationRange,
			},
			wantError: false,
		},
		{
			name: "Invalid strength range is error",
			config: ElixirConfig{
				Type:               ElixirTypeStrength,
				StrengthRange:      ElixirAttributeRange{Min: 20, Max: 5}, // Min > Max
				AgilityRange:       ElixirAttributeRange{Min: 0, Max: 0},
				DurationStepsRange: defaultDurationRange,
			},
			wantError: true,
		},
		{
			name: "Invalid agility range is error",
			config: ElixirConfig{
				Type:               ElixirTypeAgility,
				StrengthRange:      ElixirAttributeRange{Min: 0, Max: 0},
				AgilityRange:       ElixirAttributeRange{Min: 15, Max: 5}, // Min > Max
				DurationStepsRange: defaultDurationRange,
			},
			wantError: true,
		},
		{
			name: "Invalid duration range is error",
			config: ElixirConfig{
				Type:               ElixirTypeMystery,
				StrengthRange:      ElixirAttributeRange{Min: -10, Max: 10},
				AgilityRange:       ElixirAttributeRange{Min: -10, Max: 10},
				DurationStepsRange: ElixirDurationStepsRange{Min: 30, Max: 5}, // Min > Max
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

func TestElixirRegistry_AllConfigsValid(t *testing.T) {
	for elixirType, config := range ElixirRegistry {
		t.Run(string(elixirType), func(t *testing.T) {
			if err := config.Validate(); err != nil {
				t.Errorf("Registry config for %q is invalid: %v", elixirType, err)
			}

			if config.Type != elixirType {
				t.Errorf("Config type mismatch: registry key=%q, config.Type=%q", elixirType, config.Type)
			}

			if config.Description == "" {
				t.Errorf("Config for %q has empty description", elixirType)
			}
		})
	}
}

func TestElixirConfig_GenerateAttributes(t *testing.T) {
	tests := []struct {
		name   string
		seed   int64
		config ElixirConfig
	}{
		{
			name: "Strength config generates deterministic attributes",
			seed: defaultTestSeed,
			config: ElixirConfig{
				Type:          ElixirTypeStrength,
				StrengthRange: ElixirAttributeRange{Min: 5, Max: 20},
				AgilityRange:  ElixirAttributeRange{Min: 0, Max: 0},
			},
		},
		{
			name: "Mystery config generates varied attributes",
			seed: 100,
			config: ElixirConfig{
				Type:          ElixirTypeMystery,
				StrengthRange: ElixirAttributeRange{Min: -20, Max: 30},
				AgilityRange:  ElixirAttributeRange{Min: -20, Max: 30},
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

func TestElixirConfig_GenerateDuration(t *testing.T) {
	tests := []struct {
		name   string
		seed   int64
		config ElixirConfig
	}{
		{
			name: "Standard duration range",
			seed: defaultTestSeed,
			config: ElixirConfig{
				Type:               ElixirTypeStrength,
				DurationStepsRange: defaultDurationRange,
			},
		},
		{
			name: "Custom short duration range",
			seed: 50,
			config: ElixirConfig{
				Type:               "Short Elixir",
				DurationStepsRange: ElixirDurationStepsRange{Min: 1, Max: 5},
			},
		},
		{
			name: "Custom long duration range",
			seed: 75,
			config: ElixirConfig{
				Type:               "Long Elixir",
				DurationStepsRange: ElixirDurationStepsRange{Min: 50, Max: 100},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate twice with same seed
			rng1 := utils.NewRandomGeneratorWithSeed(tt.seed)
			duration1 := tt.config.GenerateDuration(rng1)

			rng2 := utils.NewRandomGeneratorWithSeed(tt.seed)
			duration2 := tt.config.GenerateDuration(rng2)

			// Should be deterministic
			if duration1 != duration2 {
				t.Errorf("Duration not deterministic: %d != %d", duration1, duration2)
			}

			// Should be within range
			if duration1 < tt.config.DurationStepsRange.Min || duration1 > tt.config.DurationStepsRange.Max {
				t.Errorf("Duration out of range: got %d, want [%d, %d]",
					duration1, tt.config.DurationStepsRange.Min, tt.config.DurationStepsRange.Max)
			}
		})
	}
}

func TestElixirConfig_GenerateAttributes_Randomness(t *testing.T) {
	config := ElixirConfig{
		Type:          ElixirTypeMystery,
		StrengthRange: ElixirAttributeRange{Min: -20, Max: 30},
		AgilityRange:  ElixirAttributeRange{Min: -20, Max: 30},
	}

	strengthValues := make(map[float64]bool)
	agilityValues := make(map[float64]bool)
	iterations := 50

	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomGeneratorWithSeed(int64(i))
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

func TestElixirConfig_GenerateDuration_Randomness(t *testing.T) {
	config := ElixirConfig{
		Type:               ElixirTypeStrength,
		DurationStepsRange: defaultDurationRange,
	}

	durationValues := make(map[uint32]bool)
	iterations := 50

	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomGeneratorWithSeed(int64(i))
		duration := config.GenerateDuration(rng)
		durationValues[duration] = true
	}

	minUniqueValues := 10
	if len(durationValues) < minUniqueValues {
		t.Errorf("Expected at least %d unique duration values, got %d", minUniqueValues, len(durationValues))
	}
}
