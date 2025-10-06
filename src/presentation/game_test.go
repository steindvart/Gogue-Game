package presentation

import (
	"testing"
)

type mockState struct {
	inputCalled  bool
	updateCalled bool
	renderCalled bool
}

func (m *mockState) Input() { m.inputCalled = true }

func (m *mockState) Update() GameStateSignal {
	m.updateCalled = true
	return StopStateSignal
}

func (m *mockState) Render() { m.renderCalled = true }

func TestGame_PushPopState(t *testing.T) {
	g := &Game{}
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

func TestGame_CurrentState(t *testing.T) {
	g := &Game{}
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

func TestGame_Run(t *testing.T) {
	g := &Game{}
	state := &mockState{}
	g.PushState(state)
	g.Run()
	if !state.inputCalled || !state.updateCalled || !state.renderCalled {
		t.Errorf("Run: expected all methods to be called on state")
	}
}
