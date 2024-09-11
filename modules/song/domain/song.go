package domain

type Song struct {
    ID     string
    Title  string
    Artist string
}

type SongRepository interface {
    GetRandomSong() (*Song, error)
}