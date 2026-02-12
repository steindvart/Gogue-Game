package save

import (
	"encoding/json"
	"gogue/internal/model/primitives"
	"gogue/internal/model/world"
	"gogue/internal/presentation/dto"
	"os"
	"sort"
)

const ScoresTop = 20

func SaveGame(level *world.Level, size primitives.Size2D[uint], filename string) error {
	dto := ConvertLevelToDto(level, size)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(dto)
}

func LoadGame(filename string) (*world.Level, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var saveDto dto.GameSaveDto
	if err := json.Unmarshal(file, &saveDto); err != nil {
		return nil, err
	}

	return RestoreLevelFromDto(saveDto)
}

func SaveScore(treasures int32, filename string) error {
	var scoreDto dto.ScoreDto

	if file, err := os.Open(filename); err == nil {
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&scoreDto); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	scoreDto.Scores = append(scoreDto.Scores, treasures)
	sort.Slice(scoreDto.Scores, func(i, j int) bool {
		return scoreDto.Scores[i] > scoreDto.Scores[j]
	})

	if len(scoreDto.Scores) > ScoresTop {
		scoreDto.Scores = scoreDto.Scores[:ScoresTop]
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(scoreDto)
}

func LoadScore(filename string) (*dto.ScoreDto, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var scoreDto dto.ScoreDto
	if err := json.Unmarshal(file, &scoreDto); err != nil {
		return nil, err
	}

	return &scoreDto, nil
}
