package routes

import (
    "net/http"
    "shufflerion/infrastructure/controllers"
    "github.com/rs/cors"
)

func RegisterRoutes(authController *controllers.AuthController, songsController *controllers.SongController) {
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
    })

    mux := http.NewServeMux()

    mux.HandleFunc("/songs/random", songsController.GetRandomSongs)
    mux.HandleFunc("/auth/tokens", authController.GetAccessTokens)

    handler := c.Handler(mux)

    http.Handle("/", handler)
}
