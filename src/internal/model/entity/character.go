package entity

import (
	"gogue/internal/model/primitive"
	"math/rand"
)

type Character struct {
	Shape      primitive.Box
	Attributes primitive.Attributes
}

func (c *Character) IsAlive() bool {
	return c.Attributes.Health > 0
}

func (c *Character) Move(delta primitive.Point2D[int]) {
	c.Shape.Move(delta)
}

func (c *Character) TakeDamage(damage float64) {
	c.Attributes.Health -= damage
	if c.Attributes.Health < 0 {
		c.Attributes.Health = 0
	}
}

func (c *Character) Heal(amount float64) {
	c.Attributes.Health += amount
	if c.Attributes.Health > c.Attributes.MaxHealth {
		c.Attributes.Health = c.Attributes.MaxHealth
	}
}

func (c *Character) Attack() float64 {
	return c.Attributes.Strength
}

// Шанс уклонения = 1 - 1/(1 + Agility/scale).
// Растёт с увеличением ловкости, но никогда не достигает 100%.
// scale регулирует скорость роста. Это обеспечивает баланс между ростом шанса и невозможностью абсолютного уклонения.
// @todo 1 - сделать настраиваемым scale? Например, для регулировки сложности игры?
// @todo 2 - сделать сравнение с учётом ловкости атакующего?
func (c *Character) CheckEvasion() bool {
	if c.Attributes.Agility == 0 {
		return false
	}

	const scale = 20.0
	chance := 1.0 - 1.0/(1.0+c.Attributes.Agility/scale)

	// @todo (Copilot):
	// Using the global rand.Float64() makes the function non-deterministic and difficult to test.
	// Consider accepting a *rand.Rand parameter or using a seeded random generator for better testability.
	roll := rand.Float64() // Случайное дробное число - [0,1)
	return roll < chance
}
