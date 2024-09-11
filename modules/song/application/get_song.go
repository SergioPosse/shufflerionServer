package application

import (
    "shufflerion/modules/song/domain"
)

type GetSongsUseCase struct {
    SongRepo domain.SongRepository
}

func NewGetSongsUseCase(repo domain.SongRepository) *GetSongsUseCase {
    return &GetSongsUseCase{SongRepo: repo}
}

func (uc *GetSongsUseCase) Execute() (*domain.Song, error) {
    return uc.SongRepo.GetRandomSong()
}