package cli

import (
	"gogue/internal/model/signals"
	"testing"

	"github.com/rivo/tview"
)

type mockState struct {
	updateCalled bool
}

func (m *mockState) Update(dt float64) signals.Type {
	m.updateCalled = true
	return signals.NoSignal
}

func (m *mockState) Primitive() tview.Primitive {
	return tview.NewBox()
}

func TestStateMachine_PushPopState(t *testing.T) {
	g := &StateMachine{}
	state1 := &mockState{}
	state2 := &mockState{}

	g.PushState(state1)
	g.PushState(state2)

	if len(g.States) != 2 {
		t.Errorf("PushState: expected 2 states, got %d", len(g.States))
	}

	g.PopState()
	if len(g.States) != 1 {
		t.Errorf("PopState: expected 1 state, got %d", len(g.States))
	}

	g.PopState()
	if len(g.States) != 0 {
		t.Errorf("PopState: expected 0 states, got %d", len(g.States))
	}

	g.PopState() // Попытка удалить из пустого
	if len(g.States) != 0 {
		t.Errorf("PopState: expected 0 states after pop from empty, got %d", len(g.States))
	}
}

func TestStateMachine_CurrentState(t *testing.T) {
	g := &StateMachine{}
	if g.CurrentState() != nil {
		t.Errorf("CurrentState: expected nil for empty stack")
	}
	state := &mockState{}
	g.PushState(state)
	if g.CurrentState() != state {
		t.Errorf("CurrentState: expected pushed state")
	}
	g.PopState()
	if g.CurrentState() != nil {
		t.Errorf("CurrentState: expected nil after pop")
	}
}
