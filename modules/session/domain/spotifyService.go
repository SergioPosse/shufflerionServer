package domain

import (
	"shufflerion/modules/song/domain"
)

type SpotifyService interface {
	FetchRandomSongs(accessToken string, quantity int) ([]domain.Song, error)
	AddUserToApp(accessToken string, email string) (bool, error)
}