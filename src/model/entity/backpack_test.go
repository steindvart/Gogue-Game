package entity

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

var elixir = Elixir{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Elixir",
	},
	EffectDuration: time.Minute,
	AffectedAttribute: Attributes{
		MaxHealth: 1,
		Agility:   0,
		Strength:  0,
	},
	Increment: 10,
}

var scroll = Scroll{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Scroll",
	},
	AffectedAttribute: Attributes{
		MaxHealth: 0,
		Agility:   1,
		Strength:  0,
	},
	Increment: 10,
}

var food = Food{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Food",
	},
	HealthRegeneration: 15,
}

var weapon = Weapon{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Weapon",
	},
	Damage: 10,
}

var elixir1 = Elixir{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Elixir 1",
	},
	EffectDuration: time.Minute,
	AffectedAttribute: Attributes{
		MaxHealth: 1,
		Agility:   0,
		Strength:  0,
	},
	Increment: 10,
}

var elixir2 = Elixir{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Elixir 2",
	},
	EffectDuration: time.Minute,
	AffectedAttribute: Attributes{
		MaxHealth: 1,
		Agility:   0,
		Strength:  0,
	},
	Increment: 10,
}

var elixir3 = Elixir{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Elixir 3",
	},
	EffectDuration: time.Minute,
	AffectedAttribute: Attributes{
		MaxHealth: 1,
		Agility:   0,
		Strength:  0,
	},
	Increment: 10,
}

var elixir4 = Elixir{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Elixir 4",
	},
	EffectDuration: time.Minute,
	AffectedAttribute: Attributes{
		MaxHealth: 1,
		Agility:   0,
		Strength:  0,
	},
	Increment: 10,
}

var elixir5 = Elixir{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Elixir 5",
	},
	EffectDuration: time.Minute,
	AffectedAttribute: Attributes{
		MaxHealth: 1,
		Agility:   0,
		Strength:  0,
	},
	Increment: 10,
}

var elixir6 = Elixir{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Elixir 6",
	},
	EffectDuration: time.Minute,
	AffectedAttribute: Attributes{
		MaxHealth: 1,
		Agility:   0,
		Strength:  0,
	},
	Increment: 10,
}

var elixir7 = Elixir{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Elixir 7",
	},
	EffectDuration: time.Minute,
	AffectedAttribute: Attributes{
		MaxHealth: 1,
		Agility:   0,
		Strength:  0,
	},
	Increment: 10,
}

var elixir8 = Elixir{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Elixir 8",
	},
	EffectDuration: time.Minute,
	AffectedAttribute: Attributes{
		MaxHealth: 1,
		Agility:   0,
		Strength:  0,
	},
	Increment: 10,
}

var elixir9 = Elixir{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Elixir 9",
	},
	EffectDuration: time.Minute,
	AffectedAttribute: Attributes{
		MaxHealth: 1,
		Agility:   0,
		Strength:  0,
	},
	Increment: 10,
}

var scroll1 = Scroll{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Scroll 1",
	},
	AffectedAttribute: Attributes{
		MaxHealth: 0,
		Agility:   1,
		Strength:  0,
	},
	Increment: 10,
}

var scroll2 = Scroll{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Scroll 2",
	},
	AffectedAttribute: Attributes{
		MaxHealth: 0,
		Agility:   1,
		Strength:  0,
	},
	Increment: 10,
}

var scroll3 = Scroll{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Scroll 3",
	},
	AffectedAttribute: Attributes{
		MaxHealth: 0,
		Agility:   1,
		Strength:  0,
	},
	Increment: 10,
}

var food1 = Food{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Food 1",
	},
	HealthRegeneration: 15,
}

var food2 = Food{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Food 2",
	},
	HealthRegeneration: 15,
}

var food3 = Food{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Food 3",
	},
	HealthRegeneration: 15,
}

var weapon1 = Weapon{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Weapon 1",
	},
	Damage: 10,
}

var weapon2 = Weapon{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Weapon 2",
	},
	Damage: 10,
}

var weapon3 = Weapon{
	Item: Item{
		Shape: Box{},
		Name:  "Awkward Weapon 3",
	},
	Damage: 10,
}

func TestBackpack_AddItem(t *testing.T) {
	tests := []struct {
		name       string
		items      []ItemLike
		want       Backpack
		wantErrors []error
	}{
		{
			name:  "add elixir",
			items: []ItemLike{&elixir},
			want: Backpack{
				Capacity: BackpackDefaultCapacity,
				ItemsNum: 1,
				Elixirs: map[string][]Elixir{
					"Awkward Elixir": {elixir},
				},
				Scrolls:   map[string][]Scroll{},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			wantErrors: []error{nil},
		},
		{
			name:  "add scroll",
			items: []ItemLike{&scroll},
			want: Backpack{
				Capacity: BackpackDefaultCapacity,
				ItemsNum: 1,
				Elixirs:  map[string][]Elixir{},
				Scrolls: map[string][]Scroll{
					"Awkward Scroll": {scroll},
				},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			wantErrors: []error{nil},
		},
		{
			name:  "add food",
			items: []ItemLike{&food},
			want: Backpack{
				Capacity: BackpackDefaultCapacity,
				ItemsNum: 1,
				Elixirs:  map[string][]Elixir{},
				Scrolls:  map[string][]Scroll{},
				Foods: map[string][]Food{
					"Awkward Food": {food},
				},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			wantErrors: []error{nil},
		},
		{
			name:  "add weapon",
			items: []ItemLike{&weapon},
			want: Backpack{
				Capacity: BackpackDefaultCapacity,
				ItemsNum: 1,
				Elixirs:  map[string][]Elixir{},
				Scrolls:  map[string][]Scroll{},
				Foods:    map[string][]Food{},
				Weapons: map[string][]Weapon{
					"Awkward Weapon": {weapon},
				},
				Treasures: 0,
			},
			wantErrors: []error{nil},
		},
		{
			name:  "add elixir, scroll, food, weapon",
			items: []ItemLike{&elixir, &scroll, &food, &weapon},
			want: Backpack{
				Capacity: BackpackDefaultCapacity,
				ItemsNum: 4,
				Elixirs: map[string][]Elixir{
					"Awkward Elixir": {elixir},
				},
				Scrolls: map[string][]Scroll{
					"Awkward Scroll": {scroll},
				},
				Foods: map[string][]Food{
					"Awkward Food": {food},
				},
				Weapons: map[string][]Weapon{
					"Awkward Weapon": {weapon},
				},
				Treasures: 0,
			},
			wantErrors: []error{nil, nil, nil, nil},
		},
		{
			name:  "add nine same scrolls",
			items: []ItemLike{&scroll, &scroll, &scroll, &scroll, &scroll, &scroll, &scroll, &scroll, &scroll},
			want: Backpack{
				Capacity: BackpackDefaultCapacity,
				ItemsNum: 9,
				Elixirs:  map[string][]Elixir{},
				Scrolls: map[string][]Scroll{
					"Awkward Scroll": {scroll, scroll, scroll, scroll, scroll, scroll, scroll, scroll, scroll},
				},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			wantErrors: []error{nil, nil, nil, nil, nil, nil, nil, nil, nil},
		},
		{
			name:  "add nine different elixirs",
			items: []ItemLike{&elixir1, &elixir2, &elixir3, &elixir4, &elixir5, &elixir6, &elixir7, &elixir8, &elixir9},
			want: Backpack{
				Capacity: BackpackDefaultCapacity,
				ItemsNum: 9,
				Elixirs: map[string][]Elixir{
					"Awkward Elixir 1": {elixir1},
					"Awkward Elixir 2": {elixir2},
					"Awkward Elixir 3": {elixir3},
					"Awkward Elixir 4": {elixir4},
					"Awkward Elixir 5": {elixir5},
					"Awkward Elixir 6": {elixir6},
					"Awkward Elixir 7": {elixir7},
					"Awkward Elixir 8": {elixir8},
					"Awkward Elixir 9": {elixir9},
				},
				Scrolls:   map[string][]Scroll{},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			wantErrors: []error{nil, nil, nil, nil, nil, nil, nil, nil, nil},
		},
		{
			name:  "add nine same scrolls and food",
			items: []ItemLike{&scroll, &scroll, &scroll, &scroll, &scroll, &scroll, &scroll, &scroll, &scroll, &food},
			want: Backpack{
				Capacity: BackpackDefaultCapacity,
				ItemsNum: 9,
				Elixirs:  map[string][]Elixir{},
				Scrolls: map[string][]Scroll{
					"Awkward Scroll": {scroll, scroll, scroll, scroll, scroll, scroll, scroll, scroll, scroll},
				},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			wantErrors: []error{nil, nil, nil, nil, nil, nil, nil, nil, nil, BackpackIsFullError{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBackpack()

			for idx := range tt.items {
				err := b.AddItem(tt.items[idx])
				errWant := tt.wantErrors[idx]

				if !errors.Is(err, errWant) {
					t.Errorf("AddItem() = %#v, want %#v", err, errWant)
				}
			}

			if !reflect.DeepEqual(*b, tt.want) {
				t.Errorf("AddItem(): Backpack got %#v, want %#v", *b, tt.want)
			}
		})
	}
}

func TestBackPack_RemoveItem(t *testing.T) {
	tests := []struct {
		name         string
		items        []ItemLike
		item         ItemLike
		wantBackpack Backpack
		want         error
	}{
		{
			name:  "remove elixir from empty backpack",
			items: []ItemLike{},
			item:  &elixir,
			wantBackpack: Backpack{
				Capacity:  BackpackDefaultCapacity,
				ItemsNum:  0,
				Elixirs:   map[string][]Elixir{},
				Scrolls:   map[string][]Scroll{},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			want: ItemIsNotInBackpackError{},
		},
		{
			name:  "remove scroll from empty backpack",
			items: []ItemLike{},
			item:  &scroll,
			wantBackpack: Backpack{
				Capacity:  BackpackDefaultCapacity,
				ItemsNum:  0,
				Elixirs:   map[string][]Elixir{},
				Scrolls:   map[string][]Scroll{},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			want: ItemIsNotInBackpackError{},
		},
		{
			name:  "remove food from empty backpack",
			items: []ItemLike{},
			item:  &food,
			wantBackpack: Backpack{
				Capacity:  BackpackDefaultCapacity,
				ItemsNum:  0,
				Elixirs:   map[string][]Elixir{},
				Scrolls:   map[string][]Scroll{},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			want: ItemIsNotInBackpackError{},
		},
		{
			name:  "remove weapon from empty backpack",
			items: []ItemLike{},
			item:  &weapon,
			wantBackpack: Backpack{
				Capacity:  BackpackDefaultCapacity,
				ItemsNum:  0,
				Elixirs:   map[string][]Elixir{},
				Scrolls:   map[string][]Scroll{},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			want: ItemIsNotInBackpackError{},
		},
		{
			name:  "remove elixir",
			items: []ItemLike{&elixir, &elixir, &elixir},
			item:  &elixir,
			wantBackpack: Backpack{
				Capacity: BackpackDefaultCapacity,
				ItemsNum: 2,
				Elixirs: map[string][]Elixir{
					"Awkward Elixir": {elixir, elixir},
				},
				Scrolls:   map[string][]Scroll{},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			want: nil,
		},
		{
			name:  "remove scroll",
			items: []ItemLike{&scroll, &scroll, &scroll},
			item:  &scroll,
			wantBackpack: Backpack{
				Capacity: BackpackDefaultCapacity,
				ItemsNum: 2,
				Elixirs:  map[string][]Elixir{},
				Scrolls: map[string][]Scroll{
					"Awkward Scroll": {scroll, scroll},
				},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			want: nil,
		},
		{
			name:  "remove food",
			items: []ItemLike{&food, &food, &food},
			item:  &food,
			wantBackpack: Backpack{
				Capacity: BackpackDefaultCapacity,
				ItemsNum: 2,
				Elixirs:  map[string][]Elixir{},
				Scrolls:  map[string][]Scroll{},
				Foods: map[string][]Food{
					"Awkward Food": {food, food},
				},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			want: nil,
		},
		{
			name:  "remove weapon",
			items: []ItemLike{&weapon, &weapon, &weapon},
			item:  &weapon,
			wantBackpack: Backpack{
				Capacity: BackpackDefaultCapacity,
				ItemsNum: 2,
				Elixirs:  map[string][]Elixir{},
				Scrolls:  map[string][]Scroll{},
				Foods:    map[string][]Food{},
				Weapons: map[string][]Weapon{
					"Awkward Weapon": {weapon, weapon},
				},
				Treasures: 0,
			},
			want: nil,
		},
		{
			name:  "remove last elixir",
			items: []ItemLike{&elixir},
			item:  &elixir,
			wantBackpack: Backpack{
				Capacity:  BackpackDefaultCapacity,
				ItemsNum:  0,
				Elixirs:   map[string][]Elixir{},
				Scrolls:   map[string][]Scroll{},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			want: nil,
		},
		{
			name:  "remove last scroll",
			items: []ItemLike{&scroll},
			item:  &scroll,
			wantBackpack: Backpack{
				Capacity:  BackpackDefaultCapacity,
				ItemsNum:  0,
				Elixirs:   map[string][]Elixir{},
				Scrolls:   map[string][]Scroll{},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			want: nil,
		},
		{
			name:  "remove last food",
			items: []ItemLike{&food},
			item:  &food,
			wantBackpack: Backpack{
				Capacity:  BackpackDefaultCapacity,
				ItemsNum:  0,
				Elixirs:   map[string][]Elixir{},
				Scrolls:   map[string][]Scroll{},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			want: nil,
		},
		{
			name:  "remove last weapon",
			items: []ItemLike{&weapon},
			item:  &weapon,
			wantBackpack: Backpack{
				Capacity:  BackpackDefaultCapacity,
				ItemsNum:  0,
				Elixirs:   map[string][]Elixir{},
				Scrolls:   map[string][]Scroll{},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBackpack()

			for idx := range tt.items {
				_ = b.AddItem(tt.items[idx])
			}

			got := b.RemoveItem(tt.item, Box{})

			if !errors.Is(got, tt.want) {
				t.Errorf("RemoveItem() = %#v, want %#v", got, tt.want)
			}

			if !reflect.DeepEqual(*b, tt.wantBackpack) {
				t.Errorf("RemoveItem(): Backpack got %#v, want %#v", *b, tt.wantBackpack)
			}
		})
	}
}

// func TestBackpack_GetElixirsList(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		items []ItemLike
// 		want  []listElement
// 	}{
// 		{
// 			name:  "zero elixirs",
// 			items: []ItemLike{},
// 			want:  []listElement{},
// 		},
// 		{
// 			name:  "multiple elixirs of same type",
// 			items: []ItemLike{&elixir, &elixir, &elixir},
// 			want: []listElement{
// 				{
// 					Name: "Awkward Elixir",
// 					Num:  3,
// 					Item: &elixir,
// 				},
// 			},
// 		},
// 		{
// 			name:  "multiple elixirs of different types",
// 			items: []ItemLike{&elixir1, &elixir2, &elixir2, &elixir3},
// 			want: []listElement{
// 				{
// 					Name: "Awkward Elixir 1",
// 					Num:  1,
// 					Item: &elixir1,
// 				},
// 				{
// 					Name: "Awkward Elixir 2",
// 					Num:  2,
// 					Item: &elixir2,
// 				},
// 				{
// 					Name: "Awkward Elixir 3",
// 					Num:  1,
// 					Item: &elixir3,
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b := NewBackpack()

// 			for idx := range tt.items {
// 				_ = b.AddItem(tt.items[idx])
// 			}

// 			got := b.GetElixirsList()

// 			for idx := range tt.want {
// 				e, ok
// 				if !slices.Contains(tt.want, got[idx]) {
// 					if !slices.Equal(got, tt.want) {
// 						t.Errorf("GetItemsList() = %#v, want in %#v", got[idx], tt.want)
// 					}
// 				}
// 			}
// 		})
// 	}
// }

// func TestBackpack_GetScrollsList(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		items []ItemLike
// 		want  []listElement
// 	}{
// 		{
// 			name:  "zero scrolls",
// 			items: []ItemLike{},
// 			want:  []listElement{},
// 		},
// 		{
// 			name:  "multiple scrolls of same type",
// 			items: []ItemLike{&scroll, &scroll, &scroll},
// 			want: []listElement{
// 				{
// 					Name: "Awkward Scroll",
// 					Num:  3,
// 					Item: &scroll,
// 				},
// 			},
// 		},
// 		{
// 			name:  "multiple scrolls of different types",
// 			items: []ItemLike{&scroll1, &scroll2, &scroll2, &scroll3},
// 			want: []listElement{
// 				{
// 					Name: "Awkward Scroll 1",
// 					Num:  1,
// 					Item: &scroll1,
// 				},
// 				{
// 					Name: "Awkward Scroll 2",
// 					Num:  2,
// 					Item: &scroll2,
// 				},
// 				{
// 					Name: "Awkward Scroll 3",
// 					Num:  1,
// 					Item: &scroll3,
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b := NewBackpack()

// 			for idx := range tt.items {
// 				_ = b.AddItem(tt.items[idx])
// 			}

// 			got := b.GetScrollsList()

// 			for idx := range got {
// 				if !slices.Contains(tt.want, got[idx]) {
// 					if !slices.Equal(got, tt.want) {
// 						t.Errorf("GetItemsList() = %#v, want in %#v", got[idx], tt.want)
// 					}
// 				}
// 			}
// 		})
// 	}
// }

// func TestBackpack_GetFoodsList(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		items []ItemLike
// 		want  []listElement
// 	}{
// 		{
// 			name:  "zero foods",
// 			items: []ItemLike{},
// 			want:  []listElement{},
// 		},
// 		{
// 			name:  "multiple foods of same type",
// 			items: []ItemLike{&food, &food, &food},
// 			want: []listElement{
// 				{
// 					Name: "Awkward Food",
// 					Num:  3,
// 					Item: &food,
// 				},
// 			},
// 		},
// 		{
// 			name:  "multiple food of different types",
// 			items: []ItemLike{&food1, &food2, &food2, &food3},
// 			want: []listElement{
// 				{
// 					Name: "Awkward Food 1",
// 					Num:  1,
// 					Item: &food1,
// 				},
// 				{
// 					Name: "Awkward Food 2",
// 					Num:  2,
// 					Item: &food2,
// 				},
// 				{
// 					Name: "Awkward Food 3",
// 					Num:  1,
// 					Item: &food3,
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b := NewBackpack()

// 			for idx := range tt.items {
// 				_ = b.AddItem(tt.items[idx])
// 			}

// 			got := b.GetFoodsList()

// 			for idx := range got {
// 				if !slices.Contains(tt.want, got[idx]) {
// 					if !slices.Equal(got, tt.want) {
// 						t.Errorf("GetItemsList() = %#v, want in %#v", got[idx], tt.want)
// 					}
// 				}
// 			}
// 		})
// 	}
// }

// func TestBackpack_GetWeaponsList(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		items []ItemLike
// 		want  []listElement
// 	}{
// 		{
// 			name:  "zero weapons",
// 			items: []ItemLike{},
// 			want:  []listElement{},
// 		},
// 		{
// 			name:  "multiple weapons of same type",
// 			items: []ItemLike{&weapon, &weapon, &weapon},
// 			want: []listElement{
// 				{
// 					Name: "Awkward Weapon",
// 					Num:  3,
// 					Item: &weapon,
// 				},
// 			},
// 		},
// 		{
// 			name:  "multiple weapons of different types",
// 			items: []ItemLike{&weapon1, &weapon2, &weapon2, &weapon3},
// 			want: []listElement{
// 				{
// 					Name: "Awkward Weapon 1",
// 					Num:  1,
// 					Item: &weapon1,
// 				},
// 				{
// 					Name: "Awkward Weapon 2",
// 					Num:  2,
// 					Item: &weapon2,
// 				},
// 				{
// 					Name: "Awkward Weapon 3",
// 					Num:  1,
// 					Item: &weapon3,
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b := NewBackpack()

// 			for idx := range tt.items {
// 				_ = b.AddItem(tt.items[idx])
// 			}

// 			got := b.GetWeaponsList()

// 			for idx := range got {
// 				if !slices.Contains(tt.want, got[idx]) {
// 					if !slices.Equal(got, tt.want) {
// 						t.Errorf("GetItemsList() = %#v, want in %#v", got[idx], tt.want)
// 					}
// 				}
// 			}
// 		})
// 	}
// }

// func TestBackpack_GetItemsList(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		items []ItemLike
// 		want  []listElement
// 	}{
// 		{
// 			name:  "elixir, scroll, food, weapon",
// 			items: []ItemLike{&elixir, &scroll, &food, &weapon},
// 			want: []listElement{
// 				{
// 					Name: "Awkward Elixir",
// 					Num:  1,
// 					Item: &elixir,
// 				},
// 				{
// 					Name: "Awkward Scroll",
// 					Num:  1,
// 					Item: &scroll,
// 				},
// 				{
// 					Name: "Awkward Food",
// 					Num:  1,
// 					Item: &food,
// 				},
// 				{
// 					Name: "Awkward Weapon",
// 					Num:  1,
// 					Item: &weapon,
// 				},
// 			},
// 		},
// 		{
// 			name:  "multiple elixirs, scrolls, foods, weapons",
// 			items: []ItemLike{&elixir1, &elixir2, &scroll1, &scroll1, &food1, &food2, &weapon1, &weapon3},
// 			want: []listElement{
// 				{
// 					Name: "Awkward Elixir 1",
// 					Num:  1,
// 					Item: &elixir1,
// 				},
// 				{
// 					Name: "Awkward Elixir 2",
// 					Num:  1,
// 					Item: &elixir2,
// 				},
// 				{
// 					Name: "Awkward Scroll 1",
// 					Num:  2,
// 					Item: &scroll1,
// 				},
// 				{
// 					Name: "Awkward Food 1",
// 					Num:  1,
// 					Item: &food1,
// 				},
// 				{
// 					Name: "Awkward Food 2",
// 					Num:  1,
// 					Item: &food2,
// 				},
// 				{
// 					Name: "Awkward Weapon 1",
// 					Num:  1,
// 					Item: &weapon1,
// 				},
// 				{
// 					Name: "Awkward Weapon 3",
// 					Num:  1,
// 					Item: &weapon3,
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b := NewBackpack()

// 			for idx := range tt.items {
// 				_ = b.AddItem(tt.items[idx])
// 			}

// 			got := b.GetItemsList()

// 			for idx := range got {
// 				if !slices.Contains(tt.want, got[idx]) {
// 					if !slices.Equal(got, tt.want) {
// 						t.Errorf("GetItemsList() = %#v, want in %#v", got[idx], tt.want)
// 					}
// 				}
// 			}
// 		})
// 	}
// }

// func TestNewBackpack(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want *Backpack
// 	}{
// 		{
// 			name: "backpack constructor",
// 			want: &Backpack{
// 				Capacity:  BackpackDefaultCapacity,
// 				ItemsNum:  0,
// 				Elixirs:   map[string][]Elixir{},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := NewBackpack()
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewBackpack() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }
