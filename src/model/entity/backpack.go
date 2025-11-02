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

func (b *Backpack) AddItem(c ConsumableLike) error {
	if b.ItemsNum < b.Capacity {
		e, ok := c.(*Elixir)
		if ok {
			b.Elixirs[e.Consumable.Name] = append(b.Elixirs[e.Consumable.Name], e)
		}

		s, ok := c.(*Scroll)
		if ok {
			b.Scrolls[s.Consumable.Name] = append(b.Scrolls[s.Consumable.Name], s)
		}

		f, ok := c.(*Food)
		if ok {
			b.Foods[f.Consumable.Name] = append(b.Foods[f.Consumable.Name], f)
		}

		w, ok := c.(*Weapon)
		if ok {
			b.Weapons[w.Consumable.Name] = append(b.Weapons[w.Consumable.Name], w)
		}

		c.Taken()
		b.ItemsNum++
	} else {
		return BackpackIsFullError{}
	}

	return nil
}

func (b *Backpack) RemoveItem(c ConsumableLike, box Box) ConsumableLike {
	e, ok := c.(*Elixir)
	if ok {
		if len(b.Elixirs[e.Consumable.Name]) > 1 {
			idx := slices.Index(b.Elixirs[e.Consumable.Name], e)
			b.Elixirs[e.Consumable.Name] = append(b.Elixirs[e.Consumable.Name][:idx], b.Elixirs[e.Consumable.Name][idx+1:]...)
		} else {
			delete(b.Elixirs, e.Consumable.Name)
		}
	}

	s, ok := c.(*Scroll)
	if ok {
		if len(b.Scrolls[s.Consumable.Name]) > 1 {
			idx := slices.Index(b.Scrolls[s.Consumable.Name], s)
			b.Scrolls[s.Consumable.Name] = append(b.Scrolls[s.Consumable.Name][:idx], b.Scrolls[s.Consumable.Name][idx+1:]...)
		} else {
			delete(b.Scrolls, s.Consumable.Name)
		}
	}

	f, ok := c.(*Food)
	if ok {
		if len(b.Foods[f.Consumable.Name]) > 1 {
			idx := slices.Index(b.Foods[f.Consumable.Name], f)
			b.Foods[f.Consumable.Name] = append(b.Foods[f.Consumable.Name][:idx], b.Foods[f.Consumable.Name][idx+1:]...)
		} else {
			delete(b.Foods, f.Consumable.Name)
		}
	}

	w, ok := c.(*Weapon)
	if ok {
		if len(b.Weapons[w.Consumable.Name]) > 1 {
			idx := slices.Index(b.Weapons[w.Consumable.Name], w)
			b.Weapons[w.Consumable.Name] = append(b.Weapons[w.Consumable.Name][:idx], b.Weapons[w.Consumable.Name][idx+1:]...)
		} else {
			delete(b.Weapons, w.Consumable.Name)
		}
	}

	b.ItemsNum--
	return c.Dropped(box)
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
