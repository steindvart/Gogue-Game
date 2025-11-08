package entity

import "math/rand"

type Attributes struct {
	MaxHealth, Agility, Strength uint
}

type CharacterLike interface {
	IsAlive() bool
	Move(delta Point2D[int])
	TakeDamage(damage float64)
	Heal(amount float64)
	Attack() uint
	CheckEvasion() bool
}

type Character struct {
	Shape     Box
	Health    float64
	MaxHealth float64
	Strength  uint
	Agility   uint
}

func (c *Character) IsAlive() bool {
	return c.Health > 0
}

func (c *Character) Move(delta Point2D[int]) {
	c.Shape.Move(delta)
}

func (c *Character) TakeDamage(damage float64) {
	c.Health -= damage
	if c.Health < 0 {
		c.Health = 0
	}
}

func (c *Character) Heal(amount float64) {
	c.Health += amount
	if c.Health > c.MaxHealth {
		c.Health = c.MaxHealth
	}
}

func (c *Character) Attack() uint {
	return c.Strength
}

// Шанс уклонения = 1 - 1/(1 + Agility/scale).
// Растёт с увеличением ловкости, но никогда не достигает 100%.
// scale регулирует скорость роста. Это обеспечивает баланс между ростом шанса и невозможностью абсолютного уклонения.
// @todo 1 - сделать настраиваемым scale? Например, для регулировки сложности игры?
// @todo 2 - сделать сравнение с учётом ловкости атакующего?
func (c *Character) CheckEvasion() bool {
	if c.Agility == 0 {
		return false
	}

	const scale = 20.0
	chance := 1.0 - 1.0/(1.0+float64(c.Agility)/scale)

	// @todo (Copilot):
	// Using the global rand.Float64() makes the function non-deterministic and difficult to test.
	// Consider accepting a *rand.Rand parameter or using a seeded random generator for better testability.
	roll := rand.Float64() // Случайное дробное число - [0,1)
	return roll < chance
}
