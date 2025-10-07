package model

import (
	"testing"
)

func TestMenu_Next(t *testing.T) {
	tests := []struct {
		name    string
		options []option
		active  int
		want    int
	}{
		{
			name:    "next increments",
			options: []option{{"A", 1}, {"B", 2}, {"C", 3}},
			active:  0,
			want:    1,
		},
		{
			name:    "next wraps",
			options: []option{{"A", 1}, {"B", 2}},
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
			m := &Menu{Options: tt.options, Active: tt.active}
			m.Next()
			if m.Active != tt.want {
				t.Errorf("Next() = %d, want %d", m.Active, tt.want)
			}
		})
	}
}

func TestMenu_Previous(t *testing.T) {
	tests := []struct {
		name    string
		options []option
		active  int
		want    int
	}{
		{
			name:    "previous decrements",
			options: []option{{"A", 1}, {"B", 2}, {"C", 3}},
			active:  2,
			want:    1,
		},
		{
			name:    "previous wraps",
			options: []option{{"A", 1}, {"B", 2}},
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
			m := &Menu{Options: tt.options, Active: tt.active}
			m.Previous()
			if m.Active != tt.want {
				t.Errorf("Previous() = %d, want %d", m.Active, tt.want)
			}
		})
	}
}

func TestMenu_Select(t *testing.T) {
	tests := []struct {
		name    string
		options []option
		active  int
		want    Signal
	}{
		{
			name:    "select valid",
			options: []option{{"A", 42}, {"B", 99}},
			active:  1,
			want:    99,
		},
		{
			name:    "select in empty menu",
			options: nil,
			active:  0,
			want:    -1,
		},
		{
			name:    "select with one option",
			options: []option{{"A", 7}},
			active:  0,
			want:    7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Menu{Options: tt.options, Active: tt.active}
			got := m.Select()
			if got != tt.want {
				t.Errorf("Select() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMenu_NextPreviousSelect(t *testing.T) {
	m := &Menu{Options: []option{{"A", 42}, {"B", 99}}, Active: 0}

	m.Next()
	got := m.Select()
	want := m.Options[1].Signal
	if got != want {
		t.Errorf("Select() = %v, want %v", got, want)
	}

	m.Previous()
	got = m.Select()
	want = m.Options[0].Signal
	if got != want {
		t.Errorf("Select() = %v, want %v", got, want)
	}

	m.Previous()
	got = m.Select()
	want = m.Options[1].Signal
	if got != want {
		t.Errorf("Select() = %v, want %v", got, want)
	}
}
