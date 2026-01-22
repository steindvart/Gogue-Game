package items

import (
	"errors"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
	"testing"
)

const (
	testBackpackSeed = int64(42)
)

func createTestElixir(rng *utils.Random, elixirType ElixirType) *Elixir {
	box := primitives.Box{
		Point: primitives.Point2D[int]{X: 0, Y: 0},
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}
	return NewElixirBuiltin(rng, box, elixirType)
}

func createTestScroll(rng *utils.Random, scrollType ScrollType) *Scroll {
	box := primitives.Box{
		Point: primitives.Point2D[int]{X: 0, Y: 0},
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}
	return NewScrollBuiltin(rng, box, scrollType)
}

func createTestFood(rng *utils.Random, foodType FoodType) *Food {
	box := primitives.Box{
		Point: primitives.Point2D[int]{X: 0, Y: 0},
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}
	return NewFoodBuiltin(rng, box, foodType)
}

func createTestWeapon(rng *utils.Random, weaponType WeaponType) *Weapon {
	box := primitives.Box{
		Point: primitives.Point2D[int]{X: 0, Y: 0},
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}
	return NewWeaponBuiltin(rng, box, weaponType)
}

func createTestTreasure(rng *utils.Random, treasureType TreasureType) *Treasure {
	box := primitives.Box{
		Point: primitives.Point2D[int]{X: 0, Y: 0},
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}
	return NewTreasureBuiltin(rng, box, treasureType)
}

func TestBackpack_NewBackpack(t *testing.T) {
	backpack := NewBackpack()

	if backpack == nil {
		t.Fatal("Expected non-nil backpack, got nil")
	}

	if backpack.Capacity != DefaultBackpackCapacity {
		t.Errorf("Expected capacity %d, got %d", DefaultBackpackCapacity, backpack.Capacity)
	}

	if backpack.ItemsNum != 0 {
		t.Errorf("Expected ItemsNum 0, got %d", backpack.ItemsNum)
	}

	if backpack.Treasures != 0 {
		t.Errorf("Expected Treasures 0, got %d", backpack.Treasures)
	}

	if backpack.Elixirs == nil {
		t.Error("Expected Elixirs map to be initialized")
	}

	if backpack.Scrolls == nil {
		t.Error("Expected Scrolls map to be initialized")
	}

	if backpack.Foods == nil {
		t.Error("Expected Foods map to be initialized")
	}

	if backpack.Weapons == nil {
		t.Error("Expected Weapons map to be initialized")
	}

	if !backpack.IsEmpty() {
		t.Error("Expected new backpack to be empty")
	}
}

func TestBackpack_IsEmpty(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*Backpack)
		wantEmpty bool
	}{
		{
			name:      "New backpack is empty",
			setupFunc: func(b *Backpack) {},
			wantEmpty: true,
		},
		{
			name: "Backpack with one item is not empty",
			setupFunc: func(b *Backpack) {
				rng := utils.NewRandomWithSeed(testBackpackSeed)
				elixir := createTestElixir(rng, ElixirTypeStrength)
				_ = b.AddElixir(elixir)
			},
			wantEmpty: false,
		},
		{
			name: "Backpack with treasure is still empty (treasures don't count as items)",
			setupFunc: func(b *Backpack) {
				rng := utils.NewRandomWithSeed(testBackpackSeed)
				treasure := createTestTreasure(rng, TreasureTypeGold)
				b.AddTreasure(treasure)
			},
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backpack := NewBackpack()
			tt.setupFunc(backpack)

			if got := backpack.IsEmpty(); got != tt.wantEmpty {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.wantEmpty)
			}
		})
	}
}

func TestBackpack_IsFull(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*Backpack)
		wantFull  bool
	}{
		{
			name:      "New backpack is not full",
			setupFunc: func(b *Backpack) {},
			wantFull:  false,
		},
		{
			name: "Backpack with capacity items is full",
			setupFunc: func(b *Backpack) {
				rng := utils.NewRandomWithSeed(testBackpackSeed)
				for i := uint(0); i < DefaultBackpackCapacity; i++ {
					elixir := createTestElixir(rng, ElixirTypeStrength)
					_ = b.AddElixir(elixir)
				}
			},
			wantFull: true,
		},
		{
			name: "Backpack with capacity-1 items is not full",
			setupFunc: func(b *Backpack) {
				rng := utils.NewRandomWithSeed(testBackpackSeed)
				for i := uint(0); i < DefaultBackpackCapacity-1; i++ {
					elixir := createTestElixir(rng, ElixirTypeStrength)
					_ = b.AddElixir(elixir)
				}
			},
			wantFull: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backpack := NewBackpack()
			tt.setupFunc(backpack)

			if got := backpack.IsFull(); got != tt.wantFull {
				t.Errorf("IsFull() = %v, want %v (ItemsNum=%d, Capacity=%d)",
					got, tt.wantFull, backpack.ItemsNum, backpack.Capacity)
			}
		})
	}
}

func TestBackpack_AddItem(t *testing.T) {
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	tests := []struct {
		name      string
		item      any
		wantError bool
	}{
		{
			name:      "Add Elixir via AddItem",
			item:      createTestElixir(rng, ElixirTypeStrength),
			wantError: false,
		},
		{
			name:      "Add Scroll via AddItem",
			item:      createTestScroll(rng, ScrollTypeStrength),
			wantError: false,
		},
		{
			name:      "Add Food via AddItem",
			item:      createTestFood(rng, FoodTypeBread),
			wantError: false,
		},
		{
			name:      "Add Weapon via AddItem",
			item:      createTestWeapon(rng, WeaponTypeSword),
			wantError: false,
		},
		{
			name:      "Add Treasure via AddItem",
			item:      createTestTreasure(rng, TreasureTypeGold),
			wantError: false,
		},
		{
			name:      "Add invalid type returns error",
			item:      "not an item",
			wantError: true,
		},
		{
			name:      "Add nil returns error",
			item:      nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backpack := NewBackpack()
			err := backpack.AddItem(tt.item)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				var wantErrType NotItemError
				if !errors.As(err, &wantErrType) {
					t.Errorf("Expected %T, got %T", wantErrType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

func TestBackpack_AddItemToFullBackpackIsError(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	for i := uint(0); i < DefaultBackpackCapacity; i++ {
		elixir := createTestElixir(rng, ElixirTypeStrength)
		if err := backpack.AddItem(elixir); err != nil {
			t.Fatalf("Failed to add item %d: %v", i, err)
		}
	}

	elixir := createTestElixir(rng, ElixirTypeAgility)
	err := backpack.AddItem(elixir)

	if err == nil {
		t.Error("Expected error when adding to full backpack, got nil")
	}

	var expectedErr BackpackIsFullError
	if !errors.As(err, &expectedErr) {
		t.Errorf("Expected BackpackIsFullError, got %T: %v", err, err)
	}
}

func TestBackpack_RemoveItem(t *testing.T) {
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	tests := []struct {
		name      string
		setupFunc func(*Backpack) any
		wantError bool
	}{
		{
			name: "Remove existing Elixir",
			setupFunc: func(b *Backpack) any {
				elixir := createTestElixir(rng, ElixirTypeStrength)
				_ = b.AddElixir(elixir)
				return elixir
			},
			wantError: false,
		},
		{
			name: "Remove existing Scroll",
			setupFunc: func(b *Backpack) any {
				scroll := createTestScroll(rng, ScrollTypeStrength)
				_ = b.AddScroll(scroll)
				return scroll
			},
			wantError: false,
		},
		{
			name: "Remove existing Food",
			setupFunc: func(b *Backpack) any {
				food := createTestFood(rng, FoodTypeBread)
				_ = b.AddFood(food)
				return food
			},
			wantError: false,
		},
		{
			name: "Remove existing Weapon",
			setupFunc: func(b *Backpack) any {
				weapon := createTestWeapon(rng, WeaponTypeSword)
				_ = b.AddWeapon(weapon)
				return weapon
			},
			wantError: false,
		},
		{
			name: "Remove non-existing item returns error",
			setupFunc: func(b *Backpack) any {
				return createTestElixir(rng, ElixirTypeStrength)
			},
			wantError: true,
		},
		{
			name: "Remove invalid type returns error",
			setupFunc: func(b *Backpack) any {
				return "not an item"
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backpack := NewBackpack()
			item := tt.setupFunc(backpack)

			err := backpack.RemoveItem(item)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

func TestBackpack_AddElixir(t *testing.T) {
	tests := []struct {
		name         string
		elixirType   ElixirType
		count        int
		wantItemsNum uint
	}{
		{
			name:         "Add single Strength elixir",
			elixirType:   ElixirTypeStrength,
			count:        1,
			wantItemsNum: 1,
		},
		{
			name:         "Add multiple Strength elixirs",
			elixirType:   ElixirTypeStrength,
			count:        3,
			wantItemsNum: 3,
		},
		{
			name:         "Add multiple different elixirs",
			elixirType:   ElixirTypeAgility,
			count:        2,
			wantItemsNum: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backpack := NewBackpack()
			rng := utils.NewRandomWithSeed(testBackpackSeed)

			for i := 0; i < tt.count; i++ {
				elixir := createTestElixir(rng, tt.elixirType)
				err := backpack.AddElixir(elixir)

				if err != nil {
					t.Fatalf("Failed to add elixir %d: %v", i, err)
				}
			}

			if backpack.ItemsNum != tt.wantItemsNum {
				t.Errorf("Expected ItemsNum %d, got %d", tt.wantItemsNum, backpack.ItemsNum)
			}

			if backpack.Elixirs.Len() != tt.count {
				t.Errorf("Expected %d elixirs in map, got %d", tt.count, backpack.Elixirs.Len())
			}
		})
	}
}

func TestBackpack_AddElixir_NilIsNothing(t *testing.T) {
	backpack := NewBackpack()
	err := backpack.AddElixir(nil)

	if err != nil {
		t.Fatalf("Add nil elixir must do nothing, without error, but got error: %v", err)
	}

	if backpack.ItemsNum != 0 {
		t.Errorf("Expected ItemsNum 0, got %d", backpack.ItemsNum)
	}

	if backpack.Elixirs.Len() != 0 {
		t.Errorf("Expected 0 elixirs in map, got %d", backpack.Elixirs.Len())
	}
}

func TestBackpack_RemoveElixir(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func(*Backpack, *utils.Random) *Elixir
		wantError    bool
		wantItemsNum uint
		wantListLen  int
	}{
		{
			name: "Remove single elixir",
			setupFunc: func(b *Backpack, rng *utils.Random) *Elixir {
				elixir := createTestElixir(rng, ElixirTypeStrength)
				_ = b.AddElixir(elixir)
				return elixir
			},
			wantError:    false,
			wantItemsNum: 0,
			wantListLen:  0,
		},
		{
			name: "Remove first of multiple elixirs",
			setupFunc: func(b *Backpack, rng *utils.Random) *Elixir {
				var firstElixir *Elixir
				for i := 0; i < 3; i++ {
					elixir := createTestElixir(rng, ElixirTypeStrength)
					if i == 0 {
						firstElixir = elixir
					}
					_ = b.AddElixir(elixir)
				}
				return firstElixir
			},
			wantError:    false,
			wantItemsNum: 2,
			wantListLen:  2,
		},
		{
			name: "Remove non-existing elixir returns error",
			setupFunc: func(b *Backpack, rng *utils.Random) *Elixir {
				elixir := createTestElixir(rng, ElixirTypeStrength)
				_ = b.AddElixir(elixir)
				// Возвращаем другой экземпляр, которого нет в рюкзаке
				return createTestElixir(rng, ElixirTypeAgility)
			},
			wantError:    true,
			wantItemsNum: 1,
			wantListLen:  1,
		},
		{
			name: "Remove from empty backpack returns error",
			setupFunc: func(b *Backpack, rng *utils.Random) *Elixir {
				return createTestElixir(rng, ElixirTypeStrength)
			},
			wantError:    true,
			wantItemsNum: 0,
			wantListLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backpack := NewBackpack()
			rng := utils.NewRandomWithSeed(testBackpackSeed)
			elixir := tt.setupFunc(backpack, rng)

			err := backpack.RemoveElixir(elixir)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				var expectedErr ItemIsNotInBackpackError
				if !errors.As(err, &expectedErr) {
					t.Errorf("Expected ItemIsNotInBackpackError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}

			if backpack.ItemsNum != tt.wantItemsNum {
				t.Errorf("Expected ItemsNum %d, got %d", tt.wantItemsNum, backpack.ItemsNum)
			}

			if backpack.Elixirs.Len() != tt.wantListLen {
				t.Errorf("Expected list length %d, got %d", tt.wantListLen, backpack.Elixirs.Len())
			}
		})
	}
}

func TestBackpack_RemoveElixir_NilIsError(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)
	_ = backpack.AddElixir(createTestElixir(rng, ElixirTypeStrength))

	err := backpack.RemoveElixir(nil)

	if err == nil {
		t.Fatalf("Remove nil elixir must return error, got nil")
	}

	if backpack.ItemsNum != 1 {
		t.Errorf("Expected ItemsNum 1, got %d", backpack.ItemsNum)
	}

	if backpack.Elixirs.Len() != 1 {
		t.Errorf("Expected 1 elixirs in map, got %d", backpack.Elixirs.Len())
	}
}

func TestBackpack_AddScroll(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	scroll := createTestScroll(rng, ScrollTypeStrength)
	err := backpack.AddScroll(scroll)

	if err != nil {
		t.Fatalf("Failed to add scroll: %v", err)
	}

	if backpack.ItemsNum != 1 {
		t.Errorf("Expected ItemsNum 1, got %d", backpack.ItemsNum)
	}

	if backpack.Scrolls.Len() != 1 {
		t.Errorf("Expected 1 scroll in list, got %d", backpack.Scrolls.Len())
	}
}

func TestBackpack_AddScroll_NilIsNothing(t *testing.T) {
	backpack := NewBackpack()
	err := backpack.AddScroll(nil)

	if err != nil {
		t.Fatalf("Add nil scroll must do nothing, without error, but got error: %v", err)
	}

	if backpack.ItemsNum != 0 {
		t.Errorf("Expected ItemsNum 0, got %d", backpack.ItemsNum)
	}

	if backpack.Scrolls.Len() != 0 {
		t.Errorf("Expected 0 scrolls in map, got %d", backpack.Scrolls.Len())
	}
}

func TestBackpack_RemoveScroll(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	scroll := createTestScroll(rng, ScrollTypeStrength)
	_ = backpack.AddScroll(scroll)

	err := backpack.RemoveScroll(scroll)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if backpack.ItemsNum != 0 {
		t.Errorf("Expected ItemsNum 0, got %d", backpack.ItemsNum)
	}

	if backpack.Scrolls.Len() != 0 {
		t.Errorf("Expected 0 scrolls in list, got %d", backpack.Scrolls.Len())
	}
}

func TestBackpack_RemoveScroll_NilIsError(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)
	_ = backpack.AddScroll(createTestScroll(rng, ScrollTypeStrength))

	err := backpack.RemoveScroll(nil)

	if err == nil {
		t.Fatalf("Remove nil scroll must return error, got nil")
	}

	if backpack.ItemsNum != 1 {
		t.Errorf("Expected ItemsNum 1, got %d", backpack.ItemsNum)
	}

	if backpack.Scrolls.Len() != 1 {
		t.Errorf("Expected 1 scroll in list, got %d", backpack.Scrolls.Len())
	}
}

func TestBackpack_AddFood(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	food := createTestFood(rng, FoodTypeBread)
	err := backpack.AddFood(food)

	if err != nil {
		t.Fatalf("Failed to add food: %v", err)
	}

	if backpack.ItemsNum != 1 {
		t.Errorf("Expected ItemsNum 1, got %d", backpack.ItemsNum)
	}

	if backpack.Foods.Len() != 1 {
		t.Errorf("Expected 1 food in list, got %d", backpack.Foods.Len())
	}
}

func TestBackpack_AddFood_NilIsNothing(t *testing.T) {
	backpack := NewBackpack()
	err := backpack.AddFood(nil)

	if err != nil {
		t.Fatalf("Add nil food must do nothing, without error, but got error: %v", err)
	}

	if backpack.ItemsNum != 0 {
		t.Errorf("Expected ItemsNum 0, got %d", backpack.ItemsNum)
	}

	if backpack.Foods.Len() != 0 {
		t.Errorf("Expected 0 foods in map, got %d", backpack.Foods.Len())
	}
}

func TestBackpack_RemoveFood(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	food := createTestFood(rng, FoodTypeBread)
	_ = backpack.AddFood(food)

	err := backpack.RemoveFood(food)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if backpack.ItemsNum != 0 {
		t.Errorf("Expected ItemsNum 0, got %d", backpack.ItemsNum)
	}

	if backpack.Foods.Len() != 0 {
		t.Errorf("Expected 0 foods in list, got %d", backpack.Foods.Len())
	}
}

func TestBackpack_RemoveFood_NilIsError(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)
	_ = backpack.AddFood(createTestFood(rng, FoodTypeBread))

	err := backpack.RemoveFood(nil)

	if err == nil {
		t.Fatalf("Remove nil food must return error, got nil")
	}

	if backpack.ItemsNum != 1 {
		t.Errorf("Expected ItemsNum 1, got %d", backpack.ItemsNum)
	}

	if backpack.Foods.Len() != 1 {
		t.Errorf("Expected 1 food in list, got %d", backpack.Foods.Len())
	}
}

func TestBackpack_AddWeapon(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	weapon := createTestWeapon(rng, WeaponTypeSword)
	err := backpack.AddWeapon(weapon)

	if err != nil {
		t.Fatalf("Failed to add weapon: %v", err)
	}

	if backpack.ItemsNum != 1 {
		t.Errorf("Expected ItemsNum 1, got %d", backpack.ItemsNum)
	}

	if backpack.Weapons.Len() != 1 {
		t.Errorf("Expected 1 weapon in list, got %d", backpack.Weapons.Len())
	}
}

func TestBackpack_AddWeapon_NilIsNothing(t *testing.T) {
	backpack := NewBackpack()
	err := backpack.AddWeapon(nil)

	if err != nil {
		t.Fatalf("Add nil weapon must do nothing, without error, but got error: %v", err)
	}

	if backpack.ItemsNum != 0 {
		t.Errorf("Expected ItemsNum 0, got %d", backpack.ItemsNum)
	}

	if backpack.Weapons.Len() != 0 {
		t.Errorf("Expected 0 weapons in map, got %d", backpack.Weapons.Len())
	}
}

func TestBackpack_RemoveWeapon(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	weapon := createTestWeapon(rng, WeaponTypeSword)
	_ = backpack.AddWeapon(weapon)

	err := backpack.RemoveWeapon(weapon)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if backpack.ItemsNum != 0 {
		t.Errorf("Expected ItemsNum 0, got %d", backpack.ItemsNum)
	}

	if backpack.Weapons.Len() != 0 {
		t.Errorf("Expected 0 weapons in list, got %d", backpack.Weapons.Len())
	}
}

func TestBackpack_RemoveWeapon_NilIsError(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)
	_ = backpack.AddWeapon(createTestWeapon(rng, WeaponTypeSword))

	err := backpack.RemoveWeapon(nil)

	if err == nil {
		t.Fatalf("Remove nil weapon must return error, got nil")
	}

	if backpack.ItemsNum != 1 {
		t.Errorf("Expected ItemsNum 1, got %d", backpack.ItemsNum)
	}

	if backpack.Weapons.Len() != 1 {
		t.Errorf("Expected 1 weapon in list, got %d", backpack.Weapons.Len())
	}
}

func TestBackpack_AddTreasure(t *testing.T) {
	tests := []struct {
		name           string
		treasures      []int32
		wantTotalValue int32
		wantItemsNum   uint
	}{
		{
			name:           "Add single treasure",
			treasures:      []int32{100},
			wantTotalValue: 100,
			wantItemsNum:   0, // Treasures don't count as items
		},
		{
			name:           "Add multiple treasures",
			treasures:      []int32{100, 200, 50},
			wantTotalValue: 350,
			wantItemsNum:   0,
		},
		{
			name:           "Add zero-value treasure",
			treasures:      []int32{0},
			wantTotalValue: 0,
			wantItemsNum:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backpack := NewBackpack()
			rng := utils.NewRandomWithSeed(testBackpackSeed)

			for _, value := range tt.treasures {
				box := primitives.Box{
					Point: primitives.Point2D[int]{X: 0, Y: 0},
					Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
				}
				config := TreasureConfig{
					Type:        "TestTreasure",
					ValueRange:  TreasureValueRange{Min: value, Max: value},
					Description: "Test treasure",
				}
				treasure, err := newTreasureByConfig(rng, box, config)
				if err != nil {
					t.Fatalf("Failed to create treasure: %v", err)
				}
				backpack.AddTreasure(treasure)
			}

			if backpack.Treasures != tt.wantTotalValue {
				t.Errorf("Expected Treasures %d, got %d", tt.wantTotalValue, backpack.Treasures)
			}

			if backpack.ItemsNum != tt.wantItemsNum {
				t.Errorf("Expected ItemsNum %d, got %d (treasures don't count)", tt.wantItemsNum, backpack.ItemsNum)
			}
		})
	}
}

func TestBackpack_AddTreasure_NilIsNothing(t *testing.T) {
	backpack := NewBackpack()
	initialTreasures := backpack.Treasures

	backpack.AddTreasure(nil)

	if backpack.Treasures != initialTreasures {
		t.Errorf("Add nil treasure must do nothing, expected Treasures=%d, got %d", initialTreasures, backpack.Treasures)
	}

	if backpack.ItemsNum != 0 {
		t.Errorf("Expected ItemsNum 0, got %d", backpack.ItemsNum)
	}
}

func TestBackpack_MixedItems(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	// Добавляем разные типы предметов
	elixir := createTestElixir(rng, ElixirTypeStrength)
	scroll := createTestScroll(rng, ScrollTypeStrength)
	food := createTestFood(rng, FoodTypeBread)
	weapon := createTestWeapon(rng, WeaponTypeSword)
	treasure := createTestTreasure(rng, TreasureTypeGold)

	if err := backpack.AddElixir(elixir); err != nil {
		t.Fatalf("Failed to add elixir: %v", err)
	}
	if err := backpack.AddScroll(scroll); err != nil {
		t.Fatalf("Failed to add scroll: %v", err)
	}
	if err := backpack.AddFood(food); err != nil {
		t.Fatalf("Failed to add food: %v", err)
	}
	if err := backpack.AddWeapon(weapon); err != nil {
		t.Fatalf("Failed to add weapon: %v", err)
	}
	backpack.AddTreasure(treasure)

	// Проверяем состояние
	if backpack.ItemsNum != 4 {
		t.Errorf("Expected ItemsNum 4, got %d (treasure doesn't count)", backpack.ItemsNum)
	}

	if backpack.Elixirs.Len() != 1 {
		t.Errorf("Expected 1 elixir in list, got %d", backpack.Elixirs.Len())
	}

	if backpack.Scrolls.Len() != 1 {
		t.Errorf("Expected 1 scroll in list, got %d", backpack.Scrolls.Len())
	}

	if backpack.Foods.Len() != 1 {
		t.Errorf("Expected 1 food in list, got %d", backpack.Foods.Len())
	}

	if backpack.Weapons.Len() != 1 {
		t.Errorf("Expected 1 weapon in list, got %d", backpack.Weapons.Len())
	}

	if backpack.Treasures <= 0 {
		t.Errorf("Expected positive treasure value, got %d", backpack.Treasures)
	}
}

func TestBackpack_FillAndEmpty(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	// Заполняем рюкзак эликсирами и сохраняем ссылки
	elixirs := make([]*Elixir, 0, DefaultBackpackCapacity)
	itemsAdded := uint(0)
	for i := uint(0); i < DefaultBackpackCapacity; i++ {
		elixir := createTestElixir(rng, ElixirTypeStrength)
		elixirs = append(elixirs, elixir)
		if err := backpack.AddElixir(elixir); err != nil {
			t.Fatalf("Failed to add item %d: %v", i, err)
		}
		itemsAdded++
	}

	if !backpack.IsFull() {
		t.Error("Expected backpack to be full")
	}

	if backpack.ItemsNum != itemsAdded {
		t.Errorf("Expected ItemsNum %d, got %d", itemsAdded, backpack.ItemsNum)
	}

	// Опустошаем рюкзак, удаляя конкретные экземпляры
	for _, elixir := range elixirs {
		err := backpack.RemoveElixir(elixir)
		if err != nil {
			t.Fatalf("Failed to remove elixir: %v", err)
		}
	}

	if !backpack.IsEmpty() {
		t.Errorf("Expected backpack to be empty, but ItemsNum=%d", backpack.ItemsNum)
	}

	if backpack.Elixirs.Len() != 0 {
		t.Errorf("Expected empty elixirs list, got length %d", backpack.Elixirs.Len())
	}
}

func TestBackpack_MultipleItemsSameType(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	// Добавляем несколько предметов одного типа и сохраняем ссылки
	count := 5
	elixirs := make([]*Elixir, 0, count)
	for i := 0; i < count; i++ {
		elixir := createTestElixir(rng, ElixirTypeStrength)
		elixirs = append(elixirs, elixir)
		if err := backpack.AddElixir(elixir); err != nil {
			t.Fatalf("Failed to add elixir %d: %v", i, err)
		}
	}

	if backpack.Elixirs.Len() != count {
		t.Errorf("Expected %d elixirs in list, got %d", count, backpack.Elixirs.Len())
	}

	// Удаляем по одному (удаляем конкретные экземпляры)
	for i := 0; i < count; i++ {
		err := backpack.RemoveElixir(elixirs[i])
		if err != nil {
			t.Fatalf("Failed to remove elixir at iteration %d: %v", i, err)
		}

		expectedCount := count - i - 1
		actualCount := backpack.Elixirs.Len()
		if actualCount != expectedCount {
			t.Errorf("After removal %d: expected %d elixirs, got %d",
				i+1, expectedCount, actualCount)
		}
	}

	// Проверяем, что список пуст
	if backpack.Elixirs.Len() != 0 {
		t.Errorf("Expected empty list, got length %d", backpack.Elixirs.Len())
	}
}

func TestBackpack_RemoveAfterListCleanup(t *testing.T) {
	backpack := NewBackpack()
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	// Добавляем и удаляем один предмет
	elixir := createTestElixir(rng, ElixirTypeStrength)
	_ = backpack.AddElixir(elixir)
	_ = backpack.RemoveElixir(elixir)

	// Попытка удалить снова должна вернуть ошибку
	err := backpack.RemoveElixir(elixir)

	if err == nil {
		t.Error("Expected error when removing non-existing item")
	}

	var expectedErr ItemIsNotInBackpackError
	if !errors.As(err, &expectedErr) {
		t.Errorf("Expected ItemIsNotInBackpackError, got %T", err)
	}

	// Проверяем, что список пуст
	if backpack.Elixirs.Len() != 0 {
		t.Errorf("Expected empty list, got length %d", backpack.Elixirs.Len())
	}
}

// Бенчмарки

func BenchmarkBackpack_AddItem(b *testing.B) {
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	b.Run("AddElixir", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			backpack := NewBackpack()
			elixir := createTestElixir(rng, ElixirTypeStrength)
			_ = backpack.AddItem(elixir)
		}
	})

	b.Run("AddScroll", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			backpack := NewBackpack()
			scroll := createTestScroll(rng, ScrollTypeStrength)
			_ = backpack.AddItem(scroll)
		}
	})

	b.Run("AddWeapon", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			backpack := NewBackpack()
			weapon := createTestWeapon(rng, WeaponTypeSword)
			_ = backpack.AddItem(weapon)
		}
	})
}

func BenchmarkBackpack_RemoveItem(b *testing.B) {
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	b.Run("RemoveElixir", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			backpack := NewBackpack()
			elixir := createTestElixir(rng, ElixirTypeStrength)
			_ = backpack.AddElixir(elixir)
			b.StartTimer()

			_ = backpack.RemoveItem(elixir)
		}
	})
}

func BenchmarkBackpack_FillBackpack(b *testing.B) {
	rng := utils.NewRandomWithSeed(testBackpackSeed)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		backpack := NewBackpack()
		for j := uint(0); j < DefaultBackpackCapacity; j++ {
			elixir := createTestElixir(rng, ElixirTypeStrength)
			_ = backpack.AddElixir(elixir)
		}
	}
}
