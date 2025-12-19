package items

import (
	"gogue/internal/model/primitives"
	"testing"
)

func TestItem_Drop(t *testing.T) {
	tests := []struct {
		name     string
		startPos primitives.Point2D[int]
		delta    primitives.Point2D[int]
		wantPos  primitives.Point2D[int]
	}{
		{
			name:     "Create and drop",
			startPos: primitives.Point2D[int]{X: 5, Y: 5},
			delta:    primitives.Point2D[int]{X: 0, Y: 0},
			wantPos:  primitives.Point2D[int]{X: 5, Y: 5},
		},
		{
			name:     "Create, move and drop",
			startPos: primitives.Point2D[int]{X: 5, Y: 5},
			delta:    primitives.Point2D[int]{X: 2, Y: 2},
			wantPos:  primitives.Point2D[int]{X: 7, Y: 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &Item{
				Box: &primitives.Box{
					Point: tt.startPos,
					Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
				},
			}

			item.Move(tt.delta)

			if item.GetPosition() != tt.wantPos {
				t.Errorf("Item position = %v, want %v", item.GetPosition(), tt.wantPos)
			}
		})
	}
}
