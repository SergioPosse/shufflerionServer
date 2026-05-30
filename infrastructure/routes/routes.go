package routes

import (
	"fmt"
	"net/http"
	"shufflerion/infrastructure/controllers"
	"shufflerion/infrastructure/server"
)

func RegisterRoutes(authController *controllers.AuthController, songsController *controllers.SongController, sessionController *controllers.SessionController, wsServer *server.WebSocketServer) *http.ServeMux {

	mux := http.NewServeMux()

	// http
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/songs/random", songsController.GetRandomSongs)
	mux.HandleFunc("/auth/tokens", authController.GetAccessTokens)
	mux.HandleFunc("/auth/refresh", authController.RefreshToken)
	mux.HandleFunc("/session/create", sessionController.CreateSession)
	mux.HandleFunc("/session/{id}", sessionController.GetSessionById)
	mux.HandleFunc("/session/update", sessionController.UpdateSession)

	// webscoket
	mux.HandleFunc("/session/socket", wsServer.HandleConnection)

	fmt.Println("🕒 Loading server routes...")
	return mux
}
