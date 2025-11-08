package entities

// import (
// 	"errors"
// 	"gogue/internal/model/entity/primitive"
// 	"reflect"
// 	"testing"
// 	"time"
// )

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

// func TestIsElixir(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		item ItemLike
// 		want *Elixir
// 	}{
// 		{
// 			name: "elixir",
// 			item: &Elixir{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Elixir",
// 				},
// 				EffectDuration: time.Minute,
// 				AffectedAttribute: primitives.Attributes{
// 					MaxHealth: 1,
// 					Agility:   0,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			want: &Elixir{},
// 		},
// 		{
// 			name: "not elixir",
// 			item: &Scroll{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Scroll",
// 				},
// 				AffectedAttribute: primitives.Attributes{
// 					MaxHealth: 0,
// 					Agility:   1,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			want: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.want != nil {
// 				ptr, _ := tt.items.(*Elixir)
// 				tt.want = ptr
// 			}

// 			got := IsElixir(tt.item)

// 			if got != tt.want {
// 				t.Errorf("IsElixir() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestIsScroll(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		item ItemLike
// 		want *Scroll
// 	}{
// 		{
// 			name: "scroll",
// 			item: &Scroll{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Scroll",
// 				},
// 				AffectedAttribute: primitives.Attributes{
// 					MaxHealth: 0,
// 					Agility:   1,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			want: &Scroll{},
// 		},
// 		{
// 			name: "not scroll",
// 			item: &Elixir{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Elixir",
// 				},
// 				EffectDuration: time.Minute,
// 				AffectedAttribute: primitives.Attributes{
// 					MaxHealth: 1,
// 					Agility:   0,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			want: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.want != nil {
// 				ptr, _ := tt.items.(*Scroll)
// 				tt.want = ptr
// 			}

// 			got := IsScroll(tt.item)

// 			if got != tt.want {
// 				t.Errorf("IsScroll() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestIsFood(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		item ItemLike
// 		want *Food
// 	}{
// 		{
// 			name: "food",
// 			item: &Food{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Food",
// 				},
// 				HealthRegeneration: 15,
// 			},
// 			want: &Food{},
// 		},
// 		{
// 			name: "not scroll",
// 			item: &Elixir{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Elixir",
// 				},
// 				EffectDuration: time.Minute,
// 				AffectedAttribute: primitives.Attributes{
// 					MaxHealth: 1,
// 					Agility:   0,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			want: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.want != nil {
// 				ptr, _ := tt.items.(*Food)
// 				tt.want = ptr
// 			}

// 			got := IsFood(tt.item)

// 			if got != tt.want {
// 				t.Errorf("IsFood() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestIsWeapon(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		item ItemLike
// 		want *Weapon
// 	}{
// 		{
// 			name: "weapon",
// 			item: &Weapon{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Weapon",
// 				},
// 				Damage: 10,
// 			},
// 			want: &Weapon{},
// 		},
// 		{
// 			name: "not weapon",
// 			item: &Elixir{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Elixir",
// 				},
// 				EffectDuration: time.Minute,
// 				AffectedAttribute: primitives.Attributes{
// 					MaxHealth: 1,
// 					Agility:   0,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			want: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.want != nil {
// 				ptr, _ := tt.items.(*Weapon)
// 				tt.want = ptr
// 			}

// 			got := IsWeapon(tt.item)

// 			if got != tt.want {
// 				t.Errorf("IsWeapon() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestBackpack_AddTreasure(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		treasure *Treasure
// 		want     Backpack
// 	}{
// 		{
// 			name: "treasure value 10",
// 			treasure: &Treasure{
// 				Shape: primitives.Box{},
// 				Name:  "Awkward Treasure",
// 				Value: 10,
// 			},
// 			want: Backpack{
// 				Capacity:  BackpackDefaultCapacity,
// 				ItemsNum:  0,
// 				Elixirs:   map[string][]Elixir{},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 10,
// 			},
// 		},
// 		{
// 			name: "treasure value 0",
// 			treasure: &Treasure{
// 				Shape: primitives.Box{},
// 				Name:  "Awkward Treasure",
// 				Value: 0,
// 			},
// 			want: Backpack{
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
// 			b := NewBackpack()

// 			b.AddTreasure(tt.treasure)

// 			if !reflect.DeepEqual(*b, tt.want) {
// 				t.Errorf("AddTreasure(): Backpack got %#v, want %#v", *b, tt.want)
// 			}
// 		})
// 	}
// }

// func TestBackpack_AddItem(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		items      []ItemLike
// 		want       Backpack
// 		wantErrors []error
// 	}{
// 		{
// 			name: "add elixir",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			want: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 1,
// 				Elixirs: map[string][]Elixir{
// 					"Awkward Elixir": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 				},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			wantErrors: []error{nil},
// 		},
// 		{
// 			name: "add scroll",
// 			items: []ItemLike{
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			want: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 1,
// 				Elixirs:  map[string][]Elixir{},
// 				Scrolls: map[string][]Scroll{
// 					"Awkward Scroll": {Scroll{
// 						Item: Item{
// 							Shape: primitives.Box{},
// 							Name:  "Awkward Scroll",
// 						},
// 						AffectedAttribute: primitives.Attributes{
// 							MaxHealth: 0,
// 							Agility:   1,
// 							Strength:  0,
// 						},
// 						Increment: 10,
// 					}},
// 				},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			wantErrors: []error{nil},
// 		},
// 		{
// 			name: "add food",
// 			items: []ItemLike{
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 			},
// 			want: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 1,
// 				Elixirs:  map[string][]Elixir{},
// 				Scrolls:  map[string][]Scroll{},
// 				Foods: map[string][]Food{
// 					"Awkward Food": {
// 						Food{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Food",
// 							},
// 							HealthRegeneration: 15,
// 						},
// 					},
// 				},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			wantErrors: []error{nil},
// 		},
// 		{
// 			name: "add weapon",
// 			items: []ItemLike{
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon",
// 					},
// 					Damage: 10,
// 				},
// 			},
// 			want: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 1,
// 				Elixirs:  map[string][]Elixir{},
// 				Scrolls:  map[string][]Scroll{},
// 				Foods:    map[string][]Food{},
// 				Weapons: map[string][]Weapon{
// 					"Awkward Weapon": {
// 						Weapon{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Weapon",
// 							},
// 							Damage: 10,
// 						},
// 					},
// 				},
// 				Treasures: 0,
// 			},
// 			wantErrors: []error{nil},
// 		},
// 		{
// 			name: "add elixir, scroll, food, weapon",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon",
// 					},
// 					Damage: 10,
// 				},
// 			},
// 			want: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 4,
// 				Elixirs: map[string][]Elixir{
// 					"Awkward Elixir": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 				},
// 				Scrolls: map[string][]Scroll{
// 					"Awkward Scroll": {Scroll{
// 						Item: Item{
// 							Shape: primitives.Box{},
// 							Name:  "Awkward Scroll",
// 						},
// 						AffectedAttribute: primitives.Attributes{
// 							MaxHealth: 0,
// 							Agility:   1,
// 							Strength:  0,
// 						},
// 						Increment: 10,
// 					}},
// 				},
// 				Foods: map[string][]Food{
// 					"Awkward Food": {
// 						Food{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Food",
// 							},
// 							HealthRegeneration: 15,
// 						},
// 					},
// 				},
// 				Weapons: map[string][]Weapon{
// 					"Awkward Weapon": {
// 						Weapon{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Weapon",
// 							},
// 							Damage: 10,
// 						},
// 					},
// 				},
// 				Treasures: 0,
// 			},
// 			wantErrors: []error{nil, nil, nil, nil},
// 		},
// 		{
// 			name: "add nine same scrolls",
// 			items: []ItemLike{
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			want: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 9,
// 				Elixirs:  map[string][]Elixir{},
// 				Scrolls: map[string][]Scroll{
// 					"Awkward Scroll": {
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 				},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			wantErrors: []error{nil, nil, nil, nil, nil, nil, nil, nil, nil},
// 		},
// 		{
// 			name: "add nine different elixirs",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 1",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 2",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 3",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 4",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 5",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 6",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 7",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 8",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 9",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			want: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 9,
// 				Elixirs: map[string][]Elixir{
// 					"Awkward Elixir 1": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Elixir 1",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 					"Awkward Elixir 2": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Elixir 2",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 					"Awkward Elixir 3": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Elixir 3",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 					"Awkward Elixir 4": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Elixir 4",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 					"Awkward Elixir 5": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Elixir 5",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 					"Awkward Elixir 6": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Elixir 6",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 					"Awkward Elixir 7": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Elixir 7",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 					"Awkward Elixir 8": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Elixir 8",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 					"Awkward Elixir 9": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Elixir 9",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 				},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			wantErrors: []error{nil, nil, nil, nil, nil, nil, nil, nil, nil},
// 		},
// 		{
// 			name: "add nine same scrolls and food",
// 			items: []ItemLike{
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			want: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 9,
// 				Elixirs:  map[string][]Elixir{},
// 				Scrolls: map[string][]Scroll{
// 					"Awkward Scroll": {
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 				},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			wantErrors: []error{nil, nil, nil, nil, nil, nil, nil, nil, nil, BackpackIsFullError{}},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b := NewBackpack()

// 			for idx := range tt.items {
// 				err := b.AddItem(tt.items[idx])
// 				errWant := tt.wantErrors[idx]

// 				if !errors.Is(err, errWant) {
// 					t.Errorf("AddItem() = %#v, want %#v", err, errWant)
// 				}
// 			}

// 			if !reflect.DeepEqual(*b, tt.want) {
// 				t.Errorf("AddItem(): Backpack got %#v, want %#v", *b, tt.want)
// 			}
// 		})
// 	}
// }

// func TestBackpack_RemoveItem(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		items        []ItemLike
// 		item         ItemLike
// 		wantBackpack Backpack
// 		want         error
// 	}{
// 		{
// 			name:  "remove elixir from empty backpack",
// 			items: []ItemLike{},
// 			item: &Elixir{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Elixir",
// 				},
// 				EffectDuration: time.Minute,
// 				AffectedAttribute: primitives.Attributes{
// 					MaxHealth: 1,
// 					Agility:   0,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			wantBackpack: Backpack{
// 				Capacity:  BackpackDefaultCapacity,
// 				ItemsNum:  0,
// 				Elixirs:   map[string][]Elixir{},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: ItemIsNotInBackpackError{},
// 		},
// 		{
// 			name:  "remove scroll from empty backpack",
// 			items: []ItemLike{},
// 			item: &Scroll{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Scroll",
// 				},
// 				AffectedAttribute: primitives.Attributes{
// 					MaxHealth: 0,
// 					Agility:   1,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			wantBackpack: Backpack{
// 				Capacity:  BackpackDefaultCapacity,
// 				ItemsNum:  0,
// 				Elixirs:   map[string][]Elixir{},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: ItemIsNotInBackpackError{},
// 		},
// 		{
// 			name:  "remove food from empty backpack",
// 			items: []ItemLike{},
// 			item: &Food{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Food",
// 				},
// 				HealthRegeneration: 15,
// 			},
// 			wantBackpack: Backpack{
// 				Capacity:  BackpackDefaultCapacity,
// 				ItemsNum:  0,
// 				Elixirs:   map[string][]Elixir{},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: ItemIsNotInBackpackError{},
// 		},
// 		{
// 			name:  "remove weapon from empty backpack",
// 			items: []ItemLike{},
// 			item: &Weapon{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Weapon",
// 				},
// 				Damage: 10,
// 			},
// 			wantBackpack: Backpack{
// 				Capacity:  BackpackDefaultCapacity,
// 				ItemsNum:  0,
// 				Elixirs:   map[string][]Elixir{},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: ItemIsNotInBackpackError{},
// 		},
// 		{
// 			name: "remove elixir",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			item: &Elixir{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Elixir",
// 				},
// 				EffectDuration: time.Minute,
// 				AffectedAttribute: primitives.Attributes{
// 					MaxHealth: 1,
// 					Agility:   0,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			wantBackpack: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 2,
// 				Elixirs: map[string][]Elixir{
// 					"Awkward Elixir": {
// 						Elixir{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Elixir{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Elixir",
// 							},
// 							EffectDuration: time.Minute,
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 1,
// 								Agility:   0,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 				},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: nil,
// 		},
// 		{
// 			name: "remove scroll",
// 			items: []ItemLike{
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			item: &Scroll{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Scroll",
// 				},
// 				AffectedAttribute: primitives.Attributes{
// 					MaxHealth: 0,
// 					Agility:   1,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			wantBackpack: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 2,
// 				Elixirs:  map[string][]Elixir{},
// 				Scrolls: map[string][]Scroll{
// 					"Awkward Scroll": {
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 						Scroll{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Scroll",
// 							},
// 							AffectedAttribute: primitives.Attributes{
// 								MaxHealth: 0,
// 								Agility:   1,
// 								Strength:  0,
// 							},
// 							Increment: 10,
// 						},
// 					},
// 				},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: nil,
// 		},
// 		{
// 			name: "remove food",
// 			items: []ItemLike{
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 			},
// 			item: &Food{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Food",
// 				},
// 				HealthRegeneration: 15,
// 			},
// 			wantBackpack: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 2,
// 				Elixirs:  map[string][]Elixir{},
// 				Scrolls:  map[string][]Scroll{},
// 				Foods: map[string][]Food{
// 					"Awkward Food": {
// 						Food{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Food",
// 							},
// 							HealthRegeneration: 15,
// 						},
// 						Food{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Food",
// 							},
// 							HealthRegeneration: 15,
// 						},
// 					},
// 				},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: nil,
// 		},
// 		{
// 			name: "remove weapon",
// 			items: []ItemLike{
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon",
// 					},
// 					Damage: 10,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon",
// 					},
// 					Damage: 10,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon",
// 					},
// 					Damage: 10,
// 				},
// 			},
// 			item: &Weapon{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Weapon",
// 				},
// 				Damage: 10,
// 			},
// 			wantBackpack: Backpack{
// 				Capacity: BackpackDefaultCapacity,
// 				ItemsNum: 2,
// 				Elixirs:  map[string][]Elixir{},
// 				Scrolls:  map[string][]Scroll{},
// 				Foods:    map[string][]Food{},
// 				Weapons: map[string][]Weapon{
// 					"Awkward Weapon": {
// 						Weapon{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Weapon",
// 							},
// 							Damage: 10,
// 						}, Weapon{
// 							Item: Item{
// 								Shape: primitives.Box{},
// 								Name:  "Awkward Weapon",
// 							},
// 							Damage: 10,
// 						},
// 					},
// 				},
// 				Treasures: 0,
// 			},
// 			want: nil,
// 		},
// 		{
// 			name: "remove last elixir",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			item: &Elixir{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Elixir",
// 				},
// 				EffectDuration: time.Minute,
// 				AffectedAttribute: primitives.Attributes{
// 					MaxHealth: 1,
// 					Agility:   0,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			wantBackpack: Backpack{
// 				Capacity:  BackpackDefaultCapacity,
// 				ItemsNum:  0,
// 				Elixirs:   map[string][]Elixir{},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: nil,
// 		},
// 		{
// 			name: "remove last scroll",
// 			items: []ItemLike{
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			item: &Scroll{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Scroll",
// 				},
// 				AffectedAttribute: primitives.Attributes{
// 					MaxHealth: 0,
// 					Agility:   1,
// 					Strength:  0,
// 				},
// 				Increment: 10,
// 			},
// 			wantBackpack: Backpack{
// 				Capacity:  BackpackDefaultCapacity,
// 				ItemsNum:  0,
// 				Elixirs:   map[string][]Elixir{},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: nil,
// 		},
// 		{
// 			name: "remove last food",
// 			items: []ItemLike{
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 			},
// 			item: &Food{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Food",
// 				},
// 				HealthRegeneration: 15,
// 			},
// 			wantBackpack: Backpack{
// 				Capacity:  BackpackDefaultCapacity,
// 				ItemsNum:  0,
// 				Elixirs:   map[string][]Elixir{},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: nil,
// 		},
// 		{
// 			name: "remove last weapon",
// 			items: []ItemLike{
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon",
// 					},
// 					Damage: 10,
// 				},
// 			},
// 			item: &Weapon{
// 				Item: Item{
// 					Shape: primitives.Box{},
// 					Name:  "Awkward Weapon",
// 				},
// 				Damage: 10,
// 			},
// 			wantBackpack: Backpack{
// 				Capacity:  BackpackDefaultCapacity,
// 				ItemsNum:  0,
// 				Elixirs:   map[string][]Elixir{},
// 				Scrolls:   map[string][]Scroll{},
// 				Foods:     map[string][]Food{},
// 				Weapons:   map[string][]Weapon{},
// 				Treasures: 0,
// 			},
// 			want: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b := NewBackpack()

// 			for idx := range tt.items {
// 				_ = b.AddItem(tt.items[idx])
// 			}

// 			got := b.RemoveItem(tt.item, primitives.Box{})

// 			if !errors.Is(got, tt.want) {
// 				t.Errorf("RemoveItem() = %#v, want %#v", got, tt.want)
// 			}

// 			if !reflect.DeepEqual(*b, tt.wantBackpack) {
// 				t.Errorf("RemoveItem(): Backpack got %#v, want %#v", *b, tt.wantBackpack)
// 			}
// 		})
// 	}
// }

// func TestBackpack_GetElixirsList(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		items []ItemLike
// 		want  ItemsList
// 	}{
// 		{
// 			name:  "zero elixirs",
// 			items: []ItemLike{},
// 			want:  ItemsList{},
// 		},
// 		{
// 			name: "multiple elixirs of same type",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			want: ItemsList{
// 				{
// 					Name: "Awkward Elixir",
// 					Num:  3,
// 				},
// 			},
// 		},
// 		{
// 			name: "multiple elixirs of different types",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 1",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 2",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 2",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 3",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			want: ItemsList{
// 				{
// 					Name: "Awkward Elixir 1",
// 					Num:  1,
// 				},
// 				{
// 					Name: "Awkward Elixir 2",
// 					Num:  2,
// 				},
// 				{
// 					Name: "Awkward Elixir 3",
// 					Num:  1,
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

// 			for idx := range tt.want {
// 				tt.want[idx].Item = &b.Elixirs[tt.want[idx].Name][0]
// 			}

// 			got := b.GetElixirsList()

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetElixirsList() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestBackpack_GetScrollsList(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		items []ItemLike
// 		want  ItemsList
// 	}{
// 		{
// 			name:  "zero scrolls",
// 			items: []ItemLike{},
// 			want:  ItemsList{},
// 		},
// 		{
// 			name: "multiple scrolls of same type",
// 			items: []ItemLike{
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			want: ItemsList{
// 				{
// 					Name: "Awkward Scroll",
// 					Num:  3,
// 				},
// 			},
// 		},
// 		{
// 			name: "multiple scrolls of different types",
// 			items: []ItemLike{
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll 1",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll 2",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll 2",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll 3",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 			},
// 			want: ItemsList{
// 				{
// 					Name: "Awkward Scroll 1",
// 					Num:  1,
// 				},
// 				{
// 					Name: "Awkward Scroll 2",
// 					Num:  2,
// 				},
// 				{
// 					Name: "Awkward Scroll 3",
// 					Num:  1,
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

// 			for idx := range tt.want {
// 				tt.want[idx].Item = &b.Scrolls[tt.want[idx].Name][0]
// 			}

// 			got := b.GetScrollsList()

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetScrollsList() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestBackpack_GetFoodsList(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		items []ItemLike
// 		want  ItemsList
// 	}{
// 		{
// 			name:  "zero foods",
// 			items: []ItemLike{},
// 			want:  ItemsList{},
// 		},
// 		{
// 			name: "multiple foods of same type",
// 			items: []ItemLike{
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 			},
// 			want: ItemsList{
// 				{
// 					Name: "Awkward Food",
// 					Num:  3,
// 				},
// 			},
// 		},
// 		{
// 			name: "multiple food of different types",
// 			items: []ItemLike{
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food 1",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food 2",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food 2",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food 3",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 			},
// 			want: ItemsList{
// 				{
// 					Name: "Awkward Food 1",
// 					Num:  1,
// 				},
// 				{
// 					Name: "Awkward Food 2",
// 					Num:  2,
// 				},
// 				{
// 					Name: "Awkward Food 3",
// 					Num:  1,
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

// 			for idx := range tt.want {
// 				tt.want[idx].Item = &b.Foods[tt.want[idx].Name][0]
// 			}

// 			got := b.GetFoodsList()

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetFoodsList() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestBackpack_GetWeaponsList(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		items []ItemLike
// 		want  ItemsList
// 	}{
// 		{
// 			name:  "zero weapons",
// 			items: []ItemLike{},
// 			want:  ItemsList{},
// 		},
// 		{
// 			name: "multiple weapons of same type",
// 			items: []ItemLike{
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon",
// 					},
// 					Damage: 10,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon",
// 					},
// 					Damage: 10,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon",
// 					},
// 					Damage: 10,
// 				},
// 			},
// 			want: ItemsList{
// 				{
// 					Name: "Awkward Weapon",
// 					Num:  3,
// 				},
// 			},
// 		},
// 		{
// 			name: "multiple weapons of different types",
// 			items: []ItemLike{
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon 1",
// 					},
// 					Damage: 10,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon 2",
// 					},
// 					Damage: 10,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon 2",
// 					},
// 					Damage: 10,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon 3",
// 					},
// 					Damage: 10,
// 				},
// 			},
// 			want: ItemsList{
// 				{
// 					Name: "Awkward Weapon 1",
// 					Num:  1,
// 				},
// 				{
// 					Name: "Awkward Weapon 2",
// 					Num:  2,
// 				},
// 				{
// 					Name: "Awkward Weapon 3",
// 					Num:  1,
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

// 			for idx := range tt.want {
// 				tt.want[idx].Item = &b.Weapons[tt.want[idx].Name][0]
// 			}

// 			got := b.GetWeaponsList()

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetWeaponsList() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestBackpack_GetItemsList(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		items []ItemLike
// 		want  ItemsList
// 	}{
// 		{
// 			name: "elixir, scroll, food, weapon",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon",
// 					},
// 					Damage: 10,
// 				},
// 			},
// 			want: ItemsList{
// 				{
// 					Name: "Awkward Elixir",
// 					Num:  1,
// 					Item: &Elixir{},
// 				},
// 				{
// 					Name: "Awkward Scroll",
// 					Num:  1,
// 					Item: &Scroll{},
// 				},
// 				{
// 					Name: "Awkward Food",
// 					Num:  1,
// 					Item: &Food{},
// 				},
// 				{
// 					Name: "Awkward Weapon",
// 					Num:  1,
// 					Item: &Weapon{},
// 				},
// 			},
// 		},
// 		{
// 			name: "multiple elixirs, scrolls, foods, weapons",
// 			items: []ItemLike{
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 1",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Elixir{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Elixir 2",
// 					},
// 					EffectDuration: time.Minute,
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 1,
// 						Agility:   0,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll 1",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Scroll{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Scroll 1",
// 					},
// 					AffectedAttribute: primitives.Attributes{
// 						MaxHealth: 0,
// 						Agility:   1,
// 						Strength:  0,
// 					},
// 					Increment: 10,
// 				},
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food 1",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 				&Food{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Food 2",
// 					},
// 					HealthRegeneration: 15,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon 1",
// 					},
// 					Damage: 10,
// 				},
// 				&Weapon{
// 					Item: Item{
// 						Shape: primitives.Box{},
// 						Name:  "Awkward Weapon 3",
// 					},
// 					Damage: 10,
// 				},
// 			},
// 			want: ItemsList{
// 				{
// 					Name: "Awkward Elixir 1",
// 					Num:  1,
// 					Item: &Elixir{},
// 				},
// 				{
// 					Name: "Awkward Elixir 2",
// 					Num:  1,
// 					Item: &Elixir{},
// 				},
// 				{
// 					Name: "Awkward Scroll 1",
// 					Num:  2,
// 					Item: &Scroll{},
// 				},
// 				{
// 					Name: "Awkward Food 1",
// 					Num:  1,
// 					Item: &Food{},
// 				},
// 				{
// 					Name: "Awkward Food 2",
// 					Num:  1,
// 					Item: &Food{},
// 				},
// 				{
// 					Name: "Awkward Weapon 1",
// 					Num:  1,
// 					Item: &Weapon{},
// 				},
// 				{
// 					Name: "Awkward Weapon 3",
// 					Num:  1,
// 					Item: &Weapon{},
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

// 			for idx := range tt.want {
// 				n := 0
// 				i := 0

// 				if IsElixir(tt.want[idx].Item) != nil {
// 					tt.want[idx].Item = &b.Elixirs[tt.want[idx].Name][0]
// 					i++
// 				}

// 				n = i

// 				if IsScroll(tt.want[idx].Item) != nil {
// 					tt.want[idx].Item = &b.Scrolls[tt.want[idx-n].Name][0]
// 					i++
// 				}

// 				n = i

// 				if IsFood(tt.want[idx].Item) != nil {
// 					tt.want[idx].Item = &b.Foods[tt.want[idx-n].Name][0]
// 					i++
// 				}

// 				n = i

// 				if IsWeapon(tt.want[idx].Item) != nil {
// 					tt.want[idx].Item = &b.Weapons[tt.want[idx-n].Name][0]
// 				}
// 			}

// 			got := b.GetItemsList()

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetItemsList() = %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestItemsList_Sort(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		elements []element
// 		want     ItemsList
// 	}{
// 		{
// 			name: "random order",
// 			elements: []element{
// 				{
// 					Name: "F",
// 					Num:  1,
// 					Item: nil,
// 				},
// 				{
// 					Name: "B",
// 					Num:  2,
// 					Item: nil,
// 				},
// 				{
// 					Name: "X",
// 					Num:  1,
// 					Item: nil,
// 				},
// 				{
// 					Name: "I",
// 					Num:  1,
// 					Item: nil,
// 				},
// 			},
// 			want: ItemsList{
// 				{
// 					Name: "B",
// 					Num:  2,
// 					Item: nil,
// 				},
// 				{
// 					Name: "F",
// 					Num:  1,
// 					Item: nil,
// 				},
// 				{
// 					Name: "I",
// 					Num:  1,
// 					Item: nil,
// 				},
// 				{
// 					Name: "X",
// 					Num:  1,
// 					Item: nil,
// 				},
// 			},
// 		},
// 		{
// 			name: "ascending order",
// 			elements: []element{
// 				{
// 					Name: "B",
// 					Num:  2,
// 					Item: nil,
// 				},
// 				{
// 					Name: "F",
// 					Num:  1,
// 					Item: nil,
// 				},
// 				{
// 					Name: "I",
// 					Num:  1,
// 					Item: nil,
// 				},
// 				{
// 					Name: "X",
// 					Num:  1,
// 					Item: nil,
// 				},
// 			},
// 			want: ItemsList{
// 				{
// 					Name: "B",
// 					Num:  2,
// 					Item: nil,
// 				},
// 				{
// 					Name: "F",
// 					Num:  1,
// 					Item: nil,
// 				},
// 				{
// 					Name: "I",
// 					Num:  1,
// 					Item: nil,
// 				},
// 				{
// 					Name: "X",
// 					Num:  1,
// 					Item: nil,
// 				},
// 			},
// 		},
// 		{
// 			name: "descending order",
// 			elements: []element{
// 				{
// 					Name: "X",
// 					Num:  1,
// 					Item: nil,
// 				},
// 				{
// 					Name: "I",
// 					Num:  1,
// 					Item: nil,
// 				},
// 				{
// 					Name: "F",
// 					Num:  1,
// 					Item: nil,
// 				},
// 				{
// 					Name: "B",
// 					Num:  2,
// 					Item: nil,
// 				},
// 			},
// 			want: ItemsList{
// 				{
// 					Name: "B",
// 					Num:  2,
// 					Item: nil,
// 				},
// 				{
// 					Name: "F",
// 					Num:  1,
// 					Item: nil,
// 				},
// 				{
// 					Name: "I",
// 					Num:  1,
// 					Item: nil,
// 				},
// 				{
// 					Name: "X",
// 					Num:  1,
// 					Item: nil,
// 				},
// 			},
// 		},
// 		{
// 			name:     "zero elements",
// 			elements: []element{},
// 			want:     ItemsList{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := ItemsList{}
// 			got = append(got, tt.elements...)
// 			got.Sort()

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Sort(): ItemsList got %#v, want %#v", got, tt.want)
// 			}
// 		})
// 	}
// }
