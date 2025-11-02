package entity

import "slices"

const (
	BackpackDefaultCapacity uint = 9
)

type BackpackIsFullError struct{}

func (BackpackIsFullError) Error() string {
	return "backpack is full, drop something"
}

type Backpack struct {
	Capacity  uint
	ItemsNum  uint
	Elixirs   map[string][]*Elixir
	Scrolls   map[string][]*Scroll
	Foods     map[string][]*Food
	Weapons   map[string][]*Weapon
	Treasures uint
}

func (b *Backpack) AddItem(i ItemLike) error {
	if b.ItemsNum < b.Capacity {
		e, ok := i.(*Elixir)
		if ok {
			b.Elixirs[e.Item.Name] = append(b.Elixirs[e.Item.Name], e)
		}

		s, ok := i.(*Scroll)
		if ok {
			b.Scrolls[s.Item.Name] = append(b.Scrolls[s.Item.Name], s)
		}

		f, ok := i.(*Food)
		if ok {
			b.Foods[f.Item.Name] = append(b.Foods[f.Item.Name], f)
		}

		w, ok := i.(*Weapon)
		if ok {
			b.Weapons[w.Item.Name] = append(b.Weapons[w.Item.Name], w)
		}

		i.Taken()
		b.ItemsNum++
	} else {
		return BackpackIsFullError{}
	}

	return nil
}

func (b *Backpack) RemoveItem(i ItemLike, box Box) ItemLike {
	e, ok := i.(*Elixir)
	if ok {
		if len(b.Elixirs[e.Item.Name]) > 1 {
			idx := slices.Index(b.Elixirs[e.Item.Name], e)
			b.Elixirs[e.Item.Name] = append(b.Elixirs[e.Item.Name][:idx], b.Elixirs[e.Item.Name][idx+1:]...)
		} else {
			delete(b.Elixirs, e.Item.Name)
		}
	}

	s, ok := i.(*Scroll)
	if ok {
		if len(b.Scrolls[s.Item.Name]) > 1 {
			idx := slices.Index(b.Scrolls[s.Item.Name], s)
			b.Scrolls[s.Item.Name] = append(b.Scrolls[s.Item.Name][:idx], b.Scrolls[s.Item.Name][idx+1:]...)
		} else {
			delete(b.Scrolls, s.Item.Name)
		}
	}

	f, ok := i.(*Food)
	if ok {
		if len(b.Foods[f.Item.Name]) > 1 {
			idx := slices.Index(b.Foods[f.Item.Name], f)
			b.Foods[f.Item.Name] = append(b.Foods[f.Item.Name][:idx], b.Foods[f.Item.Name][idx+1:]...)
		} else {
			delete(b.Foods, f.Item.Name)
		}
	}

	w, ok := i.(*Weapon)
	if ok {
		if len(b.Weapons[w.Item.Name]) > 1 {
			idx := slices.Index(b.Weapons[w.Item.Name], w)
			b.Weapons[w.Item.Name] = append(b.Weapons[w.Item.Name][:idx], b.Weapons[w.Item.Name][idx+1:]...)
		} else {
			delete(b.Weapons, w.Item.Name)
		}
	}

	b.ItemsNum--
	return i.Dropped(box)
}

func (b *Backpack) GetItemsList() map[string]int {
	list := map[string]int{}

	for key := range b.Elixirs {
		list[key] = len(b.Elixirs[key])
	}

	for key := range b.Scrolls {
		list[key] = len(b.Scrolls[key])
	}

	for key := range b.Foods {
		list[key] = len(b.Foods[key])
	}

	for key := range b.Weapons {
		list[key] = len(b.Weapons[key])
	}

	return list
}

func (b *Backpack) GetElixirsList() map[string]int {
	list := map[string]int{}

	for key := range b.Elixirs {
		list[key] = len(b.Elixirs[key])
	}

	return list
}

func (b *Backpack) GetScrollsList() map[string]int {
	list := map[string]int{}

	for key := range b.Scrolls {
		list[key] = len(b.Scrolls[key])
	}

	return list
}

func (b *Backpack) GetFoodsList() map[string]int {
	list := map[string]int{}

	for key := range b.Foods {
		list[key] = len(b.Foods[key])
	}

	return list
}

func (b *Backpack) GetWeaponsList() map[string]int {
	list := map[string]int{}

	for key := range b.Weapons {
		list[key] = len(b.Weapons[key])
	}

	return list
}

func NewBackpack() *Backpack {
	return &Backpack{
		Capacity:  BackpackDefaultCapacity,
		ItemsNum:  0,
		Elixirs:   map[string][]*Elixir{},
		Scrolls:   map[string][]*Scroll{},
		Foods:     map[string][]*Food{},
		Weapons:   map[string][]*Weapon{},
		Treasures: 0,
	}
}
