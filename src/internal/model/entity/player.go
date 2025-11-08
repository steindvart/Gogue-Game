package entity

import (
	"gogue/internal/model/items"
	"gogue/internal/model/primitive"
)

type Player struct {
	Character      Character
	Experience     uint
	CharacterLevel uint
	Backpack       *Backpack
	Weapon         *items.Weapon
}

func (p *Player) IsAlive() bool {
	return p.Character.IsAlive()
}

func (p *Player) Move(delta primitive.Point2D[int]) {
	p.Character.Move(delta)
}

func (p *Player) TakeDamage(damage float64) {
	p.Character.TakeDamage(damage)
}

func (p *Player) Heal(amount float64) {
	p.Character.Heal(amount)
}

func (p *Player) Attack() float64 {
	damage := p.Character.Attack()
	if p.Weapon != nil {
		damage += p.Weapon.Damage
	}
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

func NewPlayer(box primitive.Box) *Player {
	return &Player{
		Character: Character{
			Shape: box,
			Attributes: primitive.Attributes{
				Health:    float64(AttributeRateAverage),
				MaxHealth: float64(AttributeRateAverage),
				Strength:  float64(AttributeRateAverage),
				Agility:   float64(AttributeRateAverage),
			},
		},
		Experience:     0,
		CharacterLevel: 1,
		Backpack:       NewBackpack(),
		Weapon:         nil,
	}
}
