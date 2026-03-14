package entities

import (
	"math"

	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type CharacterProvider interface {
	GetCharacter() *Character
}

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

func (c *Character) GetCharacter() *Character {
	return c
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

// damageVariation - максимальное отклонение урона вниз (30%).
// Например, при Strength = 10 урон будет в диапазоне [7, 10].
const damageVariation = 0.3

// applyDamageVariation применяет случайное отклонение к базовому урону.
// Урон может быть снижен до (1 - damageVariation) * baseDamage, но не увеличен.
// Результат округляется по математическим правилам (math.Round).
func applyDamageVariation(baseDamage float64, rnd utils.Randomizer) float64 {
	if baseDamage <= 0 {
		return 0
	}

	// Множитель в диапазоне [1 - damageVariation, 1.0], т.е. [0.8, 1.0]
	multiplier := 1.0 - rnd.Float64()*damageVariation
	return math.Round(baseDamage * multiplier)
}

func (c *Character) MakeDamage(rnd utils.Randomizer) float64 {
	return applyDamageVariation(c.Attributes.Strength, rnd)
}

func (c *Character) Attack(defender *Character, rnd utils.Randomizer) AttackResult {
	result := AttackResult{
		DefenderMaxHealth: defender.Attributes.MaxHealth,
	}

	if defender.CheckEvasion(c.Attributes.Agility, rnd) {
		result.Evaded = true
		result.DefenderHealthAfter = defender.Attributes.Health
		return result
	}

	damage := c.MakeDamage(rnd)
	defender.TakeDamage(damage)

	result.Damage = damage
	result.DefenderHealthAfter = defender.Attributes.Health
	result.DefenderKilled = !defender.IsAlive()

	return result
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
		if e.Duration.Turns > turns {
			e.Duration.Turns -= turns
		} else {
			e.Duration.Turns = 0
		}

		if e.Duration.Turns == 0 {
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

	// Если какой-то из атрибутов стал отрицательным, устанавливаем его в 0.
	if c.Attributes.Strength < 0 {
		c.Attributes.Strength = 0
	} else if c.Attributes.Agility < 0 {
		c.Attributes.Agility = 0
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

// CheckEvasion вычисляет шанс уклонения защитника с учётом ловкости атакующего.
//
// Формула: chance = 1 - 1/(1 + effectiveAgility/scale),
// где effectiveAgility = max(0, defenderAgility - attackerAgility * attackerWeight).
//
// attackerWeight (0.5) определяет, насколько сильно ловкость атакующего
// снижает эффективность уклонения. Значение < 1 означает, что уклоняться
// легче, чем попадать по ловкому противнику - это даёт преимущество защитнику.
//
// scale (50) регулирует скорость роста шанса уклонения.
// Шанс растёт с увеличением ловкости, но никогда не достигает 100%.
// @todo - сделать настраиваемым scale, чтобы ещё лучше настраивать сложность игры?
func (c *Character) CheckEvasion(attackerAgility float64, rnd utils.Randomizer) bool {
	if c.Attributes.Agility <= 0 {
		return false
	}

	const (
		scale          = 50.0
		attackerWeight = 0.5 // Ловкость атакующего влияет в 2 раза слабее
	)

	effectiveAgility := c.Attributes.Agility - attackerAgility*attackerWeight
	if effectiveAgility <= 0 {
		return false
	}

	chance := 1.0 - 1.0/(1.0+effectiveAgility/scale)

	roll := rnd.Float64() // Случайное дробное число - [0,1)
	return roll < chance
}
