package application

import (
    "shufflerion/modules/song/domain"
)

type GetRandomSongsUseCase struct {
    SongsRepo domain.SongsRepository
}

func NewGetSongsUseCase(repo domain.SongsRepository) *GetRandomSongsUseCase {
    return &GetRandomSongsUseCase{SongsRepo: repo}
}

func (uc *GetRandomSongsUseCase) Execute(accessToken1 string, accessToken2 string) ([]domain.Song, error) {
    return uc.SongsRepo.GetRandomSongs(accessToken1, accessToken2)
}