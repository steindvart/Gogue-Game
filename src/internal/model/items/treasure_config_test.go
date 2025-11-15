package items

import (
	"gogue/internal/utils"
	"testing"
)

func TestGetTreasureConfig(t *testing.T) {
	tests := []struct {
		name           string
		treasureType   TreasureType
		wantType       TreasureType
		wantValueRange TreasureValueRange
	}{
		{
			name:           "Get Gold config",
			treasureType:   TreasureTypeGold,
			wantType:       TreasureTypeGold,
			wantValueRange: TreasureValueRange{Min: 5, Max: 20},
		},
		{
			name:           "Get Gem config",
			treasureType:   TreasureTypeGem,
			wantType:       TreasureTypeGem,
			wantValueRange: TreasureValueRange{Min: 50, Max: 200},
		},
		{
			name:           "Get Artifact config",
			treasureType:   TreasureTypeArtifact,
			wantType:       TreasureTypeArtifact,
			wantValueRange: TreasureValueRange{Min: 200, Max: 1000},
		},
		{
			name:           "Get Mystery config",
			treasureType:   TreasureTypeMystery,
			wantType:       TreasureTypeMystery,
			wantValueRange: TreasureValueRange{Min: 1, Max: 500},
		},
		{
			name:           "Unknown type returns Mystery",
			treasureType:   "Unknown Treasure",
			wantType:       TreasureTypeMystery,
			wantValueRange: TreasureValueRange{Min: 1, Max: 500},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := GetTreasureConfig(tt.treasureType)

			if config.Type != tt.wantType {
				t.Errorf("Expected type %q, got %q", tt.wantType, config.Type)
			}

			if config.ValueRange != tt.wantValueRange {
				t.Errorf("Expected ValueRange %+v, got %+v", tt.wantValueRange, config.ValueRange)
			}

			// Validate the config
			if err := config.Validate(); err != nil {
				t.Errorf("Config validation failed: %v", err)
			}

			// Check description is not empty
			if config.Description == "" {
				t.Errorf("Config for %q has empty description", tt.treasureType)
			}
		})
	}
}

func TestTreasureConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    TreasureConfig
		wantError bool
	}{
		{
			name: "Valid config with positive range",
			config: TreasureConfig{
				Type:        TreasureTypeGold,
				ValueRange:  TreasureValueRange{Min: 5, Max: 20},
				Description: "Valid gold treasure",
			},
			wantError: false,
		},
		{
			name: "Valid config with large range",
			config: TreasureConfig{
				Type:        TreasureTypeArtifact,
				ValueRange:  TreasureValueRange{Min: 100, Max: 1000},
				Description: "Valid artifact",
			},
			wantError: false,
		},
		{
			name: "Valid config with single value",
			config: TreasureConfig{
				Type:        "Fixed Treasure",
				ValueRange:  TreasureValueRange{Min: 10, Max: 10},
				Description: "Always 10 value",
			},
			wantError: false,
		},
		{
			name: "Invalid range - Min > Max",
			config: TreasureConfig{
				Type:        "Invalid Treasure",
				ValueRange:  TreasureValueRange{Min: 100, Max: 10}, // Min > Max
				Description: "Invalid range",
			},
			wantError: true,
		},
		{
			name: "Invalid range - large difference",
			config: TreasureConfig{
				Type:        "Invalid Treasure",
				ValueRange:  TreasureValueRange{Min: 1000, Max: 1}, // Min > Max
				Description: "Invalid range",
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

func TestTreasureRegistry_AllConfigsValid(t *testing.T) {
	if len(TreasureRegistry) == 0 {
		t.Fatal("TreasureRegistry is empty")
	}

	expectedCount := 4
	if len(TreasureRegistry) != expectedCount {
		t.Errorf("Expected %d treasure types in registry, got %d", expectedCount, len(TreasureRegistry))
	}

	for treasureType, config := range TreasureRegistry {
		t.Run(string(treasureType), func(t *testing.T) {
			// Validate config
			if err := config.Validate(); err != nil {
				t.Errorf("Registry config for %q is invalid: %v", treasureType, err)
			}

			// Check type matches registry key
			if config.Type != treasureType {
				t.Errorf("Config type mismatch: registry key=%q, config.Type=%q", treasureType, config.Type)
			}

			// Check description is not empty
			if config.Description == "" {
				t.Errorf("Config for %q has empty description", treasureType)
			}

			// Check range is reasonable
			if config.ValueRange.Min > config.ValueRange.Max {
				t.Errorf("Config for %q has invalid value range: Min=%d > Max=%d",
					treasureType, config.ValueRange.Min, config.ValueRange.Max)
			}

			// Check values are non-negative
			if config.ValueRange.Min < 0 {
				t.Errorf("Config for %q has negative Min value: %d", treasureType, config.ValueRange.Min)
			}
		})
	}
}

func TestTreasureRegistry_ExpectedTypes(t *testing.T) {
	expectedTypes := []TreasureType{
		TreasureTypeGold,
		TreasureTypeGem,
		TreasureTypeArtifact,
		TreasureTypeMystery,
	}

	for _, treasureType := range expectedTypes {
		t.Run(string(treasureType), func(t *testing.T) {
			config, exists := TreasureRegistry[treasureType]
			if !exists {
				t.Errorf("Expected treasure type %q not found in registry", treasureType)
				return
			}

			if config.Type != treasureType {
				t.Errorf("Expected Type=%q, got %q", treasureType, config.Type)
			}
		})
	}
}

func TestTreasureConfig_GenerateValue(t *testing.T) {
	tests := []struct {
		name   string
		seed   int64
		config TreasureConfig
	}{
		{
			name: "Gold config generates deterministic value",
			seed: defaultItemsTestSeed,
			config: TreasureConfig{
				Type:       TreasureTypeGold,
				ValueRange: TreasureValueRange{Min: 5, Max: 20},
			},
		},
		{
			name: "Gem config generates deterministic value",
			seed: 100,
			config: TreasureConfig{
				Type:       TreasureTypeGem,
				ValueRange: TreasureValueRange{Min: 50, Max: 200},
			},
		},
		{
			name: "Artifact config generates deterministic value",
			seed: 200,
			config: TreasureConfig{
				Type:       TreasureTypeArtifact,
				ValueRange: TreasureValueRange{Min: 200, Max: 1000},
			},
		},
		{
			name: "Mystery config generates deterministic value",
			seed: 300,
			config: TreasureConfig{
				Type:       TreasureTypeMystery,
				ValueRange: TreasureValueRange{Min: 1, Max: 500},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate twice with same seed
			rng1 := utils.NewRandomGeneratorWithSeed(tt.seed)
			value1 := tt.config.GenerateValue(rng1)

			rng2 := utils.NewRandomGeneratorWithSeed(tt.seed)
			value2 := tt.config.GenerateValue(rng2)

			// Should be deterministic
			if value1 != value2 {
				t.Errorf("Value not deterministic: %d != %d", value1, value2)
			}

			// Should be within range
			if value1 < tt.config.ValueRange.Min || value1 > tt.config.ValueRange.Max {
				t.Errorf("Value out of range: got %d, want [%d, %d]",
					value1, tt.config.ValueRange.Min, tt.config.ValueRange.Max)
			}

			// Value should be non-negative
			if value1 < 0 {
				t.Errorf("Expected non-negative value, got %d", value1)
			}
		})
	}
}

func TestTreasureConfig_GenerateValue_Randomness(t *testing.T) {
	config := TreasureConfig{
		Type:       TreasureTypeMystery,
		ValueRange: TreasureValueRange{Min: 1, Max: 500},
	}

	values := make(map[int32]bool)
	iterations := 50

	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomGeneratorWithSeed(int64(i))
		value := config.GenerateValue(rng)
		values[value] = true
	}

	minUniqueValues := 10
	if len(values) < minUniqueValues {
		t.Errorf("Expected at least %d unique values, got %d", minUniqueValues, len(values))
	}

	// Verify all values are within range
	for value := range values {
		if value < config.ValueRange.Min || value > config.ValueRange.Max {
			t.Errorf("Value %d is out of range [%d, %d]",
				value, config.ValueRange.Min, config.ValueRange.Max)
		}
	}
}

func TestTreasureConfig_GenerateValue_FixedRange(t *testing.T) {
	config := TreasureConfig{
		Type:       "Fixed Treasure",
		ValueRange: TreasureValueRange{Min: 10, Max: 10},
	}

	iterations := 20
	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomGeneratorWithSeed(int64(i))
		value := config.GenerateValue(rng)

		// Fixed range should always produce same value
		if value != 10 {
			t.Errorf("Expected value=10 for fixed range, got %d", value)
		}
	}
}

func TestTreasureConfig_GenerateValue_LargeRange(t *testing.T) {
	config := TreasureConfig{
		Type:       "Large Range Treasure",
		ValueRange: TreasureValueRange{Min: 1, Max: 10000},
	}

	values := make(map[int32]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomGeneratorWithSeed(int64(i))
		value := config.GenerateValue(rng)
		values[value] = true

		// Verify within range
		if value < config.ValueRange.Min || value > config.ValueRange.Max {
			t.Errorf("Value out of range: got %d, want [%d, %d]",
				value, config.ValueRange.Min, config.ValueRange.Max)
		}
	}

	// Should have good distribution
	minUniqueValues := 50
	if len(values) < minUniqueValues {
		t.Errorf("Expected at least %d unique values in large range, got %d", minUniqueValues, len(values))
	}
}

// Benchmarks
func BenchmarkGetTreasureConfig(b *testing.B) {
	b.Run("Known type", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = GetTreasureConfig(TreasureTypeGold)
		}
	})

	b.Run("Unknown type", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = GetTreasureConfig("Unknown")
		}
	})
}

func BenchmarkTreasureConfig_Validate(b *testing.B) {
	config := TreasureConfig{
		Type:        TreasureTypeMystery,
		ValueRange:  TreasureValueRange{Min: 1, Max: 500},
		Description: "Test treasure",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}

func BenchmarkTreasureConfig_GenerateValue(b *testing.B) {
	config := GetTreasureConfig(TreasureTypeMystery)
	rng := utils.NewRandomGeneratorWithSeed(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.GenerateValue(rng)
	}
}
