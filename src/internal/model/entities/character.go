package entities

import (
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Character struct {
	*primitives.Box
	*primitives.Attributes
	TemporaryEffects []*primitives.Effect
}

func NewCharacter(box primitives.Box, attrs primitives.Attributes) *Character {
	return &Character{
		Box:              &box,
		Attributes:       &attrs,
		TemporaryEffects: nil,
	}
}

func (c *Character) GetPosition() primitives.Point2D[int] {
	return c.Box.Point
}

func (c *Character) SetPosition(pos primitives.Point2D[int]) {
	c.Box.Point = pos
}

func (c *Character) IsAlive() bool {
	return c.Attributes.Health > 0
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

func (c *Character) Use(usable items.Usable) {
	c.ApplyEffect(usable.Use())
}

func (c *Character) ProcessTurns(turns uint32) {
	c.ProcessTemporaryEffects(turns)
}

func (c *Character) ProcessTemporaryEffects(turns uint32) {
	// Важно начинать с конца среза, чтобы при удалении не сбивался индекс.
	for i := len(c.TemporaryEffects) - 1; i >= 0; i-- {
		e := c.TemporaryEffects[i]
		if e.Duration.Steps > turns {
			e.Duration.Steps -= turns
		} else {
			e.Duration.Steps = 0
		}

		if e.Duration.Steps == 0 {
			c.removeTemporaryEffectByIndex(i)
		}
	}
}

func (c *Character) ApplyEffect(effect *primitives.Effect) {
	if effect == nil {
		return
	}

	if effect.Duration.Type == primitives.EffectDurationTypeAllTemporary ||
		effect.Duration.Type == primitives.EffectDurationTypeAllTemporaryHealPermanent {
		c.TemporaryEffects = append(c.TemporaryEffects, effect)
	}

	// При изменении максимального здоровья, текущее здоровье также изменяется на ту же величину.
	if effect.Attributes.MaxHealth != 0 {
		effect.Attributes.Health = effect.Attributes.MaxHealth
	}

	c.Attributes.Affect(effect.Attributes)

	// Если здоровье превысило максимальное, устанавливаем его в максимум.
	// Если здоровье стало отрицательным, устанавливаем его в 1, чтобы персонаж не умер от эффекта.
	if c.Attributes.Health > c.Attributes.MaxHealth {
		c.Attributes.Health = c.Attributes.MaxHealth
	} else if c.Attributes.Health <= 0 {
		c.Attributes.Health = 1
	}
}

func (c *Character) RemoveTemporaryEffect(effect *primitives.Effect) {
	if effect == nil {
		return
	}

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
func (c *Character) CheckEvasion(rnd utils.Randomizer) bool {
	if c.Attributes.Agility == 0 {
		return false
	}

	const scale = 20.0
	chance := 1.0 - 1.0/(1.0+c.Attributes.Agility/scale)

	roll := rnd.Float64() // Случайное дробное число - [0,1)
	return roll < chance
}
