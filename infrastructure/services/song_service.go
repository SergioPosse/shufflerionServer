package services

import (
	"fmt"
	"shufflerion/modules/song/domain"
	domainSession "shufflerion/modules/session/domain"
)

type SongsService struct{
	spotifyService domainSession.SpotifyService
}

func NewSongsService(spotifyService domainSession.SpotifyService) domain.SongsRepository {
	return &SongsService{spotifyService: spotifyService}
}

func (s *SongsService) GetRandomSongs(accessToken1 string, accessToken2 string) ([]domain.Song, error) {

	tracks1, err := s.spotifyService.FetchRandomSongs(accessToken1, 10)
	if len(accessToken2) > 20 {
		tracks2, err2 := s.spotifyService.FetchRandomSongs(accessToken2, 10)

		if err != nil || err2 != nil {
			var combinedErr error
			if err != nil && err2 != nil {
				combinedErr = fmt.Errorf("error getting track: %v, %v", err, err2)
			} else if err != nil {
				combinedErr = fmt.Errorf("error getting track: %v", err)
			} else {
				combinedErr = fmt.Errorf("error getting track: %v", err2)
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
	totalLength := len(arrayA) + len(arrayB)
	response := make([]domain.Song, 0, totalLength)

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
