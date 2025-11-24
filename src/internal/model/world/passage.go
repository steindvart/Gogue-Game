package world

import (
	"errors"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Passage struct {
	Way     []primitives.Point2D[int]
	DoorOne primitives.Point2D[int]
	DoorTwo primitives.Point2D[int]
}

func NewPassageOnX(doorOne primitives.Point2D[int], doorTwo primitives.Point2D[int], random utils.RandomSource) (*Passage, error) {
	if doorOne.X == doorTwo.X {
		return nil, errors.New("doors cannot be positioned on the same x axis")
	}

	if doorOne.X > doorTwo.X {
		doorTemp := doorOne
		doorOne = doorTwo
		doorTwo = doorTemp
	}

	// прокладываю первый кубик(?) тоннеля от двери
	passage := []primitives.Point2D[int]{
		{X: doorOne.X + 1, Y: doorOne.Y},
	}

	// pointKinkOnX - рандомная точка, на которой тоннель свернёт
	// иду по прямой до этой точки
	pointKinkOnX := random.Intn(doorTwo.X-doorOne.X-1) + doorOne.X
	lastElement := passage[len(passage)-1]
	for x := lastElement.X; x < pointKinkOnX; x++ {
		lastElement = passage[len(passage)-1]
		passage = append(passage, primitives.Point2D[int]{X: lastElement.X + 1, Y: lastElement.Y})
	}

	// поднимаюсь или опускаюсь, пока не буду на одной линии с финишной дверью
	lastElement = passage[len(passage)-1]
	if doorOne.Y < doorTwo.Y {
		for y := lastElement.Y; y < doorTwo.Y; y++ {
			lastElement = passage[len(passage)-1]
			passage = append(passage, primitives.Point2D[int]{X: lastElement.X, Y: lastElement.Y + 1})
		}
	} else if doorOne.Y > doorTwo.Y {
		for y := lastElement.Y; y > doorTwo.Y; y-- {
			lastElement = passage[len(passage)-1]
			passage = append(passage, primitives.Point2D[int]{X: lastElement.X, Y: lastElement.Y - 1})
		}
	}

	// иду по прямой до финишной двери
	lastElement = passage[len(passage)-1]
	for x := lastElement.X; x < doorTwo.X-1; x++ {
		lastElement = passage[len(passage)-1]
		passage = append(passage, primitives.Point2D[int]{X: lastElement.X + 1, Y: lastElement.Y})
	}

	return &Passage{DoorOne: doorOne, Way: passage, DoorTwo: doorTwo}, nil
}

func NewPassageOnY(doorOne primitives.Point2D[int], doorTwo primitives.Point2D[int], random utils.RandomSource) (*Passage, error) {
	if doorOne.Y == doorTwo.Y {
		return nil, errors.New("doors cannot be positioned on the same y axis")
	}

	if doorOne.Y > doorTwo.Y {
		doorTemp := doorOne
		doorOne = doorTwo
		doorTwo = doorTemp
	}

	// прокладываю первый кубик(?) тоннеля от двери
	passage := []primitives.Point2D[int]{
		{X: doorOne.X, Y: doorOne.Y + 1},
	}

	// pointKinkOnX - рандомная точка, на которой тоннель свернёт
	// иду по прямой до этой точки
	pointKinkOnY := random.Intn(doorTwo.Y-doorOne.Y-1) + doorOne.Y
	lastElement := passage[len(passage)-1]
	for y := lastElement.Y; y < pointKinkOnY; y++ {
		lastElement = passage[len(passage)-1]
		passage = append(passage, primitives.Point2D[int]{X: lastElement.X, Y: lastElement.Y + 1})
	}

	// двигаюсь влево или вправо, пока не буду на одной линии с финишной дверью
	lastElement = passage[len(passage)-1]
	if doorOne.X < doorTwo.X {
		for x := lastElement.X; x < doorTwo.X; x++ {
			lastElement = passage[len(passage)-1]
			passage = append(passage, primitives.Point2D[int]{X: lastElement.X + 1, Y: lastElement.Y})
		}
	} else if doorOne.X > doorTwo.X {
		for x := lastElement.X; x > doorTwo.X; x-- {
			lastElement = passage[len(passage)-1]
			passage = append(passage, primitives.Point2D[int]{X: lastElement.X - 1, Y: lastElement.Y})
		}
	}

	// иду по прямой до финишной двери
	lastElement = passage[len(passage)-1]
	for y := lastElement.Y; y < doorTwo.Y-1; y++ {
		lastElement = passage[len(passage)-1]
		passage = append(passage, primitives.Point2D[int]{X: lastElement.X, Y: lastElement.Y + 1})
	}

	return &Passage{DoorOne: doorOne, Way: passage, DoorTwo: doorTwo}, nil
}
