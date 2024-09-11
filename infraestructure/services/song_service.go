package services

import (
    "errors"
    "math/rand"
    "shufflerion/modules/song/domain"
)

type SongService struct {
    Songs []domain.Song
}

func NewSongService() *SongService {
    return &SongService{
        Songs: []domain.Song{
            {ID: "1", Title: "Song One", Artist: "Artist A"},
            {ID: "2", Title: "Song Two", Artist: "Artist B"},
        },
    }
}

func (s *SongService) GetRandomSong() (*domain.Song, error) {
    if len(s.Songs) == 0 {
        return nil, errors.New("no songs available")
    }
    randomIndex := rand.Intn(len(s.Songs))
    return &s.Songs[randomIndex], nil
}