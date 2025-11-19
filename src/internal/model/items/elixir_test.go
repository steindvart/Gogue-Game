package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
	"testing"
)

func TestElixir_NewElixir_BuiltinConfig(t *testing.T) {
	tests := []struct {
		name         string
		elixirType   ElixirType
		expectedName string
	}{
		{
			name:         "Strength type uses Strength elixir config",
			elixirType:   ElixirTypeStrength,
			expectedName: string(ElixirTypeStrength),
		},
		{
			name:         "Agility type uses Agility elixir config",
			elixirType:   ElixirTypeAgility,
			expectedName: string(ElixirTypeAgility),
		},
		{
			name:         "PhantomBreath type uses PhantomBreath elixir config",
			elixirType:   ElixirTypeDwarfism,
			expectedName: string(ElixirTypeDwarfism),
		},
		{
			name:         "FrozenStar type uses FrozenStar elixir config",
			elixirType:   ElixirTypeGiantism,
			expectedName: string(ElixirTypeGiantism),
		},
		{
			name:         "Mystery type uses Mystery elixir config",
			elixirType:   ElixirTypeMystery,
			expectedName: string(ElixirTypeMystery),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomGeneratorWithSeed(defaultItemsTestSeed)
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

			elixir := NewElixir(rng, box, tt.elixirType)

			if elixir == nil {
				t.Fatal("Expected valid elixir, got nil")
			}

			if elixir.Type != tt.elixirType {
				t.Errorf("Expected type %q, got %q", tt.expectedName, elixir.Type)
			}

			if elixir.Name != tt.expectedName {
				t.Errorf("Expected name %q, got %q", tt.expectedName, elixir.Name)
			}

			// Verify attributes match the config's ranges
			config := GetElixirConfig(tt.elixirType)
			if elixir.Effect.Attributes.Strength < config.StrengthRange.Min ||
				elixir.Effect.Attributes.Strength > config.StrengthRange.Max {
				t.Errorf("Strength out of config range: got %.2f, want [%.2f, %.2f]",
					elixir.Effect.Attributes.Strength, config.StrengthRange.Min, config.StrengthRange.Max)
			}

			if elixir.Effect.Attributes.Agility < config.AgilityRange.Min ||
				elixir.Effect.Attributes.Agility > config.AgilityRange.Max {
				t.Errorf("Agility out of config range: got %.2f, want [%.2f, %.2f]",
					elixir.Effect.Attributes.Agility, config.AgilityRange.Min, config.AgilityRange.Max)
			}

			if elixir.Duration.Steps < config.DurationStepsRange.Min ||
				elixir.Duration.Steps > config.DurationStepsRange.Max {
				t.Errorf("Duration out of config range: got %d, want [%d, %d]",
					elixir.Duration.Steps, config.DurationStepsRange.Min, config.DurationStepsRange.Max)
			}

			if elixir.Name != string(config.Type) {
				t.Errorf("Expected name %q, got %q", string(config.Type), elixir.Name)
			}
		})
	}
}

func TestElixir_NewElixir_CustomConfig(t *testing.T) {
	tests := []struct {
		name      string
		seed      int64
		config    ElixirConfig
		box       primitives.Box
		wantError bool
	}{
		{
			name: "Valid custom config",
			seed: defaultItemsTestSeed,
			config: ElixirConfig{
				Type:               "Custom Elixir",
				StrengthRange:      primitives.AttributeRange{Min: 10, Max: 50},
				AgilityRange:       primitives.AttributeRange{Min: -5, Max: 15},
				DurationStepsRange: ElixirDurationStepsRange{Min: 10, Max: 20},
				Description:        "A custom test elixir",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 3, Y: 3}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: false,
		},
		{
			name: "Invalid strength range returns error",
			seed: defaultItemsTestSeed,
			config: ElixirConfig{
				Type:               "Invalid Elixir",
				StrengthRange:      primitives.AttributeRange{Min: 50, Max: 10}, // Min > Max
				AgilityRange:       primitives.AttributeRange{Min: 0, Max: 10},
				DurationStepsRange: defaultDurationRange,
				Description:        "Invalid config",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: true,
		},
		{
			name: "Invalid agility range returns error",
			seed: defaultItemsTestSeed,
			config: ElixirConfig{
				Type:               "Invalid Elixir",
				StrengthRange:      primitives.AttributeRange{Min: 0, Max: 10},
				AgilityRange:       primitives.AttributeRange{Min: 20, Max: 5}, // Min > Max
				DurationStepsRange: defaultDurationRange,
				Description:        "Invalid config",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: true,
		},
		{
			name: "Invalid duration range returns error",
			seed: defaultItemsTestSeed,
			config: ElixirConfig{
				Type:               "Invalid Elixir",
				StrengthRange:      primitives.AttributeRange{Min: 0, Max: 10},
				AgilityRange:       primitives.AttributeRange{Min: 0, Max: 10},
				DurationStepsRange: ElixirDurationStepsRange{Min: 50, Max: 10}, // Min > Max
				Description:        "Invalid config",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomGeneratorWithSeed(tt.seed)
			elixir, err := NewElixirByConfig(rng, tt.box, tt.config)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if elixir != nil {
					t.Error("Expected nil elixir on error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if elixir == nil {
				t.Fatal("Expected valid elixir, got nil")
			}

			if elixir.Type != ElixirTypeCustom {
				t.Errorf("Expected type %q, got %q", ElixirTypeCustom, elixir.Type)
			}

			// Verify attributes are within config ranges
			if elixir.Effect.Attributes.Strength < tt.config.StrengthRange.Min ||
				elixir.Effect.Attributes.Strength > tt.config.StrengthRange.Max {
				t.Errorf("Strength out of range: got %.2f, want [%.2f, %.2f]",
					elixir.Effect.Attributes.Strength, tt.config.StrengthRange.Min, tt.config.StrengthRange.Max)
			}

			if elixir.Effect.Attributes.Agility < tt.config.AgilityRange.Min ||
				elixir.Effect.Attributes.Agility > tt.config.AgilityRange.Max {
				t.Errorf("Agility out of range: got %.2f, want [%.2f, %.2f]",
					elixir.Effect.Attributes.Agility, tt.config.AgilityRange.Min, tt.config.AgilityRange.Max)
			}

			if elixir.Duration.Steps < tt.config.DurationStepsRange.Min ||
				elixir.Duration.Steps > tt.config.DurationStepsRange.Max {
				t.Errorf("Duration out of range: got %d, want [%d, %d]",
					elixir.Duration.Steps, tt.config.DurationStepsRange.Min, tt.config.DurationStepsRange.Max)
			}

			// Verify name matches config type
			if elixir.Name != string(tt.config.Type) {
				t.Errorf("Expected name %q, got %q", tt.config.Type, elixir.Name)
			}
		})
	}
}

func TestElixir_NewElixir_NoRandom(t *testing.T) {
	tests := []struct {
		name       string
		seed       int64
		elixirType ElixirType
	}{
		{
			name:       "Same seed produces same Strength elixir",
			seed:       defaultItemsTestSeed,
			elixirType: ElixirTypeStrength,
		},
		{
			name:       "Same seed produces same Agility elixir",
			seed:       defaultItemsTestSeed,
			elixirType: ElixirTypeAgility,
		},
		{
			name:       "Same seed produces same Mystery elixir",
			seed:       defaultItemsTestSeed,
			elixirType: ElixirTypeMystery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

			// Create first elixir
			rng1 := utils.NewRandomGeneratorWithSeed(tt.seed)
			elixir1 := NewElixir(rng1, box, tt.elixirType)

			// Create second elixir with same seed
			rng2 := utils.NewRandomGeneratorWithSeed(tt.seed)
			elixir2 := NewElixir(rng2, box, tt.elixirType)

			// Verify they are identical
			if elixir1.Effect.Attributes.Strength != elixir2.Effect.Attributes.Strength {
				t.Errorf("Strength mismatch: %f != %f", elixir1.Effect.Attributes.Strength, elixir2.Effect.Attributes.Strength)
			}

			if elixir1.Effect.Attributes.Agility != elixir2.Effect.Attributes.Agility {
				t.Errorf("Agility mismatch: %f != %f", elixir1.Effect.Attributes.Agility, elixir2.Effect.Attributes.Agility)
			}

			if elixir1.Duration.Steps != elixir2.Duration.Steps {
				t.Errorf("Duration mismatch: %d != %d", elixir1.Duration.Steps, elixir2.Duration.Steps)
			}
		})
	}
}

func TestElixir_NewElixir_Randomness(t *testing.T) {
	tests := []struct {
		name       string
		elixirType ElixirType
		iterations int
	}{
		{
			name:       "Strength elixir produces varied values",
			elixirType: ElixirTypeStrength,
			iterations: 50,
		},
		{
			name:       "Mystery elixir produces varied values",
			elixirType: ElixirTypeMystery,
			iterations: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
			strengthValues := make(map[float64]bool)
			agilityValues := make(map[float64]bool)
			durationValues := make(map[uint32]bool)

			// Generate multiple elixirs with different seeds
			for i := 0; i < tt.iterations; i++ {
				rng := utils.NewRandomGeneratorWithSeed(int64(i))
				elixir := NewElixir(rng, box, tt.elixirType)

				strengthValues[elixir.Effect.Attributes.Strength] = true
				agilityValues[elixir.Effect.Attributes.Agility] = true
				durationValues[elixir.Duration.Steps] = true
			}

			// Check that we got varied values (at least 5 different values)
			minUniqueValues := 5
			if len(durationValues) < minUniqueValues {
				t.Errorf("Expected at least %d unique duration values, got %d", minUniqueValues, len(durationValues))
			}

			// For attributes, check based on type
			switch tt.elixirType {
			case ElixirTypeStrength:
				if len(strengthValues) < minUniqueValues {
					t.Errorf("Expected at least %d unique strength values, got %d", minUniqueValues, len(strengthValues))
				}
			case ElixirTypeMystery:
				if len(strengthValues) < minUniqueValues {
					t.Errorf("Expected at least %d unique strength values, got %d", minUniqueValues, len(strengthValues))
				}
				if len(agilityValues) < minUniqueValues {
					t.Errorf("Expected at least %d unique agility values, got %d", minUniqueValues, len(agilityValues))
				}
			}
		})
	}
}

func TestElixir_NewElixir_ZeroSizedBoxIsOk(t *testing.T) {
	rng := utils.NewRandomGeneratorWithSeed(defaultItemsTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 0, Height: 0}}
	elixir := NewElixir(rng, box, ElixirTypeStrength)

	if elixir == nil {
		t.Fatal("Expected valid elixir with zero-sized box")
	}

	if elixir.Shape != box {
		t.Errorf("Expected box %v, got %v", box, elixir.Shape)
	}
}

func TestElixir_Drop(t *testing.T) {
	rng := utils.NewRandomGeneratorWithSeed(defaultItemsTestSeed)
	initialBox := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	elixir := NewElixir(rng, initialBox, ElixirTypeStrength)

	if elixir.Shape != initialBox {
		t.Errorf("Expected initial box %v, got %v", initialBox, elixir.Shape)
	}

	newPosition := primitives.Point2D[int]{X: 10, Y: 10}
	resultBox := elixir.Drop(newPosition)

	if elixir.Shape.Point != newPosition {
		t.Errorf("Expected position to be updated to %v, got %v", newPosition, elixir.Shape.Point)
	}

	if resultBox.Point != newPosition {
		t.Errorf("Expected returned box to have position %v, got %v", newPosition, resultBox.Point)
	}
}

func TestElixir_Use(t *testing.T) {
	tests := []struct {
		name       string
		seed       int64
		elixirType ElixirType
	}{
		{
			name:       "Use Strength elixir returns correct attributes",
			seed:       defaultItemsTestSeed,
			elixirType: ElixirTypeStrength,
		},
		{
			name:       "Use Agility elixir returns correct attributes",
			seed:       1,
			elixirType: ElixirTypeAgility,
		},
		{
			name:       "Use Mystery elixir returns correct attributes",
			seed:       2,
			elixirType: ElixirTypeMystery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomGeneratorWithSeed(tt.seed)
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
			elixir := NewElixir(rng, box, tt.elixirType)

			effect := elixir.Use()

			// Verify returned attributes match stored attributes
			if effect.Attributes.Strength != elixir.Effect.Attributes.Strength {
				t.Errorf("Strength mismatch: got %f, want %f", effect.Attributes.Strength, elixir.Effect.Attributes.Strength)
			}

			if effect.Attributes.Agility != elixir.Effect.Attributes.Agility {
				t.Errorf("Agility mismatch: got %f, want %f", effect.Attributes.Agility, elixir.Effect.Attributes.Agility)
			}
		})
	}
}

func TestElixir_UseMultipleTimesIsOk(t *testing.T) {
	rng := utils.NewRandomGeneratorWithSeed(defaultItemsTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	elixir := NewElixir(rng, box, ElixirTypeStrength)

	effect1 := elixir.Use()
	effect2 := elixir.Use()

	if effect1.Attributes.Strength != effect2.Attributes.Strength {
		t.Error("Use() should return consistent attributes on multiple calls")
	}
	if effect1.Attributes.Agility != effect2.Attributes.Agility {
		t.Error("Use() should return consistent attributes on multiple calls")
	}
}

func TestElixir_AsElixir_ValidPointer(t *testing.T) {
	input := &Elixir{
		Item: &Item{
			Shape: primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			Name:  "Test Elixir",
		},
		Effect: &primitives.Effect{
			Attributes: &primitives.Attributes{Strength: 10},
			Duration: &primitives.EffectDuration{
				Type:  primitives.EffectDurationTypeTemporary,
				Steps: 20,
			},
		},
	}

	result := AsElixir(input)

	if result == nil {
		t.Fatal("Expected non-nil result, got nil")
	}

	if result.Effect.Attributes.Strength != 10 {
		t.Errorf("Expected Strength 10, got %f", result.Effect.Attributes.Strength)
	}
	if result.Duration.Steps != 20 {
		t.Errorf("Expected Duration 20, got %d", result.Duration.Steps)
	}
	if result != input {
		t.Error("Expected same pointer to be returned")
	}
}

func TestElixir_AsElixir_NonElixirTypeIsNil(t *testing.T) {
	result := AsElixir("not an elixir")

	if result != nil {
		t.Errorf("Expected nil for non-Elixir type, got %v", result)
	}
}

func TestElixir_AsElixir_NilInputIsNil(t *testing.T) {
	result := AsElixir(nil)

	if result != nil {
		t.Errorf("Expected nil for nil input, got %v", result)
	}
}

func TestElixir_AsElixir_DifferentStructTypeIsNil(t *testing.T) {
	input := &Item{
		Shape: primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 1}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		Name:  "Just an item",
	}

	result := AsElixir(input)

	if result != nil {
		t.Errorf("Expected nil for Item type, got %v", result)
	}
}

// Benchmarks
func BenchmarkElixir_NewElixir(b *testing.B) {
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

	b.Run("Strength", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomGeneratorWithSeed(int64(i))
			_ = NewElixir(rng, box, ElixirTypeStrength)
		}
	})

	b.Run("Mystery", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomGeneratorWithSeed(int64(i))
			_ = NewElixir(rng, box, ElixirTypeMystery)
		}
	})
}

func BenchmarkElixir_Use(b *testing.B) {
	rng := utils.NewRandomGeneratorWithSeed(defaultItemsTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	elixir := NewElixir(rng, box, ElixirTypeStrength)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = elixir.Use()
	}
}
