package cli

import (
	"gogue/view/action"
	"testing"

	gc "github.com/rthornton128/goncurses"
)

func TestAction(t *testing.T) {
	tests := []struct {
		name string
		key  gc.Key
		want action.Type
	}{
		{"KEY_UP", gc.KEY_UP, action.MoveUp},
		{"KEY_DOWN", gc.KEY_DOWN, action.MoveDown},
		{"KEY_LEFT", gc.KEY_LEFT, action.MoveLeft},
		{"KEY_RIGHT", gc.KEY_RIGHT, action.MoveRight},
		{"KEY_ENTER", gc.KEY_ENTER, action.Select},
		{"KEY_RETURN", gc.KEY_RETURN, action.Select},
		// буквы
		{"w", gc.Key('w'), action.MoveUp},
		{"W", gc.Key('W'), action.MoveUp},
		{"s", gc.Key('s'), action.MoveDown},
		{"a", gc.Key('a'), action.MoveLeft},
		{"d", gc.Key('d'), action.MoveRight},
		{"e", gc.Key('e'), action.Select},
		// невалидные
		{"no action", gc.Key('x'), action.NoAction},
		{"empty", gc.Key(0), action.NoAction},
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
