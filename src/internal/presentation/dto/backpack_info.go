package dto

import (
	"container/list"
	"gogue/internal/model/items"
)

// BackpackInfo содержит информацию о содержимом рюкзака для отображения в UI
type BackpackInfo struct {
	Capacity  uint
	ItemsNum  uint
	Elixirs   []*ItemInfo
	Scrolls   []*ItemInfo
	Foods     []*ItemInfo
	Weapons   []*ItemInfo
	Treasures int32
}

// ConvertBackpackToDto преобразует модель рюкзака в DTO для отображения
func ConvertBackpackToDto(backpack *items.Backpack) *BackpackInfo {
	if backpack == nil {
		return nil
	}

	return &BackpackInfo{
		Capacity:  backpack.Capacity,
		ItemsNum:  backpack.ItemsNum,
		Elixirs:   convertElixirListToDto(backpack.Elixirs),
		Scrolls:   convertScrollListToDto(backpack.Scrolls),
		Foods:     convertFoodListToDto(backpack.Foods),
		Weapons:   convertWeaponListToDto(backpack.Weapons),
		Treasures: backpack.Treasures,
	}
}

// convertElixirListToDto преобразует список эликсиров в слайс DTO
func convertElixirListToDto(itemList *list.List) []*ItemInfo {
	if itemList == nil {
		return []*ItemInfo{}
	}

	result := make([]*ItemInfo, 0, itemList.Len())
	for elem := itemList.Front(); elem != nil; elem = elem.Next() {
		item := elem.Value.(*items.Elixir)
		info := &ItemInfo{
			Type:       "Elixir",
			Name:       item.Name,
			EffectInfo: ConvertEffectToDto(item.Effect),
			CanUse:     true,
			CanTake:    true,
		}
		result = append(result, info)
	}
	return result
}

// convertScrollListToDto преобразует список свитков в слайс DTO
func convertScrollListToDto(itemList *list.List) []*ItemInfo {
	if itemList == nil {
		return []*ItemInfo{}
	}

	result := make([]*ItemInfo, 0, itemList.Len())
	for elem := itemList.Front(); elem != nil; elem = elem.Next() {
		item := elem.Value.(*items.Scroll)
		info := &ItemInfo{
			Type:       "Scroll",
			Name:       item.Name,
			EffectInfo: ConvertEffectToDto(item.Effect),
			CanUse:     true,
			CanTake:    true,
		}
		result = append(result, info)
	}
	return result
}

// convertFoodListToDto преобразует список еды в слайс DTO
func convertFoodListToDto(itemList *list.List) []*ItemInfo {
	if itemList == nil {
		return []*ItemInfo{}
	}

	result := make([]*ItemInfo, 0, itemList.Len())
	for elem := itemList.Front(); elem != nil; elem = elem.Next() {
		item := elem.Value.(*items.Food)
		info := &ItemInfo{
			Type:       "Food",
			Name:       item.Name,
			EffectInfo: ConvertEffectToDto(item.Effect),
			CanUse:     true,
			CanTake:    true,
		}
		result = append(result, info)
	}
	return result
}

// convertWeaponListToDto преобразует список оружия в слайс DTO
func convertWeaponListToDto(itemList *list.List) []*ItemInfo {
	if itemList == nil {
		return []*ItemInfo{}
	}

	result := make([]*ItemInfo, 0, itemList.Len())
	for elem := itemList.Front(); elem != nil; elem = elem.Next() {
		item := elem.Value.(*items.Weapon)
		info := &ItemInfo{
			Type:       "Weapon",
			Name:       item.Name,
			EffectInfo: ConvertEffectToDto(item.Effect),
			CanUse:     true,
			CanTake:    true,
		}
		result = append(result, info)
	}
	return result
}
