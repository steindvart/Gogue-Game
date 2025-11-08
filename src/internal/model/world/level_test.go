package world

import (
	"gogue/internal/model/primitive"
	"strings"
	"testing"
)

func TestLevel_GenerateRoomsOnLevel(t *testing.T) {
	tests := []struct {
		name      string
		sizeMap   primitive.Size2D[uint]
		wantErr   bool
		errorText string
	}{
		{
			name:    "Valid map size 15x15",
			sizeMap: primitive.Size2D[uint]{Height: 30, Width: 90},
			wantErr: false,
		},
		{
			name:      "Map too small width",
			sizeMap:   primitive.Size2D[uint]{Height: 6, Width: 8},
			wantErr:   true,
			errorText: "map size is too small",
		},
		{
			name:      "Map too small height",
			sizeMap:   primitive.Size2D[uint]{Height: 6, Width: 8},
			wantErr:   true,
			errorText: "map size is too small",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := &Level{}

			err := level.GenerateNineRooms(tt.sizeMap)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("Expected error containing %q, but got nil", tt.errorText)
				}
				if !strings.Contains(err.Error(), tt.errorText) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorText, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
