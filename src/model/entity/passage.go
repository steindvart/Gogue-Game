package entity

import "math/rand"

type Passage struct {
	DoorOne Point2D[int]
	Passage []Point2D[int]
	DoorTwo Point2D[int]
}

func NewPassageOnX(doorOne Point2D[int], doorTwo Point2D[int]) *Passage {
	// если будут введены координаты дверей справа налево, а не слева направо, я их переставляю будто слева направо
	if doorOne.X > doorTwo.X {
		doorTemp := doorOne
		doorOne = doorTwo
		doorTwo = doorTemp
	}

	var passage []Point2D[int]
	passage = append(passage, Point2D[int]{X: doorOne.X + 1, Y: doorOne.Y})

	pontKinkOnX := rand.Intn(doorTwo.X-doorOne.X-1) + doorOne.X
	lastElement := passage[len(passage)-1]
	for x := lastElement.X; x < pontKinkOnX; x++ {
		lastElement = passage[len(passage)-1]
		passage = append(passage, Point2D[int]{X: lastElement.X + 1, Y: lastElement.Y})
	}

	lastElement = passage[len(passage)-1]
	if doorOne.Y < doorTwo.Y {
		for y := lastElement.Y; y < doorTwo.Y; y++ {
			lastElement = passage[len(passage)-1]
			passage = append(passage, Point2D[int]{X: lastElement.X, Y: lastElement.Y + 1})
		}
	} else if doorOne.Y > doorTwo.Y {
		for y := lastElement.Y; y > doorTwo.Y; y-- {
			lastElement = passage[len(passage)-1]
			passage = append(passage, Point2D[int]{X: lastElement.X, Y: lastElement.Y - 1})
		}
	}

	lastElement = passage[len(passage)-1]
	for x := lastElement.X; x < doorTwo.X; x++ {
		lastElement = passage[len(passage)-1]
		passage = append(passage, Point2D[int]{X: lastElement.X + 1, Y: lastElement.Y})
	}

	return &Passage{DoorOne: doorOne, Passage: passage, DoorTwo: doorTwo}
}

func NewPassageOnY(doorOne Point2D[int], doorTwo Point2D[int]) *Passage {
	if doorOne.Y > doorTwo.Y {
		doorTemp := doorOne
		doorOne = doorTwo
		doorTwo = doorTemp
	}

	var passage []Point2D[int]
	passage = append(passage, Point2D[int]{X: doorOne.X, Y: doorOne.Y + 1})

	pointKinkOnY := rand.Intn(doorTwo.Y-doorOne.Y-1) + doorOne.Y
	lastElement := passage[len(passage)-1]
	for y := lastElement.Y; y < pointKinkOnY; y++ {
		lastElement = passage[len(passage)-1]
		passage = append(passage, Point2D[int]{X: lastElement.X, Y: lastElement.Y + 1})
	}

	lastElement = passage[len(passage)-1]
	if doorOne.X < doorTwo.X {
		for x := lastElement.X; x < doorTwo.X; x++ {
			lastElement = passage[len(passage)-1]
			passage = append(passage, Point2D[int]{X: lastElement.X + 1, Y: lastElement.Y})
		}
	} else if doorOne.X > doorTwo.X {
		for x := lastElement.X; x > doorTwo.X; x-- {
			lastElement = passage[len(passage)-1]
			passage = append(passage, Point2D[int]{X: lastElement.X - 1, Y: lastElement.Y})
		}
	}

	lastElement = passage[len(passage)-1]
	for y := lastElement.Y; y < doorTwo.Y; y++ {
		lastElement = passage[len(passage)-1]
		passage = append(passage, Point2D[int]{X: lastElement.X, Y: lastElement.Y + 1})
	}

	return &Passage{DoorOne: doorOne, Passage: passage, DoorTwo: doorTwo}
}
