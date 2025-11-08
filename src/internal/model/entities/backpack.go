package entities

import "gogue/internal/model/items"

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
	Elixirs   map[string][]items.Elixir
	Scrolls   map[string][]items.Scroll
	Foods     map[string][]items.Food
	Weapons   map[string][]items.Weapon
	Treasures uint
}

// type element struct {
// 	Name string
// 	Num  int
// 	Item items.Type
// }

// type ItemsList []element

func NewBackpack() *Backpack {
	return &Backpack{
		Capacity:  BackpackDefaultCapacity,
		ItemsNum:  0,
		Elixirs:   map[string][]items.Elixir{},
		Scrolls:   map[string][]items.Scroll{},
		Foods:     map[string][]items.Food{},
		Weapons:   map[string][]items.Weapon{},
		Treasures: 0,
	}
}

func (b *Backpack) AddTreasure(t *items.Treasure) {
	b.Treasures += t.Value
}

func (b *Backpack) AddItem(item any) error {
	if b.ItemsNum >= b.Capacity {
		return BackpackIsFullError{}
	}

	if e := items.AsElixir(item); e != nil {
		b.Elixirs[e.Item.Name] = append(b.Elixirs[e.Item.Name], *e)
	}

	if s := items.AsScroll(item); s != nil {
		b.Scrolls[s.Item.Name] = append(b.Scrolls[s.Item.Name], *s)
	}

	if f := items.AsFood(item); f != nil {
		b.Foods[f.Item.Name] = append(b.Foods[f.Item.Name], *f)
	}

	if w := items.AsWeapon(item); w != nil {
		b.Weapons[w.Item.Name] = append(b.Weapons[w.Item.Name], *w)
	}

	b.ItemsNum++

	return nil
}

func (b *Backpack) RemoveItem(item any) error {
	if e := items.AsElixir(item); e != nil {
		removeItemFromMap(e.Item.Name, b.Elixirs)
	} else if s := items.AsScroll(item); s != nil {
		removeItemFromMap(s.Item.Name, b.Scrolls)
	} else if f := items.AsFood(item); f != nil {
		removeItemFromMap(f.Item.Name, b.Foods)
	} else if w := items.AsWeapon(item); w != nil {
		removeItemFromMap(w.Item.Name, b.Weapons)
	} else {
		return ItemIsNotInBackpackError{}
	}

	b.ItemsNum--
	return nil
}

func removeItemFromMap[V items.Type](key string, itemsMap map[string][]V) bool {
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

// func (b *Backpack) GetItemsList() ItemsList {
// 	list := ItemsList{}

// 	list = append(list, b.GetElixirsList()...)
// 	list = append(list, b.GetScrollsList()...)
// 	list = append(list, b.GetFoodsList()...)
// 	list = append(list, b.GetWeaponsList()...)

// 	return list
// }

// func (b *Backpack) GetElixirsList() ItemsList {
// 	return appendItemsList(b.Elixirs)
// }

// func (b *Backpack) GetScrollsList() ItemsList {
// 	return appendItemsList(b.Scrolls)
// }

// func (b *Backpack) GetFoodsList() ItemsList {
// 	return appendItemsList(b.Foods)
// }

// func (b *Backpack) GetWeaponsList() ItemsList {
// 	return appendItemsList(b.Weapons)
// }

// func appendItemsList[V items.Elixir | items.Scroll | items.Food | items.Weapon](itemsMap map[string][]V) ItemsList {
// 	list := ItemsList{}

// 	for key := range itemsMap {
// 		item := itemsMap[key][0]
// 		ptr := any(&item).(ItemLike)
// 		list = append(list, element{Name: key, Num: len(itemsMap[key]), Item: ptr})
// 	}

// 	list.Sort()

// 	return list
// }

// func (l *ItemsList) Sort() {
// 	sort.Slice(*l, func(i int, j int) bool {
// 		return (*l)[i].Name < (*l)[j].Name
// 	})
// }

// func (b *Backpack) weaponIsInBackpack(w *items.Weapon) error {
// 	if _, ok := b.Weapons[w.Item.Name]; !ok {
// 		return ItemIsNotInBackpackError{}
// 	}

// 	return nil
// }
