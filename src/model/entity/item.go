package entity

import (
	"math/rand"
)

const (
	IncreaseAttributeBaseParcentage int = 5
	IncreaseAttributeMaxPercentage  int = 20
	AttributesNum                   int = 3
)

type AttributeType int

const (
	AttributeTypeMaxHealth AttributeType = iota
	AttributeTypeAgility
	AttributeTypeStrength
)

type ItemLike interface {
	Takeable
	Droppable
	Usable
}

type Takeable interface {
	Take()
}

type Droppable interface {
	Drop(box Box)
}

type Usable interface {
	Use(p *Player)
}

type Item struct {
	Shape Box
	Name  string
}

func (i *Item) Take() {
	i.Shape = Box{}
}

func (i *Item) Drop(box Box) {
	i.Shape = box
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
	attribute := rand.Intn(AttributesNum)
	attributes := Attributes{}

	switch attribute {
	case int(AttributeTypeMaxHealth):
		attributes.MaxHealth = 1
	case int(AttributeTypeAgility):
		attributes.Agility = 1
	case int(AttributeTypeStrength):
		attributes.Strength = 1
	}

	return attributes
}
