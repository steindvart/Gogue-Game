package entity

import (
	"slices"
	"testing"
)

func TestGetAttributeRandomPercentIncrease(t *testing.T) {
	tests := []struct {
		name                        string
		wantIncreaseLessThan        uint
		wantIncreaseMoreOrEqualThan uint
	}{
		{
			name:                        "get attribute random percent increase",
			wantIncreaseLessThan:        uint(IncreaseAttributeBaseParcentage + IncreaseAttributeMaxPercentage + 1),
			wantIncreaseMoreOrEqualThan: uint(IncreaseAttributeBaseParcentage),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAttributeRandomPercentIncrease()
			if !(got < tt.wantIncreaseLessThan && got >= tt.wantIncreaseMoreOrEqualThan) {
				t.Errorf("getAttributeRandomPercentIncrease() = %#v, want less than %#v and more or equal than %#v", got, tt.wantIncreaseLessThan, tt.wantIncreaseMoreOrEqualThan)
			}
		})
	}
}

func TestGetAttributeRandomName(t *testing.T) {
	tests := []struct {
		name      string
		namesList []string
		wantNames []string
	}{
		{
			name: "get attribute random name",
			namesList: []string{
				"a",
				"b",
				"c",
				"d",
			},
			wantNames: []string{
				"a",
				"b",
				"c",
				"d",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAttributeRandomName(tt.namesList)
			if !(slices.Contains(tt.wantNames, got)) {
				t.Errorf("getAttributeRandomName() = %v, want in %v", got, tt.wantNames)
			}
		})
	}
}

func TestGetRandomAttribute(t *testing.T) {
	tests := []struct {
		name                       string
		wantAffectedAttributesList []Attributes
	}{
		{
			name: "get random attribute",
			wantAffectedAttributesList: []Attributes{
				{
					MaxHealth: 1,
					Agility:   0,
					Strength:  0,
				},
				{
					MaxHealth: 0,
					Agility:   1,
					Strength:  0,
				},
				{
					MaxHealth: 0,
					Agility:   0,
					Strength:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getRandomAttribute()
			if !(slices.Contains(tt.wantAffectedAttributesList, got)) {
				t.Errorf("getRandomAttribute() = %#v, want in %#v", got, tt.wantAffectedAttributesList)
			}
		})
	}
}

func TestConsumable_Taken(t *testing.T) {
	tests := []struct {
		name string
		want Box
	}{
		{
			name: "consumable is taken",
			want: Box{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumable := &Consumable{
				Shape: Box{
					Point: Point2D[int]{X: 1, Y: 2},
					Size:  Size2D[uint]{Height: 1, Width: 1},
				},
			}
			consumable.Taken()
			if consumable.Shape != tt.want {
				t.Errorf("Taken() = (%v), want (%v)", consumable.Shape, tt.want)
			}
		})
	}
}

func TestConsumable_Dropped(t *testing.T) {
	tests := []struct {
		name string
		box  Box
		want Box
	}{
		{
			name: "consumable is dropped",
			box:  Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
			want: Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumable := &Consumable{
				Shape: Box{
					Point: Point2D[int]{X: 1, Y: 2},
					Size:  Size2D[uint]{Height: 1, Width: 1},
				},
			}
			consumable.Dropped(tt.box)
			if consumable.Shape != tt.want {
				t.Errorf("Dropped() = (%v), want (%v)", consumable.Shape, tt.want)
			}
		})
	}
}

func TestConsumable_Use(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "use consumable",
			want: "Awkward Consumable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumable := &Consumable{
				Name: "Awkward Consumable",
			}
			consumable.Taken()
			if consumable.Name != tt.want {
				t.Errorf("Use() = (%v), want (%v)", consumable.Name, tt.want)
			}
		})
	}
}