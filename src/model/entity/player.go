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
	p.Backpack.Treasures += t.Value
}

func (p *Player) TakeItem(i ItemLike) error {
	return p.Backpack.AddItem(i)
}

func (p *Player) DropItem(i ItemLike) error {
	return p.Backpack.RemoveItem(i, p.Character.Shape)
}

func (p *Player) UseItem(i ItemLike) (string, error) {
	w, ok := i.(*Weapon)
	if ok {
		return w.Use(p), nil
	}

	err := p.Backpack.RemoveItem(i, p.Character.Shape)
	if err == nil {
		return i.Use(p), nil
	}

	return "", ItemIsNotInBackpackError{}
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
