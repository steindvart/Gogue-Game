package entity

const (
	BackpackDefaultCapacity uint = 9
)

type Backpack struct {
	Capacity    uint
	ItemsNum    uint
	Consumables []ConsumableLike
	Treasures   uint
}

// func (b *Backpack) GetItemsNamesList() []string {
// 	list := []string{}

// 	for _, item := range b.Consumables {
// 		// list = append(list, item.(Consumable).Name)
// 	}

// 	return list
// }

func NewBackpack() *Backpack {
	return &Backpack{
		Capacity:    BackpackDefaultCapacity,
		ItemsNum:    0,
		Consumables: []ConsumableLike{},
		Treasures:   0,
	}
}
