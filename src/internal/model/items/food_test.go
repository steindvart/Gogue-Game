package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
	"testing"
)

const defaultFoodTestSeed int64 = 42

func TestFood_NewFood_BuiltinConfig(t *testing.T) {
	tests := []struct {
		name         string
		foodType     FoodType
		expectedName string
		wantMinHP    float64
		wantMaxHP    float64
	}{
		{
			name:         "Potatoes type uses Potatoes config",
			foodType:     FoodTypePotatoes,
			expectedName: string(FoodTypePotatoes),
		},
		{
			name:         "Bread type uses Bread config",
			foodType:     FoodTypeBread,
			expectedName: string(FoodTypeBread),
		},
		{
			name:         "Meat type uses Meat config",
			foodType:     FoodTypeMeat,
			expectedName: string(FoodTypeMeat),
		},
		{
			name:         "Beer type uses Beer config - fixed value",
			foodType:     FoodTypeBeer,
			expectedName: string(FoodTypeBeer),
		},
		{
			name:         "Mistery type uses Mistery config",
			foodType:     FoodTypeMistery,
			expectedName: string(FoodTypeMistery),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomGeneratorWithSeed(defaultFoodTestSeed)
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

			food := NewFood(rng, box, tt.foodType)

			if food == nil {
				t.Fatal("Expected valid food, got nil")
			}

			if food.Type != tt.foodType {
				t.Errorf("Expected type %q, got %q", tt.foodType, food.Type)
			}

			if food.Name != tt.expectedName {
				t.Errorf("Expected name %q, got %q", tt.expectedName, food.Name)
			}

			// Verify attributes match the config's ranges
			config := GetFoodConfig(tt.foodType)
			if food.Effect.Attributes.Health < config.HealthRange.Min ||
				food.Effect.Attributes.Health > config.HealthRange.Max {
				t.Errorf("Health out of config range: got %.2f, want [%.2f, %.2f]",
					food.Effect.Attributes.Health, config.HealthRange.Min, config.HealthRange.Max)
			}

			if food.Name != string(config.Type) {
				t.Errorf("Expected name %q, got %q", string(config.Type), food.Name)
			}

			// Check box is correctly set
			if food.Shape != box {
				t.Errorf("Expected box %v, got %v", box, food.Shape)
			}

			// Other attributes should be zero
			if food.Effect.Attributes.Strength != 0 {
				t.Errorf("Expected Strength=0, got %f", food.Effect.Attributes.Strength)
			}
			if food.Effect.Attributes.Agility != 0 {
				t.Errorf("Expected Agility=0, got %f", food.Effect.Attributes.Agility)
			}
		})
	}
}

func TestFood_NewFoodByConfig(t *testing.T) {
	tests := []struct {
		name      string
		seed      int64
		config    FoodConfig
		box       primitives.Box
		wantError bool
	}{
		{
			name: "Valid custom config",
			seed: defaultFoodTestSeed,
			config: FoodConfig{
				Type:        "Custom Food",
				HealthRange: primitives.AttributeRange{Min: 15, Max: 35},
				Description: "A custom test food",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 3, Y: 3}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: false,
		},
		{
			name: "Valid config with negative health",
			seed: defaultFoodTestSeed,
			config: FoodConfig{
				Type:        "Poison Food",
				HealthRange: primitives.AttributeRange{Min: -20, Max: -5},
				Description: "Damages health",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 1}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: false,
		},
		{
			name: "Valid config with fixed value",
			seed: defaultFoodTestSeed,
			config: FoodConfig{
				Type:        "Fixed Health Food",
				HealthRange: primitives.AttributeRange{Min: 25, Max: 25},
				Description: "Always gives 25 HP",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 2, Y: 2}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: false,
		},
		{
			name: "Invalid health range returns error",
			seed: defaultFoodTestSeed,
			config: FoodConfig{
				Type:        "Invalid Food",
				HealthRange: primitives.AttributeRange{Min: 50, Max: 10}, // Min > Max
				Description: "Invalid config",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: true,
		},
		{
			name: "Invalid negative to positive inverted range",
			seed: defaultFoodTestSeed,
			config: FoodConfig{
				Type:        "Invalid Food 2",
				HealthRange: primitives.AttributeRange{Min: 20, Max: -10}, // Min > Max
				Description: "Invalid config",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomGeneratorWithSeed(tt.seed)
			food, err := NewFoodByConfig(rng, tt.box, tt.config)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if food != nil {
					t.Error("Expected nil food on error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if food == nil {
				t.Fatal("Expected valid food, got nil")
			}

			if food.Type != FoodTypeCustom {
				t.Errorf("Expected type %q, got %q", FoodTypeCustom, food.Type)
			}

			// Verify attributes are within config ranges
			if food.Effect.Attributes.Health < tt.config.HealthRange.Min ||
				food.Effect.Attributes.Health > tt.config.HealthRange.Max {
				t.Errorf("Health out of range: got %.2f, want [%.2f, %.2f]",
					food.Effect.Attributes.Health, tt.config.HealthRange.Min, tt.config.HealthRange.Max)
			}

			// Verify name matches config type
			if food.Name != string(tt.config.Type) {
				t.Errorf("Expected name %q, got %q", tt.config.Type, food.Name)
			}

			// Verify box is set correctly
			if food.Shape != tt.box {
				t.Errorf("Expected box %v, got %v", tt.box, food.Shape)
			}

			// For fixed-value range, verify exact value
			if tt.config.HealthRange.Min == tt.config.HealthRange.Max {
				if food.Effect.Attributes.Health != tt.config.HealthRange.Min {
					t.Errorf("Expected exact Health=%f for fixed range, got %f",
						tt.config.HealthRange.Min, food.Effect.Attributes.Health)
				}
			}
		})
	}
}

func TestFood_NewFood_Deterministic(t *testing.T) {
	tests := []struct {
		name     string
		seed     int64
		foodType FoodType
	}{
		{
			name:     "Same seed produces same Potatoes",
			seed:     defaultFoodTestSeed,
			foodType: FoodTypePotatoes,
		},
		{
			name:     "Same seed produces same Meat",
			seed:     defaultFoodTestSeed,
			foodType: FoodTypeMeat,
		},
		{
			name:     "Same seed produces same Mistery",
			seed:     defaultFoodTestSeed,
			foodType: FoodTypeMistery,
		},
		{
			name:     "Same seed produces same Beer - fixed value",
			seed:     defaultFoodTestSeed,
			foodType: FoodTypeBeer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

			// Create first food
			rng1 := utils.NewRandomGeneratorWithSeed(tt.seed)
			food1 := NewFood(rng1, box, tt.foodType)

			// Create second food with same seed
			rng2 := utils.NewRandomGeneratorWithSeed(tt.seed)
			food2 := NewFood(rng2, box, tt.foodType)

			// Verify they are identical
			if food1.Effect.Attributes.Health != food2.Effect.Attributes.Health {
				t.Errorf("Health mismatch: %f != %f", food1.Effect.Attributes.Health, food2.Effect.Attributes.Health)
			}

			if food1.Name != food2.Name {
				t.Errorf("Name mismatch: %q != %q", food1.Name, food2.Name)
			}
		})
	}
}

func TestFood_NewFood_Randomness(t *testing.T) {
	tests := []struct {
		name       string
		foodType   FoodType
		iterations int
	}{
		{
			name:       "Potatoes produces varied values",
			foodType:   FoodTypePotatoes,
			iterations: 50,
		},
		{
			name:       "Mistery produces varied values",
			foodType:   FoodTypeMistery,
			iterations: 50,
		},
		{
			name:       "Meat produces varied values",
			foodType:   FoodTypeMeat,
			iterations: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
			healthValues := make(map[float64]bool)

			// Generate multiple foods with different seeds
			for i := 0; i < tt.iterations; i++ {
				rng := utils.NewRandomGeneratorWithSeed(int64(i))
				food := NewFood(rng, box, tt.foodType)
				healthValues[food.Effect.Attributes.Health] = true
			}

			// Check that we got varied values (at least 5 different values)
			minUniqueValues := 5
			if len(healthValues) < minUniqueValues {
				t.Errorf("Expected at least %d unique health values, got %d", minUniqueValues, len(healthValues))
			}
		})
	}
}

func TestFood_NewFood_BeerProducesFixedValue(t *testing.T) {
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	expectedHealth := 20.0
	iterations := 30

	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomGeneratorWithSeed(int64(i))
		food := NewFood(rng, box, FoodTypeBeer)

		if food.Effect.Attributes.Health != expectedHealth {
			t.Errorf("Iteration %d: Expected Health=%f, got %f",
				i, expectedHealth, food.Effect.Attributes.Health)
		}
	}
}

func TestFood_NewFood_UnknownTypeDefaultsToMistery(t *testing.T) {
	rng := utils.NewRandomGeneratorWithSeed(defaultFoodTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

	food := NewFood(rng, box, "Unknown Food Type")

	if food == nil {
		t.Fatal("Expected valid food, got nil")
	}

	if food.Name != string(FoodTypeMistery) {
		t.Errorf("Expected name %q for unknown type, got %q", FoodTypeMistery, food.Name)
	}

	// Should use Mistery config ranges
	misteryConfig := GetFoodConfig(FoodTypeMistery)
	if food.Effect.Attributes.Health < misteryConfig.HealthRange.Min ||
		food.Effect.Attributes.Health > misteryConfig.HealthRange.Max {
		t.Errorf("Health out of Mistery config range: got %.2f, want [%.2f, %.2f]",
			food.Effect.Attributes.Health, misteryConfig.HealthRange.Min, misteryConfig.HealthRange.Max)
	}
}

func TestFood_NewFood_ZeroSizedBoxIsOk(t *testing.T) {
	rng := utils.NewRandomGeneratorWithSeed(defaultFoodTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 0, Height: 0}}
	food := NewFood(rng, box, FoodTypeMeat)

	if food == nil {
		t.Fatal("Expected valid food with zero-sized box")
	}

	if food.Shape != box {
		t.Errorf("Expected box %v, got %v", box, food.Shape)
	}
}

func TestFood_NewFood_VariousPositions(t *testing.T) {
	tests := []struct {
		name     string
		position primitives.Point2D[int]
		size     primitives.Size2D[uint]
	}{
		{
			name:     "Origin position",
			position: primitives.Point2D[int]{X: 0, Y: 0},
			size:     primitives.Size2D[uint]{Width: 1, Height: 1},
		},
		{
			name:     "Positive position",
			position: primitives.Point2D[int]{X: 100, Y: 200},
			size:     primitives.Size2D[uint]{Width: 2, Height: 2},
		},
		{
			name:     "Negative position",
			position: primitives.Point2D[int]{X: -10, Y: -20},
			size:     primitives.Size2D[uint]{Width: 1, Height: 1},
		},
		{
			name:     "Large size",
			position: primitives.Point2D[int]{X: 5, Y: 5},
			size:     primitives.Size2D[uint]{Width: 10, Height: 15},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomGeneratorWithSeed(defaultFoodTestSeed)
			box := primitives.Box{Point: tt.position, Size: tt.size}
			food := NewFood(rng, box, FoodTypeBread)

			if food == nil {
				t.Fatal("Expected valid food")
			}

			if food.Shape.Point != tt.position {
				t.Errorf("Expected position %v, got %v", tt.position, food.Shape.Point)
			}

			if food.Shape.Size != tt.size {
				t.Errorf("Expected size %v, got %v", tt.size, food.Shape.Size)
			}
		})
	}
}

func TestFood_AsFood_ValidPointer(t *testing.T) {
	input := &Food{
		Item: &Item{
			Shape: primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			Name:  "Test Food",
		},
		Effect: &primitives.Effect{
			Attributes: &primitives.Attributes{Health: 25},
		},
	}

	result := AsFood(input)

	if result == nil {
		t.Fatal("Expected non-nil result, got nil")
	}

	if result.Effect.Attributes.Health != 25 {
		t.Errorf("Expected Health 25, got %f", result.Effect.Attributes.Health)
	}

	if result != input {
		t.Error("Expected same pointer to be returned")
	}
}

func TestFood_AsFood_NonFoodTypeIsNil(t *testing.T) {
	result := AsFood("not a food")

	if result != nil {
		t.Errorf("Expected nil for non-Food type, got %v", result)
	}
}

func TestFood_AsFood_NilInputIsNil(t *testing.T) {
	result := AsFood(nil)

	if result != nil {
		t.Errorf("Expected nil for nil input, got %v", result)
	}
}

func TestFood_AsFood_DifferentStructTypeIsNil(t *testing.T) {
	input := &Item{
		Shape: primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 1}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		Name:  "Just an item",
	}

	result := AsFood(input)

	if result != nil {
		t.Errorf("Expected nil for Item type, got %v", result)
	}
}

func TestFood_AsFood_ElixirTypeIsNil(t *testing.T) {
	input := &Elixir{
		Item: &Item{
			Shape: primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			Name:  "Test Elixir",
		},
		Effect: &primitives.Effect{
			Attributes: &primitives.Attributes{Strength: 10},
			Duration:   20,
		},
	}

	result := AsFood(input)

	if result != nil {
		t.Errorf("Expected nil for Elixir type, got %v", result)
	}
}

// Benchmarks
func BenchmarkFood_NewFood(b *testing.B) {
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

	b.Run("Potatoes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomGeneratorWithSeed(int64(i))
			_ = NewFood(rng, box, FoodTypePotatoes)
		}
	})

	b.Run("Meat", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomGeneratorWithSeed(int64(i))
			_ = NewFood(rng, box, FoodTypeMeat)
		}
	})

	b.Run("Mistery", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomGeneratorWithSeed(int64(i))
			_ = NewFood(rng, box, FoodTypeMistery)
		}
	})

	b.Run("Beer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomGeneratorWithSeed(int64(i))
			_ = NewFood(rng, box, FoodTypeBeer)
		}
	})
}

func BenchmarkFood_NewFoodByConfig(b *testing.B) {
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	config := FoodConfig{
		Type:        "Benchmark Food",
		HealthRange: primitives.AttributeRange{Min: 10, Max: 50},
		Description: "Benchmark test food",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng := utils.NewRandomGeneratorWithSeed(int64(i))
		_, _ = NewFoodByConfig(rng, box, config)
	}
}

func BenchmarkFood_AsFood(b *testing.B) {
	food := &Food{
		Item: &Item{
			Shape: primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			Name:  "Test Food",
		},
		Effect: &primitives.Effect{
			Attributes: &primitives.Attributes{Health: 25},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AsFood(food)
	}
}
