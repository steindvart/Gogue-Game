package entities

import (
	"testing"

	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

const defaultPlayerTestSeed int64 = 42

func TestPlayer_NewPlayer_BasicInit(t *testing.T) {
	box := &primitives.Box{Point: primitives.Point2D[int]{X: 10, Y: 20}, Size: primitives.Size2D[uint]{Width: 2, Height: 3}}
	p := NewPlayer(box)

	if p == nil || p.Character == nil {
		t.Fatalf("NewPlayer should initialize Character, got: p=%v, Character=%v", p, p.Character)
	}
	if p.Shape != box { // Shape points to the same box pointer passed in
		t.Errorf("Shape should reference provided box pointer; got %p want %p", p.Shape, box)
	}

	// Default attributes
	if p.Attributes.Health != 100 || p.Attributes.MaxHealth != 100 || p.Attributes.Strength != 10 || p.Attributes.Agility != 5 {
		t.Errorf("Unexpected default attributes: got (H %.1f/MH %.1f, S %.1f, A %.1f)", p.Attributes.Health, p.Attributes.MaxHealth, p.Attributes.Strength, p.Attributes.Agility)
	}

	// Backpack, Weapon, progression
	if p.Backpack == nil {
		t.Fatalf("Backpack must be initialized")
	}
	if !p.Backpack.IsEmpty() || p.Backpack.ItemsNum != 0 {
		t.Errorf("Backpack should be empty on init; got items=%d", p.Backpack.ItemsNum)
	}
	if p.Weapon != nil {
		t.Errorf("Weapon should be nil on init")
	}
	if p.Experience != 0 {
		t.Errorf("Experience should start at 0, got %d", p.Experience)
	}
	if p.CharacterLevel != 1 {
		t.Errorf("CharacterLevel should start at 1, got %d", p.CharacterLevel)
	}
}

func TestPlayer_EquipWeapon_AppliesEffectAndStoresWeapon(t *testing.T) {
	box := &primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	p := NewPlayer(box)
	base := p.Attributes

	rnd := utils.NewRandomGeneratorWithSeed(defaultPlayerTestSeed)
	w := items.NewWeaponBuiltin(rnd, primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 1}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}, items.WeaponTypeSword)
	if w == nil || w.Effect == nil {
		t.Fatalf("Weapon or its effect/attributes should be initialized")
	}
	delta := w.Effect.Attributes

	p.EquipWeapon(w)

	if p.Weapon != w {
		t.Errorf("EquipWeapon should set current weapon")
	}

	// Attributes should be increased by weapon effect
	if p.Attributes.Strength != base.Strength+delta.Strength || p.Attributes.Agility != base.Agility+delta.Agility {
		t.Errorf("Attributes not updated by weapon: got (S %.1f, A %.1f) want (S %.1f, A %.1f)", p.Attributes.Strength, p.Attributes.Agility, base.Strength+delta.Strength, base.Agility+delta.Agility)
	}
}

func TestPlayer_EquipWeapon_NilIsNothing(t *testing.T) {
	box := &primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	p := NewPlayer(box)
	base := p.Attributes

	p.EquipWeapon(nil)

	if p.Weapon != nil {
		t.Errorf("Weapon should remain nil when equipping nil")
	}

	if p.Attributes.Strength != base.Strength || p.Attributes.Agility != base.Agility {
		t.Errorf("Attributes should not change when equipping nil: got (S %.1f, A %.1f) want (S %.1f, A %.1f)", p.Attributes.Strength, p.Attributes.Agility, base.Strength, base.Agility)
	}

	rnd := utils.NewRandomGeneratorWithSeed(defaultPlayerTestSeed)
	w := items.NewWeaponBuiltin(rnd, primitives.Box{Point: primitives.Point2D[int]{X: 1, Y: 1}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}, items.WeaponTypeSword)
	if w == nil || w.Effect == nil {
		t.Fatalf("Weapon or its effect/attributes should be initialized")
	}
	delta := w.Effect.Attributes

	p.EquipWeapon(w)
	p.EquipWeapon(nil) // Equip nil after valid weapon

	if p.Weapon != w {
		t.Errorf("Weapon should remain unchanged when equipping nil")
	}

	// Attributes should remain as after valid equip
	if p.Attributes.Strength != base.Strength+delta.Strength || p.Attributes.Agility != base.Agility+delta.Agility {
		t.Errorf("Attributes should remain unchanged when equipping nil: got (S %.1f, A %.1f) want (S %.1f, A %.1f)", p.Attributes.Strength, p.Attributes.Agility, base.Strength+delta.Strength, base.Agility+delta.Agility)
	}
}

func TestPlayer_UnequipWeapon_RevertsEffectAndUnsetsWeapon(t *testing.T) {
	box := &primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	p := NewPlayer(box)
	base := p.Attributes

	rnd := utils.NewRandomGeneratorWithSeed(defaultPlayerTestSeed)
	w := items.NewWeaponBuiltin(rnd, primitives.Box{Point: primitives.Point2D[int]{X: 2, Y: 2}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}, items.WeaponTypeDagger)
	delta := w.Effect.Attributes

	p.EquipWeapon(w)
	if p.Attributes.Strength != base.Strength+delta.Strength || p.Attributes.Agility != base.Agility+delta.Agility {
		t.Fatalf("Precondition failed after equip: got (S %.1f, A %.1f)", p.Attributes.Strength, p.Attributes.Agility)
	}

	ret := p.UnequipWeapon()
	if ret != w {
		t.Errorf("UnequipWeapon should return previously equipped weapon")
	}
	if p.Weapon != nil {
		t.Errorf("Weapon should be nil after unequip")
	}
	// Attributes should be reverted to base
	if p.Attributes.Strength != base.Strength || p.Attributes.Agility != base.Agility {
		t.Errorf("Attributes should revert after unequip: got (S %.1f, A %.1f) want (S %.1f, A %.1f)", p.Attributes.Strength, p.Attributes.Agility, base.Strength, base.Agility)
	}
}

func TestPlayer_UnequipWeapon_NoWeapon(t *testing.T) {
	box := &primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	p := NewPlayer(box)
	base := p.Attributes

	ret := p.UnequipWeapon()
	if ret != nil {
		t.Errorf("UnequipWeapon should return nil when no weapon equipped")
	}
	// Nothing should change
	if p.Attributes != base {
		t.Errorf("Attributes should not change when unequipping without weapon")
	}
}

// @todo - пока такое поведение, но в будущем нужно пересмотреть
func TestPlayer_EquipWeapon_Twice_StacksByDesign(t *testing.T) {
	// Document current behavior: equipping a second weapon does not auto-revert the first; effects stack.
	box := &primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	p := NewPlayer(box)
	base := p.Attributes

	rnd1 := utils.NewRandomGeneratorWithSeed(100)
	w1 := items.NewWeaponBuiltin(rnd1, primitives.Box{}, items.WeaponTypeDagger)
	d1 := w1.Effect.Attributes
	p.EquipWeapon(w1)

	rnd2 := utils.NewRandomGeneratorWithSeed(200)
	w2 := items.NewWeaponBuiltin(rnd2, primitives.Box{}, items.WeaponTypeAxe)
	d2 := w2.Effect.Attributes
	p.EquipWeapon(w2)

	// Last equipped weapon reference is stored
	if p.Weapon != w2 {
		t.Errorf("Expected weapon reference to point to last equipped weapon")
	}
	// Attributes include both effects
	if p.Attributes.Strength != base.Strength+d1.Strength+d2.Strength || p.Attributes.Agility != base.Agility+d1.Agility+d2.Agility {
		t.Errorf("Attributes should stack with multiple equips: got (S %.1f, A %.1f) want (S %.1f, A %.1f)", p.Attributes.Strength, p.Attributes.Agility, base.Strength+d1.Strength+d2.Strength, base.Agility+d1.Agility+d2.Agility)
	}
}
