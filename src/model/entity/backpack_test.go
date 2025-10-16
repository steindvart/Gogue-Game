package entity

import (
	"reflect"
	"testing"
)

func TestNewBackpack(t *testing.T) {
	tests := []struct {
		name string
		want *Backpack
	}{
		{
			name: "backpack constructor",
			want: &Backpack{
				Capacity:  BackpackDefaultCapacity,
				ItemsNum:  0,
				Foods:     []Food{},
				Elixirs:   []Elixir{},
				Scrolls:   []Scroll{},
				Weapon:    []Weapon{},
				Treasures: []Treasure{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBackpack()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBackpack() = %#v, want %#v", got, tt.want)
			}
		})
	}
}