package dto

import (
	"container/list"
	"gogue/internal/model/items"
)

type BackpackInfo struct {
	Capacity  uint
	ItemsNum  uint
	Elixirs   []*ItemInfo
	Scrolls   []*ItemInfo
	Foods     []*ItemInfo
	Weapons   []*ItemInfo
	Treasures int32
}

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

// @todo - подумать как избежать дублирования и обобщить код

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
