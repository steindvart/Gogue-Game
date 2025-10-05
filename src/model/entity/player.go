package entity

type Player struct {
	Character Character
	Backpack  *Backpack
	Weapon    *Weapon
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

func (p *Player) Heal(amount uint) {
	p.Character.Heal(amount)
}

func (p *Player) Attack() uint {
	damage := p.Character.Attack()
	// @todo - логика атаки с оружием
	// if p.Weapon != nil {
	// 	damage += p.Weapon.Damage
	// }

	return damage
}

func (p *Player) CheckEvasion() bool {
	return p.Character.CheckEvasion()
}
