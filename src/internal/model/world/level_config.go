package world

import (
	"math"

	"gogue/internal/model/primitives"
)

type LevelConfig struct {
	MapSize primitives.Size2D[uint]

	// Параметры комнат
	RoomMinWidth   int
	RoomMinHeight  int
	MapMaxWidth    uint
	MapMaxHeight   uint
	MinRoomPadding int

	// Параметры генерации мира
	MaxExtraPassageCount int
	NumberXYSections     int
	RoomsCount           int
	ItemsSpawnConfig     ItemSpawnConfig
	EnemiesSpawnConfig   EnemySpawnConfig

	// Конфигурация сложности
	DifficultyScaling DifficultyScaling
}

type ItemSpawnConfig struct {
	FoodsQuantity   uint
	ElixirsQuantity uint
	ScrollsQuantity uint
	WeaponsQuantity uint
}

type EnemySpawnConfig struct {
	Quantity uint

	// Множитель атрибутов врагов (1.0 = базовые значения).
	// Увеличивается с каждым уровнем подземелья.
	AttributeMultiplier float64
}

// TreasureDropConfig определяет проценты выпадения типов сокровищ с побежденных врагов.
// Сумма всех полей должна быть равна 100.
type TreasureDropConfig struct {
	GoldPercent     uint
	GemPercent      uint
	ArtifactPercent uint
	MysteryPercent  uint
}

// DifficultyScaling содержит параметры, управляющие тем, как сложность
// масштабируется с ростом номера уровня подземелья.
type DifficultyScaling struct {
	// Враги: базовое количество + прирост за уровень
	BaseEnemyCount       uint
	EnemyCountPerLevel   uint
	MaxEnemyCount        uint
	EnemyAttributeGrowth float64 // Прирост множителя атрибутов за уровень (например, 0.15 => +15% за уровень)

	// Предметы: базовое количество - убыль за уровень (с минимальным порогом)
	BaseFoods     uint
	BaseElixirs   uint
	BaseScrolls   uint
	BaseWeapons   uint
	ItemDecayRate float64 // Доля убыли предметов за уровень (0.08 => -8% от базы за уровень)
	MinFoods      uint
	MinElixirs    uint
	MinScrolls    uint
	MinWeapons    uint

	// Сокровища: сдвиг процентов от Gold к более ценным типам
	BaseTreasureDrop       TreasureDropConfig
	GemGrowthPerLevel      float64 // Прирост % гемов за уровень
	ArtifactGrowthPerLevel float64 // Прирост % артефактов за уровень
}

func DefaultDifficultyScaling() DifficultyScaling {
	return DifficultyScaling{
		// Враги
		BaseEnemyCount:       8,
		EnemyCountPerLevel:   1,
		MaxEnemyCount:        20,
		EnemyAttributeGrowth: 0.08,

		// Предметы
		BaseFoods:     7,
		BaseElixirs:   5,
		BaseScrolls:   3,
		BaseWeapons:   2,
		ItemDecayRate: 0.05,
		MinFoods:      3,
		MinElixirs:    2,
		MinScrolls:    1,
		MinWeapons:    1,

		// Сокровища
		BaseTreasureDrop: TreasureDropConfig{
			GoldPercent:     65,
			GemPercent:      20,
			ArtifactPercent: 5,
			MysteryPercent:  10,
		},
		GemGrowthPerLevel:      1.5,
		ArtifactGrowthPerLevel: 0.8,
	}
}

func DefaultLevelConfig(mapSize primitives.Size2D[uint]) LevelConfig {
	scaling := DefaultDifficultyScaling()
	return LevelConfig{
		MapSize:              mapSize,
		RoomMinWidth:         3,
		RoomMinHeight:        3,
		MapMaxWidth:          200,
		MapMaxHeight:         150,
		MinRoomPadding:       1,
		MaxExtraPassageCount: 2,
		NumberXYSections:     3,
		RoomsCount:           9,
		ItemsSpawnConfig: ItemSpawnConfig{
			FoodsQuantity:   scaling.BaseFoods,
			ElixirsQuantity: scaling.BaseElixirs,
			ScrollsQuantity: scaling.BaseScrolls,
			WeaponsQuantity: scaling.BaseWeapons,
		},
		EnemiesSpawnConfig: EnemySpawnConfig{
			Quantity:            scaling.BaseEnemyCount,
			AttributeMultiplier: 1.0,
		},
		DifficultyScaling: scaling,
	}
}

// ScaleForLevel пересчитывает конфигурацию спавна врагов и предметов
// в соответствии с номером уровня подземелья.
// levelNumber начинается с 1 (первый уровень - базовая сложность).
func (c *LevelConfig) ScaleForLevel(levelNumber uint) {
	if levelNumber == 0 {
		levelNumber = 1
	}

	s := c.DifficultyScaling
	progression := float64(levelNumber - 1) // 0 для первого уровня

	// --- Враги ---
	enemyCount := s.BaseEnemyCount + s.EnemyCountPerLevel*uint(progression)
	if enemyCount > s.MaxEnemyCount {
		enemyCount = s.MaxEnemyCount
	}
	c.EnemiesSpawnConfig.Quantity = enemyCount
	c.EnemiesSpawnConfig.AttributeMultiplier = 1.0 + s.EnemyAttributeGrowth*progression

	// --- Предметы (экспоненциальное затухание, чтобы не обнулялось слишком резко) ---
	decayMultiplier := math.Pow(1.0-s.ItemDecayRate, progression)

	c.ItemsSpawnConfig.FoodsQuantity = clampUint(
		uint(math.Round(float64(s.BaseFoods)*decayMultiplier)),
		s.MinFoods, s.BaseFoods,
	)
	c.ItemsSpawnConfig.ElixirsQuantity = clampUint(
		uint(math.Round(float64(s.BaseElixirs)*decayMultiplier)),
		s.MinElixirs, s.BaseElixirs,
	)
	c.ItemsSpawnConfig.ScrollsQuantity = clampUint(
		uint(math.Round(float64(s.BaseScrolls)*decayMultiplier)),
		s.MinScrolls, s.BaseScrolls,
	)
	c.ItemsSpawnConfig.WeaponsQuantity = clampUint(
		uint(math.Round(float64(s.BaseWeapons)*decayMultiplier)),
		s.MinWeapons, s.BaseWeapons,
	)
}

// GetTreasureDropConfig возвращает конфигурацию выпадения сокровищ
// для указанного номера уровня. С ростом уровня Gold-процент уменьшается
// в пользу Gem и Artifact.
//
// Целочисленные проценты распределяются по методу наибольших остатков
// (Largest Remainder / Hare–Niemeyer), который гарантирует сумму = 100
// при минимальной погрешности округления.
func (c *LevelConfig) GetTreasureDropConfig(levelNumber uint) TreasureDropConfig {
	if levelNumber == 0 {
		levelNumber = 1
	}

	s := c.DifficultyScaling
	progression := float64(levelNumber - 1)
	base := s.BaseTreasureDrop

	// Вычисляем дробные проценты
	gemF := float64(base.GemPercent) + s.GemGrowthPerLevel*progression
	artifactF := float64(base.ArtifactPercent) + s.ArtifactGrowthPerLevel*progression
	mysteryF := float64(base.MysteryPercent)
	goldF := 100.0 - gemF - artifactF - mysteryF

	// Ограничиваем минимальный Gold
	const minGoldPercent = 5.0
	if goldF < minGoldPercent {
		excess := minGoldPercent - goldF
		goldF = minGoldPercent
		// Срезаем excess пропорционально из gem и artifact
		sum := gemF + artifactF
		if sum > 0 {
			gemF -= excess * (gemF / sum)
			artifactF -= excess * (artifactF / sum)
		}
	}

	// Защита от отрицательных дробных значений
	if gemF < 0 {
		gemF = 0
	}
	if artifactF < 0 {
		artifactF = 0
	}

	// Распределяем 100 целых процентов методом наибольших остатков (Hare–Niemeyer).
	// Порядок элементов: gold, gem, artifact, mystery.
	fractions := [4]float64{goldF, gemF, artifactF, mysteryF}
	result := distributeByLargestRemainder(fractions, 100)

	return TreasureDropConfig{
		GoldPercent:     uint(result[0]),
		GemPercent:      uint(result[1]),
		ArtifactPercent: uint(result[2]),
		MysteryPercent:  uint(result[3]),
	}
}

// distributeByLargestRemainder распределяет total целых единиц между N категориями
// пропорционально дробным значениям fractions, используя метод наибольших остатков.
//
// Алгоритм:
//  1. Каждой категории присваивается floor(fraction).
//  2. Оставшиеся единицы (total - сумма floor) раздаются категориям
//     с наибольшими дробными остатками (fraction - floor).
//
// Это стандартный метод пропорционального представительства (Hare–Niemeyer),
// гарантирующий, что сумма результатов = total.
func distributeByLargestRemainder(fractions [4]float64, total int) [4]int {
	var result [4]int
	var remainders [4]float64
	allocated := 0

	for i, f := range fractions {
		floored := int(math.Floor(f))
		if floored < 0 {
			floored = 0
		}
		result[i] = floored
		remainders[i] = f - float64(floored)
		allocated += floored
	}

	// Раздаём оставшиеся единицы по убыванию дробного остатка
	remaining := total - allocated
	for remaining > 0 {
		bestIdx := -1
		bestRem := -1.0
		for i, r := range remainders {
			if r > bestRem {
				bestRem = r
				bestIdx = i
			}
		}
		if bestIdx < 0 {
			break
		}
		result[bestIdx]++
		remainders[bestIdx] = -1.0 // Исключаем из дальнейшего выбора
		remaining--
	}

	return result
}

// clampUint ограничивает значение v диапазоном [minVal, maxVal].
func clampUint(v, minVal, maxVal uint) uint {
	if v < minVal {
		return minVal
	}
	if v > maxVal {
		return maxVal
	}
	return v
}
