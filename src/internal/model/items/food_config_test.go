package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
	"testing"
)

func TestGetFoodConfig(t *testing.T) {
	tests := []struct {
		name            string
		foodType        FoodType
		wantType        FoodType
		wantHealthRange primitives.AttributeRange
	}{
		{
			name:            "Get Potatoes config",
			foodType:        FoodTypePotatoes,
			wantType:        FoodTypePotatoes,
			wantHealthRange: primitives.AttributeRange{Min: 1, Max: 10},
		},
		{
			name:            "Get Bread config",
			foodType:        FoodTypeBread,
			wantType:        FoodTypeBread,
			wantHealthRange: primitives.AttributeRange{Min: 10, Max: 20},
		},
		{
			name:            "Get Meat config",
			foodType:        FoodTypeMeat,
			wantType:        FoodTypeMeat,
			wantHealthRange: primitives.AttributeRange{Min: 20, Max: 30},
		},
		{
			name:            "Get Beer config - fixed value",
			foodType:        FoodTypeBeer,
			wantType:        FoodTypeBeer,
			wantHealthRange: primitives.AttributeRange{Min: 20, Max: 20},
		},
		{
			name:            "Get Mistery config - wide range",
			foodType:        FoodTypeMistery,
			wantType:        FoodTypeMistery,
			wantHealthRange: primitives.AttributeRange{Min: -10, Max: 40},
		},
		{
			name:            "Unknown type returns Mistery",
			foodType:        "Unknown Food",
			wantType:        FoodTypeMistery,
			wantHealthRange: primitives.AttributeRange{Min: -10, Max: 40},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := GetFoodConfig(tt.foodType)

			if config.Type != tt.wantType {
				t.Errorf("Expected type %q, got %q", tt.wantType, config.Type)
			}

			if config.HealthRange != tt.wantHealthRange {
				t.Errorf("Expected HealthRange=%v, got %v", tt.wantHealthRange, config.HealthRange)
			}

			// Validate the config
			if err := config.Validate(); err != nil {
				t.Errorf("Config validation failed: %v", err)
			}

			// Check description is not empty
			if config.Description == "" {
				t.Errorf("Config for %q has empty description", tt.foodType)
			}
		})
	}
}

func TestFoodConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    FoodConfig
		wantError bool
	}{
		{
			name: "Valid config with positive range",
			config: FoodConfig{
				Type:        FoodTypeMeat,
				HealthRange: primitives.AttributeRange{Min: 20, Max: 30},
				Description: "Valid meat",
			},
			wantError: false,
		},
		{
			name: "Valid config with negative to positive range",
			config: FoodConfig{
				Type:        FoodTypeMistery,
				HealthRange: primitives.AttributeRange{Min: -10, Max: 40},
				Description: "Valid mystery",
			},
			wantError: false,
		},
		{
			name: "Valid config with zero-width range",
			config: FoodConfig{
				Type:        FoodTypeBeer,
				HealthRange: primitives.AttributeRange{Min: 20, Max: 20},
				Description: "Fixed health",
			},
			wantError: false,
		},
		{
			name: "Valid config with zero values",
			config: FoodConfig{
				Type:        "Zero Food",
				HealthRange: primitives.AttributeRange{Min: 0, Max: 0},
				Description: "No health change",
			},
			wantError: false,
		},
		{
			name: "Invalid health range - Min > Max",
			config: FoodConfig{
				Type:        "Invalid Food",
				HealthRange: primitives.AttributeRange{Min: 30, Max: 10},
				Description: "Invalid range",
			},
			wantError: true,
		},
		{
			name: "Invalid health range - positive to negative",
			config: FoodConfig{
				Type:        "Invalid Food 2",
				HealthRange: primitives.AttributeRange{Min: 20, Max: -5},
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

func TestFoodRegistry_AllConfigsValid(t *testing.T) {
	if len(FoodRegistry) == 0 {
		t.Fatal("FoodRegistry is empty")
	}

	expectedCount := 5
	if len(FoodRegistry) != expectedCount {
		t.Errorf("Expected %d food types in registry, got %d", expectedCount, len(FoodRegistry))
	}

	for foodType, config := range FoodRegistry {
		t.Run(string(foodType), func(t *testing.T) {
			// Validate config
			if err := config.Validate(); err != nil {
				t.Errorf("Registry config for %q is invalid: %v", foodType, err)
			}

			// Check type matches registry key
			if config.Type != foodType {
				t.Errorf("Config type mismatch: registry key=%q, config.Type=%q", foodType, config.Type)
			}

			// Check description is not empty
			if config.Description == "" {
				t.Errorf("Config for %q has empty description", foodType)
			}

			// Check health range is reasonable
			if config.HealthRange.Min > config.HealthRange.Max {
				t.Errorf("Config for %q has invalid health range: Min=%f > Max=%f",
					foodType, config.HealthRange.Min, config.HealthRange.Max)
			}
		})
	}
}

func TestFoodRegistry_ExpectedTypes(t *testing.T) {
	expectedTypes := []FoodType{
		FoodTypePotatoes,
		FoodTypeBread,
		FoodTypeMeat,
		FoodTypeBeer,
		FoodTypeMistery,
	}

	for _, foodType := range expectedTypes {
		t.Run(string(foodType), func(t *testing.T) {
			config, exists := FoodRegistry[foodType]
			if !exists {
				t.Errorf("Expected food type %q not found in registry", foodType)
				return
			}

			if config.Type != foodType {
				t.Errorf("Expected Type=%q, got %q", foodType, config.Type)
			}
		})
	}
}

func TestFoodConfig_GenerateAttributes(t *testing.T) {
	tests := []struct {
		name   string
		seed   int64
		config FoodConfig
	}{
		{
			name: "Potatoes config generates deterministic health",
			seed: defaultTestSeed,
			config: FoodConfig{
				Type:        FoodTypePotatoes,
				HealthRange: primitives.AttributeRange{Min: 1, Max: 10},
			},
		},
		{
			name: "Meat config generates higher health",
			seed: 100,
			config: FoodConfig{
				Type:        FoodTypeMeat,
				HealthRange: primitives.AttributeRange{Min: 20, Max: 30},
			},
		},
		{
			name: "Mistery config with negative range",
			seed: 200,
			config: FoodConfig{
				Type:        FoodTypeMistery,
				HealthRange: primitives.AttributeRange{Min: -10, Max: 40},
			},
		},
		{
			name: "Beer config with fixed value",
			seed: 300,
			config: FoodConfig{
				Type:        FoodTypeBeer,
				HealthRange: primitives.AttributeRange{Min: 20, Max: 20},
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
			if attrs1.Health != attrs2.Health {
				t.Errorf("Health not deterministic: %f != %f", attrs1.Health, attrs2.Health)
			}

			// Should be within primitives.AttributeRange
			if attrs1.Health < tt.config.HealthRange.Min || attrs1.Health > tt.config.HealthRange.Max {
				t.Errorf("Health out of AttributeRange: got %.2f, want [%.2f, %.2f]",
					attrs1.Health, tt.config.HealthRange.Min, tt.config.HealthRange.Max)
			}

			// Other attributes should be zero
			if attrs1.Strength != 0 {
				t.Errorf("Expected Strength=0, got %f", attrs1.Strength)
			}
			if attrs1.Agility != 0 {
				t.Errorf("Expected Agility=0, got %f", attrs1.Agility)
			}

			// For fixed-value range, should always return exact value
			if tt.config.HealthRange.Min == tt.config.HealthRange.Max {
				if attrs1.Health != tt.config.HealthRange.Min {
					t.Errorf("Expected exact value %f for fixed range, got %f",
						tt.config.HealthRange.Min, attrs1.Health)
				}
			}
		})
	}
}

func TestFoodConfig_GenerateAttributes_Randomness(t *testing.T) {
	config := FoodConfig{
		Type:        FoodTypeMistery,
		HealthRange: primitives.AttributeRange{Min: -10, Max: 40},
	}

	healthValues := make(map[float64]bool)
	iterations := 50

	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomGeneratorWithSeed(int64(i))
		attrs := config.GenerateAttributes(rng)
		healthValues[attrs.Health] = true
	}

	minUniqueValues := 10
	if len(healthValues) < minUniqueValues {
		t.Errorf("Expected at least %d unique health values, got %d", minUniqueValues, len(healthValues))
	}

	// Verify all values are within range
	for health := range healthValues {
		if health < config.HealthRange.Min || health > config.HealthRange.Max {
			t.Errorf("Health value %f is out of range [%f, %f]",
				health, config.HealthRange.Min, config.HealthRange.Max)
		}
	}
}

func TestFoodConfig_GenerateAttributes_FixedValue(t *testing.T) {
	config := FoodConfig{
		Type:        FoodTypeBeer,
		HealthRange: primitives.AttributeRange{Min: 20, Max: 20},
	}

	expectedHealth := 20.0
	iterations := 20

	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomGeneratorWithSeed(int64(i))
		attrs := config.GenerateAttributes(rng)

		if attrs.Health != expectedHealth {
			t.Errorf("Iteration %d: Expected Health=%f, got %f", i, expectedHealth, attrs.Health)
		}
	}
}

// Benchmarks
func BenchmarkGetFoodConfig(b *testing.B) {
	b.Run("Known type", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = GetFoodConfig(FoodTypeMeat)
		}
	})

	b.Run("Unknown type", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = GetFoodConfig("Unknown")
		}
	})
}

func BenchmarkFoodConfig_Validate(b *testing.B) {
	config := FoodConfig{
		Type:        FoodTypeMeat,
		HealthRange: primitives.AttributeRange{Min: 20, Max: 30},
		Description: "Test food",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}

func BenchmarkFoodConfig_GenerateAttributes(b *testing.B) {
	config := GetFoodConfig(FoodTypeMistery)
	rng := utils.NewRandomGeneratorWithSeed(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.GenerateAttributes(rng)
	}
}
