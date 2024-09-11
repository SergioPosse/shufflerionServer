package controllers

import (
    "encoding/json"
    "net/http"
    "shufflerion/modules/song/application"
)

type SongController struct {
    GetSongsUC *application.GetSongsUseCase
}

func NewSongController(usecase *application.GetSongsUseCase) *SongController {
    return &SongController{GetSongsUC: usecase}
}

func (c *SongController) GetRandomSong(w http.ResponseWriter, r *http.Request) {
    song, err := c.GetSongsUC.Execute()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(song)
}