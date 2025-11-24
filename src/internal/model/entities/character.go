package entities

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Character struct {
	Shape            *primitives.Box
	Attributes       primitives.Attributes
	TemporaryEffects []*primitives.Effect
}

func NewCharacter(box primitives.Box, attrs primitives.Attributes) *Character {
	return &Character{
		Shape:            &box,
		Attributes:       attrs,
		TemporaryEffects: nil,
	}
}

func (c *Character) IsAlive() bool {
	return c.Attributes.Health > 0
}

func (c *Character) Move(delta primitives.Point2D[int]) {
	c.Shape.Move(delta)
}

func (c *Character) TakeDamage(damage float64) {
	c.Attributes.Health -= damage
	if c.Attributes.Health < 0 {
		c.Attributes.Health = 0
	}
}

func (c *Character) Attack() float64 {
	return c.Attributes.Strength
}

type Usable interface {
	Use() *primitives.Effect
}

func (c *Character) Use(usable Usable) {
	c.ApplyEffect(usable.Use())
}

func (c *Character) ProcessTemporaryEffects(steps uint32) {
	// Важно начинать с конца среза, чтобы при удалении не сбивался индекс.
	for i := len(c.TemporaryEffects) - 1; i >= 0; i-- {
		e := c.TemporaryEffects[i]
		if e.Duration.Steps > steps {
			e.Duration.Steps -= steps
		} else {
			e.Duration.Steps = 0
		}

		if e.Duration.Steps == 0 {
			c.removeTemporaryEffectByIndex(i)
		}
	}
}

func (c *Character) ApplyEffect(effect *primitives.Effect) {
	// @todo - сделать обработку nil значений

	if effect.Duration.Type == primitives.EffectDurationTypeAllTemporary ||
		effect.Duration.Type == primitives.EffectDurationTypeAllTemporaryHealPermanent {
		c.TemporaryEffects = append(c.TemporaryEffects, effect)
	}

	// Apply attribute changes (permanent or immediate part of temporary)
	c.Attributes.Affect(effect.Attributes)

	if c.Attributes.Health > c.Attributes.MaxHealth {
		c.Attributes.Health = c.Attributes.MaxHealth
	}
}

func (c *Character) RemoveTemporaryEffect(effect *primitives.Effect) {
	// @todo - сделать обработку nil значений

	for i, e := range c.TemporaryEffects {
		if e == effect {
			c.removeTemporaryEffectByIndex(i)
			return
		}
	}
}

func (c *Character) removeTemporaryEffectByIndex(idx int) {
	e := c.TemporaryEffects[idx]
	c.TemporaryEffects = append(c.TemporaryEffects[:idx], c.TemporaryEffects[idx+1:]...)

	// Если эффект временный, но воздействие на здоровье было мгновенным, то не отменяем его.
	// Например, зелье лечения с мгновенным восстановлением здоровья, но временным увеличением силы.
	// При снятии эффекта здоровье не должно уменьшаться.
	if e.Duration.Type == primitives.EffectDurationTypeAllTemporaryHealPermanent {
		e.Attributes.Health = 0
	}

	c.Attributes.Affect(primitives.Inverse(e.Attributes))
}

// Шанс уклонения = 1 - 1/(1 + Agility/scale).
// Растёт с увеличением ловкости, но никогда не достигает 100%.
// scale регулирует скорость роста. Это обеспечивает баланс между ростом шанса и невозможностью абсолютного уклонения.
// @todo 1 - сделать настраиваемым scale? Например, для регулировки сложности игры?
// @todo 2 - сделать сравнение с учётом ловкости атакующего?
func (c *Character) CheckEvasion(rnd utils.RandomSource) bool {
	if c.Attributes.Agility == 0 {
		return false
	}

	const scale = 20.0
	chance := 1.0 - 1.0/(1.0+c.Attributes.Agility/scale)

	roll := rnd.Float64() // Случайное дробное число - [0,1)
	return roll < chance
}
