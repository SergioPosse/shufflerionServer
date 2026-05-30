package services

import (
	"fmt"
	"shufflerion/modules/song/domain"
	domainSession "shufflerion/modules/session/domain"
)

type SongsService struct {
	spotifyService domainSession.SpotifyService
}

func NewSongsService(spotifyService domainSession.SpotifyService) domain.SongsRepository {
	return &SongsService{spotifyService: spotifyService}
}

func (s *SongsService) GetRandomSongs(accessToken1 string, accessToken2 string) ([]domain.Song, error) {
	fmt.Printf("🎲 GetRandomSongs: fetching songs for user1\n")
	tracks1, err := s.spotifyService.FetchRandomSongs(accessToken1, 15)
	if err != nil {
		fmt.Printf("❌ GetRandomSongs: error fetching user1 songs: %v\n", err)
	} else {
		fmt.Printf("✅ GetRandomSongs: got %d songs from user1\n", len(tracks1))
	}

	if accessToken2 != "" {
		fmt.Printf("🎲 GetRandomSongs: fetching songs for user2\n")
		tracks2, err2 := s.spotifyService.FetchRandomSongs(accessToken2, 15)
		if err2 != nil {
			fmt.Printf("❌ GetRandomSongs: error fetching user2 songs: %v\n", err2)
		} else {
			fmt.Printf("✅ GetRandomSongs: got %d songs from user2\n", len(tracks2))
		}

		if err != nil && err2 != nil {
			return nil, fmt.Errorf("error getting tracks from both users: user1=%v, user2=%v", err, err2)
		}

		if err2 == nil {
			mixed := interleaveArrays(tracks1, tracks2)
			fmt.Printf("🎵 GetRandomSongs: mixed playlist has %d songs (%d user1 + %d user2)\n", len(mixed), len(tracks1), len(tracks2))
			return mixed, nil
		}
	}

	if err != nil {
		return nil, err
	}

	fmt.Printf("🎵 GetRandomSongs: single-user playlist has %d songs\n", len(tracks1))
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
