package entity

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestPlayer_IsAlive(t *testing.T) {
	tests := []struct {
		name   string
		health float64
		want   bool
	}{
		{
			name:   "alive",
			health: 10,
			want:   true,
		},
		{
			name:   "dead",
			health: 0,
			want:   false,
		},
		{
			name:   "negative health",
			health: -5,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Character: Character{Health: tt.health}}
			got := p.IsAlive()
			if got != tt.want {
				t.Errorf("IsAlive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_Move(t *testing.T) {
	tests := []struct {
		name  string
		start Point2D[int]
		delta Point2D[int]
		want  Point2D[int]
	}{
		{
			name:  "move positive",
			start: Point2D[int]{X: 0, Y: 0},
			delta: Point2D[int]{X: 2, Y: 3},
			want:  Point2D[int]{X: 2, Y: 3},
		},
		{
			name:  "move negative",
			start: Point2D[int]{X: 5, Y: 5},
			delta: Point2D[int]{X: -2, Y: -3},
			want:  Point2D[int]{X: 3, Y: 2},
		},
		{
			name:  "move zero",
			start: Point2D[int]{X: 1, Y: 1},
			delta: Point2D[int]{X: 0, Y: 0},
			want:  Point2D[int]{X: 1, Y: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Character: Character{Shape: Box{Point: tt.start}}}
			p.Move(tt.delta)
			if p.Character.Shape.Point != tt.want {
				t.Errorf("Move() = (%v), want (%v)", p.Character.Shape.Point, tt.want)
			}
		})
	}
}

func TestPlayer_TakeDamage(t *testing.T) {
	tests := []struct {
		name   string
		health float64
		damage float64
		want   float64
	}{
		{
			name:   "normal damage",
			health: 10,
			damage: 4,
			want:   6,
		},
		{
			name:   "overkill",
			health: 5,
			damage: 10,
			want:   0,
		},
		{
			name:   "zero damage",
			health: 7,
			damage: 0,
			want:   7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Character: Character{Health: tt.health}}
			p.TakeDamage(tt.damage)
			if p.Character.Health != tt.want {
				t.Errorf("TakeDamage() = %v, want %v", p.Character.Health, tt.want)
			}
		})
	}
}

func TestPlayer_Heal(t *testing.T) {
	tests := []struct {
		name      string
		health    float64
		maxHealth float64
		heal      float64
		want      float64
	}{
		{
			name:      "normal heal",
			health:    5,
			maxHealth: 10,
			heal:      3,
			want:      8,
		},
		{
			name:      "overheal",
			health:    8,
			maxHealth: 10,
			heal:      5,
			want:      10,
		},
		{
			name:      "zero heal",
			health:    7,
			maxHealth: 10,
			heal:      0,
			want:      7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Character: Character{Health: tt.health, MaxHealth: tt.maxHealth}}
			p.Heal(tt.heal)
			if p.Character.Health != tt.want {
				t.Errorf("Heal() = %v, want %v", p.Character.Health, tt.want)
			}
		})
	}
}

func TestPlayer_Attack(t *testing.T) {
	tests := []struct {
		name     string
		strength uint
		want     uint
	}{
		{
			name:     "normal attack",
			strength: 7,
			want:     7,
		},
		{
			name:     "zero strength",
			strength: 0,
			want:     0,
		},
		{
			name:     "high strength",
			strength: 100,
			want:     100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Character: Character{Strength: tt.strength}}
			got := p.Attack()
			if got != tt.want {
				t.Errorf("Attack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_AttackWithWeapon(t *testing.T) {
	tests := []struct {
		name     string
		strength uint
		want     uint
	}{
		{
			name:     "normal attack",
			strength: 7,
			want:     7,
		},
		{
			name:     "zero strength",
			strength: 0,
			want:     0,
		},
		{
			name:     "high strength",
			strength: 100,
			want:     100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Character: Character{Strength: tt.strength}}
			p.Weapon = NewWeapon(Box{})

			got := p.Attack()
			if got != tt.want+p.Weapon.Damage {
				t.Errorf("Attack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_CheckEvasion(t *testing.T) {
	tests := []struct {
		name         string
		agility      uint
		wantAllFalse bool
		wantHighRate bool
	}{
		{
			name:         "zero agility",
			agility:      0,
			wantAllFalse: true,
			wantHighRate: false,
		},
		{
			name:         "high agility",
			agility:      1000,
			wantAllFalse: false,
			wantHighRate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Character: Character{Agility: tt.agility}}
			tries := 1000
			count := 0
			for i := 0; i < tries; i++ {
				if p.CheckEvasion() {
					count++
				}
			}
			if tt.wantAllFalse && count != 0 {
				t.Errorf("CheckEvasion() with agility 0 should always be false, got %d/%d", count, tries)
			}
			if tt.wantHighRate && count < tries/2 {
				t.Errorf("CheckEvasion() with high agility should be high rate, got %d/%d", count, tries)
			}
			if count == tries {
				t.Errorf("CheckEvasion() should never be 100%%")
			}
		})
	}
}

func TestPlayer_TakeTreasure(t *testing.T) {
	tests := []struct {
		name     string
		treasure Treasure
		want     uint
	}{
		{
			name: "take treasure cost 10",
			treasure: Treasure{
				Shape: Box{},
				Name:  "Gold",
				Value: 10,
			},
			want: 10,
		},
		{
			name: "take treasure cost 0",
			treasure: Treasure{
				Shape: Box{},
				Name:  "Gold",
				Value: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(Box{})

			p.TakeTreasure(&tt.treasure)
			got := p.Backpack.Treasures
			if got != tt.want {
				t.Errorf("TakeTreasure(): Backpack Treasures got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_TakeItem(t *testing.T) {
	tests := []struct {
		name         string
		items        []ItemLike
		wantBackpack Backpack
		want         []error
	}{
		{
			name:  "take elixir",
			items: []ItemLike{&elixir},
			wantBackpack: Backpack{
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
			want: []error{nil},
		},
		{
			name:  "take nine elixirs",
			items: []ItemLike{&elixir, &elixir, &elixir, &elixir, &elixir, &elixir, &elixir, &elixir, &elixir},
			wantBackpack: Backpack{
				Capacity: BackpackDefaultCapacity,
				ItemsNum: 9,
				Elixirs: map[string][]Elixir{
					"Awkward Elixir": {elixir, elixir, elixir, elixir, elixir, elixir, elixir, elixir, elixir},
				},
				Scrolls:   map[string][]Scroll{},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			want: []error{nil, nil, nil, nil, nil, nil, nil, nil, nil},
		},
		{
			name:  "take ten elixirs",
			items: []ItemLike{&elixir, &elixir, &elixir, &elixir, &elixir, &elixir, &elixir, &elixir, &elixir, &elixir},
			wantBackpack: Backpack{
				Capacity: BackpackDefaultCapacity,
				ItemsNum: 9,
				Elixirs: map[string][]Elixir{
					"Awkward Elixir": {elixir, elixir, elixir, elixir, elixir, elixir, elixir, elixir, elixir},
				},
				Scrolls:   map[string][]Scroll{},
				Foods:     map[string][]Food{},
				Weapons:   map[string][]Weapon{},
				Treasures: 0,
			},
			want: []error{nil, nil, nil, nil, nil, nil, nil, nil, nil, BackpackIsFullError{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(Box{})
			for idx := range tt.items {
				err := p.TakeItem(tt.items[idx])
				errWant := tt.want[idx]

				if !errors.Is(err, errWant) {
					t.Errorf("TakeItem() = %#v, want %#v", err, errWant)
				}
			}

			if !reflect.DeepEqual(*p.Backpack, tt.wantBackpack) {
				t.Errorf("TakeItem(): Backpack got %#v, want %#v", *p.Backpack, tt.wantBackpack)
			}

		})
	}
}

func TestPlayer_DropItem(t *testing.T) {
	tests := []struct {
		name         string
		items        []ItemLike
		item         ItemLike
		wantBackpack Backpack
		want         error
	}{
		{
			name:  "drop elixir from empty backpack",
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
			name:  "drop elixir",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(Box{})

			for idx := range tt.items {
				_ = p.Backpack.AddItem(tt.items[idx])
			}

			got := p.DropItem(tt.item)

			if !errors.Is(got, tt.want) {
				t.Errorf("DropItem() = %#v, want %#v", got, tt.want)
			}

			if !reflect.DeepEqual(*p.Backpack, tt.wantBackpack) {
				t.Errorf("DropItem(): Backpack got %#v, want %#v", *p.Backpack, tt.wantBackpack)
			}
		})
	}
}

func TestPlayer_UseItem(t *testing.T) {
	tests := []struct {
		name       string
		items      []ItemLike
		item       ItemLike
		wantPlayer Player
		wantString string
		wantError  error
	}{
		{
			name:  "use elixir, which is not in backpack",
			items: []ItemLike{},
			item:  &elixir,
			wantPlayer: Player{
				Character: Character{
					Shape:     Box{},
					Health:    float64(AttributeRateAverage),
					MaxHealth: float64(AttributeRateAverage),
					Strength:  uint(AttributeRateAverage),
					Agility:   uint(AttributeRateAverage),
				},
				Experience:     0,
				CharacterLevel: 1,
				Backpack:       NewBackpack(),
				Weapon:         nil,
			},
			wantString: "",
			wantError:  ItemIsNotInBackpackError{},
		},
		{
			name:  "use elixir",
			items: []ItemLike{&elixir},
			item:  &elixir,
			wantPlayer: Player{
				Character: Character{
					Shape:     Box{},
					Health:    float64(AttributeRateAverage),
					MaxHealth: float64(AttributeRateAverage) + float64(10),
					Strength:  uint(AttributeRateAverage),
					Agility:   uint(AttributeRateAverage),
				},
				Experience:     0,
				CharacterLevel: 1,
				Backpack:       NewBackpack(),
				Weapon:         nil,
			},
			wantString: "You drank the Awkward Elixir, your MaxHealth has increased by 10",
			wantError:  nil,
		},
		{
			name:  "use scroll",
			items: []ItemLike{&scroll},
			item:  &scroll,
			wantPlayer: Player{
				Character: Character{
					Shape:     Box{},
					Health:    float64(AttributeRateAverage),
					MaxHealth: float64(AttributeRateAverage),
					Strength:  uint(AttributeRateAverage),
					Agility:   uint(AttributeRateAverage) + 10,
				},
				Experience:     0,
				CharacterLevel: 1,
				Backpack:       NewBackpack(),
				Weapon:         nil,
			},
			wantString: "You read the Awkward Scroll, your Agility has increased by 10",
			wantError:  nil,
		},
		{
			name:  "use food",
			items: []ItemLike{&food},
			item:  &food,
			wantPlayer: Player{
				Character: Character{
					Shape:     Box{},
					Health:    float64(AttributeRateAverage) + float64(15),
					MaxHealth: float64(AttributeRateAverage),
					Strength:  uint(AttributeRateAverage),
					Agility:   uint(AttributeRateAverage),
				},
				Experience:     0,
				CharacterLevel: 1,
				Backpack:       NewBackpack(),
				Weapon:         nil,
			},
			wantString: "You ate the Awkward Food, your Health has increased by 15",
			wantError:  nil,
		},
		{
			name:  "use weapon",
			items: []ItemLike{&weapon},
			item:  &weapon,
			wantPlayer: Player{
				Character: Character{
					Shape:     Box{},
					Health:    float64(AttributeRateAverage),
					MaxHealth: float64(AttributeRateAverage),
					Strength:  uint(AttributeRateAverage),
					Agility:   uint(AttributeRateAverage),
				},
				Experience:     0,
				CharacterLevel: 1,
				Backpack:       NewBackpack(),
				Weapon:         &weapon,
			},
			wantError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(Box{})
			p.Backpack = tt.wantPlayer.Backpack

			for idx := range tt.items {
				_ = p.Backpack.AddItem(tt.items[idx])
			}

			gotErr := p.UseItem(tt.item)

			time.Sleep(1 * time.Millisecond)

			if !errors.Is(gotErr, tt.wantError) {
				t.Errorf("UseItem(): error got %#v, want %#v", gotErr, tt.wantError)
			}

			if !reflect.DeepEqual(*p, tt.wantPlayer) {
				t.Errorf("UseItem(): Player got %#v, want %#v", *p, tt.wantPlayer)
			}
		})
	}
}

func TestNewPlayer(t *testing.T) {
	tests := []struct {
		name string
		box  Box
		want *Player
	}{
		{
			name: "player constructor",
			box:  Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}},
			want: &Player{
				Character: Character{
					Shape:     Box{Point: Point2D[int]{X: 1, Y: 2}, Size: Size2D[uint]{Height: 1, Width: 2}},
					Health:    float64(AttributeRateAverage),
					MaxHealth: float64(AttributeRateAverage),
					Strength:  uint(AttributeRateAverage),
					Agility:   uint(AttributeRateAverage),
				},
				Experience:     0,
				CharacterLevel: 1,
				Backpack:       NewBackpack(),
				Weapon:         nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPlayer(tt.box)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPlayer() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
