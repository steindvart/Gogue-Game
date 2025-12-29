package entities

import (
	"container/list"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
)

type Player struct {
	*Character
	*items.Backpack
	*items.Weapon
	// @todo - убрать поля уровня и опыта. Это не нужно
	Experience uint
	Level      uint
	ViewRadius int
}

type BackpackItemType int

const (
	BackpackItemTypeWeapon BackpackItemType = iota
	BackpackItemTypeFood
	BackpackItemTypeElixir
	BackpackItemTypeScroll
)

func NewPlayer(box primitives.Box) *Player {
	return &Player{
		Character: &Character{
			Box: &box,
			Attributes: &primitives.Attributes{
				Health:    100,
				MaxHealth: 100,
				Strength:  10,
				Agility:   5,
			},
		},
		Experience: 0,
		Level:      1,
		ViewRadius: 3,
		Backpack:   items.NewBackpack(),
		Weapon:     nil,
	}
}

func (p *Player) DropEquipWeaponToBackpack() error {
	if p.Weapon == nil {
		return nil
	}

	if err := p.Backpack.AddItem(p.Weapon); err != nil {
		return err
	}

	_ = p.UnequipWeapon()

	return nil
}

func (p *Player) EquipWeapon(w *items.Weapon) error {
	if w == nil {
		return nil
	}

	// Если уже есть экипированный предмет, пытаемся положить его в рюкзак.
	if err := p.DropEquipWeaponToBackpack(); err != nil {
		return err
	}

	p.Weapon = w
	p.Character.Use(w)

	return nil
}

func (p *Player) UnequipWeapon() *items.Weapon {
	w := p.Weapon
	if w != nil {
		p.ApplyEffect(&primitives.Effect{
			Attributes: primitives.Inverse(w.Effect.Attributes),
		})
		p.Weapon = nil
	}
	return w
}

// GetItemFromBackpackByIndex возвращает предмет из рюкзака по индексу в зависимости от типа
func (p *Player) GetItemFromBackpackByIndex(itemType BackpackItemType, index int) any {
	var targetList *list.List

	switch itemType {
	case BackpackItemTypeWeapon:
		targetList = p.Backpack.Weapons
	case BackpackItemTypeFood:
		targetList = p.Backpack.Foods
	case BackpackItemTypeElixir:
		targetList = p.Backpack.Elixirs
	case BackpackItemTypeScroll:
		targetList = p.Backpack.Scrolls
	default:
		return nil
	}

	if index < 0 || index >= targetList.Len() {
		return nil
	}

	i := 0
	for elem := targetList.Front(); elem != nil; elem = elem.Next() {
		if i == index {
			return elem.Value
		}
		i++
	}

	return nil
}

func (p *Player) UseItemFromBackpack(item any) error {
	if item == nil {
		return items.ItemIsNotInBackpackError{}
	}

	// Обрабатываем оружие отдельно
	if weapon, ok := item.(*items.Weapon); ok {
		if err := p.EquipWeapon(weapon); err != nil {
			return err
		}
	} else if usableItem, ok := item.(items.Usable); ok {
		p.Character.Use(usableItem)
	}

	return p.Backpack.RemoveItem(item)
}
