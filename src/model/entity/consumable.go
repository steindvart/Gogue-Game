package entity

import (
	"math/rand"
)

const (
	IncreaseAttributeBaseParcentage int = 5
	IncreaseAttributeMaxPercentage  int = 20
)

type ConsumableLike interface {
	Taken()
	Dropped(box Box) ConsumableLike
	Use(p *Player) (string, ConsumableLike)
}

type Consumable struct {
	Shape Box
	Name  string
}

type Attributes struct {
	MaxHealth, Agility, Strength uint
}

func (c *Consumable) Taken() {
	c.Shape = Box{}
}

func (c *Consumable) Dropped(box Box) {
	c.Shape = box
}

func (a *Attributes) GetAffectedAttributeName() string {
	if a.MaxHealth == 1 && a.Agility == 0 && a.Strength == 0 {
		return "MaxHealth"
	} else if a.Agility == 1 && a.MaxHealth == 0 && a.Strength == 0 {
		return "Agility"
	} else if a.Strength == 1 && a.MaxHealth == 0 && a.Agility == 0 {
		return "Strength"
	} else {
		return "None"
	}
}

func getAttributeRandomPercentIncrease() uint {
	return uint(IncreaseAttributeBaseParcentage + rand.Intn(IncreaseAttributeMaxPercentage+1))
}

func getAttributeRandomName(names []string) string {
	return names[rand.Intn(len(names))]
}

func getRandomAttribute() Attributes {
	attributesNum := 3
	attribute := rand.Intn(attributesNum)

	attributes := Attributes{}

	switch attribute {
	case 0:
		attributes.MaxHealth = 1
	case 1:
		attributes.Agility = 1
	case 2:
		attributes.Strength = 1
	}

	return attributes
}
