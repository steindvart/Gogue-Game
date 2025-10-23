package entity

import "math/rand"

type Passage struct {
	DoorOne Point2D[int]
	Passage []Point2D[int]
	DoorTwo Point2D[int]
}

func (p *Passage) NewPassageX(doorOne Point2D[int], doorTwo Point2D[int]) {
	// если будут введены координаты дверей справа налево, а не слева направо, я их переставляю будто слева направо
	if doorOne.X > doorTwo.X {
		doorTemp := doorOne
		doorOne = doorTwo
		doorTwo = doorTemp
	}

	p.DoorOne = doorOne
	p.DoorTwo = doorTwo

	p.Passage = append(p.Passage, Point2D[int]{X: p.DoorOne.X + 1, Y: p.DoorOne.Y})

	pontKinkX := rand.Intn(doorTwo.X-doorOne.X) + doorOne.X
	lastElement := p.Passage[len(p.Passage)-1]
	for x := lastElement.X; x < pontKinkX; x++ {
		lastElement = p.Passage[len(p.Passage)-1]
		p.Passage = append(p.Passage, Point2D[int]{X: lastElement.X + 1, Y: lastElement.Y})
	}

	lastElement = p.Passage[len(p.Passage)-1]
	if p.DoorOne.Y < p.DoorTwo.Y {
		for y := lastElement.Y; y <= p.DoorTwo.Y; y++ {
			lastElement = p.Passage[len(p.Passage)-1]
			p.Passage = append(p.Passage, Point2D[int]{X: lastElement.X, Y: lastElement.Y + 1})
		}
	} else if p.DoorTwo.Y > p.DoorOne.Y {
		for y := lastElement.Y; y <= p.DoorTwo.Y; y-- {
			lastElement = p.Passage[len(p.Passage)-1]
			p.Passage = append(p.Passage, Point2D[int]{X: lastElement.X, Y: lastElement.Y - 1})
		}
	}

	lastElement = p.Passage[len(p.Passage)-1]
	for x := lastElement.X; x < doorTwo.X; x++ {
		lastElement = p.Passage[len(p.Passage)-1]
		p.Passage = append(p.Passage, Point2D[int]{X: lastElement.X + 1, Y: lastElement.Y})
	}
}
