package cli

import (
	"gogue/view"
	"testing"

	gc "github.com/rthornton128/goncurses"
)

func TestAction(t *testing.T) {
	tests := []struct {
		name string
		key  gc.Key
		want view.ActionType
	}{
		{"KEY_UP", gc.KEY_UP, view.MoveUp},
		{"KEY_DOWN", gc.KEY_DOWN, view.MoveDown},
		{"KEY_LEFT", gc.KEY_LEFT, view.MoveLeft},
		{"KEY_RIGHT", gc.KEY_RIGHT, view.MoveRight},
		{"KEY_ENTER", gc.KEY_ENTER, view.Select},
		// буквы
		{"w", gc.Key('w'), view.MoveUp},
		{"W", gc.Key('W'), view.MoveUp},
		{"s", gc.Key('s'), view.MoveDown},
		{"a", gc.Key('a'), view.MoveLeft},
		{"d", gc.Key('d'), view.MoveRight},
		{"e", gc.Key('e'), view.Select},
		// невалидные
		{"no action", gc.Key('x'), view.NoAction},
		{"empty", gc.Key(0), view.NoAction},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Action(tt.key)
			if got != tt.want {
				t.Errorf("Action(%v) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}
