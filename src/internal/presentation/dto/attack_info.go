package dto

import "gogue/internal/model/entities"

// AttackInfo - DTO для передачи информации об атаке в view-слой.
type AttackInfo struct {
	AttackerName        string
	DefenderName        string
	Damage              float64
	Evaded              bool
	DefenderHealthAfter float64
	DefenderMaxHealth   float64
	DefenderKilled      bool

	// Guaranteed - атака была гарантированной (без проверки evasion).
	Guaranteed bool
	// AppliedStun - атака наложила оглушение.
	AppliedStun bool

	// Loot - информация о добыче при убийстве врага (опционально).
	LootName  string
	LootValue int32
}

func ConvertAttackResultToDto(result *entities.AttackResult) *AttackInfo {
	if result == nil {
		return nil
	}

	return &AttackInfo{
		AttackerName:        result.AttackerName,
		DefenderName:        result.DefenderName,
		Damage:              result.Damage,
		Evaded:              result.Evaded,
		DefenderHealthAfter: result.DefenderHealthAfter,
		DefenderMaxHealth:   result.DefenderMaxHealth,
		DefenderKilled:      result.DefenderKilled,
		Guaranteed:          result.Guaranteed,
		AppliedStun:         result.AppliedStun,
	}
}

func ConvertAttackResultsToDto(results []entities.AttackResult) []*AttackInfo {
	infos := make([]*AttackInfo, 0, len(results))
	for i := range results {
		infos = append(infos, ConvertAttackResultToDto(&results[i]))
	}
	return infos
}
