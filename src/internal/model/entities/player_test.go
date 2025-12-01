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

	err := p.EquipWeapon(w)
	if err != nil {
		t.Fatalf("EquipWeapon returned unexpected error: %v", err)
	}

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

	err := p.EquipWeapon(nil)
	if err != nil {
		t.Fatalf("EquipWeapon returned unexpected error: %v", err)
	}

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

	err = p.EquipWeapon(w)
	if err != nil {
		t.Fatalf("EquipWeapon returned unexpected error: %v", err)
	}
	err = p.EquipWeapon(nil) // Equip nil after valid weapon
	if err != nil {
		t.Fatalf("EquipWeapon returned unexpected error: %v", err)
	}

	if p.Weapon != w {
		t.Errorf("Weapon should remain unchanged when equipping nil")
	}

	// Attributes should remain as after valid equip
	if p.Attributes.Strength != base.Strength+delta.Strength || p.Attributes.Agility != base.Agility+delta.Agility {
		t.Errorf("Attributes should remain unchanged when equipping nil: got (S %.1f, A %.1f) want (S %.1f, A %.1f)", p.Attributes.Strength, p.Attributes.Agility, base.Strength+delta.Strength, base.Agility+delta.Agility)
	}
}

func TestPlayer_EquipWeapon_NoAddToBackpackIfPreviousWeaponIsNil(t *testing.T) {
	box := &primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	p := NewPlayer(box)
	base := p.Attributes

	if p.Weapon != nil {
		t.Fatalf("Precondition failed: expected no weapon equipped")
	}

	rnd1 := utils.NewRandomGeneratorWithSeed(100)
	w1 := items.NewWeaponBuiltin(rnd1, primitives.Box{}, items.WeaponTypeDagger)
	delta := w1.Effect.Attributes
	err := p.EquipWeapon(w1)
	if err != nil {
		t.Fatalf("EquipWeapon returned unexpected error: %v", err)
	}

	if p.Weapon != w1 {
		t.Errorf("Expected weapon reference to point to last equipped weapon")
	}

	if p.Attributes.Strength != base.Strength+delta.Strength || p.Attributes.Agility != base.Agility+delta.Agility {
		t.Errorf("Attributes should reflect last equipped weapon only: got (S %.1f, A %.1f) want (S %.1f, A %.1f)", p.Attributes.Strength, p.Attributes.Agility, base.Strength+delta.Strength, base.Agility+delta.Agility)
	}

	// Проверяем что рюкзак пуст в случае экипировки на пустой слот
	if !p.Backpack.IsEmpty() || p.Backpack.ItemsNum != 0 {
		t.Fatalf("Expected backpack to have 0 items after equipping weapon on free weapon slot, got %d items", p.Backpack.ItemsNum)
	}
}

func TestPlayer_EquipWeapon_MovePreviousWeaponToBackpack(t *testing.T) {
	box := &primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	p := NewPlayer(box)
	base := p.Attributes

	rnd1 := utils.NewRandomGeneratorWithSeed(100)
	w1 := items.NewWeaponBuiltin(rnd1, primitives.Box{}, items.WeaponTypeDagger)
	delta := w1.Effect.Attributes
	p.EquipWeapon(w1)

	rnd2 := utils.NewRandomGeneratorWithSeed(200)
	w2 := items.NewWeaponBuiltin(rnd2, primitives.Box{}, items.WeaponTypeAxe)
	delta = w2.Effect.Attributes
	p.EquipWeapon(w2)

	if p.Weapon != w2 {
		t.Errorf("Expected weapon reference to point to last equipped weapon")
	}

	if p.Attributes.Strength != base.Strength+delta.Strength || p.Attributes.Agility != base.Agility+delta.Agility {
		t.Errorf("Attributes should reflect last equipped weapon only: got (S %.1f, A %.1f) want (S %.1f, A %.1f)", p.Attributes.Strength, p.Attributes.Agility, base.Strength+delta.Strength, base.Agility+delta.Agility)
	}

	// Проверяем что первый предмет теперь в рюкзаке
	if p.Backpack.IsEmpty() || p.Backpack.ItemsNum != 1 {
		t.Fatalf("Expected backpack to have 1 item after equipping second weapon, got %d items", p.Backpack.ItemsNum)
	}

	wInBackpack, ok := p.Backpack.Weapons.Front().Value.(*items.Weapon)
	if !ok || wInBackpack != w1 {
		t.Errorf("Expected first equipped weapon to be in backpack, got %v", wInBackpack)
	}
}

func TestPlayer_EquipWeapon_PreviousWeaponIsNotNilAndBackpackIsFull_IsError(t *testing.T) {
	box := &primitives.Box{Point: primitives.Point2D[int]{X: 0, Y: 0}, Size: primitives.Size2D[uint]{Width: 1, Height: 1}}
	p := NewPlayer(box)

	rnd1 := utils.NewRandomGeneratorWithSeed(100)
	w1 := items.NewWeaponBuiltin(rnd1, primitives.Box{}, items.WeaponTypeDagger)
	p.EquipWeapon(w1)

	p.Backpack.Capacity = 0 // Уменьшаем вместимость рюкзака для теста

	rnd2 := utils.NewRandomGeneratorWithSeed(200)
	w2 := items.NewWeaponBuiltin(rnd2, primitives.Box{}, items.WeaponTypeAxe)
	err := p.EquipWeapon(w2)
	if err == nil {
		t.Fatalf("Expected error when equipping weapon with full backpack, got nil")
	}

	if p.Weapon != w1 {
		t.Errorf("Weapon should remain unchanged after failed equip attempt")
	}

	// Проверяем что рюкзак остался с 0 предметов
	if !p.Backpack.IsEmpty() || p.Backpack.ItemsNum != 0 {
		t.Fatalf("Expected backpack to have 0 items after failed equip attempt, got %d items", p.Backpack.ItemsNum)
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
