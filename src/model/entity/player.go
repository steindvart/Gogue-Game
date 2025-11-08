package entity

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

func (p *Player) TakeTreasure(t *Treasure) {
	p.Backpack.AddTreasure(t)
}

func (p *Player) TakeItem(item ItemLike) error {
	return p.Backpack.AddItem(item)
}

func (p *Player) DropItem(item ItemLike) error {
	return p.Backpack.RemoveItem(item, p.Character.Shape)
}

func (p *Player) UseItem(item ItemLike) error {
	if w := IsWeapon(item); w != nil {
		err := p.Backpack.weaponIsInBackpack(w)
		if err != nil {
			return err
		}

		if p.Weapon != nil && p.Weapon != w {
			_ = p.Backpack.RemoveItem(p.Weapon, p.Character.Shape)
		}

		item.Use(p)
		return nil
	}

	err := p.Backpack.RemoveItem(item, p.Character.Shape)
	if err == nil {
		item.Use(p)
		return nil
	}

	return err
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
