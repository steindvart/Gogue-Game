package entity

const (
	BackpackDefaultCapacity uint = 9
)

type BackpackIsFullError struct{}

func (BackpackIsFullError) Error() string {
	return "Backpack is full, drop something"
}

type ItemIsNotInBackpackError struct{}

func (ItemIsNotInBackpackError) Error() string {
	return "Item is not found in the Backpack"
}

type Backpack struct {
	Capacity  uint
	ItemsNum  uint
	Elixirs   map[string][]Elixir
	Scrolls   map[string][]Scroll
	Foods     map[string][]Food
	Weapons   map[string][]Weapon
	Treasures uint
}

type listElement struct {
	Name string
	Num  int
	Item ItemLike
}

func (b *Backpack) AddItem(i ItemLike) error {
	if b.ItemsNum < b.Capacity {
		e, ok := i.(*Elixir)
		if ok {
			b.Elixirs[e.Item.Name] = append(b.Elixirs[e.Item.Name], *e)
		}

		s, ok := i.(*Scroll)
		if ok {
			b.Scrolls[s.Item.Name] = append(b.Scrolls[s.Item.Name], *s)
		}

		f, ok := i.(*Food)
		if ok {
			b.Foods[f.Item.Name] = append(b.Foods[f.Item.Name], *f)
		}

		w, ok := i.(*Weapon)
		if ok {
			b.Weapons[w.Item.Name] = append(b.Weapons[w.Item.Name], *w)
		}

		i.Taken()
		b.ItemsNum++
	} else {
		return BackpackIsFullError{}
	}

	return nil
}

func (b *Backpack) RemoveItem(i ItemLike, box Box) error {
	e, ok := i.(*Elixir)
	if ok {
		if len(b.Elixirs[e.Item.Name]) > 1 {
			b.Elixirs[e.Item.Name] = b.Elixirs[e.Item.Name][1:]
		} else {
			_, ok := b.Elixirs[e.Item.Name]
			if ok {
				delete(b.Elixirs, e.Item.Name)
			} else {
				return ItemIsNotInBackpackError{}
			}
		}

		b.ItemsNum--
		i.Dropped(box)
		return nil
	}

	s, ok := i.(*Scroll)
	if ok {
		if len(b.Scrolls[s.Item.Name]) > 1 {
			b.Scrolls[s.Item.Name] = b.Scrolls[s.Item.Name][1:]
		} else {
			_, ok := b.Scrolls[s.Item.Name]
			if ok {
				delete(b.Scrolls, s.Item.Name)
			} else {
				return ItemIsNotInBackpackError{}
			}
		}

		b.ItemsNum--
		i.Dropped(box)
		return nil
	}

	f, ok := i.(*Food)
	if ok {
		if len(b.Foods[f.Item.Name]) > 1 {
			b.Foods[f.Item.Name] = b.Foods[f.Item.Name][1:]
		} else {
			_, ok := b.Foods[f.Item.Name]
			if ok {
				delete(b.Foods, f.Item.Name)
			} else {
				return ItemIsNotInBackpackError{}
			}
		}

		b.ItemsNum--
		i.Dropped(box)
		return nil
	}

	w, ok := i.(*Weapon)
	if ok {
		if len(b.Weapons[w.Item.Name]) > 1 {
			b.Weapons[w.Item.Name] = b.Weapons[w.Item.Name][1:]
		} else {
			_, ok := b.Weapons[w.Item.Name]
			if ok {
				delete(b.Weapons, w.Item.Name)
			} else {
				return ItemIsNotInBackpackError{}
			}
		}

		b.ItemsNum--
		i.Dropped(box)
		return nil
	}

	return ItemIsNotInBackpackError{}
}

func (b *Backpack) GetItemsList() []listElement {
	list := []listElement{}

	list = append(list, b.GetElixirsList()...)
	list = append(list, b.GetScrollsList()...)
	list = append(list, b.GetFoodsList()...)
	list = append(list, b.GetWeaponsList()...)

	return list
}

func (b *Backpack) GetElixirsList() []listElement {
	list := []listElement{}

	for key := range b.Elixirs {
		list = append(list, listElement{Name: key, Num: len(b.Elixirs[key]), Item: &b.Elixirs[key][0]})
	}

	return list
}

func (b *Backpack) GetScrollsList() []listElement {
	list := []listElement{}

	for key := range b.Scrolls {
		list = append(list, listElement{Name: key, Num: len(b.Scrolls[key]), Item: &b.Scrolls[key][0]})
	}

	return list
}

func (b *Backpack) GetFoodsList() []listElement {
	list := []listElement{}

	for key := range b.Foods {
		list = append(list, listElement{Name: key, Num: len(b.Foods[key]), Item: &b.Foods[key][0]})
	}

	return list
}

func (b *Backpack) GetWeaponsList() []listElement {
	list := []listElement{}

	for key := range b.Weapons {
		list = append(list, listElement{Name: key, Num: len(b.Weapons[key]), Item: &b.Weapons[key][0]})
	}

	return list
}

func NewBackpack() *Backpack {
	return &Backpack{
		Capacity:  BackpackDefaultCapacity,
		ItemsNum:  0,
		Elixirs:   map[string][]Elixir{},
		Scrolls:   map[string][]Scroll{},
		Foods:     map[string][]Food{},
		Weapons:   map[string][]Weapon{},
		Treasures: 0,
	}
}
