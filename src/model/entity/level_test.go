package entity

import (
	"strings"
	"testing"
)

func TestLevel_GenerateRoomsOnLevel(t *testing.T) {
	tests := []struct {
		name          string
		sizeMap       Size2D[uint]
		wantErr       bool
		errorContains string
	}{
		{
			name:    "Valid map size 15x15",
			sizeMap: Size2D[uint]{Height: 30, Width: 90},
			wantErr: false,
		},
		{
			name:          "Map too small width",
			sizeMap:       Size2D[uint]{Height: 6, Width: 8},
			wantErr:       true,
			errorContains: "map size is too small",
		},
		{
			name:          "Map too small height",
			sizeMap:       Size2D[uint]{Height: 6, Width: 8},
			wantErr:       true,
			errorContains: "map size is too small",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := &Level{}

			err := level.GenerateNineRooms(tt.sizeMap)

			if tt.wantErr {
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
