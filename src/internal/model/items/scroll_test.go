package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
	"testing"
)

func TestScroll_NewScroll_BuiltinConfig(t *testing.T) {
	tests := []struct {
		name         string
		scrollType   ScrollType
		expectedName string
	}{
		{
			name:         "Strength type uses Strength scroll config",
			scrollType:   ScrollTypeStrength,
			expectedName: string(ScrollTypeStrength),
		},
		{
			name:         "Agility type uses Agility scroll config",
			scrollType:   ScrollTypeAgility,
			expectedName: string(ScrollTypeAgility),
		},
		{
			name:         "MaxHealth type uses MaxHealth scroll config",
			scrollType:   ScrollTypeMaxHealth,
			expectedName: string(ScrollTypeMaxHealth),
		},
		{
			name:         "Mystery type uses Mystery scroll config",
			scrollType:   ScrollTypeMystery,
			expectedName: string(ScrollTypeMystery),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomGeneratorWithSeed(defaultItemsTestSeed)
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

			scroll := NewScroll(rng, box, tt.scrollType)

			if scroll == nil {
				t.Fatal("Expected valid scroll, got nil")
			}

			if scroll.Type != tt.scrollType {
				t.Errorf("Expected type %q, got %q", tt.scrollType, scroll.Type)
			}

			if scroll.Name != tt.expectedName {
				t.Errorf("Expected name %q, got %q", tt.expectedName, scroll.Name)
			}

			// Verify attributes match the config's ranges
			config := GetScrollConfig(tt.scrollType)
			if scroll.Effect.Attributes.Strength < config.StrengthRange.Min ||
				scroll.Effect.Attributes.Strength > config.StrengthRange.Max {
				t.Errorf("Strength out of config range: got %.2f, want [%.2f, %.2f]",
					scroll.Effect.Attributes.Strength, config.StrengthRange.Min, config.StrengthRange.Max)
			}

			if scroll.Effect.Attributes.Agility < config.AgilityRange.Min ||
				scroll.Effect.Attributes.Agility > config.AgilityRange.Max {
				t.Errorf("Agility out of config range: got %.2f, want [%.2f, %.2f]",
					scroll.Effect.Attributes.Agility, config.AgilityRange.Min, config.AgilityRange.Max)
			}

			if scroll.Effect.Attributes.MaxHealth < config.MaxHealthRange.Min ||
				scroll.Effect.Attributes.MaxHealth > config.MaxHealthRange.Max {
				t.Errorf("MaxHealth out of config range: got %.2f, want [%.2f, %.2f]",
					scroll.Effect.Attributes.MaxHealth, config.MaxHealthRange.Min, config.MaxHealthRange.Max)
			}

			if scroll.Name != string(config.Type) {
				t.Errorf("Expected name %q, got %q", string(config.Type), scroll.Name)
			}

			// Check box is correctly set
			if scroll.Shape != box {
				t.Errorf("Expected box %v, got %v", box, scroll.Shape)
			}

			// Health should always be zero for scrolls
			if scroll.Effect.Attributes.Health != 0 {
				t.Errorf("Expected Health=0, got %f", scroll.Effect.Attributes.Health)
			}
		})
	}
}

func TestScroll_NewScrollByConfig(t *testing.T) {
	tests := []struct {
		name      string
		seed      int64
		config    ScrollConfig
		box       primitives.Box
		wantError bool
	}{
		{
			name: "Valid custom config with all attributes",
			seed: defaultItemsTestSeed,
			config: ScrollConfig{
				Type:           "Custom Scroll",
				StrengthRange:  primitives.AttributeRange{Min: 5, Max: 15},
				AgilityRange:   primitives.AttributeRange{Min: 3, Max: 12},
				MaxHealthRange: primitives.AttributeRange{Min: 10, Max: 25},
				Description:    "A custom test scroll",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 3, Y: 3}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: false,
		},
		{
			name: "Valid config with single attribute",
			seed: defaultItemsTestSeed,
			config: ScrollConfig{
				Type:           "Single Attribute Scroll",
				StrengthRange:  primitives.AttributeRange{Min: 10, Max: 20},
				AgilityRange:   primitives.AttributeRange{Min: 0, Max: 0},
				MaxHealthRange: primitives.AttributeRange{Min: 0, Max: 0},
				Description:    "Only strength",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 1}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: false,
		},
		{
			name: "Valid config with negative values",
			seed: defaultItemsTestSeed,
			config: ScrollConfig{
				Type:           "Curse Scroll",
				StrengthRange:  primitives.AttributeRange{Min: -10, Max: -2},
				AgilityRange:   primitives.AttributeRange{Min: -8, Max: -1},
				MaxHealthRange: primitives.AttributeRange{Min: -15, Max: -5},
				Description:    "Decreases all attributes",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 2, Y: 2}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: false,
		},
		{
			name: "Invalid strength range returns error",
			seed: defaultItemsTestSeed,
			config: ScrollConfig{
				Type:           "Invalid Scroll",
				StrengthRange:  primitives.AttributeRange{Min: 50, Max: 10}, // Min > Max
				AgilityRange:   primitives.AttributeRange{Min: 0, Max: 10},
				MaxHealthRange: primitives.AttributeRange{Min: 0, Max: 10},
				Description:    "Invalid config",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: true,
		},
		{
			name: "Invalid agility range returns error",
			seed: defaultItemsTestSeed,
			config: ScrollConfig{
				Type:           "Invalid Scroll",
				StrengthRange:  primitives.AttributeRange{Min: 0, Max: 10},
				AgilityRange:   primitives.AttributeRange{Min: 20, Max: 5}, // Min > Max
				MaxHealthRange: primitives.AttributeRange{Min: 0, Max: 10},
				Description:    "Invalid config",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: true,
		},
		{
			name: "Invalid MaxHealth range returns error",
			seed: defaultItemsTestSeed,
			config: ScrollConfig{
				Type:           "Invalid Scroll",
				StrengthRange:  primitives.AttributeRange{Min: 0, Max: 10},
				AgilityRange:   primitives.AttributeRange{Min: 0, Max: 10},
				MaxHealthRange: primitives.AttributeRange{Min: 50, Max: 10}, // Min > Max
				Description:    "Invalid config",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomGeneratorWithSeed(tt.seed)
			scroll, err := NewScrollByConfig(rng, tt.box, tt.config)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if scroll != nil {
					t.Error("Expected nil scroll on error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if scroll == nil {
				t.Fatal("Expected valid scroll, got nil")
			}

			if scroll.Type != ScrollTypeCustom {
				t.Errorf("Expected type %q, got %q", ScrollTypeCustom, scroll.Type)
			}

			// Verify attributes are within config ranges
			if scroll.Effect.Attributes.Strength < tt.config.StrengthRange.Min ||
				scroll.Effect.Attributes.Strength > tt.config.StrengthRange.Max {
				t.Errorf("Strength out of range: got %.2f, want [%.2f, %.2f]",
					scroll.Effect.Attributes.Strength, tt.config.StrengthRange.Min, tt.config.StrengthRange.Max)
			}

			if scroll.Effect.Attributes.Agility < tt.config.AgilityRange.Min ||
				scroll.Effect.Attributes.Agility > tt.config.AgilityRange.Max {
				t.Errorf("Agility out of range: got %.2f, want [%.2f, %.2f]",
					scroll.Effect.Attributes.Agility, tt.config.AgilityRange.Min, tt.config.AgilityRange.Max)
			}

			if scroll.Effect.Attributes.MaxHealth < tt.config.MaxHealthRange.Min ||
				scroll.Effect.Attributes.MaxHealth > tt.config.MaxHealthRange.Max {
				t.Errorf("MaxHealth out of range: got %.2f, want [%.2f, %.2f]",
					scroll.Effect.Attributes.MaxHealth, tt.config.MaxHealthRange.Min, tt.config.MaxHealthRange.Max)
			}

			// Verify name matches config type
			if scroll.Name != string(tt.config.Type) {
				t.Errorf("Expected name %q, got %q", tt.config.Type, scroll.Name)
			}

			// Verify box is set correctly
			if scroll.Shape != tt.box {
				t.Errorf("Expected box %v, got %v", tt.box, scroll.Shape)
			}

			// Health should always be zero
			if scroll.Effect.Attributes.Health != 0 {
				t.Errorf("Expected Health=0, got %f", scroll.Effect.Attributes.Health)
			}
		})
	}
}

func TestScroll_NewScrollByConfig_FixedValues(t *testing.T) {
	tests := []struct {
		name      string
		seed      int64
		config    ScrollConfig
		box       primitives.Box
		wantError bool
	}{
		{
			name: "Valid custom config with all fixed attributes",
			seed: defaultItemsTestSeed,
			config: ScrollConfig{
				Type:           "Custom Scroll",
				StrengthRange:  primitives.AttributeRange{Min: 5, Max: 5},
				AgilityRange:   primitives.AttributeRange{Min: 5, Max: 5},
				MaxHealthRange: primitives.AttributeRange{Min: 5, Max: 5},
				Description:    "A custom test scroll",
			},
			box:       primitives.Box{Point: primitives.Point2D[int]{X: 3, Y: 3}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomGeneratorWithSeed(tt.seed)
			scroll, err := NewScrollByConfig(rng, tt.box, tt.config)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if scroll != nil {
					t.Error("Expected nil scroll on error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if scroll == nil {
				t.Fatal("Expected valid scroll, got nil")
			}

			// Verify attributes are within config ranges
			if scroll.Effect.Attributes.Strength < tt.config.StrengthRange.Min ||
				scroll.Effect.Attributes.Strength > tt.config.StrengthRange.Max {
				t.Errorf("Strength out of range: got %.2f, want [%.2f, %.2f]",
					scroll.Effect.Attributes.Strength, tt.config.StrengthRange.Min, tt.config.StrengthRange.Max)
			}

			if scroll.Effect.Attributes.Agility < tt.config.AgilityRange.Min ||
				scroll.Effect.Attributes.Agility > tt.config.AgilityRange.Max {
				t.Errorf("Agility out of range: got %.2f, want [%.2f, %.2f]",
					scroll.Effect.Attributes.Agility, tt.config.AgilityRange.Min, tt.config.AgilityRange.Max)
			}

			if scroll.Effect.Attributes.MaxHealth < tt.config.MaxHealthRange.Min ||
				scroll.Effect.Attributes.MaxHealth > tt.config.MaxHealthRange.Max {
				t.Errorf("MaxHealth out of range: got %.2f, want [%.2f, %.2f]",
					scroll.Effect.Attributes.MaxHealth, tt.config.MaxHealthRange.Min, tt.config.MaxHealthRange.Max)
			}

			// Verify name matches config type
			if scroll.Name != string(tt.config.Type) {
				t.Errorf("Expected name %q, got %q", tt.config.Type, scroll.Name)
			}

			// Verify box is set correctly
			if scroll.Shape != tt.box {
				t.Errorf("Expected box %v, got %v", tt.box, scroll.Shape)
			}

			// Health should always be zero
			if scroll.Effect.Attributes.Health != 0 {
				t.Errorf("Expected Health=0, got %f", scroll.Effect.Attributes.Health)
			}
		})
	}
}

func TestScroll_NewScroll_Deterministic(t *testing.T) {
	tests := []struct {
		name       string
		seed       int64
		scrollType ScrollType
	}{
		{
			name:       "Same seed produces same Strength scroll",
			seed:       defaultItemsTestSeed,
			scrollType: ScrollTypeStrength,
		},
		{
			name:       "Same seed produces same Agility scroll",
			seed:       defaultItemsTestSeed,
			scrollType: ScrollTypeAgility,
		},
		{
			name:       "Same seed produces same MaxHealth scroll",
			seed:       defaultItemsTestSeed,
			scrollType: ScrollTypeMaxHealth,
		},
		{
			name:       "Same seed produces same Mystery scroll",
			seed:       defaultItemsTestSeed,
			scrollType: ScrollTypeMystery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

			// Create first scroll
			rng1 := utils.NewRandomGeneratorWithSeed(tt.seed)
			scroll1 := NewScroll(rng1, box, tt.scrollType)

			// Create second scroll with same seed
			rng2 := utils.NewRandomGeneratorWithSeed(tt.seed)
			scroll2 := NewScroll(rng2, box, tt.scrollType)

			// Verify they are identical
			if scroll1.Effect.Attributes.Strength != scroll2.Effect.Attributes.Strength {
				t.Errorf("Strength mismatch: %f != %f", scroll1.Effect.Attributes.Strength, scroll2.Effect.Attributes.Strength)
			}

			if scroll1.Effect.Attributes.Agility != scroll2.Effect.Attributes.Agility {
				t.Errorf("Agility mismatch: %f != %f", scroll1.Effect.Attributes.Agility, scroll2.Effect.Attributes.Agility)
			}

			if scroll1.Effect.Attributes.MaxHealth != scroll2.Effect.Attributes.MaxHealth {
				t.Errorf("MaxHealth mismatch: %f != %f", scroll1.Effect.Attributes.MaxHealth, scroll2.Effect.Attributes.MaxHealth)
			}

			if scroll1.Name != scroll2.Name {
				t.Errorf("Name mismatch: %q != %q", scroll1.Name, scroll2.Name)
			}
		})
	}
}

func TestScroll_NewScroll_Randomness(t *testing.T) {
	tests := []struct {
		name       string
		scrollType ScrollType
		iterations int
	}{
		{
			name:       "Strength scroll produces varied values",
			scrollType: ScrollTypeStrength,
			iterations: 50,
		},
		{
			name:       "Mystery scroll produces varied values",
			scrollType: ScrollTypeMystery,
			iterations: 50,
		},
		{
			name:       "MaxHealth scroll produces varied values",
			scrollType: ScrollTypeMaxHealth,
			iterations: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
			strengthValues := make(map[float64]bool)
			agilityValues := make(map[float64]bool)
			maxHealthValues := make(map[float64]bool)

			// Generate multiple scrolls with different seeds
			for i := 0; i < tt.iterations; i++ {
				rng := utils.NewRandomGeneratorWithSeed(int64(i))
				scroll := NewScroll(rng, box, tt.scrollType)
				strengthValues[scroll.Effect.Attributes.Strength] = true
				agilityValues[scroll.Effect.Attributes.Agility] = true
				maxHealthValues[scroll.Effect.Attributes.MaxHealth] = true
			}

			// Check that we got varied values (at least 5 different values)
			minUniqueValues := 5
			config := GetScrollConfig(tt.scrollType)

			// Only check non-zero ranges
			if config.StrengthRange.Min != 0 || config.StrengthRange.Max != 0 {
				if len(strengthValues) < minUniqueValues {
					t.Errorf("Expected at least %d unique strength values, got %d", minUniqueValues, len(strengthValues))
				}
			}

			if config.AgilityRange.Min != 0 || config.AgilityRange.Max != 0 {
				if len(agilityValues) < minUniqueValues {
					t.Errorf("Expected at least %d unique agility values, got %d", minUniqueValues, len(agilityValues))
				}
			}

			if config.MaxHealthRange.Min != 0 || config.MaxHealthRange.Max != 0 {
				if len(maxHealthValues) < minUniqueValues {
					t.Errorf("Expected at least %d unique MaxHealth values, got %d", minUniqueValues, len(maxHealthValues))
				}
			}
		})
	}
}

func TestScroll_NewScroll_UnknownTypeDefaultsToMystery(t *testing.T) {
	rng := utils.NewRandomGeneratorWithSeed(defaultItemsTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

	scroll := NewScroll(rng, box, "Unknown Scroll Type")

	if scroll == nil {
		t.Fatal("Expected valid scroll, got nil")
	}

	if scroll.Name != string(ScrollTypeMystery) {
		t.Errorf("Expected name %q for unknown type, got %q", ScrollTypeMystery, scroll.Name)
	}

	// Should use Mystery config ranges
	mysteryConfig := GetScrollConfig(ScrollTypeMystery)
	if scroll.Effect.Attributes.Strength < mysteryConfig.StrengthRange.Min ||
		scroll.Effect.Attributes.Strength > mysteryConfig.StrengthRange.Max {
		t.Errorf("Strength out of Mystery config range: got %.2f, want [%.2f, %.2f]",
			scroll.Effect.Attributes.Strength, mysteryConfig.StrengthRange.Min, mysteryConfig.StrengthRange.Max)
	}
}

func TestScroll_NewScroll_ZeroSizedBoxIsOk(t *testing.T) {
	rng := utils.NewRandomGeneratorWithSeed(defaultItemsTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 0, Height: 0}}
	scroll := NewScroll(rng, box, ScrollTypeStrength)

	if scroll == nil {
		t.Fatal("Expected valid scroll with zero-sized box")
	}

	if scroll.Shape != box {
		t.Errorf("Expected box %v, got %v", box, scroll.Shape)
	}
}

func TestScroll_NewScroll_VariousPositions(t *testing.T) {
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
			rng := utils.NewRandomGeneratorWithSeed(defaultItemsTestSeed)
			box := primitives.Box{Point: tt.position, Size: tt.size}
			scroll := NewScroll(rng, box, ScrollTypeAgility)

			if scroll == nil {
				t.Fatal("Expected valid scroll")
			}

			if scroll.Shape.Point != tt.position {
				t.Errorf("Expected position %v, got %v", tt.position, scroll.Shape.Point)
			}

			if scroll.Shape.Size != tt.size {
				t.Errorf("Expected size %v, got %v", tt.size, scroll.Shape.Size)
			}
		})
	}
}

func TestScroll_Use(t *testing.T) {
	tests := []struct {
		name       string
		seed       int64
		scrollType ScrollType
	}{
		{
			name:       "Use Strength scroll returns correct attributes",
			seed:       defaultItemsTestSeed,
			scrollType: ScrollTypeStrength,
		},
		{
			name:       "Use Agility scroll returns correct attributes",
			seed:       1,
			scrollType: ScrollTypeAgility,
		},
		{
			name:       "Use MaxHealth scroll returns correct attributes",
			seed:       2,
			scrollType: ScrollTypeMaxHealth,
		},
		{
			name:       "Use Mystery scroll returns correct attributes",
			seed:       3,
			scrollType: ScrollTypeMystery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := utils.NewRandomGeneratorWithSeed(tt.seed)
			box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
			scroll := NewScroll(rng, box, tt.scrollType)

			effect := scroll.Use()

			// Verify returned attributes match stored attributes
			if effect.Attributes.Strength != scroll.Effect.Attributes.Strength {
				t.Errorf("Strength mismatch: got %f, want %f", effect.Attributes.Strength, scroll.Effect.Attributes.Strength)
			}

			if effect.Attributes.Agility != scroll.Effect.Attributes.Agility {
				t.Errorf("Agility mismatch: got %f, want %f", effect.Attributes.Agility, scroll.Effect.Attributes.Agility)
			}

			if effect.Attributes.MaxHealth != scroll.Effect.Attributes.MaxHealth {
				t.Errorf("MaxHealth mismatch: got %f, want %f", effect.Attributes.MaxHealth, scroll.Effect.Attributes.MaxHealth)
			}

			// Health should always be zero
			if effect.Attributes.Health != 0 {
				t.Errorf("Expected Health=0, got %f", effect.Attributes.Health)
			}
		})
	}
}

func TestScroll_UseMultipleTimesIsOk(t *testing.T) {
	rng := utils.NewRandomGeneratorWithSeed(defaultItemsTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	scroll := NewScroll(rng, box, ScrollTypeStrength)

	effect1 := scroll.Use()
	effect2 := scroll.Use()

	if effect1.Attributes.Strength != effect2.Attributes.Strength {
		t.Error("Use() should return consistent Strength on multiple calls")
	}
	if effect1.Attributes.Agility != effect2.Attributes.Agility {
		t.Error("Use() should return consistent Agility on multiple calls")
	}
	if effect1.Attributes.MaxHealth != effect2.Attributes.MaxHealth {
		t.Error("Use() should return consistent MaxHealth on multiple calls")
	}
}

func TestScroll_AsScroll_ValidPointer(t *testing.T) {
	input := &Scroll{
		Item: &Item{
			Shape: primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			Name:  "Test Scroll",
		},
		Effect: &primitives.Effect{
			Attributes: primitives.Attributes{Strength: 10, Agility: 5, MaxHealth: 15},
		},
	}

	result := AsScroll(input)

	if result == nil {
		t.Fatal("Expected non-nil result, got nil")
	}

	if result.Effect.Attributes.Strength != 10 {
		t.Errorf("Expected Strength 10, got %f", result.Effect.Attributes.Strength)
	}
	if result.Effect.Attributes.Agility != 5 {
		t.Errorf("Expected Agility 5, got %f", result.Effect.Attributes.Agility)
	}
	if result.Effect.Attributes.MaxHealth != 15 {
		t.Errorf("Expected MaxHealth 15, got %f", result.Effect.Attributes.MaxHealth)
	}
	if result != input {
		t.Error("Expected same pointer to be returned")
	}
}

func TestScroll_AsScroll_NonScrollTypeIsNil(t *testing.T) {
	result := AsScroll("not a scroll")

	if result != nil {
		t.Errorf("Expected nil for non-Scroll type, got %v", result)
	}
}

func TestScroll_AsScroll_NilInputIsNil(t *testing.T) {
	result := AsScroll(nil)

	if result != nil {
		t.Errorf("Expected nil for nil input, got %v", result)
	}
}

func TestScroll_AsScroll_DifferentStructTypeIsNil(t *testing.T) {
	input := &Item{
		Shape: primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 1}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		Name:  "Just an item",
	}

	result := AsScroll(input)

	if result != nil {
		t.Errorf("Expected nil for Item type, got %v", result)
	}
}

func TestScroll_AsScroll_ElixirTypeIsNil(t *testing.T) {
	input := &Elixir{
		Item: &Item{
			Shape: primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			Name:  "Test Elixir",
		},
		Effect: &primitives.Effect{
			Attributes: primitives.Attributes{Strength: 10},
			Duration: primitives.EffectDuration{
				Type:  primitives.EffectDurationTypeAllTemporary,
				Steps: 20,
			},
		},
	}

	result := AsScroll(input)

	if result != nil {
		t.Errorf("Expected nil for Elixir type, got %v", result)
	}
}

func TestScroll_AsScroll_FoodTypeIsNil(t *testing.T) {
	input := &Food{
		Item: &Item{
			Shape: primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			Name:  "Test Food",
		},
		Effect: &primitives.Effect{
			Attributes: primitives.Attributes{Health: 25},
		},
	}

	result := AsScroll(input)

	if result != nil {
		t.Errorf("Expected nil for Food type, got %v", result)
	}
}

// Benchmarks
func BenchmarkScroll_NewScroll(b *testing.B) {
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}

	b.Run("Strength", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomGeneratorWithSeed(int64(i))
			_ = NewScroll(rng, box, ScrollTypeStrength)
		}
	})

	b.Run("Agility", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomGeneratorWithSeed(int64(i))
			_ = NewScroll(rng, box, ScrollTypeAgility)
		}
	})

	b.Run("MaxHealth", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomGeneratorWithSeed(int64(i))
			_ = NewScroll(rng, box, ScrollTypeMaxHealth)
		}
	})

	b.Run("Mystery", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rng := utils.NewRandomGeneratorWithSeed(int64(i))
			_ = NewScroll(rng, box, ScrollTypeMystery)
		}
	})
}

func BenchmarkScroll_NewScrollByConfig(b *testing.B) {
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	config := ScrollConfig{
		Type:           "Benchmark Scroll",
		StrengthRange:  primitives.AttributeRange{Min: 5, Max: 15},
		AgilityRange:   primitives.AttributeRange{Min: 3, Max: 12},
		MaxHealthRange: primitives.AttributeRange{Min: 10, Max: 25},
		Description:    "Benchmark test scroll",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng := utils.NewRandomGeneratorWithSeed(int64(i))
		_, _ = NewScrollByConfig(rng, box, config)
	}
}

func BenchmarkScroll_Use(b *testing.B) {
	rng := utils.NewRandomGeneratorWithSeed(defaultItemsTestSeed)
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	scroll := NewScroll(rng, box, ScrollTypeStrength)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = scroll.Use()
	}
}

func BenchmarkScroll_AsScroll(b *testing.B) {
	scroll := &Scroll{
		Item: &Item{
			Shape: primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
			Name:  "Test Scroll",
		},
		Effect: &primitives.Effect{
			Attributes: primitives.Attributes{Strength: 10, Agility: 5, MaxHealth: 15},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AsScroll(scroll)
	}
}
