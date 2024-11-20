package services

import (
	"fmt"
	"shufflerion/modules/song/domain"
)

type SongsService struct{}

func NewSongsService() domain.SongsRepository {
	return &SongsService{}
}

func (s *SongsService) GetRandomSongs(accessToken1 string, accessToken2 string) ([]domain.Song, error) {

	fmt.Println("Token:", accessToken1)
	fmt.Println("Token2:", accessToken2)

	tracks1, err := FetchRandomSongs(accessToken1, 3)
	if len(accessToken2) > 20 {
		tracks2, err2 := FetchRandomSongs(accessToken2, 3)

		if err != nil || err2 != nil {
			var combinedErr error
			if err != nil && err2 != nil {
				combinedErr = fmt.Errorf("error al obtener tracks: %v, %v", err, err2)
			} else if err != nil {
				combinedErr = fmt.Errorf("error al obtener tracks: %v", err)
			} else {
				combinedErr = fmt.Errorf("error al obtener tracks: %v", err2)
			}
			return nil, combinedErr
		}
		responseIntercalated := interleaveArrays(tracks1, tracks2)
		return responseIntercalated, nil
	}

	if err != nil {
		return nil, err
	}

	return tracks1, nil
}

func interleaveArrays(arrayA, arrayB []domain.Song) []domain.Song {
	// Determinamos el tamaño total del nuevo array intercalado
	totalLength := len(arrayA) + len(arrayB)
	response := make([]domain.Song, 0, totalLength)

	// Iteramos mientras haya elementos en ambos arrays
	i, j := 0, 0
	for i < len(arrayA) || j < len(arrayB) {
		if i < len(arrayA) {
			response = append(response, arrayA[i])
			i++
		}
		if j < len(arrayB) {
			response = append(response, arrayB[j])
			j++
		}
	}
	return response
}
