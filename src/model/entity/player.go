package entity

type BackpackIsFullError struct{}

func (BackpackIsFullError) Error() string {
	return "backpack is full, drop something"
}

type Player struct {
	Character      Character
	Experience     uint
	CharacterLevel uint
	Backpack       *Backpack
	Weapon         *Weapon
}

func (p *Player) IsAlive() bool {
	return p.Character.IsAlive()
}

func (p *Player) Move(delta Point2D[int]) {
	p.Character.Move(delta)
}

func (p *Player) TakeDamage(damage float64) {
	p.Character.TakeDamage(damage)
}

func (p *Player) Heal(amount float64) {
	p.Character.Heal(amount)
}

func (p *Player) Attack() uint {
	damage := p.Character.Attack()
	if p.Weapon != nil {
		damage += p.Weapon.Damage
	}
	return damage
}

func (p *Player) CheckEvasion() bool {
	return p.Character.CheckEvasion()
}

func (p *Player) TakeTreasure(treasure *Treasure) {
	p.Backpack.Treasures += treasure.Value
}

func (p *Player) TakeConsumableLike(consumableLike ConsumableLike) error {
	if p.Backpack.ItemsNum < p.Backpack.Capacity {
		p.Backpack.Consumables = append(p.Backpack.Consumables, consumableLike)
		p.Backpack.ItemsNum++

		consumableLike.Taken()
		return nil
	} else {
		return BackpackIsFullError{}
	}
}

func (p *Player) DropConsumableLike(consumableLike ConsumableLike, box Box) ConsumableLike {
	for idx, item := range p.Backpack.Consumables {
		if item == consumableLike {
			p.Backpack.Consumables = append(p.Backpack.Consumables[:idx], p.Backpack.Consumables[idx+1:]...)
			p.Backpack.ItemsNum--
		}
	}
	return consumableLike.Dropped(box)
}

func (p *Player) UseConsumableLike(consumableLike ConsumableLike) (string, ConsumableLike) {
	for idx, item := range p.Backpack.Consumables {
		if item == consumableLike {
			p.Backpack.Consumables = append(p.Backpack.Consumables[:idx], p.Backpack.Consumables[idx+1:]...)
			p.Backpack.ItemsNum--
		}
	}
	return consumableLike.Use(p)
}

func NewPlayer(box Box) *Player {
	return &Player{
		Character: Character{
			Shape:     box,
			Health:    float64(AttributeRateAverage),
			MaxHealth: float64(AttributeRateAverage),
			Strength:  uint(AttributeRateAverage),
			Agility:   uint(AttributeRateAverage),
		},
		Experience:     0,
		CharacterLevel: 1,
		Backpack:       NewBackpack(),
		Weapon:         nil,
	}
}
