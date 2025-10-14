package entity

import (
	"strings"
	"testing"
)

func TestLevel_GenerateRoomsOnLevel(t *testing.T) {
	tests := []struct {
		name          string
		mapWidth      int
		mapHeight     int
		want          bool
		errorContains string
	}{
		{
			name:      "Valid map size 15x15",
			mapWidth:  15,
			mapHeight: 15,
			want:      false,
		},
		{
			name:      "Minimal valid size 9x9",
			mapWidth:  9,
			mapHeight: 9,
			want:      false,
		},
		{
			name:          "Map too small width",
			mapWidth:      8,
			mapHeight:     6,
			want:          true,
			errorContains: "map size is too small",
		},
		{
			name:          "Map too small height",
			mapWidth:      6,
			mapHeight:     8,
			want:          true,
			errorContains: "map size is too small",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := &Level{}

			err := level.GenerateRoomsOnLevel(tt.mapWidth, tt.mapHeight)

			if tt.want {
				if err == nil {
					t.Fatalf("Expected error containing %q, but got nil", tt.errorContains)
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
