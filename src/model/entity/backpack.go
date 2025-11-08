package entity

import (
	"sort"
)

const (
	BackpackDefaultCapacity uint = 9
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
	Elixirs   map[string][]Elixir
	Scrolls   map[string][]Scroll
	Foods     map[string][]Food
	Weapons   map[string][]Weapon
	Treasures uint
}

type element struct {
	Name string
	Num  int
	Item ItemLike
}

type ItemsList []element

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

func (b *Backpack) AddTreasure(t *Treasure) {
	b.Treasures += t.Value
}

func (b *Backpack) AddItem(item Takeable) error {
	if b.ItemsNum >= b.Capacity {
		return BackpackIsFullError{}
	}

	if e := IsElixir(item); e != nil {
		b.Elixirs[e.Item.Name] = append(b.Elixirs[e.Item.Name], *e)
	}

	if s := IsScroll(item); s != nil {
		b.Scrolls[s.Item.Name] = append(b.Scrolls[s.Item.Name], *s)
	}

	if f := IsFood(item); f != nil {
		b.Foods[f.Item.Name] = append(b.Foods[f.Item.Name], *f)
	}

	if w := IsWeapon(item); w != nil {
		b.Weapons[w.Item.Name] = append(b.Weapons[w.Item.Name], *w)
	}

	item.Take()
	b.ItemsNum++

	return nil
}

func (b *Backpack) RemoveItem(item Droppable, box Box) error {
	var removed bool

	if e := IsElixir(item); e != nil {
		removed = removeItemFromMap(e.Item.Name, b.Elixirs)
	}

	if s := IsScroll(item); s != nil {
		removed = removeItemFromMap(s.Item.Name, b.Scrolls)
	}

	if f := IsFood(item); f != nil {
		removed = removeItemFromMap(f.Item.Name, b.Foods)
	}

	if w := IsWeapon(item); w != nil {
		removed = removeItemFromMap(w.Item.Name, b.Weapons)
	}

	if removed {
		item.Drop(box)
		b.ItemsNum--

		return nil
	}

	return ItemIsNotInBackpackError{}
}

func IsElixir(item any) *Elixir {
	e, ok := item.(*Elixir)
	if ok {
		return e
	}
	return nil
}

func IsScroll(item any) *Scroll {
	s, ok := item.(*Scroll)
	if ok {
		return s
	}
	return nil
}

func IsFood(item any) *Food {
	f, ok := item.(*Food)
	if ok {
		return f
	}
	return nil
}

func IsWeapon(item any) *Weapon {
	w, ok := item.(*Weapon)
	if ok {
		return w
	}
	return nil
}

func removeItemFromMap[V Elixir | Scroll | Food | Weapon](key string, itemsMap map[string][]V) bool {
	if _, ok := itemsMap[key]; !ok {
		return false
	}

	if len(itemsMap[key]) > 1 {
		itemsMap[key] = itemsMap[key][1:]
	} else {
		delete(itemsMap, key)
	}

	return true
}

func (b *Backpack) GetItemsList() ItemsList {
	list := ItemsList{}

	list = append(list, b.GetElixirsList()...)
	list = append(list, b.GetScrollsList()...)
	list = append(list, b.GetFoodsList()...)
	list = append(list, b.GetWeaponsList()...)

	return list
}

func (b *Backpack) GetElixirsList() ItemsList {
	return appendItemsList(b.Elixirs)
}

func (b *Backpack) GetScrollsList() ItemsList {
	return appendItemsList(b.Scrolls)
}

func (b *Backpack) GetFoodsList() ItemsList {
	return appendItemsList(b.Foods)
}

func (b *Backpack) GetWeaponsList() ItemsList {
	return appendItemsList(b.Weapons)
}

func appendItemsList[V Elixir | Scroll | Food | Weapon](itemsMap map[string][]V) ItemsList {
	list := ItemsList{}

	for key := range itemsMap {
		item := itemsMap[key][0]
		ptr := any(&item).(ItemLike)
		list = append(list, element{Name: key, Num: len(itemsMap[key]), Item: ptr})
	}

	list.Sort()

	return list
}

func (l *ItemsList) Sort() {
	sort.Slice(*l, func(i int, j int) bool {
		return (*l)[i].Name < (*l)[j].Name
	})
}

func (b *Backpack) weaponIsInBackpack(w *Weapon) error {
	if _, ok := b.Weapons[w.Item.Name]; !ok {
		return ItemIsNotInBackpackError{}
	}

	return nil
}
