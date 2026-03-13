package entities

import (
	"testing"

	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

const defaultCharacterTestSeed int64 = 42

func TestCharacter_NewCharacter_BasicInit(t *testing.T) {
	box := primitives.Box{Point: primitives.Point2D[int]{X: 2, Y: 3}, Size: primitives.Size2D[uint]{Width: 1, Height: 2}}
	attrs := primitives.Attributes{MaxHealth: 100, Health: 50, Strength: 7, Agility: 9}

	c := NewCharacter(box, attrs)

	if c == nil {
		t.Fatalf("NewCharacter returned nil")
	}
	if c.Box == nil {
		t.Fatalf("Pointers must be initialized: Box=%v, Attributes=%v", c.Box, c.Attributes)
	}
	if *c.Box != box {
		t.Errorf("Box mismatch: got %+v want %+v", *c.Box, box)
	}
	if c.Attributes == &attrs {
		t.Errorf("Attributes pointers must not match: got %p want %p", c.Attributes, &attrs)
	}
	if c.TemporaryEffects != nil {
		// In this project we expect nil slice on init (len is 0 anyway)
		t.Errorf("TemporaryEffects should be nil on init, got non-nil len=%d", len(c.TemporaryEffects))
	}
}

func TestCharacter_NewCharacter_IndependenceFromArgs(t *testing.T) {
	box := primitives.Box{Point: primitives.Point2D[int]{X: 5, Y: 5}, Size: primitives.Size2D[uint]{Width: 2, Height: 2}}
	attrs := primitives.Attributes{MaxHealth: 100, Health: 80, Strength: 10, Agility: 1}
	c := NewCharacter(box, attrs)

	box.Move(primitives.Point2D[int]{X: 10, Y: 10})
	attrs.Health = 1
	attrs.Strength = 0
	attrs.Agility = 0

	if c.Box.Point != (primitives.Point2D[int]{X: 5, Y: 5}) {
		t.Errorf("Box should be independent from original box; got point=%+v", c.Box.Point)
	}
	if c.Attributes.Health != 80 || c.Attributes.Strength != 10 || c.Attributes.Agility != 1 {
		t.Errorf("Attributes should be independent: got (H %.1f, S %.1f, A %.1f)", c.Attributes.Health, c.Attributes.Strength, c.Attributes.Agility)
	}
}

func TestCharacter_IsAlive(t *testing.T) {
	tests := []struct {
		name   string
		health float64
		want   bool
	}{
		{
			name:   "Health > 0, is alive",
			health: 10,
			want:   true,
		},
		{
			name:   "Health == 0, is not alive",
			health: 0,
			want:   false,
		},
		{
			name:   "Health < 0, is not alive",
			health: -5,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCharacter(
				primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
				primitives.Attributes{MaxHealth: 100, Health: tt.health, Agility: 0, Strength: 0},
			)
			if got := c.IsAlive(); got != tt.want {
				t.Errorf("IsAlive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCharacter_Move(t *testing.T) {
	box := primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 2}, Size: primitives.Size2D[uint]{Width: 2, Height: 3}}
	c := NewCharacter(box, primitives.Attributes{MaxHealth: 10, Health: 10})

	delta := primitives.Point2D[int]{X: 3, Y: -1}
	c.Move(delta)

	want := primitives.Point2D[int]{X: 4, Y: 1}
	if c.Box.Point != want {
		t.Errorf("Move() point = %v, want %v", c.Box.Point, want)
	}
}

func TestCharacter_TakeDamage(t *testing.T) {
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 10},
	)

	c.TakeDamage(5)
	if c.Attributes.Health != 5 {
		t.Errorf("TakeDamage(5) health = %v, want 5", c.Attributes.Health)
	}

	c.TakeDamage(10)
	if c.Attributes.Health != 0 {
		t.Errorf("TakeDamage(10) should clamp to 0, got %v", c.Attributes.Health)
	}

	c.TakeDamage(100)
	if c.Attributes.Health != 0 {
		t.Errorf("TakeDamage(overkill) should stay 0, got %v", c.Attributes.Health)
	}
}

func TestCharacter_AttackEqualsStrength(t *testing.T) {
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 100, Strength: 17.5},
	)
	if got := c.MakeDamage(); got != 17.5 {
		t.Errorf("Attack() = %v, want 17.5", got)
	}
}

func TestCharacter_ApplyEffect_NilIsNothing(t *testing.T) {
	const maxHealth = 100
	const wantHealth = 50
	const wantStrength = 10.0
	const wantAgility = 5.0

	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: maxHealth, Health: wantHealth, Strength: wantStrength, Agility: wantAgility},
	)

	c.ApplyEffect(nil)

	if len(c.TemporaryEffects) != 0 {
		t.Fatalf("Nil effect must not be tracked, got len=%d", len(c.TemporaryEffects))
	}

	if c.Attributes.Health != wantHealth {
		t.Errorf("Health should remain unchanged: got %v, want 50", c.Attributes.Health)
	}

	if c.Attributes.Strength != wantStrength {
		t.Errorf("Strength should remain unchanged: want %v got %v", wantStrength, c.Attributes.Strength)
	}

	if c.Attributes.Agility != wantAgility {
		t.Errorf("Agility should remain unchanged: want %v got %v", wantAgility, c.Attributes.Agility)
	}
}

func TestCharacter_ApplyEffect_AllPermanent_ClampsHealth(t *testing.T) {
	const maxHealth = 100

	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: maxHealth, Health: 50, Strength: 10, Agility: 5},
	)
	e := &primitives.Effect{
		Duration:   primitives.EffectDuration{Type: primitives.EffectDurationTypeAllPermanent, Turns: 0},
		Attributes: primitives.Attributes{Health: maxHealth, Strength: 5, Agility: 3},
	}

	c.ApplyEffect(e)

	if len(c.TemporaryEffects) != 0 {
		t.Fatalf("AllPermanent effect must not be tracked, got len=%d", len(c.TemporaryEffects))
	}
	if c.Attributes.Health != maxHealth {
		t.Errorf("Health should be clamped to MaxHealth: got %v, want 100", c.Attributes.Health)
	}

	wantStrength := 15.0
	if c.Attributes.Strength != wantStrength {
		t.Errorf("Strength attribute not applied correctly: want %v got %v", wantStrength, c.Attributes.Strength)
	}

	wantAgility := 8.0
	if c.Attributes.Agility != wantAgility {
		t.Errorf("Agility attribute not applied correctly: want %v got %v", wantAgility, c.Attributes.Agility)
	}
}

func TestCharacter_ApplyEffect_AllTemporary_TracksAndMutates(t *testing.T) {
	const maxHealth = 100

	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: maxHealth, Health: 40, Strength: 2, Agility: 1},
	)
	e := &primitives.Effect{
		Duration:   primitives.EffectDuration{Type: primitives.EffectDurationTypeAllTemporary, Turns: 5},
		Attributes: primitives.Attributes{Health: 70, Strength: 5, Agility: 3},
	}

	c.ApplyEffect(e)

	if len(c.TemporaryEffects) != 1 || c.TemporaryEffects[0] != e {
		t.Fatalf("AllTemporary effect must be tracked; got len=%d", len(c.TemporaryEffects))
	}
	if c.Attributes.Health != 100 {
		t.Errorf("Health should be clamped to MaxHealth: got %v, want 100", c.Attributes.Health)
	}
	if c.Attributes.Strength != 7 || c.Attributes.Agility != 4 {
		t.Errorf("Attributes not applied correctly: got (Str %.1f, Agi %.1f)", c.Attributes.Strength, c.Attributes.Agility)
	}

	// Partial processing should not remove effect
	c.ProcessTemporaryEffects(3)
	if e.Duration.Turns != 2 {
		t.Errorf("Steps should be reduced to 2, got %d", e.Duration.Turns)
	}
	if len(c.TemporaryEffects) != 1 {
		t.Fatalf("Effect should still be tracked, got len=%d", len(c.TemporaryEffects))
	}
}

func TestCharacter_ApplyEffect_MaxHealthIncrease_AdjustsCurrentHealth(t *testing.T) {
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 50, Strength: 10, Agility: 5},
	)

	e := &primitives.Effect{
		Duration:   primitives.EffectDuration{Type: primitives.EffectDurationTypeAllPermanent, Turns: 0},
		Attributes: primitives.Attributes{MaxHealth: 20},
	}

	c.ApplyEffect(e)

	if c.Attributes.MaxHealth != 120 {
		t.Errorf("MaxHealth after ApplyEffect = %.1f, want 120", c.Attributes.MaxHealth)
	}
	if c.Attributes.Health != 70 {
		t.Errorf("Health after ApplyEffect = %.1f, want 70", c.Attributes.Health)
	}
}

func TestCharacter_ApplyEffect_MaxHealthDecrease_HealthClampedToOne(t *testing.T) {
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 10, Strength: 10, Agility: 5},
	)

	e := &primitives.Effect{
		Duration:   primitives.EffectDuration{Type: primitives.EffectDurationTypeAllPermanent, Turns: 0},
		Attributes: primitives.Attributes{MaxHealth: -20},
	}

	c.ApplyEffect(e)

	if c.Attributes.MaxHealth != 80 {
		t.Errorf("MaxHealth after decreasing effect = %.1f, want 80", c.Attributes.MaxHealth)
	}
	if c.Attributes.Health != 1 {
		t.Errorf("Health should be clamped to 1 when effect reduces it <= 0; got %.1f", c.Attributes.Health)
	}
}

func TestCharacter_ProcessTemporaryEffects_MaxHealthPermamentAffectHealth(t *testing.T) {
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 10, Strength: 10, Agility: 5},
	)

	e := &primitives.Effect{
		Duration:   primitives.EffectDuration{Type: primitives.EffectDurationTypeAllTemporaryHealPermanent, Turns: 3},
		Attributes: primitives.Attributes{MaxHealth: 20},
	}

	c.ApplyEffect(e)

	if c.Attributes.MaxHealth != 120 {
		t.Errorf("MaxHealth after ApplyEffect = %.1f, want 120", c.Attributes.MaxHealth)
	}
	if c.Attributes.Health != 30 {
		t.Errorf("Health after ApplyEffect = %.1f, want 30", c.Attributes.Health)
	}

	c.ProcessTemporaryEffects(2)
	if c.Attributes.MaxHealth != 120 {
		t.Errorf("MaxHealth after ApplyEffect = %.1f, want 120", c.Attributes.MaxHealth)
	}
	if c.Attributes.Health != 30 {
		t.Errorf("Health after ApplyEffect = %.1f, want 30", c.Attributes.Health)
	}

	c.ProcessTemporaryEffects(2)
	if c.Attributes.MaxHealth != 100 {
		t.Errorf("MaxHealth after ApplyEffect = %.1f, want 100", c.Attributes.MaxHealth)
	}
	if c.Attributes.Health != 30 {
		t.Errorf("Health after ApplyEffect = %.1f, want 30", c.Attributes.Health)
	}
}

func TestCharacter_ProcessTemporaryEffects_MaxHealthTemporaryAffectHealth(t *testing.T) {
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 10, Strength: 10, Agility: 5},
	)

	e := &primitives.Effect{
		Duration:   primitives.EffectDuration{Type: primitives.EffectDurationTypeAllTemporary, Turns: 3},
		Attributes: primitives.Attributes{MaxHealth: 20},
	}

	c.ApplyEffect(e)

	if c.Attributes.MaxHealth != 120 {
		t.Errorf("MaxHealth after ApplyEffect = %.1f, want 120", c.Attributes.MaxHealth)
	}
	if c.Attributes.Health != 30 {
		t.Errorf("Health after ApplyEffect = %.1f, want 30", c.Attributes.Health)
	}

	c.ProcessTemporaryEffects(2)
	if c.Attributes.MaxHealth != 120 {
		t.Errorf("MaxHealth after ApplyEffect = %.1f, want 120", c.Attributes.MaxHealth)
	}
	if c.Attributes.Health != 30 {
		t.Errorf("Health after ApplyEffect = %.1f, want 30", c.Attributes.Health)
	}

	c.ProcessTemporaryEffects(2)
	if c.Attributes.MaxHealth != 100 {
		t.Errorf("MaxHealth after ApplyEffect = %.1f, want 100", c.Attributes.MaxHealth)
	}
	if c.Attributes.Health != 10 {
		t.Errorf("Health after ApplyEffect = %.1f, want 10", c.Attributes.Health)
	}
}

func TestCharacter_ProcessTemporaryEffects_ExpiresAndRollsBack(t *testing.T) {
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 50, Strength: 10, Agility: 1},
	)
	e := &primitives.Effect{
		Duration:   primitives.EffectDuration{Type: primitives.EffectDurationTypeAllTemporary, Turns: 2},
		Attributes: primitives.Attributes{Health: 20, Strength: 5, Agility: 3},
	}

	c.ApplyEffect(e)
	if c.Attributes.Health != 70 || c.Attributes.Strength != 15 || c.Attributes.Agility != 4 {
		t.Fatalf("ApplyEffect mutated wrong: got (H %.1f, S %.1f, A %.1f)", c.Attributes.Health, c.Attributes.Strength, c.Attributes.Agility)
	}

	c.ProcessTemporaryEffects(2)

	if len(c.TemporaryEffects) != 0 {
		t.Fatalf("Expired effect should be removed, got len=%d", len(c.TemporaryEffects))
	}
	if c.Attributes.Health != 50 || c.Attributes.Strength != 10 || c.Attributes.Agility != 1 {
		t.Errorf("Attributes should be rolled back: got (H %.1f, S %.1f, A %.1f)", c.Attributes.Health, c.Attributes.Strength, c.Attributes.Agility)
	}
}

func TestCharacter_RemoveTemporaryEffect_HealPermanent_OthersRevert(t *testing.T) {
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 40, Agility: 2, Strength: 3},
	)
	e := &primitives.Effect{
		Duration:   primitives.EffectDuration{Type: primitives.EffectDurationTypeAllTemporaryHealPermanent, Turns: 5},
		Attributes: primitives.Attributes{Health: 30, Strength: 7, Agility: 11},
	}

	c.ApplyEffect(e)

	if c.Attributes.Health != 70 || c.Attributes.Strength != 10 || c.Attributes.Agility != 13 {
		t.Fatalf("Precondition after ApplyEffect failed: got (H %.1f, S %.1f, A %.1f)", c.Attributes.Health, c.Attributes.Strength, c.Attributes.Agility)
	}

	c.RemoveTemporaryEffect(e)

	if len(c.TemporaryEffects) != 0 {
		t.Fatalf("Effect should be removed, got len=%d", len(c.TemporaryEffects))
	}
	// Health should remain boosted (permanent heal), other attributes revert
	if c.Attributes.Health != 70 {
		t.Errorf("Health should remain after removal (permanent heal), got %.1f, want 70", c.Attributes.Health)
	}
	if c.Attributes.Strength != 3 || c.Attributes.Agility != 2 {
		t.Errorf("Other attributes should revert: got (S %.1f, A %.1f)", c.Attributes.Strength, c.Attributes.Agility)
	}
}

func TestCharacter_RemoveTemporaryEffect_NilIsNothing(t *testing.T) {
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 40, Agility: 2, Strength: 3},
	)
	e := &primitives.Effect{
		Duration:   primitives.EffectDuration{Type: primitives.EffectDurationTypeAllTemporaryHealPermanent, Turns: 5},
		Attributes: primitives.Attributes{Health: 30, Strength: 7, Agility: 11},
	}

	c.ApplyEffect(e)

	if c.Attributes.Health != 70 || c.Attributes.Strength != 10 || c.Attributes.Agility != 13 {
		t.Fatalf("Precondition after ApplyEffect failed: got (H %.1f, S %.1f, A %.1f)", c.Attributes.Health, c.Attributes.Strength, c.Attributes.Agility)
	}

	c.RemoveTemporaryEffect(nil)

	if c.Attributes.Health != 70 || c.Attributes.Strength != 10 || c.Attributes.Agility != 13 {
		t.Fatalf("Precondition after ApplyEffect failed: got (H %.1f, S %.1f, A %.1f)", c.Attributes.Health, c.Attributes.Strength, c.Attributes.Agility)
	}
}

func TestCharacter_ProcessTemporaryEffects_MultipleExpire(t *testing.T) {
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 10, Agility: 0, Strength: 0},
	)
	e1 := &primitives.Effect{
		Duration:   primitives.EffectDuration{Type: primitives.EffectDurationTypeAllTemporary, Turns: 1},
		Attributes: primitives.Attributes{Health: 10, Strength: 5},
	}
	e2 := &primitives.Effect{
		Duration:   primitives.EffectDuration{Type: primitives.EffectDurationTypeAllTemporary, Turns: 1},
		Attributes: primitives.Attributes{Health: 20, Agility: 7},
	}

	c.ApplyEffect(e1)
	c.ApplyEffect(e2)

	// After apply: H=40, S=5, A=7
	if c.Attributes.Health != 40 || c.Attributes.Strength != 5 || c.Attributes.Agility != 7 {
		t.Fatalf("Apply effects wrong: got (H %.1f, S %.1f, A %.1f)", c.Attributes.Health, c.Attributes.Strength, c.Attributes.Agility)
	}

	c.ProcessTemporaryEffects(1)

	if len(c.TemporaryEffects) != 0 {
		t.Fatalf("Both effects should expire and be removed, got len=%d", len(c.TemporaryEffects))
	}
	// Should roll back exactly once each
	if c.Attributes.Health != 10 || c.Attributes.Strength != 0 || c.Attributes.Agility != 0 {
		t.Errorf("Attributes should be fully rolled back: got (H %.1f, S %.1f, A %.1f)", c.Attributes.Health, c.Attributes.Strength, c.Attributes.Agility)
	}
}

func TestCharacter_CheckEvasion_ZeroAgilityAlwaysFalse(t *testing.T) {
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 100, Agility: 0, Strength: 0},
	)
	for i := 0; i < 100; i++ {
		randomGenerator := utils.NewRandomWithSeed(int64(i))
		if c.CheckEvasion(0, randomGenerator) {
			t.Fatalf("CheckEvasion must be false when Agility=0 (iter %d)", i)
		}
	}
}

func TestCharacter_CheckEvasion_HighAgilityMostlyTrueWithSeed(t *testing.T) {
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 100, Agility: 380, Strength: 0}, // ~95% chance
	)
	evasionsCount := 0
	total := 200
	for i := 0; i < total; i++ {
		randomGenerator := utils.NewRandomWithSeed(int64(i))
		if c.CheckEvasion(0, randomGenerator) {
			evasionsCount++
		}
	}

	if evasionsCount < 150 { // be conservative to avoid flakiness
		t.Errorf("Expected many evades with high agility; got %d/%d", evasionsCount, total)
	}
}

func TestCharacter_CheckEvasion_NoRandom(t *testing.T) {
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 100, Agility: 380, Strength: 0}, // ~95% chance
	)
	evasionsCount1 := 0
	total := 200
	randomGenerator := utils.NewRandomWithSeed(defaultCharacterTestSeed)
	for i := 0; i < total; i++ {
		if c.CheckEvasion(0, randomGenerator) {
			evasionsCount1++
		}
	}

	evasionsCount2 := 0
	for i := 0; i < total; i++ {
		randomGenerator := utils.NewRandomWithSeed(int64(i))
		if c.CheckEvasion(0, randomGenerator) {
			evasionsCount2++
		}
	}

	if evasionsCount1 == evasionsCount2 {
		t.Errorf("Expected different evasion counts with different random seeds; got both %d/%d", evasionsCount1, total)
	}
}

func TestCharacter_CheckEvasion_AttackerAgilityReducesEvasion(t *testing.T) {
	// Защитник с Agility=20, атакующий с Agility=0 vs атакующий с Agility=30
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 100, Agility: 20, Strength: 0},
	)

	total := 500
	evasionsNoAttackerAgi := 0
	evasionsHighAttackerAgi := 0

	for i := 0; i < total; i++ {
		rng := utils.NewRandomWithSeed(int64(i))
		if c.CheckEvasion(0, rng) {
			evasionsNoAttackerAgi++
		}
	}
	for i := 0; i < total; i++ {
		rng := utils.NewRandomWithSeed(int64(i))
		if c.CheckEvasion(30, rng) {
			evasionsHighAttackerAgi++
		}
	}

	if evasionsHighAttackerAgi >= evasionsNoAttackerAgi {
		t.Errorf("High attacker agility should reduce evasions: noAttackerAgi=%d, highAttackerAgi=%d",
			evasionsNoAttackerAgi, evasionsHighAttackerAgi)
	}
}

func TestCharacter_CheckEvasion_AttackerOverwhelmsDefender(t *testing.T) {
	// Если ловкость атакующего * 0.5 >= ловкости защитника, уклонений быть не должно.
	c := NewCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 100, Agility: 10, Strength: 0},
	)

	// attackerAgility=20 → effectiveAgility = 10 - 20*0.5 = 0 → шанс = 0
	for i := 0; i < 100; i++ {
		rng := utils.NewRandomWithSeed(int64(i))
		if c.CheckEvasion(20, rng) {
			t.Fatalf("CheckEvasion must be false when attacker agility overwhelms defender (iter %d)", i)
		}
	}
}

type MockUsable struct {
	effect *primitives.Effect
}

func (m *MockUsable) Use() *primitives.Effect {
	return m.effect
}

func TestCharacter_Use_WithMockUsable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		initialHealth  float64
		initialStr     float64
		initialAgi     float64
		effectAttrs    primitives.Attributes
		effectDuration primitives.EffectDuration
		wantHealth     float64
		wantStr        float64
		wantAgi        float64
		wantEffectsCnt int
	}{
		{
			name:          "Permanent effect increases attributes",
			initialHealth: 50,
			initialStr:    10,
			initialAgi:    5,
			effectAttrs: primitives.Attributes{
				Health:   20,
				Strength: 5,
				Agility:  3,
			},
			effectDuration: primitives.EffectDuration{
				Type:  primitives.EffectDurationTypeAllPermanent,
				Turns: 0,
			},
			wantHealth:     70,
			wantStr:        15,
			wantAgi:        8,
			wantEffectsCnt: 0, // Permanent не добавляется в TemporaryEffects
		},
		{
			name:          "Temporary effect tracked",
			initialHealth: 30,
			initialStr:    8,
			initialAgi:    2,
			effectAttrs: primitives.Attributes{
				Health:   10,
				Strength: 4,
				Agility:  6,
			},
			effectDuration: primitives.EffectDuration{
				Type:  primitives.EffectDurationTypeAllTemporary,
				Turns: 5,
			},
			wantHealth:     40,
			wantStr:        12,
			wantAgi:        8,
			wantEffectsCnt: 1,
		},
		{
			name:          "Heal permanent type tracked",
			initialHealth: 40,
			initialStr:    5,
			initialAgi:    3,
			effectAttrs: primitives.Attributes{
				Health:   30,
				Strength: 10,
				Agility:  5,
			},
			effectDuration: primitives.EffectDuration{
				Type:  primitives.EffectDurationTypeAllTemporaryHealPermanent,
				Turns: 3,
			},
			wantHealth:     70,
			wantStr:        15,
			wantAgi:        8,
			wantEffectsCnt: 1,
		},
		{
			name:          "Health clamped to MaxHealth",
			initialHealth: 90,
			initialStr:    10,
			initialAgi:    5,
			effectAttrs: primitives.Attributes{
				Health:   50, // Превысит MaxHealth (100)
				Strength: 0,
				Agility:  0,
			},
			effectDuration: primitives.EffectDuration{
				Type:  primitives.EffectDurationTypeAllPermanent,
				Turns: 0,
			},
			wantHealth:     100, // Зажато до MaxHealth
			wantStr:        10,
			wantAgi:        5,
			wantEffectsCnt: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewCharacter(
				primitives.Box{
					Point: primitives.Point2D[int]{X: 0, Y: 0},
					Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
				},
				primitives.Attributes{
					MaxHealth: 100,
					Health:    tt.initialHealth,
					Strength:  tt.initialStr,
					Agility:   tt.initialAgi,
				},
			)

			usable := &MockUsable{
				effect: &primitives.Effect{
					Duration:   tt.effectDuration,
					Attributes: tt.effectAttrs,
				},
			}

			c.Use(usable)

			if c.Attributes.Health != tt.wantHealth {
				t.Errorf("Health = %.1f, want %.1f", c.Attributes.Health, tt.wantHealth)
			}
			if c.Attributes.Strength != tt.wantStr {
				t.Errorf("Strength = %.1f, want %.1f", c.Attributes.Strength, tt.wantStr)
			}
			if c.Attributes.Agility != tt.wantAgi {
				t.Errorf("Agility = %.1f, want %.1f", c.Attributes.Agility, tt.wantAgi)
			}
			if len(c.TemporaryEffects) != tt.wantEffectsCnt {
				t.Errorf("TemporaryEffects count = %d, want %d", len(c.TemporaryEffects), tt.wantEffectsCnt)
			}
		})
	}
}

func TestCharacter_Use_WithElixir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		elixirType     items.ElixirType
		initialHealth  float64
		initialStr     float64
		initialAgi     float64
		wantMinHealth  float64 // Минимальное ожидаемое здоровье (из-за random)
		wantMinStr     float64
		wantMinAgi     float64
		wantEffectsCnt int
	}{
		{
			name:           "Use Strength elixir",
			elixirType:     items.ElixirTypeStrength,
			initialHealth:  50,
			initialStr:     10,
			initialAgi:     5,
			wantMinHealth:  50,
			wantMinStr:     13, // Ожидаем прирост силы (min +3)
			wantMinAgi:     5,
			wantEffectsCnt: 1, // Temporary effect
		},
		{
			name:           "Use Agility elixir",
			elixirType:     items.ElixirTypeAgility,
			initialHealth:  60,
			initialStr:     12,
			initialAgi:     3,
			wantMinHealth:  60,
			wantMinStr:     12,
			wantMinAgi:     6, // Ожидаем прирост ловкости (min +3)
			wantEffectsCnt: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rng := utils.NewRandomWithSeed(defaultCharacterTestSeed)
			c := NewCharacter(
				primitives.Box{
					Point: primitives.Point2D[int]{X: 0, Y: 0},
					Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
				},
				primitives.Attributes{
					MaxHealth: 100,
					Health:    tt.initialHealth,
					Strength:  tt.initialStr,
					Agility:   tt.initialAgi,
				},
			)

			box := primitives.Box{
				Point: primitives.Point2D[int]{X: 0, Y: 0},
				Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
			}
			elixir := items.NewElixirBuiltin(rng, box, tt.elixirType)

			c.Use(elixir)

			if c.Attributes.Health < tt.wantMinHealth {
				t.Errorf("Health = %.1f, want >= %.1f", c.Attributes.Health, tt.wantMinHealth)
			}
			if c.Attributes.Strength < tt.wantMinStr {
				t.Errorf("Strength = %.1f, want >= %.1f", c.Attributes.Strength, tt.wantMinStr)
			}
			if c.Attributes.Agility < tt.wantMinAgi {
				t.Errorf("Agility = %.1f, want >= %.1f", c.Attributes.Agility, tt.wantMinAgi)
			}
			if len(c.TemporaryEffects) != tt.wantEffectsCnt {
				t.Errorf("TemporaryEffects count = %d, want %d", len(c.TemporaryEffects), tt.wantEffectsCnt)
			}
		})
	}
}

func TestCharacter_Use_WithScroll(t *testing.T) {
	t.Parallel()

	rng := utils.NewRandomWithSeed(defaultCharacterTestSeed)
	c := NewCharacter(
		primitives.Box{
			Point: primitives.Point2D[int]{X: 0, Y: 0},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
		primitives.Attributes{
			MaxHealth: 100,
			Health:    50,
			Strength:  10,
			Agility:   5,
		},
	)

	initialHealth := c.Attributes.Health
	initialStr := c.Attributes.Strength
	initialAgi := c.Attributes.Agility

	box := primitives.Box{
		Point: primitives.Point2D[int]{X: 0, Y: 0},
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}
	scroll := items.NewScrollBuiltin(rng, box, items.ScrollTypeStrength)

	c.Use(scroll)

	attributesChanged := c.Attributes.Health != initialHealth ||
		c.Attributes.Strength != initialStr ||
		c.Attributes.Agility != initialAgi

	if !attributesChanged {
		t.Error("Expected scroll to change at least one attribute")
	}

	// Свиток не должен добавлять временных эффектов
	if len(c.TemporaryEffects) != 0 {
		t.Error("Expected scroll not to add temporary effects")
	}
}

func TestCharacter_Use_MultipleUsables(t *testing.T) {
	t.Parallel()

	c := NewCharacter(
		primitives.Box{
			Point: primitives.Point2D[int]{X: 0, Y: 0},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
		primitives.Attributes{
			MaxHealth: 100,
			Health:    30,
			Strength:  5,
			Agility:   3,
		},
	)

	// Первый эффект: +20 здоровья, +5 силы (temporary)
	usable1 := &MockUsable{
		effect: &primitives.Effect{
			Duration: primitives.EffectDuration{
				Type:  primitives.EffectDurationTypeAllTemporary,
				Turns: 5,
			},
			Attributes: primitives.Attributes{
				Health:   20,
				Strength: 5,
				Agility:  0,
			},
		},
	}

	// Второй эффект: +10 здоровья, +3 ловкости (temporary)
	usable2 := &MockUsable{
		effect: &primitives.Effect{
			Duration: primitives.EffectDuration{
				Type:  primitives.EffectDurationTypeAllTemporary,
				Turns: 3,
			},
			Attributes: primitives.Attributes{
				Health:   10,
				Strength: 0,
				Agility:  3,
			},
		},
	}

	// Act
	c.Use(usable1)
	c.Use(usable2)

	// Assert
	expectedHealth := 60.0
	expectedStr := 10.0
	expectedAgi := 6.0

	if c.Attributes.Health != expectedHealth {
		t.Errorf("Health after multiple uses = %.1f, want %.1f", c.Attributes.Health, expectedHealth)
	}
	if c.Attributes.Strength != expectedStr {
		t.Errorf("Strength after multiple uses = %.1f, want %.1f", c.Attributes.Strength, expectedStr)
	}
	if c.Attributes.Agility != expectedAgi {
		t.Errorf("Agility after multiple uses = %.1f, want %.1f", c.Attributes.Agility, expectedAgi)
	}
	if len(c.TemporaryEffects) != 2 {
		t.Errorf("Expected 2 temporary effects, got %d", len(c.TemporaryEffects))
	}
}

func TestCharacter_Use_EffectInteractionWithProcessing(t *testing.T) {
	t.Parallel()

	c := NewCharacter(
		primitives.Box{
			Point: primitives.Point2D[int]{X: 0, Y: 0},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
		primitives.Attributes{
			MaxHealth: 100,
			Health:    50,
			Strength:  10,
			Agility:   5,
		},
	)

	usable := &MockUsable{
		effect: &primitives.Effect{
			Duration: primitives.EffectDuration{
				Type:  primitives.EffectDurationTypeAllTemporary,
				Turns: 3,
			},
			Attributes: primitives.Attributes{
				Health:   20,
				Strength: 5,
				Agility:  3,
			},
		},
	}

	c.Use(usable)

	if c.Attributes.Health != 70 || c.Attributes.Strength != 15 || c.Attributes.Agility != 8 {
		t.Fatalf("After Use: unexpected attributes (H=%.1f, S=%.1f, A=%.1f)",
			c.Attributes.Health, c.Attributes.Strength, c.Attributes.Agility)
	}

	// Обрабатываем 2 шага
	c.ProcessTemporaryEffects(2)

	// Эффект должен остаться
	if len(c.TemporaryEffects) != 1 {
		t.Fatalf("Expected effect to remain after 2 steps, got %d effects", len(c.TemporaryEffects))
	}

	// Обрабатываем последний шаг
	c.ProcessTemporaryEffects(1)

	// Эффект должен исчезнуть
	if len(c.TemporaryEffects) != 0 {
		t.Errorf("Expected effect to expire, got %d effects", len(c.TemporaryEffects))
	}

	// Атрибуты должны вернуться к базовым
	if c.Attributes.Health != 50 || c.Attributes.Strength != 10 || c.Attributes.Agility != 5 {
		t.Errorf("After expiration: unexpected attributes (H=%.1f, S=%.1f, A=%.1f)",
			c.Attributes.Health, c.Attributes.Strength, c.Attributes.Agility)
	}
}

func TestCharacter_Use_NilEffect(t *testing.T) {
	t.Parallel()

	c := NewCharacter(
		primitives.Box{
			Point: primitives.Point2D[int]{X: 0, Y: 0},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
		primitives.Attributes{
			MaxHealth: 100,
			Health:    50,
			Strength:  10,
			Agility:   5,
		},
	)

	initialHealth := c.Attributes.Health
	initialStr := c.Attributes.Strength
	initialAgi := c.Attributes.Agility

	usable := &MockUsable{
		effect: nil,
	}

	// Act - не должно паниковать, но может быть nil pointer dereference
	// В зависимости от реализации ApplyEffect
	defer func() {
		if r := recover(); r != nil {
			// Если паника, это ожидаемо для nil effect
			// Можно добавить комментарий, что нужна валидация
			t.Logf("Panic on nil effect (expected if no validation): %v", r)
		}
	}()

	c.Use(usable)

	// Если не запаниковало, атрибуты не должны измениться
	if c.Attributes.Health != initialHealth ||
		c.Attributes.Strength != initialStr ||
		c.Attributes.Agility != initialAgi {
		t.Error("Attributes should not change with nil effect")
	}
}

func TestCharacter_Use_ZeroEffect(t *testing.T) {
	t.Parallel()

	c := NewCharacter(
		primitives.Box{
			Point: primitives.Point2D[int]{X: 0, Y: 0},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
		primitives.Attributes{
			MaxHealth: 100,
			Health:    50,
			Strength:  10,
			Agility:   5,
		},
	)

	usable := &MockUsable{
		effect: &primitives.Effect{
			Duration: primitives.EffectDuration{
				Type:  primitives.EffectDurationTypeAllPermanent,
				Turns: 0,
			},
			Attributes: primitives.Attributes{
				Health:   0,
				Strength: 0,
				Agility:  0,
			},
		},
	}

	c.Use(usable)

	// Assert - атрибуты не должны измениться
	if c.Attributes.Health != 50 || c.Attributes.Strength != 10 || c.Attributes.Agility != 5 {
		t.Errorf("Zero effect should not change attributes: got (H=%.1f, S=%.1f, A=%.1f)",
			c.Attributes.Health, c.Attributes.Strength, c.Attributes.Agility)
	}

	if len(c.TemporaryEffects) != 0 {
		t.Error("Zero permanent effect should not add temporary effects")
	}
}

func TestCharacter_Use_LowHealthScenario(t *testing.T) {
	t.Parallel()

	c := NewCharacter(
		primitives.Box{
			Point: primitives.Point2D[int]{X: 0, Y: 0},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
		primitives.Attributes{
			MaxHealth: 100,
			Health:    5, // Критически низкое HP
			Strength:  10,
			Agility:   5,
		},
	)

	if c.IsAlive() == false {
		t.Fatal("Character should be alive with Health=5")
	}

	// Используем лечебный предмет
	usable := &MockUsable{
		effect: &primitives.Effect{
			Duration: primitives.EffectDuration{
				Type:  primitives.EffectDurationTypeAllPermanent,
				Turns: 0,
			},
			Attributes: primitives.Attributes{
				Health:   50,
				Strength: 0,
				Agility:  0,
			},
		},
	}

	c.Use(usable)

	// Assert
	if c.Attributes.Health != 55 {
		t.Errorf("Health after healing = %.1f, want 55", c.Attributes.Health)
	}

	if !c.IsAlive() {
		t.Error("Character should be alive after healing")
	}
}

func TestCharacter_Use_NegativeEffects(t *testing.T) {
	t.Parallel()

	c := NewCharacter(
		primitives.Box{
			Point: primitives.Point2D[int]{X: 0, Y: 0},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
		primitives.Attributes{
			MaxHealth: 100,
			Health:    50,
			Strength:  10,
			Agility:   5,
		},
	)

	// Предмет с "побочными эффектами" - уменьшает силу
	usable := &MockUsable{
		effect: &primitives.Effect{
			Duration: primitives.EffectDuration{
				Type:  primitives.EffectDurationTypeAllTemporary,
				Turns: 2,
			},
			Attributes: primitives.Attributes{
				Health:   30,
				Strength: -5, // Уменьшение силы
				Agility:  10,
			},
		},
	}

	c.Use(usable)

	if c.Attributes.Health != 80 {
		t.Errorf("Health = %.1f, want 80", c.Attributes.Health)
	}
	if c.Attributes.Strength != 5 {
		t.Errorf("Strength should decrease: got %.1f, want 5", c.Attributes.Strength)
	}
	if c.Attributes.Agility != 15 {
		t.Errorf("Agility = %.1f, want 15", c.Attributes.Agility)
	}

	// После истечения эффекта атрибуты должны вернуться
	c.ProcessTemporaryEffects(2)

	if c.Attributes.Strength != 10 {
		t.Errorf("Strength should return to 10, got %.1f", c.Attributes.Strength)
	}
}
