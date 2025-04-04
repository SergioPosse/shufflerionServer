package controllers

import (
    "encoding/json"
    "net/http"
    songsModule "shufflerion/modules/song/application"
)

type SongControllerRequestBody struct {
    AccessToken1 string `json:"access_token1"`
    AccessToken2 string `json:"access_token2"`
}

type SongController struct {
    GetSongsUC         *songsModule.GetRandomSongsUseCase
}

func NewSongsController(getSongsUC *songsModule.GetRandomSongsUseCase) *SongController {
    return &SongController{GetSongsUC: getSongsUC}
}

func (c *SongController) GetRandomSongs(w http.ResponseWriter, r *http.Request) {

    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var requestBody SongControllerRequestBody

    err := json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil {
        http.Error(w, "Error reading request body", http.StatusBadRequest)
        return
    }

    accessToken1 := requestBody.AccessToken1
    accessToken2 := requestBody.AccessToken2

    if accessToken1 == "" {
        http.Error(w, "ccessToken1 is required", http.StatusBadRequest)
        return
    }

    songs, err2 := c.GetSongsUC.Execute(accessToken1, accessToken2)
    if err2 != nil {
        http.Error(w, err2.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(songs)
}