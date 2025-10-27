// entity_test.go
package entity

import (
	"strings"
	"testing"
)

func TestPassage_NewPassageOnX(t *testing.T) {
	cases := []struct {
		name      string
		doorOne   Point2D[int]
		doorTwo   Point2D[int]
		wantErr   bool
		errorText string
	}{
		{
			name:    "DoorOne is earlier than doorTwo",
			doorOne: Point2D[int]{X: 2, Y: 1},
			doorTwo: Point2D[int]{X: 6, Y: 4},
			wantErr: false,
		},
		{
			name:    "DoorOne is further than doorTwo",
			doorOne: Point2D[int]{X: 6, Y: 4},
			doorTwo: Point2D[int]{X: 2, Y: 1},
			wantErr: false,
		},
		{
			name:    "Passage on one axis Y",
			doorOne: Point2D[int]{X: 1, Y: 6},
			doorTwo: Point2D[int]{X: 5, Y: 6},
			wantErr: false,
		},
		{
			name:      "Err: Passage on one axis X",
			doorOne:   Point2D[int]{X: 6, Y: 1},
			doorTwo:   Point2D[int]{X: 6, Y: 5},
			wantErr:   true,
			errorText: "NewPassageOnX: doors cannot be positioned on the same x axis",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			passage, err := NewPassageOnX(tt.doorOne, tt.doorTwo)

			if !tt.wantErr {
				if passage == nil {
					t.Errorf("Expected a valid Passage, got nil")
				}

				// Проверка, что каждая точка соединена с предыдущей с шагом 1
				for i := 1; i < len(passage.Passage); i++ {
					current := passage.Passage[i]
					previous := passage.Passage[i-1]
					dx := current.X - previous.X
					dy := current.Y - previous.Y

					if !((abs(dx) == 1 && dy == 0) || (dx == 0 && abs(dy) == 1)) {
						t.Errorf("Points are not connected properly: %d -> %d", previous, current)
					}
				}
			}
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

func TestPassage_NewPassageOnY(t *testing.T) {
	cases := []struct {
		name      string
		doorOne   Point2D[int]
		doorTwo   Point2D[int]
		wantErr   bool
		errorText string
	}{
		{
			name:    "DoorOne is earlier than doorTwo",
			doorOne: Point2D[int]{X: 2, Y: 1},
			doorTwo: Point2D[int]{X: 6, Y: 4},
			wantErr: false,
		},
		{
			name:    "DoorOne is further than doorTwo",
			doorOne: Point2D[int]{X: 6, Y: 4},
			doorTwo: Point2D[int]{X: 2, Y: 1},
			wantErr: false,
		},
		{
			name:    "Passage on one axis X",
			doorOne: Point2D[int]{X: 6, Y: 1},
			doorTwo: Point2D[int]{X: 6, Y: 5},
			wantErr: false,
		},
		{
			name:      "Err: Passage on one axis Y",
			doorOne:   Point2D[int]{X: 1, Y: 6},
			doorTwo:   Point2D[int]{X: 5, Y: 6},
			wantErr:   true,
			errorText: "NewPassageOnY: doors cannot be positioned on the same y axis",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			passage, err := NewPassageOnY(tt.doorOne, tt.doorTwo)

			if !tt.wantErr {
				if passage == nil {
					t.Errorf("Expected a valid Passage, got nil")
				}

				// Проверка, что каждая точка соединена с предыдущей с шагом 1
				for i := 1; i < len(passage.Passage); i++ {
					current := passage.Passage[i]
					previous := passage.Passage[i-1]
					dx := current.X - previous.X
					dy := current.Y - previous.Y

					if !((abs(dx) == 1 && dy == 0) || (dx == 0 && abs(dy) == 1)) {
						t.Errorf("Points are not connected properly: %d -> %d", previous, current)
					}
				}
			}
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
