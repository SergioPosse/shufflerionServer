package routes

import (
	"log"
	"net/http"
	"shufflerion/infrastructure/controllers"
	"shufflerion/infrastructure/server"

	"github.com/rs/cors"
)

func RegisterRoutes(authController *controllers.AuthController, songsController *controllers.SongController, sessionController *controllers.SessionController, wsServer *server.WebSocketServer) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	mux := http.NewServeMux()

	// Rutas HTTP
	mux.HandleFunc("/songs/random", songsController.GetRandomSongs)
	mux.HandleFunc("/auth/tokens", authController.GetAccessTokens)
	mux.HandleFunc("/session/create", sessionController.CreateSession)
	mux.HandleFunc("/session/{id}", sessionController.GetSessionById)
	mux.HandleFunc("/session/update", sessionController.UpdateSession)

	// Ruta para WebSocket
	mux.HandleFunc("/session/socket", wsServer.HandleConnection)

	handler := c.Handler(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	log.Println("Server started at http://localhost:8080")
	server.ListenAndServe()
}