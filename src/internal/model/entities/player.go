package entities

import (
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
)

type Player struct {
	Character      Character
	Experience     uint
	CharacterLevel uint
	Backpack       *items.Backpack
	Weapon         *items.Weapon
}

func (p *Player) IsAlive() bool {
	return p.Character.IsAlive()
}

func (p *Player) Move(delta primitives.Point2D[int]) {
	p.Character.Move(delta)
}

func (p *Player) TakeDamage(damage float64) {
	p.Character.TakeDamage(damage)
}

func (p *Player) Heal(amount float64) {
	p.Character.Heal(amount)
}

func (p *Player) Attack() float64 {
	damage := 0.0
	if p.Weapon != nil {
		damage += p.Weapon.AffectedAttributes.Strength
	}

	damage += p.Character.Attack()
	return damage
}

func (p *Player) CheckEvasion() bool {
	return p.Character.CheckEvasion()
}

func (p *Player) TakeTreasure(t *items.Treasure) {
	p.Backpack.AddTreasure(t)
}

func (p *Player) TakeItemToBackpack(item any) error {
	return p.Backpack.AddItem(item)
}

// func (p *Player) DropItemFromBackpack(item any) error {
// 	return p.Backpack.RemoveItem(item, p.Character.Shape)
// }

// func (p *Player) UseItem(item any) error {
// 	if w := AsWeapon(item); w != nil {
// 		err := p.Backpack.weaponIsInBackpack(w)
// 		if err != nil {
// 			return err
// 		}

// 		if p.Weapon != nil && p.Weapon != w {
// 			_ = p.Backpack.RemoveItem(p.Weapon, p.Character.Shape)
// 		}

// 		items.Use(p)
// 		return nil
// 	}

// 	err := p.Backpack.RemoveItem(item, p.Character.Shape)
// 	if err == nil {
// 		items.Use(p)
// 		return nil
// 	}

// 	return err
// }

func NewPlayer(box primitives.Box) *Player {
	return &Player{
		Character: Character{
			Shape: box,
			Attributes: primitives.Attributes{
				Health:    100,
				MaxHealth: 100,
				Strength:  10,
				Agility:   5,
			},
		},
		Experience:     0,
		CharacterLevel: 1,
		Backpack:       items.NewBackpack(),
		Weapon:         nil,
	}
}
