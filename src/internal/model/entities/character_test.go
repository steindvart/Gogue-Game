package entities

import (
	"testing"

	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

const defaultCharacterTestSeed int64 = 42

func newCharacter(box primitives.Box, attrs primitives.Attributes) *Character {
	return &Character{
		Shape:            &box,
		Attributes:       &attrs,
		TemporaryEffects: nil,
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
			c := newCharacter(
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
	c := newCharacter(box, primitives.Attributes{MaxHealth: 10, Health: 10})

	delta := primitives.Point2D[int]{X: 3, Y: -1}
	c.Move(delta)

	want := primitives.Point2D[int]{X: 4, Y: 1}
	if c.Shape.Point != want {
		t.Errorf("Move() point = %v, want %v", c.Shape.Point, want)
	}
}

func TestCharacter_TakeDamage(t *testing.T) {
	c := newCharacter(
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
	c := newCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 100, Strength: 17.5},
	)
	if got := c.Attack(); got != 17.5 {
		t.Errorf("Attack() = %v, want 17.5", got)
	}
}

func TestCharacter_ApplyEffect_AllPermanent_ClampsHealth(t *testing.T) {
	const maxHealth = 100

	c := newCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: maxHealth, Health: 50, Strength: 10, Agility: 5},
	)
	e := &primitives.Effect{
		Duration:   &primitives.EffectDuration{Type: primitives.EffectDurationTypeAllPermanent, Steps: 0},
		Attributes: &primitives.Attributes{Health: maxHealth, Strength: 5, Agility: 3},
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

	c := newCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: maxHealth, Health: 40, Strength: 2, Agility: 1},
	)
	e := &primitives.Effect{
		Duration:   &primitives.EffectDuration{Type: primitives.EffectDurationTypeAllTemporary, Steps: 5},
		Attributes: &primitives.Attributes{Health: 70, Strength: 5, Agility: 3},
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
	if e.Duration.Steps != 2 {
		t.Errorf("Steps should be reduced to 2, got %d", e.Duration.Steps)
	}
	if len(c.TemporaryEffects) != 1 {
		t.Fatalf("Effect should still be tracked, got len=%d", len(c.TemporaryEffects))
	}
}

func TestCharacter_ProcessTemporaryEffects_ExpiresAndRollsBack(t *testing.T) {
	c := newCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 50, Strength: 10, Agility: 1},
	)
	e := &primitives.Effect{
		Duration:   &primitives.EffectDuration{Type: primitives.EffectDurationTypeAllTemporary, Steps: 2},
		Attributes: &primitives.Attributes{Health: 20, Strength: 5, Agility: 3},
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
	c := newCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 40, Agility: 2, Strength: 3},
	)
	e := &primitives.Effect{
		Duration:   &primitives.EffectDuration{Type: primitives.EffectDurationTypeAllTemporaryHealPermanent, Steps: 5},
		Attributes: &primitives.Attributes{Health: 30, Strength: 7, Agility: 11},
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

func TestCharacter_ProcessTemporaryEffects_MultipleExpire(t *testing.T) {
	c := newCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 10, Agility: 0, Strength: 0},
	)
	e1 := &primitives.Effect{
		Duration:   &primitives.EffectDuration{Type: primitives.EffectDurationTypeAllTemporary, Steps: 1},
		Attributes: &primitives.Attributes{Health: 10, Strength: 5},
	}
	e2 := &primitives.Effect{
		Duration:   &primitives.EffectDuration{Type: primitives.EffectDurationTypeAllTemporary, Steps: 1},
		Attributes: &primitives.Attributes{Health: 20, Agility: 7},
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
	c := newCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 100, Agility: 0, Strength: 0},
	)
	for i := 0; i < 100; i++ {
		randomGenerator := utils.NewRandomGeneratorWithSeed(int64(i))
		if c.CheckEvasion(randomGenerator) {
			t.Fatalf("CheckEvasion must be false when Agility=0 (iter %d)", i)
		}
	}
}

func TestCharacter_CheckEvasion_HighAgilityMostlyTrueWithSeed(t *testing.T) {
	c := newCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 100, Agility: 380, Strength: 0}, // ~95% chance
	)
	evasionsCount := 0
	total := 200
	for i := 0; i < total; i++ {
		randomGenerator := utils.NewRandomGeneratorWithSeed(int64(i))
		if c.CheckEvasion(randomGenerator) {
			evasionsCount++
		}
	}

	if evasionsCount < 150 { // be conservative to avoid flakiness
		t.Errorf("Expected many evades with high agility; got %d/%d", evasionsCount, total)
	}
}

func TestCharacter_CheckEvasion_NoRandom(t *testing.T) {
	c := newCharacter(
		primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}},
		primitives.Attributes{MaxHealth: 100, Health: 100, Agility: 380, Strength: 0}, // ~95% chance
	)
	evasionsCount1 := 0
	total := 200
	randomGenerator := utils.NewRandomGeneratorWithSeed(defaultCharacterTestSeed)
	for i := 0; i < total; i++ {
		if c.CheckEvasion(randomGenerator) {
			evasionsCount1++
		}
	}

	evasionsCount2 := 0
	for i := 0; i < total; i++ {
		randomGenerator := utils.NewRandomGeneratorWithSeed(int64(i))
		if c.CheckEvasion(randomGenerator) {
			evasionsCount2++
		}
	}

	if evasionsCount1 == evasionsCount2 {
		t.Errorf("Expected different evasion counts with different random seeds; got both %d/%d", evasionsCount1, total)
	}
}
