package routes

import (
    "net/http"
    "shufflerion/infrastructure/controllers"
    "github.com/rs/cors"
)

func RegisterRoutes(authController *controllers.AuthController, songsController *controllers.SongController, sessionController *controllers.SessionController) {
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
    })

    mux := http.NewServeMux()

    mux.HandleFunc("/songs/random", songsController.GetRandomSongs)
    mux.HandleFunc("/auth/tokens", authController.GetAccessTokens)
    mux.HandleFunc("/session/create", sessionController.CreateSession)
    mux.HandleFunc("/session/{id}", sessionController.GetSessionById)
    mux.HandleFunc("/session/update", sessionController.UpdateSession)

    handler := c.Handler(mux)

    server := &http.Server{
        Addr:    ":8080",
        Handler: handler,
    }
    server.ListenAndServe()
}
