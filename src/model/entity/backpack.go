package entity

const (
	BackpackDefaultCapacity uint = 9
)

type BackpackIsFullError struct{}

func (BackpackIsFullError) Error() string {
	return "Backpack is full, drop something"
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

type listElement struct {
	Name string
	Num int
	Item ItemLike
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
			b.Elixirs[e.Item.Name] = b.Elixirs[e.Item.Name][1:]
		} else {
			delete(b.Elixirs, e.Item.Name)
		}
	}

	s, ok := i.(*Scroll)
	if ok {
		if len(b.Scrolls[s.Item.Name]) > 1 {
			b.Scrolls[e.Item.Name] = b.Scrolls[e.Item.Name][1:]
		} else {
			delete(b.Scrolls, s.Item.Name)
		}
	}

	f, ok := i.(*Food)
	if ok {
		if len(b.Foods[f.Item.Name]) > 1 {
			b.Foods[e.Item.Name] = b.Foods[e.Item.Name][1:]
		} else {
			delete(b.Foods, f.Item.Name)
		}
	}

	w, ok := i.(*Weapon)
	if ok {
		if len(b.Weapons[w.Item.Name]) > 1 {
			b.Weapons[e.Item.Name] = b.Weapons[e.Item.Name][1:]
		} else {
			delete(b.Weapons, w.Item.Name)
		}
	}

	b.ItemsNum--
	return i.Dropped(box)
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
		list = append(list, listElement{Name: key, Num: len(b.Elixirs[key]), Item: b.Elixirs[key][0]})
	}

	return list
}

func (b *Backpack) GetScrollsList() []listElement {
	list := []listElement{}

	for key := range b.Scrolls {
		list = append(list, listElement{Name: key, Num: len(b.Scrolls[key]), Item: b.Elixirs[key][0]})
	}

	return list
}

func (b *Backpack) GetFoodsList() []listElement {
	list := []listElement{}

	for key := range b.Foods {
		list = append(list, listElement{Name: key, Num: len(b.Foods[key]), Item: b.Elixirs[key][0]})
	}

	return list
}

func (b *Backpack) GetWeaponsList() []listElement {
	list := []listElement{}

	for key := range b.Weapons {
		list = append(list, listElement{Name: key, Num: len(b.Weapons[key]), Item: b.Elixirs[key][0]})
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
