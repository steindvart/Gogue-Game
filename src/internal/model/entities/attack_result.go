package entities

// AttackResult описывает результат одной атаки.
type AttackResult struct {
	// AttackerName - имя атакующего (например, "Player", "Zombie").
	AttackerName string
	// DefenderName - имя защищающегося.
	DefenderName string
	// Damage - нанесённый урон. 0, если атака была уклонена.
	Damage float64
	// Evaded - true, если защитник уклонился от атаки.
	Evaded bool
	// DefenderHealthAfter - здоровье защитника после атаки.
	DefenderHealthAfter float64
	// DefenderMaxHealth - максимальное здоровье защитника.
	DefenderMaxHealth float64
	// DefenderKilled - true, если защитник убит в результате атаки.
	DefenderKilled bool
	// Guaranteed - true, если атака была гарантированной (без проверки уклонения).
	Guaranteed bool
	// AppliedStun - true, если атака наложила оглушение на защитника.
	AppliedStun bool
}
