package items

import (
	"gogue/internal/model/primitive"
	"testing"
)

func TestNewTreasure(t *testing.T) {
	tests := []struct {
		name              string
		box               primitive.Box
		wantShape         primitive.Box
		wantName          string
		wantValueLessThan uint
	}{
		{
			name:              "treasure constructor",
			box:               primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 1}},
			wantShape:         primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 1}},
			wantName:          "Gold",
			wantValueLessThan: uint(TreasureBaseValue + TreasureMaxValue),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTreasure(tt.box)
			if !(got.Shape == tt.wantShape) {
				t.Errorf("NewTreasure(): got Shape %v, want %v", got.Shape, tt.wantShape)
			}
			if !(got.Value < tt.wantValueLessThan) {
				t.Errorf("NewTreasure(): got Value %v, want less than %v", got.Value, tt.wantValueLessThan)
			}
			if !(got.Name == tt.wantName) {
				t.Errorf("NewTreasure(): got Name %v, want %v", got.Name, tt.wantName)
			}
		})
	}
}
