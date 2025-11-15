package utils

import (
	"math"
	"math/rand"
)

type EntityType int

const (
	EntityTypePlayer EntityType = iota + 1
	EntityTypeHorizontalWall
	EntityTypeVerticalWall
	EntityTypePortal
	EntityTypePassage
	EntityTypeDoorOne
	EntityTypeDoorTwo
	EntityTypeZombie
	EntityTypeVampire
	EntityTypeGhost
	EntityTypeOgre
	EntityTypeSnakeMage
)

type RandomSource interface {
	Intn(n int) int
	Float64() float64
	Perm(n int) []int
	Shuffle(n int, swap func(i, j int))
	Uint32() uint32
	Seed(seed int64)
}

type RandomGenerator struct {
	rng *rand.Rand
}

func NewRandomGeneratorWithSeed(seed int64) *RandomGenerator {
	return &RandomGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

func (r *RandomGenerator) Intn(n int) int {
	return r.rng.Intn(n)
}

func (r *RandomGenerator) Float64() float64 {
	return r.rng.Float64()
}

func (r *RandomGenerator) Perm(n int) []int {
	return r.rng.Perm(n)
}

func (r *RandomGenerator) Shuffle(n int, swap func(i, j int)) {
	r.rng.Shuffle(n, swap)
}

func (r *RandomGenerator) Uint32() uint32 {
	return r.rng.Uint32()
}

func (r *RandomGenerator) Seed(seed int64) {
	r.rng.Seed(seed)
}

func RandomFloatInRange(random RandomSource, min, max float64) float64 {
	return min + random.Float64()*(max-min)
}

func RandomRoundedFloatInRange(random RandomSource, min, max float64) int {
	return int(math.Round(RandomFloatInRange(random, min, max)))
}

func RandomIntInRange(random RandomSource, min, max int) int {
	return min + random.Intn(max-min+1)
}
