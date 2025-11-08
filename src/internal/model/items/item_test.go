package items

// import (
// 	"gogue/internal/model/primitive"
// 	"slices"
// 	"testing"
// )

// func TestGetAttributeRandomPercentIncrease(t *testing.T) {
// 	tests := []struct {
// 		name                        string
// 		wantIncreaseLessThan        uint
// 		wantIncreaseMoreOrEqualThan uint
// 	}{
// 		{
// 			name:                        "get attribute random percent increase",
// 			wantIncreaseLessThan:        uint(IncreaseAttributeBaseParcentage + IncreaseAttributeMaxPercentage + 1),
// 			wantIncreaseMoreOrEqualThan: uint(IncreaseAttributeBaseParcentage),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := getAttributeRandomPercentIncrease()
// 			if !(got < tt.wantIncreaseLessThan && got >= tt.wantIncreaseMoreOrEqualThan) {
// 				t.Errorf("getAttributeRandomPercentIncrease() = %#v, want less than %#v and more or equal than %#v", got, tt.wantIncreaseLessThan, tt.wantIncreaseMoreOrEqualThan)
// 			}
// 		})
// 	}
// }

// func TestGetAttributeRandomName(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		namesList []string
// 		wantNames []string
// 	}{
// 		{
// 			name: "get attribute random name",
// 			namesList: []string{
// 				"a",
// 				"b",
// 				"c",
// 				"d",
// 			},
// 			wantNames: []string{
// 				"a",
// 				"b",
// 				"c",
// 				"d",
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := getAttributeRandomName(tt.namesList)
// 			if !(slices.Contains(tt.wantNames, got)) {
// 				t.Errorf("getAttributeRandomName() = %v, want in %v", got, tt.wantNames)
// 			}
// 		})
// 	}
// }

// func TestGetRandomAttribute(t *testing.T) {
// 	tests := []struct {
// 		name                       string
// 		wantAffectedAttributesList []primitive.Attributes
// 	}{
// 		{
// 			name: "get random attribute",
// 			wantAffectedAttributesList: []primitive.Attributes{
// 				{
// 					MaxHealth: 1,
// 					Agility:   0,
// 					Strength:  0,
// 				},
// 				{
// 					MaxHealth: 0,
// 					Agility:   1,
// 					Strength:  0,
// 				},
// 				{
// 					MaxHealth: 0,
// 					Agility:   0,
// 					Strength:  1,
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := getRandomAttribute()
// 			if !(slices.Contains(tt.wantAffectedAttributesList, got)) {
// 				t.Errorf("getRandomAttribute() = %#v, want in %#v", got, tt.wantAffectedAttributesList)
// 			}
// 		})
// 	}
// }

// func TestItem_Taken(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want primitive.Box
// 	}{
// 		{
// 			name: "item is taken",
// 			want: primitive.Box{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			item := &Item{
// 				Shape: primitive.Box{
// 					Point: primitive.Point2D[int]{X: 1, Y: 2},
// 					Size:  primitive.Size2D[uint]{Height: 1, Width: 1},
// 				},
// 			}
// 			items.Take()
// 			if items.Shape != tt.want {
// 				t.Errorf("Taken() = (%v), want (%v)", items.Shape, tt.want)
// 			}
// 		})
// 	}
// }

// func TestItem_Dropped(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		box  primitive.Box
// 		want primitive.Box
// 	}{
// 		{
// 			name: "item is dropped",
// 			box:  primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 1}},
// 			want: primitive.Box{Point: primitive.Point2D[int]{X: 1, Y: 2}, Size: primitive.Size2D[uint]{Height: 1, Width: 1}},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			item := &Item{
// 				Shape: primitive.Box{},
// 			}
// 			items.Drop(tt.box)
// 			if items.Shape != tt.want {
// 				t.Errorf("Dropped() = (%v), want (%v)", items.Shape, tt.want)
// 			}
// 		})
// 	}
// }

// func TestAttribute_GetAffectedAttributeName(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		attributes primitive.Attributes
// 		want       string
// 	}{
// 		{
// 			name: "get MaxHealth",
// 			attributes: primitive.Attributes{
// 				MaxHealth: 1,
// 				Agility:   0,
// 				Strength:  0,
// 			},
// 			want: "MaxHealth",
// 		},
// 		{
// 			name: "get Agility",
// 			attributes: primitive.Attributes{
// 				MaxHealth: 0,
// 				Agility:   1,
// 				Strength:  0,
// 			},
// 			want: "Agility",
// 		},
// 		{
// 			name: "get Strength",
// 			attributes: primitive.Attributes{
// 				MaxHealth: 0,
// 				Agility:   0,
// 				Strength:  1,
// 			},
// 			want: "Strength",
// 		},
// 		{
// 			name: "no affected attributes",
// 			attributes: primitive.Attributes{
// 				MaxHealth: 0,
// 				Agility:   0,
// 				Strength:  0,
// 			},
// 			want: "None",
// 		},
// 		{
// 			name: "multiple affected attributes",
// 			attributes: primitive.Attributes{
// 				MaxHealth: 1,
// 				Agility:   1,
// 				Strength:  0,
// 			},
// 			want: "None",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := tt.attributes.GetAffectedAttributeName()

// 			if got != tt.want {
// 				t.Errorf("GetAffectedAttributeName(): got %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
