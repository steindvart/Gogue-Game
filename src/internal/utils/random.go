package utils

import (
	"math"
	"math/rand"
)

type Randomizer interface {
	Intn(n int) int
	Float64() float64
	Perm(n int) []int
	Shuffle(n int, swap func(i, j int))
	Uint32() uint32
	Seed(seed int64)
}

type Random struct {
	rng *rand.Rand
}

func NewRandomWithSeed(seed int64) *Random {
	return &Random{
		rng: rand.New(rand.NewSource(seed)),
	}
}

func (r *Random) Intn(n int) int {
	return r.rng.Intn(n)
}

func (r *Random) Float64() float64 {
	return r.rng.Float64()
}

func (r *Random) Perm(n int) []int {
	return r.rng.Perm(n)
}

func (r *Random) Shuffle(n int, swap func(i, j int)) {
	r.rng.Shuffle(n, swap)
}

func (r *Random) Uint32() uint32 {
	return r.rng.Uint32()
}

func (r *Random) Seed(seed int64) {
	r.rng.Seed(seed)
}

func RandomFloatInRange(random Randomizer, min, max float64) float64 {
	return min + random.Float64()*(max-min)
}

func RandomRoundedFloatInRange(random Randomizer, min, max float64) int {
	return int(math.Round(RandomFloatInRange(random, min, max)))
}

func RandomIntInRange(random Randomizer, min, max int) int {
	return min + random.Intn(max-min+1)
}

func GetRandomElement[T any](random Randomizer, slice []T) T {
	idx := random.Intn(len(slice))
	return slice[idx]
}

// WeightedChoice - обобщённый элемент со значением и весом.
type WeightedChoice[T any] struct {
	Value  T
	Weight int
}

// GetWeightedRandomElement выбирает случайный элемент из слайса
// с учётом весов. Веса задают относительную вероятность выбора.
// Паникует, если суммарный вес <= 0 или слайс пуст.
func GetWeightedRandomElement[T any](random Randomizer, choices []WeightedChoice[T]) T {
	totalWeight := 0
	for _, c := range choices {
		totalWeight += c.Weight
	}

	roll := random.Intn(totalWeight)
	cumulative := 0
	for _, c := range choices {
		cumulative += c.Weight
		if roll < cumulative {
			return c.Value
		}
	}

	// Не должно произойти при totalWeight > 0
	return choices[len(choices)-1].Value
}
