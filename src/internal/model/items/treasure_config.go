package items

import (
	"errors"
	"fmt"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type TreasureType string

const (
	TreasureTypeGold     TreasureType = "Gold"
	TreasureTypeGem      TreasureType = "Gem"
	TreasureTypeArtifact TreasureType = "Artifact"
	TreasureTypeMystery  TreasureType = "Mystery"
	TreasureTypeCustom   TreasureType = "Custom" // For custom treasure - dynamicly or from external data created
)

// TreasureValueRange is an integer range for treasure value
type TreasureValueRange = primitives.Range[int32]

type TreasureConfig struct {
	Type        TreasureType
	ValueRange  TreasureValueRange
	Description string
}

var TreasureRegistry = map[TreasureType]TreasureConfig{
	TreasureTypeGold: {
		Type:        TreasureTypeGold,
		ValueRange:  TreasureValueRange{Min: 5, Max: 20},
		Description: "Common gold coins",
	},
	TreasureTypeGem: {
		Type:        TreasureTypeGem,
		ValueRange:  TreasureValueRange{Min: 50, Max: 200},
		Description: "Precious gem",
	},
	TreasureTypeArtifact: {
		Type:        TreasureTypeArtifact,
		ValueRange:  TreasureValueRange{Min: 200, Max: 1000},
		Description: "Powerful ancient artifact",
	},
	TreasureTypeMystery: {
		Type:        TreasureTypeMystery,
		ValueRange:  TreasureValueRange{Min: 1, Max: 500},
		Description: "Unknown treasure with random value",
	},
}

func GetTreasureConfig(t TreasureType) TreasureConfig {
	if cfg, exists := TreasureRegistry[t]; exists {
		return cfg
	}
	return TreasureRegistry[TreasureTypeMystery]
}

func (cfg *TreasureConfig) GenerateValue(rng utils.Randomizer) int32 {
	return int32(utils.RandomIntInRange(rng, int(cfg.ValueRange.Min), int(cfg.ValueRange.Max)))
}

func (cfg *TreasureConfig) Validate() error {
	if cfg.ValueRange.Min > cfg.ValueRange.Max {
		return fmt.Errorf("invalid value range: min=%d > max=%d", cfg.ValueRange.Min, cfg.ValueRange.Max)
	}
	return nil
}

// getRandomTreasureType выбирает случайный тип сокровища по заданным процентам.
//
// Если сумма процентов не равна 100, проценты нормализуются пропорционально.
// Это гарантирует корректную работу даже при ошибках округления
// на стороне вызывающего кода.
func getRandomTreasureType(rnd utils.Randomizer, goldPercent, gemPercent, artifactPercent, mysteryPercent uint) (TreasureType, error) {
	sum := goldPercent + gemPercent + artifactPercent + mysteryPercent
	if sum == 0 {
		return TreasureTypeGold, errors.New("all treasure percentages are zero")
	}

	// Fail-safe: нормализуем к 100, если сумма отклонилась из-за округления
	if sum != 100 {
		gold64 := float64(goldPercent) * 100.0 / float64(sum)
		gem64 := float64(gemPercent) * 100.0 / float64(sum)
		artifact64 := float64(artifactPercent) * 100.0 / float64(sum)

		goldPercent = uint(gold64)
		gemPercent = uint(gem64)
		artifactPercent = uint(artifact64)
	}

	percent := rnd.Intn(100)
	switch {
	case percent < int(goldPercent):
		return TreasureTypeGold, nil
	case percent < int(goldPercent)+int(gemPercent):
		return TreasureTypeGem, nil
	case percent < int(goldPercent)+int(gemPercent)+int(artifactPercent):
		return TreasureTypeArtifact, nil
	}
	return TreasureTypeMystery, nil
}
