package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
	"testing"
)

func TestTreasure_NewTreasure_BuiltinConfig(t *testing.T) {
	tests := []struct {
		name         string
		treasureType TreasureType
		expectedName string
	}{
		{
			name:         "Gold type uses Gold treasure config",
			treasureType: TreasureTypeGold,
			expectedName: string(TreasureTypeGold),
		},
		{
			name:         "Gem type uses Gem treasure config",
			treasureType: TreasureTypeGem,
			expectedName: string(TreasureTypeGem),
		},
		{
			name:         "Artifact type uses Artifact treasure config",
			treasureType: TreasureTypeArtifact,
			expectedName: string(TreasureTypeArtifact),
		},
		{
			name:         "Mystery type uses Mystery treasure config",
			treasureType: TreasureTypeMystery,
			expectedName: string(TreasureTypeMystery),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomWithSeed(defaultItemsTestSeed)
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

			treasure := NewTreasureBuiltin(rng, box, tt.treasureType)

			if treasure == nil {
				t.Fatal("Expected valid treasure, got nil")
			}

			if treasure.Type != tt.treasureType {
				t.Errorf("Expected type %q, got %q", tt.treasureType, treasure.Type)
			}

			if treasure.Name != tt.expectedName {
				t.Errorf("Expected name %q, got %q", tt.expectedName, treasure.Name)
			}

			// Verify value is within config's range
			config := GetTreasureConfig(tt.treasureType)
			if treasure.Value < config.ValueRange.Min || treasure.Value > config.ValueRange.Max {
				t.Errorf("Value out of config range: got %d, want [%d, %d]",
					treasure.Value, config.ValueRange.Min, config.ValueRange.Max)
			}

			if treasure.Name != string(config.Type) {
				t.Errorf("Expected name %q, got %q", string(config.Type), treasure.Name)
			}

			// Check box is correctly set
			if *treasure.Box != box {
				t.Errorf("Expected box %v, got %v", box, treasure.Box)
			}

			// Value should be positive
			if treasure.Value <= 0 {
				t.Errorf("Expected positive value, got %d", treasure.Value)
			}
		})
	}
}

func TestTreasure_NewTreasureByConfig(t *testing.T) {
	tests := []struct {
		name      string
		seed      int64
		config    TreasureConfig
		box       primitives.Box
		wantError bool
	}{
		{
			name: "Valid custom config",
			seed: defaultItemsTestSeed,
			config: TreasureConfig{
				Type:        "Custom Treasure",
				ValueRange:  TreasureValueRange{Min: 100, Max: 500},
				Description: "A custom test treasure",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 3, Y: 3}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: false,
		},
		{
			name: "Valid config with small range",
			seed: defaultItemsTestSeed,
			config: TreasureConfig{
				Type:        "Small Treasure",
				ValueRange:  TreasureValueRange{Min: 1, Max: 10},
				Description: "Small value treasure",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 1}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: false,
		},
		{
			name: "Valid config with large range",
			seed: defaultItemsTestSeed,
			config: TreasureConfig{
				Type:        "Large Treasure",
				ValueRange:  TreasureValueRange{Min: 1000, Max: 10000},
				Description: "Large value treasure",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 2, Y: 2}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: false,
		},
		{
			name: "Invalid value range returns error",
			seed: defaultItemsTestSeed,
			config: TreasureConfig{
				Type:        "Invalid Treasure",
				ValueRange:  TreasureValueRange{Min: 500, Max: 10}, // Min > Max
				Description: "Invalid config",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomWithSeed(tt.seed)
			treasure, err := NewTreasureByConfig(rng, tt.box, tt.config)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if treasure != nil {
					t.Error("Expected nil treasure on error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if treasure == nil {
				t.Fatal("Expected valid treasure, got nil")
			}

			if treasure.Type != tt.config.Type {
				t.Errorf("Expected type %q, got %q", tt.config.Type, treasure.Type)
			}

			// Verify value is within config range
			if treasure.Value < tt.config.ValueRange.Min || treasure.Value > tt.config.ValueRange.Max {
				t.Errorf("Value out of range: got %d, want [%d, %d]",
					treasure.Value, tt.config.ValueRange.Min, tt.config.ValueRange.Max)
			}

			// Verify name matches config type
			if treasure.Name != string(tt.config.Type) {
				t.Errorf("Expected name %q, got %q", tt.config.Type, treasure.Name)
			}

			// Verify box is set correctly
			if *treasure.Box != tt.box {
				t.Errorf("Expected box %v, got %v", tt.box, treasure.Box)
			}
		})
	}
}

func TestTreasure_NewTreasureByConfig_FixedValue(t *testing.T) {
	config := TreasureConfig{
		Type:        "Fixed Treasure",
		ValueRange:  TreasureValueRange{Min: 100, Max: 100},
		Description: "Always 100 value",
	}

	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	iterations := 20

	for i := 0; i < iterations; i++ {
		rng := utils.NewRandomWithSeed(int64(i))
		treasure, err := NewTreasureByConfig(rng, box, config)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if treasure.Value != 100 {
			t.Errorf("Expected fixed value=100, got %d", treasure.Value)
		}
	}
}

func TestTreasure_NewTreasure_Deterministic(t *testing.T) {
	tests := []struct {
		name         string
		seed         int64
		treasureType TreasureType
	}{
		{
			name:         "Same seed produces same Gold treasure",
			seed:         defaultItemsTestSeed,
			treasureType: TreasureTypeGold,
		},
		{
			name:         "Same seed produces same Gem treasure",
			seed:         defaultItemsTestSeed,
			treasureType: TreasureTypeGem,
		},
		{
			name:         "Same seed produces same Artifact treasure",
			seed:         defaultItemsTestSeed,
			treasureType: TreasureTypeArtifact,
		},
		{
			name:         "Same seed produces same Mystery treasure",
			seed:         defaultItemsTestSeed,
			treasureType: TreasureTypeMystery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

			// Create first treasure
			rng1 := utils.NewRandomWithSeed(tt.seed)
			treasure1 := NewTreasureBuiltin(rng1, box, tt.treasureType)

			// Create second treasure with same seed
			rng2 := utils.NewRandomWithSeed(tt.seed)
			treasure2 := NewTreasureBuiltin(rng2, box, tt.treasureType)

			// Verify they are identical
			if treasure1.Value != treasure2.Value {
				t.Errorf("Value mismatch: %d != %d", treasure1.Value, treasure2.Value)
			}

			if treasure1.Name != treasure2.Name {
				t.Errorf("Name mismatch: %q != %q", treasure1.Name, treasure2.Name)
			}
		})
	}
}

func TestTreasure_NewTreasure_Randomness(t *testing.T) {
	tests := []struct {
		name         string
		treasureType TreasureType
		iterations   int
	}{
		{
			name:         "Gold treasure produces varied values",
			treasureType: TreasureTypeGold,
			iterations:   50,
		},
		{
			name:         "Mystery treasure produces varied values",
			treasureType: TreasureTypeMystery,
			iterations:   50,
		},
		{
			name:         "Artifact treasure produces varied values",
			treasureType: TreasureTypeArtifact,
			iterations:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
			values := make(map[int32]bool)

			// Generate multiple treasures with different seeds
			for i := 0; i < tt.iterations; i++ {
				rng := utils.NewRandomWithSeed(int64(i))
				treasure := NewTreasureBuiltin(rng, box, tt.treasureType)
				values[treasure.Value] = true
			}

			// Check that we got varied values (at least 5 different values)
			minUniqueValues := 5
			config := GetTreasureConfig(tt.treasureType)

			// For ranges larger than minUniqueValues, expect variation
			rangeSize := config.ValueRange.Max - config.ValueRange.Min + 1
			if rangeSize > int32(minUniqueValues) {
				if len(values) < minUniqueValues {
					t.Errorf("Expected at least %d unique values, got %d", minUniqueValues, len(values))
				}
			}
		})
	}
}

func TestTreasure_NewTreasure_UnknownTypeDefaultsToMystery(t *testing.T) {
	rng := utils.NewRandomWithSeed(defaultItemsTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

	treasure := NewTreasureBuiltin(rng, box, "Unknown Treasure Type")

	if treasure == nil {
		t.Fatal("Expected valid treasure, got nil")
	}

	if treasure.Name != string(TreasureTypeMystery) {
		t.Errorf("Expected name %q for unknown type, got %q", TreasureTypeMystery, treasure.Name)
	}

	// Should use Mystery config range
	mysteryConfig := GetTreasureConfig(TreasureTypeMystery)
	if treasure.Value < mysteryConfig.ValueRange.Min || treasure.Value > mysteryConfig.ValueRange.Max {
		t.Errorf("Value out of Mystery config range: got %d, want [%d, %d]",
			treasure.Value, mysteryConfig.ValueRange.Min, mysteryConfig.ValueRange.Max)
	}
}

func TestTreasure_NewTreasure_ZeroSizedBoxIsOk(t *testing.T) {
	rng := utils.NewRandomWithSeed(defaultItemsTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 0, Height: 0}}
	treasure := NewTreasureBuiltin(rng, box, TreasureTypeGold)

	if treasure == nil {
		t.Fatal("Expected valid treasure with zero-sized box")
	}

	if *treasure.Box != box {
		t.Errorf("Expected box %v, got %v", box, treasure.Box)
	}
}

func TestTreasure_NewTreasure_VariousPositions(t *testing.T) {
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
			rng := utils.NewRandomWithSeed(defaultItemsTestSeed)
			box := primitives.Box{Point: tt.position, Size: tt.size}
			treasure := NewTreasureBuiltin(rng, box, TreasureTypeGem)

			if treasure == nil {
				t.Fatal("Expected valid treasure")
			}

			if treasure.Box.Point != tt.position {
				t.Errorf("Expected position %v, got %v", tt.position, treasure.Box.Point)
			}

			if treasure.Box.Size != tt.size {
				t.Errorf("Expected size %v, got %v", tt.size, treasure.Box.Size)
			}
		})
	}
}

func TestTreasure_Take(t *testing.T) {
	tests := []struct {
		name         string
		seed         int64
		treasureType TreasureType
	}{
		{
			name:         "Take Gold treasure returns correct value",
			seed:         defaultItemsTestSeed,
			treasureType: TreasureTypeGold,
		},
		{
			name:         "Take Gem treasure returns correct value",
			seed:         1,
			treasureType: TreasureTypeGem,
		},
		{
			name:         "Take Artifact treasure returns correct value",
			seed:         2,
			treasureType: TreasureTypeArtifact,
		},
		{
			name:         "Take Mystery treasure returns correct value",
			seed:         3,
			treasureType: TreasureTypeMystery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomWithSeed(tt.seed)
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
			treasure := NewTreasureBuiltin(rng, box, tt.treasureType)

			value := treasure.Take()

			// Verify returned value matches stored value
			if value != treasure.Value {
				t.Errorf("Value mismatch: got %d, want %d", value, treasure.Value)
			}

			// Value should be positive
			if value <= 0 {
				t.Errorf("Expected positive value, got %d", value)
			}
		})
	}
}

func TestTreasure_TakeMultipleTimesIsOk(t *testing.T) {
	rng := utils.NewRandomWithSeed(defaultItemsTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	treasure := NewTreasureBuiltin(rng, box, TreasureTypeGold)

	value1 := treasure.Take()
	value2 := treasure.Take()

	if value1 != value2 {
		t.Error("Take() should return consistent value on multiple calls")
	}
}

// Benchmarks
func BenchmarkTreasure_NewTreasure(b *testing.B) {
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

	b.Run("Gold", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomWithSeed(int64(i))
			_ = NewTreasureBuiltin(rng, box, TreasureTypeGold)
		}
	})

	b.Run("Gem", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomWithSeed(int64(i))
			_ = NewTreasureBuiltin(rng, box, TreasureTypeGem)
		}
	})

	b.Run("Artifact", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomWithSeed(int64(i))
			_ = NewTreasureBuiltin(rng, box, TreasureTypeArtifact)
		}
	})

	b.Run("Mystery", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomWithSeed(int64(i))
			_ = NewTreasureBuiltin(rng, box, TreasureTypeMystery)
		}
	})
}

func BenchmarkTreasure_NewTreasureByConfig(b *testing.B) {
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	config := TreasureConfig{
		Type:        "Benchmark Treasure",
		ValueRange:  TreasureValueRange{Min: 100, Max: 500},
		Description: "Benchmark test treasure",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng := utils.NewRandomWithSeed(int64(i))
		_, _ = NewTreasureByConfig(rng, box, config)
	}
}

func BenchmarkTreasure_Take(b *testing.B) {
	rng := utils.NewRandomWithSeed(defaultItemsTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	treasure := NewTreasureBuiltin(rng, box, TreasureTypeGold)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = treasure.Take()
	}
}
