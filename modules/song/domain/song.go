package domain

type Song struct {
	Title    string
	Artist   string
	Url      string
	Image    string
	Duration int
    Explicit bool
}

type SongsRepository interface {
	GetRandomSongs(accessToken1 string, accessToken2 string) ([]Song, error)
}