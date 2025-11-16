package items

const (
	DefaultBackpackCapacity uint = 9
)

type BackpackIsFullError struct{}

func (BackpackIsFullError) Error() string {
	return "backpack is full, drop something"
}

type ItemIsNotInBackpackError struct{}

func (ItemIsNotInBackpackError) Error() string {
	return "item is not found in the Backpack"
}

type Backpack struct {
	Capacity  uint
	ItemsNum  uint
	Elixirs   map[ElixirType][]*Elixir
	Scrolls   map[ScrollType][]*Scroll
	Foods     map[FoodType][]*Food
	Weapons   map[WeaponType][]*Weapon
	Treasures int32
}

func NewBackpack() *Backpack {
	return &Backpack{
		Capacity:  DefaultBackpackCapacity,
		ItemsNum:  0,
		Elixirs:   make(map[ElixirType][]*Elixir),
		Scrolls:   make(map[ScrollType][]*Scroll),
		Foods:     make(map[FoodType][]*Food),
		Weapons:   make(map[WeaponType][]*Weapon),
		Treasures: 0,
	}
}

func (b *Backpack) AddItem(item any) error {
	switch v := item.(type) {
	case *Elixir:
		return b.AddElixir(v)
	case *Scroll:
		return b.AddScroll(v)
	case *Food:
		return b.AddFood(v)
	case *Weapon:
		return b.AddWeapon(v)
	case *Treasure:
		b.AddTreasure(v)
		return nil
	default:
		return ItemIsNotInBackpackError{}
	}
}

func (b *Backpack) RemoveItem(item any) error {
	switch v := item.(type) {
	case *Elixir:
		return b.RemoveElixir(v.Type)
	case *Scroll:
		return b.RemoveScroll(v.Type)
	case *Food:
		return b.RemoveFood(v.Type)
	case *Weapon:
		return b.RemoveWeapon(v.Type)
	default:
		return ItemIsNotInBackpackError{}
	}
}

func (b *Backpack) IsFull() bool {
	return b.ItemsNum >= b.Capacity
}

func (b *Backpack) IsEmpty() bool {
	return b.ItemsNum == 0
}

// Treasures

func (b *Backpack) AddTreasure(t *Treasure) {
	b.Treasures += t.Value
}

// Elixirs

func (b *Backpack) AddElixir(e *Elixir) error {
	if b.ItemsNum >= b.Capacity {
		return BackpackIsFullError{}
	}

	b.Elixirs[e.Type] = append(b.Elixirs[e.Type], e)
	b.ItemsNum++
	return nil
}

func (b *Backpack) RemoveElixir(t ElixirType) error {
	if items, ok := b.Elixirs[t]; !ok || len(items) == 0 {
		return ItemIsNotInBackpackError{}
	}

	if len(b.Elixirs[t]) > 1 {
		b.Elixirs[t] = b.Elixirs[t][1:]
	} else {
		delete(b.Elixirs, t)
	}

	b.ItemsNum--
	return nil
}

// Scrolls

func (b *Backpack) AddScroll(s *Scroll) error {
	if b.ItemsNum >= b.Capacity {
		return BackpackIsFullError{}
	}

	t := s.Type
	b.Scrolls[t] = append(b.Scrolls[t], s)
	b.ItemsNum++
	return nil
}

func (b *Backpack) RemoveScroll(t ScrollType) error {
	if items, ok := b.Scrolls[t]; !ok || len(items) == 0 {
		return ItemIsNotInBackpackError{}
	}

	if len(b.Scrolls[t]) > 1 {
		b.Scrolls[t] = b.Scrolls[t][1:]
	} else {
		delete(b.Scrolls, t)
	}

	b.ItemsNum--
	return nil
}

// Foods

func (b *Backpack) AddFood(f *Food) error {
	if b.ItemsNum >= b.Capacity {
		return BackpackIsFullError{}
	}

	t := f.Type
	b.Foods[t] = append(b.Foods[t], f)
	b.ItemsNum++
	return nil
}

func (b *Backpack) RemoveFood(t FoodType) error {
	if items, ok := b.Foods[t]; !ok || len(items) == 0 {
		return ItemIsNotInBackpackError{}
	}

	if len(b.Foods[t]) > 1 {
		b.Foods[t] = b.Foods[t][1:]
	} else {
		delete(b.Foods, t)
	}

	b.ItemsNum--
	return nil
}

// Weapons

func (b *Backpack) AddWeapon(w *Weapon) error {
	if b.ItemsNum >= b.Capacity {
		return BackpackIsFullError{}
	}

	t := w.Type
	b.Weapons[t] = append(b.Weapons[t], w)
	b.ItemsNum++
	return nil
}

func (b *Backpack) RemoveWeapon(t WeaponType) error {
	if items, ok := b.Weapons[t]; !ok || len(items) == 0 {
		return ItemIsNotInBackpackError{}
	}

	if len(b.Weapons[t]) > 1 {
		b.Weapons[t] = b.Weapons[t][1:]
	} else {
		delete(b.Weapons, t)
	}

	b.ItemsNum--
	return nil
}
