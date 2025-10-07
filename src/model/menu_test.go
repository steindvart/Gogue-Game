package model

import (
	"testing"
)

func TestMenu_Next(t *testing.T) {
	tests := []struct {
		name    string
		options []MenuOption
		active  int
		want    int
	}{
		{
			name:    "next increments",
			options: []MenuOption{{"A", 1}, {"B", 2}, {"C", 3}},
			active:  0,
			want:    1,
		},
		{
			name:    "next wraps",
			options: []MenuOption{{"A", 1}, {"B", 2}},
			active:  1,
			want:    0,
		},
		{
			name:    "empty options",
			options: nil,
			active:  0,
			want:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Menu{options: tt.options, active: tt.active}
			m.Next()
			if m.active != tt.want {
				t.Errorf("Next() = %d, want %d", m.active, tt.want)
			}
		})
	}
}

func TestMenu_Previous(t *testing.T) {
	tests := []struct {
		name    string
		options []MenuOption
		active  int
		want    int
	}{
		{
			name:    "previous decrements",
			options: []MenuOption{{"A", 1}, {"B", 2}, {"C", 3}},
			active:  2,
			want:    1,
		},
		{
			name:    "previous wraps",
			options: []MenuOption{{"A", 1}, {"B", 2}},
			active:  0,
			want:    1,
		},
		{
			name:    "empty options",
			options: nil,
			active:  0,
			want:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Menu{options: tt.options, active: tt.active}
			m.Previous()
			if m.active != tt.want {
				t.Errorf("Previous() = %d, want %d", m.active, tt.want)
			}
		})
	}
}

func TestMenu_Select(t *testing.T) {
	tests := []struct {
		name    string
		options []MenuOption
		active  int
		want    MenuOption
	}{
		{
			name:    "select valid",
			options: []MenuOption{{"A", 42}, {"B", 99}},
			active:  1,
			want:    MenuOption{"B", 99},
		},
		{
			name:    "select in empty menu",
			options: nil,
			active:  0,
			want:    MenuOption{Label: "", Signal: -1},
		},
		{
			name:    "select with one option",
			options: []MenuOption{{"A", 7}},
			active:  0,
			want:    MenuOption{"A", 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Menu{options: tt.options, active: tt.active}
			got := m.Select()
			if got != tt.want {
				t.Errorf("Select() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMenu_NextPreviousSelect(t *testing.T) {
	m := &Menu{options: []MenuOption{{"A", 42}, {"B", 99}}, active: 0}

	m.Next()
	got := m.Select()
	want := m.options[1]
	if got != want {
		t.Errorf("Select() = %v, want %v", got, want)
	}

	m.Previous()
	got = m.Select()
	want = m.options[0]
	if got != want {
		t.Errorf("Select() = %v, want %v", got, want)
	}

	m.Previous()
	got = m.Select()
	want = m.options[1]
	if got != want {
		t.Errorf("Select() = %v, want %v", got, want)
	}
}

func TestMenu_GetActive(t *testing.T) {
	tests := []struct {
		name   string
		active int
		want   int
	}{
		{
			name:   "zero active",
			active: 0,
			want:   0,
		},
		{
			name:   "positive active",
			active: 2,
			want:   2,
		},
		{
			name:   "negative active",
			active: -1,
			want:   -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Menu{active: tt.active}
			got := m.GetActive()
			if got != tt.want {
				t.Errorf("GetActive() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMenu_GetOptionsLabels(t *testing.T) {
	tests := []struct {
		name    string
		options []MenuOption
		want    []string
	}{
		{
			name:    "empty menu",
			options: nil,
			want:    []string{},
		},
		{
			name:    "single option",
			options: []MenuOption{{Label: "Start", Signal: 1}},
			want:    []string{"Start"},
		},
		{
			name:    "multiple options",
			options: []MenuOption{{Label: "Start", Signal: 1}, {Label: "Exit", Signal: 2}},
			want:    []string{"Start", "Exit"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Menu{options: tt.options}
			got := m.GetOptionsLabels()
			if len(got) != len(tt.want) {
				t.Errorf("GetOptionsLabels() length = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("GetOptionsLabels()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
