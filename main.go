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
	config "shufflerion/infrastructure/server/config"

	"github.com/joho/godotenv"
)

func main() {

	// load .env
	if os.Getenv("ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env")
		}
	}

	// load config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	// spotify service
	spotifyService := services.NewSpotifyService(cfg)

	// auth module injection
	authService := services.NewAuthService(cfg)
	getAccessTokensUC := auth.NewGetAccessTokensUseCase(authService)
	authController := controllers.NewAuthController(getAccessTokensUC)

	// song module injection
	songsService := services.NewSongsService(spotifyService)
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

	// session module injection
	sessionRepo := repository.NewMongoSessionRepository(mongoDB.DB)
	sessionUseCase := session.NewSessionUseCase(sessionRepo, spotifyService)
	sessionController := controllers.NewSessionController(sessionUseCase)

	// create webSocket server
	wsServer := server.NewWebSocketServer(mongoDB.Client)

	// register routes
	serverWithRoutes := routes.RegisterRoutes(authController, songsController, sessionController, wsServer)

	// start server
	server.StartServer(serverWithRoutes)
}
