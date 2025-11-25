package items

import "container/list"

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
	Elixirs   *list.List
	Scrolls   *list.List
	Foods     *list.List
	Weapons   *list.List
	Treasures int32
}

func NewBackpack() *Backpack {
	return &Backpack{
		Capacity: DefaultBackpackCapacity,
		Elixirs:  list.New(),
		Scrolls:  list.New(),
		Foods:    list.New(),
		Weapons:  list.New(),
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
		return NotItemError{}
	}
}

func (b *Backpack) RemoveItem(item any) error {
	switch i := item.(type) {
	case *Elixir:
		return b.RemoveElixir(i)
	case *Scroll:
		return b.RemoveScroll(i)
	case *Food:
		return b.RemoveFood(i)
	case *Weapon:
		return b.RemoveWeapon(i)
	default:
		return NotItemError{}
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

	b.Elixirs.PushBack(e)
	b.ItemsNum++
	return nil
}

func (b *Backpack) RemoveElixir(e *Elixir) error {
	for elem := b.Elixirs.Front(); elem != nil; elem = elem.Next() {
		if elem.Value.(*Elixir) == e {
			b.Elixirs.Remove(elem)
			b.ItemsNum--
			return nil
		}
	}

	return ItemIsNotInBackpackError{}
}

// Scrolls

func (b *Backpack) AddScroll(s *Scroll) error {
	if b.ItemsNum >= b.Capacity {
		return BackpackIsFullError{}
	}

	b.Scrolls.PushBack(s)
	b.ItemsNum++
	return nil
}

func (b *Backpack) RemoveScroll(t *Scroll) error {
	for elem := b.Scrolls.Front(); elem != nil; elem = elem.Next() {
		if elem.Value.(*Scroll) == t {
			b.Scrolls.Remove(elem)
			b.ItemsNum--
			return nil
		}
	}

	return ItemIsNotInBackpackError{}
}

// Foods

func (b *Backpack) AddFood(f *Food) error {
	if b.ItemsNum >= b.Capacity {
		return BackpackIsFullError{}
	}

	b.Foods.PushBack(f)
	b.ItemsNum++
	return nil
}

func (b *Backpack) RemoveFood(t *Food) error {
	for elem := b.Foods.Front(); elem != nil; elem = elem.Next() {
		if elem.Value.(*Food) == t {
			b.Foods.Remove(elem)
			b.ItemsNum--
			return nil
		}
	}

	return ItemIsNotInBackpackError{}
}

// Weapons

func (b *Backpack) AddWeapon(w *Weapon) error {
	if b.ItemsNum >= b.Capacity {
		return BackpackIsFullError{}
	}

	b.Weapons.PushBack(w)
	b.ItemsNum++
	return nil
}

func (b *Backpack) RemoveWeapon(w *Weapon) error {
	for elem := b.Weapons.Front(); elem != nil; elem = elem.Next() {
		if elem.Value.(*Weapon) == w {
			b.Weapons.Remove(elem)
			b.ItemsNum--
			return nil
		}
	}

	return ItemIsNotInBackpackError{}
}
