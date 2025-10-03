package entity

import (
	"time"
)

const RoomsInWidth, RoomsInHeight = 3, 3
const RoomsNum = RoomsInWidth * RoomsInHeight

const MaxPassageParts = 3
const MaxPassagesNum = RoomsNum - 1

const MaxMonstersPerRoom = 2
const MaxConsumablesPerRoom = 3

type LevelNum uint

const (
	One LevelNum = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Eleven
	Twelve
	Thirteen
	Fourteen
	Fifteen
	Sixteen
	Seventeen
	Eighteen
	Nineteen
	Twenty
	TwentyOne
)

type Dimension uint

const (
	X Dimension = iota
	Y
	CoordinatesNum
)

type Measure uint

const (
	Height Measure = iota
	Width
	MeasuresNum
)

type Statistics uint

const (
	HealthStat Statistics = iota
	AgilityStat
	StrengthStat
)

type EnemyType uint

const (
	Zombie EnemyType = iota
	Vampire
	Ghost
	Ogre
	Snake
)

type HostilityType uint

const (
	Low HostilityType = iota
	Average
	High
)

type DirectionType uint

const (
	Forward DirectionType = iota
	Back
	Left
	Right
	DiagonallyForwardLeft
	DiagonallyForwardRight
	DiagonallyBackLeft
	DiagonallyBackRight
	Stop
)

type SessionEntity struct {
	Treasures,
	Level,
	Enemies,
	Food,
	Lixirs,
	Scrolls,
	Attacks,
	Missed,
	Moves uint
}

type PlayerEntity struct {
	BaseCharacter CharacterEntity
	MaxHealth     uint
	Backpack      *BackpackEntity
	Weapon        WeaponEntity
}

type BackpackEntity struct {
	Size          uint
	FoodsOnHand   []FoodEntity
	ElixirsOnHand []ElixirEntity
	ScrollsOnHand []ScrollEntity
	WeaponOnHand  []WeaponEntity
	Treasures     []TreasureEntity
}

type LevelEntity struct {
	Geometry    ObjectEntity
	Rooms       [RoomsNum]RoomEntity
	Passages    [MaxPassagesNum]PassageEntity
	LevelNumber LevelNum
	LevelEnd    ObjectEntity
}

type RoomEntity struct {
	Geometry    ObjectEntity
	Consumables ConsumablesEntity
	Enemies     []EnemyEntity
}

type PassageEntity struct {
	Geometries []ObjectEntity
}

type ConsumablesEntity struct {
	Foods   []FoodEntity
	Elixirs []ElixirEntity
	Scrolls []ScrollEntity
	Weapons []WeaponEntity
}

type FoodEntity struct {
	Geometry           ObjectEntity
	HealthRegeneration uint
	Name               string
}

type ElixirEntity struct {
	Geometry       ObjectEntity
	EffectDuration time.Duration
	AffectedStat   Statistics
	Increment      uint
	Name           string
}

type ScrollEntity struct {
	Geometry     ObjectEntity
	AffectedStat Statistics
	Increment    uint
	Name         string
}

type WeaponEntity struct {
	Geometry     ObjectEntity
	StrengthBuff uint
	Name         string
}

type TreasureEntity struct {
	Value uint
	Name  string
}

type EnemyEntity struct {
	BaseCharacter CharacterEntity
	Type          EnemyType
	Hostility     HostilityType
	IsChasing     bool
	Direction     DirectionType
}

type CharacterEntity struct {
	Geometry ObjectEntity
	Health   float64
	Strength uint
	Agility  uint
}

type ObjectEntity struct {
	Coordinates [CoordinatesNum]int
	Sizes       [MeasuresNum]uint
}
