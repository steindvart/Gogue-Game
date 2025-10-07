package entity

import (
	"testing"
)

func TestPoint2D_Move(t *testing.T) {
	tests := []struct {
		name  string
		start Point2D[int]
		delta Point2D[int]
		want  Point2D[int]
	}{
		{
			name:  "move positive",
			start: Point2D[int]{X: 1, Y: 2},
			delta: Point2D[int]{X: 3, Y: 4},
			want:  Point2D[int]{X: 4, Y: 6},
		},
		{
			name:  "move negative",
			start: Point2D[int]{X: 5, Y: 5},
			delta: Point2D[int]{X: -2, Y: -3},
			want:  Point2D[int]{X: 3, Y: 2},
		},
		{
			name:  "move zero",
			start: Point2D[int]{X: 1, Y: 1},
			delta: Point2D[int]{X: 0, Y: 0},
			want:  Point2D[int]{X: 1, Y: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.start
			p.Move(tt.delta)
			if p != tt.want {
				t.Errorf("Move() = %v, want %v", p, tt.want)
			}
		})
	}
}

func TestBox_Move(t *testing.T) {
	tests := []struct {
		name  string
		start Point2D[int]
		delta Point2D[int]
		want  Point2D[int]
	}{
		{
			name:  "box move positive",
			start: Point2D[int]{X: 0, Y: 0},
			delta: Point2D[int]{X: 5, Y: 7},
			want:  Point2D[int]{X: 5, Y: 7},
		},
		{
			name:  "box move negative",
			start: Point2D[int]{X: 10, Y: 10},
			delta: Point2D[int]{X: -3, Y: -4},
			want:  Point2D[int]{X: 7, Y: 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Box{Point: tt.start}
			b.Move(tt.delta)
			if b.Point != tt.want {
				t.Errorf("Box.Move() = %v, want %v", b.Point, tt.want)
			}
		})
	}
}
