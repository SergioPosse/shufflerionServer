package routes

import (
    "net/http"
    "shufflerion/infrastructure/controllers"
)

func RegisterRoutes(authController *controllers.AuthController, songsController *controllers.SongController) {
    http.HandleFunc("/songs/random", songsController.GetRandomSongs)
    http.HandleFunc("/auth/tokens", authController.GetAccessTokens)
}