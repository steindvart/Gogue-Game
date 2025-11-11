package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
	"testing"
)

const defaultTestSeed int64 = 42

func TestElixir_NewElixir(t *testing.T) {
	tests := []struct {
		name               string
		seed               int64
		elixirType         ElixirType
		box                primitives.Box
		wantStrengthMin    float64
		wantStrengthMax    float64
		wantAgilityMin     float64
		wantAgilityMax     float64
		wantDurationMin    uint32
		wantDurationMax    uint32
		checkName          bool
		expectedNamePrefix string
	}{
		{
			name:               "Elixir of Strength - positive strength only",
			seed:               1,
			elixirType:         ElixirTypeStrength,
			box:                primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantStrengthMin:    5,
			wantStrengthMax:    20,
			wantAgilityMin:     0,
			wantAgilityMax:     0,
			wantDurationMin:    uint32(ElixirDurationMin),
			wantDurationMax:    uint32(ElixirMaxDurationMax),
			checkName:          true,
			expectedNamePrefix: string(ElixirTypeStrength),
		},
		{
			name:               "Elixir of Agility - positive agility only",
			seed:               2,
			elixirType:         ElixirTypeAgility,
			box:                primitives.Box{Point: primitives.Point2D[int]{X: 10, Y: 10}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantStrengthMin:    0,
			wantStrengthMax:    0,
			wantAgilityMin:     5,
			wantAgilityMax:     20,
			wantDurationMin:    uint32(ElixirDurationMin),
			wantDurationMax:    uint32(ElixirMaxDurationMax),
			checkName:          true,
			expectedNamePrefix: string(ElixirTypeAgility),
		},
		{
			name:               "Phantom's Breath - positive agility, negative strength",
			seed:               3,
			elixirType:         ElixirTypePhantomBreath,
			box:                primitives.Box{Point: primitives.Point2D[int]{X: 3, Y: 7}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantStrengthMin:    -10,
			wantStrengthMax:    -2,
			wantAgilityMin:     10,
			wantAgilityMax:     30,
			wantDurationMin:    uint32(ElixirDurationMin),
			wantDurationMax:    uint32(ElixirMaxDurationMax),
			checkName:          true,
			expectedNamePrefix: string(ElixirTypePhantomBreath),
		},
		{
			name:               "Frozen Star - positive strength, negative agility",
			seed:               4,
			elixirType:         ElixirTypeFrozenStar,
			box:                primitives.Box{Point: primitives.Point2D[int]{X: 15, Y: 20}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantStrengthMin:    10,
			wantStrengthMax:    30,
			wantAgilityMin:     -10,
			wantAgilityMax:     -2,
			wantDurationMin:    uint32(ElixirDurationMin),
			wantDurationMax:    uint32(ElixirMaxDurationMax),
			checkName:          true,
			expectedNamePrefix: string(ElixirTypeFrozenStar),
		},
		{
			name:               "Elixir of Mystery - random attributes",
			seed:               5,
			elixirType:         ElixirTypeMystery,
			box:                primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantStrengthMin:    -20,
			wantStrengthMax:    30,
			wantAgilityMin:     -20,
			wantAgilityMax:     30,
			wantDurationMin:    uint32(ElixirDurationMin),
			wantDurationMax:    uint32(ElixirMaxDurationMax),
			checkName:          true,
			expectedNamePrefix: string(ElixirTypeMystery),
		},
		{
			name:               "Unknown type - defaults to Mystery",
			seed:               6,
			elixirType:         "Unknown Type",
			box:                primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 1}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantStrengthMin:    -20,
			wantStrengthMax:    30,
			wantAgilityMin:     -20,
			wantAgilityMax:     30,
			wantDurationMin:    uint32(ElixirDurationMin),
			wantDurationMax:    uint32(ElixirMaxDurationMax),
			checkName:          true,
			expectedNamePrefix: string(ElixirTypeMystery),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomGeneratorWithSeed(tt.seed)
			elixir := NewElixir(*rng, tt.box, tt.elixirType)

			if elixir == nil {
				t.Fatal("Expected a valid Elixir, got nil")
			}

			// Check attributes are within expected ranges
			if elixir.AffectedAttributes.Strength < tt.wantStrengthMin || elixir.AffectedAttributes.Strength > tt.wantStrengthMax {
				t.Errorf("Strength out of range: got %.2f, want [%.2f, %.2f]",
					elixir.AffectedAttributes.Strength, tt.wantStrengthMin, tt.wantStrengthMax)
			}

			if elixir.AffectedAttributes.Agility < tt.wantAgilityMin || elixir.AffectedAttributes.Agility > tt.wantAgilityMax {
				t.Errorf("Agility out of range: got %.2f, want [%.2f, %.2f]",
					elixir.AffectedAttributes.Agility, tt.wantAgilityMin, tt.wantAgilityMax)
			}

			// Check duration is within expected range
			if elixir.EffectDuration < tt.wantDurationMin || elixir.EffectDuration > tt.wantDurationMax {
				t.Errorf("EffectDuration out of range: got %d, want [%d, %d]",
					elixir.EffectDuration, tt.wantDurationMin, tt.wantDurationMax)
			}

			// Check name
			if tt.checkName && elixir.Name != tt.expectedNamePrefix {
				t.Errorf("Expected name %q, got %q", tt.expectedNamePrefix, elixir.Name)
			}

			// Check box is correctly set
			if elixir.Shape != tt.box {
				t.Errorf("Expected box %v, got %v", tt.box, elixir.Shape)
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
			seed:       defaultTestSeed,
			elixirType: ElixirTypeStrength,
		},
		{
			name:       "Same seed produces same Agility elixir",
			seed:       defaultTestSeed,
			elixirType: ElixirTypeAgility,
		},
		{
			name:       "Same seed produces same Mystery elixir",
			seed:       defaultTestSeed,
			elixirType: ElixirTypeMystery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

			// Create first elixir
			rng1 := utils.NewRandomGeneratorWithSeed(tt.seed)
			elixir1 := NewElixir(*rng1, box, tt.elixirType)

			// Create second elixir with same seed
			rng2 := utils.NewRandomGeneratorWithSeed(tt.seed)
			elixir2 := NewElixir(*rng2, box, tt.elixirType)

			// Verify they are identical
			if elixir1.AffectedAttributes.Strength != elixir2.AffectedAttributes.Strength {
				t.Errorf("Strength mismatch: %f != %f", elixir1.AffectedAttributes.Strength, elixir2.AffectedAttributes.Strength)
			}

			if elixir1.AffectedAttributes.Agility != elixir2.AffectedAttributes.Agility {
				t.Errorf("Agility mismatch: %f != %f", elixir1.AffectedAttributes.Agility, elixir2.AffectedAttributes.Agility)
			}

			if elixir1.EffectDuration != elixir2.EffectDuration {
				t.Errorf("Duration mismatch: %d != %d", elixir1.EffectDuration, elixir2.EffectDuration)
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
				elixir := NewElixir(*rng, box, tt.elixirType)

				strengthValues[elixir.AffectedAttributes.Strength] = true
				agilityValues[elixir.AffectedAttributes.Agility] = true
				durationValues[elixir.EffectDuration] = true
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
	rng := utils.NewRandomGeneratorWithSeed(defaultTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 0, Height: 0}}
	elixir := NewElixir(*rng, box, ElixirTypeStrength)

	if elixir == nil {
		t.Fatal("Expected valid elixir with zero-sized box")
	}

	if elixir.Shape != box {
		t.Errorf("Expected box %v, got %v", box, elixir.Shape)
	}
}

func TestElixir_Drop(t *testing.T) {
	rng := utils.NewRandomGeneratorWithSeed(defaultTestSeed)
	initialBox := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	elixir := NewElixir(*rng, initialBox, ElixirTypeStrength)

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
			seed:       defaultTestSeed,
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
			elixir := NewElixir(*rng, box, tt.elixirType)

			attrs := elixir.Use()

			// Verify returned attributes match stored attributes
			if attrs.Strength != elixir.AffectedAttributes.Strength {
				t.Errorf("Strength mismatch: got %f, want %f", attrs.Strength, elixir.AffectedAttributes.Strength)
			}

			if attrs.Agility != elixir.AffectedAttributes.Agility {
				t.Errorf("Agility mismatch: got %f, want %f", attrs.Agility, elixir.AffectedAttributes.Agility)
			}
		})
	}
}

func TestElixir_UseMultipleTimesIsOk(t *testing.T) {
	rng := utils.NewRandomGeneratorWithSeed(defaultTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	elixir := NewElixir(*rng, box, ElixirTypeStrength)

	attrs1 := elixir.Use()
	attrs2 := elixir.Use()

	if attrs1.Strength != attrs2.Strength {
		t.Error("Use() should return consistent attributes on multiple calls")
	}
	if attrs1.Agility != attrs2.Agility {
		t.Error("Use() should return consistent attributes on multiple calls")
	}
}

func TestElixir_AsElixir_ValidPointerIsElixir(t *testing.T) {
	input := &Elixir{
		Item: &Item{
			Shape: primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			Name:  "Test Elixir",
		},
		AffectedAttributes: primitives.Attributes{Strength: 10},
		EffectDuration:     20,
	}

	result := AsElixir(input)

	if result == nil {
		t.Fatal("Expected non-nil result, got nil")
	}

	if result.AffectedAttributes.Strength != 10 {
		t.Errorf("Expected Strength 10, got %f", result.AffectedAttributes.Strength)
	}
	if result.EffectDuration != 20 {
		t.Errorf("Expected Duration 20, got %d", result.EffectDuration)
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
			_ = NewElixir(*rng, box, ElixirTypeStrength)
		}
	})

	b.Run("Mystery", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomGeneratorWithSeed(int64(i))
			_ = NewElixir(*rng, box, ElixirTypeMystery)
		}
	})
}

func BenchmarkElixir_Use(b *testing.B) {
	rng := utils.NewRandomGeneratorWithSeed(defaultTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	elixir := NewElixir(*rng, box, ElixirTypeStrength)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = elixir.Use()
	}
}
