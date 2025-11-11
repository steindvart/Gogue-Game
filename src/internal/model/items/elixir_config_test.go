package items

import (
	"testing"
)

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
