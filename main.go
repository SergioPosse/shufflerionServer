package main

import (
	"fmt"
	"log"
	"os"
	"shufflerion/infrastructure/controllers"
	"shufflerion/infrastructure/db"
	repository "shufflerion/infrastructure/repository/session"
	"shufflerion/infrastructure/routes"
	"shufflerion/infrastructure/server"
	"shufflerion/infrastructure/services"
	auth "shufflerion/modules/auth/application"
	session "shufflerion/modules/session/application"
	songs "shufflerion/modules/song/application"

	"github.com/joho/godotenv"
)

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
			fmt.Println("Error loading .env")
	}
	// auth injection
	authService := services.NewAuthService()
	getAccessTokensUC := auth.NewGetAccessTokensUseCase(authService)
	authController := controllers.NewAuthController(getAccessTokensUC)

	// song injection
	songsService := services.NewSongsService()
	getRandomSongUC := songs.NewGetSongsUseCase(songsService)
	songsController := controllers.NewSongsController(getRandomSongUC)

	// mongo config
	mongoURI := fmt.Sprintf(
		"mongodb+srv://%s:%s@clusterfree.x3n59lo.mongodb.net/?retryWrites=true&w=majority&appName=ClusterFree",
		os.Getenv("DB_MONGO_USER"),
		os.Getenv("DB_MONGO_PASSWORD"),
	)

	mongoDB, err := db.NewMongoDB(mongoURI, "shufflerion")
	if err != nil {
		log.Fatal(err)
	}
	defer mongoDB.Close()

	// session injection
	sessionRepo := repository.NewMongoSessionRepository(mongoDB.DB)
	sessionUseCase := session.NewSessionUseCase(sessionRepo)
	sessionController := controllers.NewSessionController(sessionUseCase)

	// Crear servidor WebSocket
	wsServer := server.NewWebSocketServer(mongoDB.Client)

	// register routes
	routes.RegisterRoutes(authController, songsController, sessionController, wsServer)

	// start server
	server.StartServer()
}
