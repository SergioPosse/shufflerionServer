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
	err := godotenv.Load()
if err != nil {
    fmt.Println("Error cargando el archivo .env")
}
	// Inyección de dependencias para AuthService
	authService := services.NewAuthService()
	getAccessTokensUC := auth.NewGetAccessTokensUseCase(authService)
	authController := controllers.NewAuthController(getAccessTokensUC)

	// inyeccion para song controller
	songsService := services.NewSongsService()
	getRandomSongUC := songs.NewGetSongsUseCase(songsService)
	songsController := controllers.NewSongsController(getRandomSongUC)

	// Configurar MongoDB
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

	// Inyectar repositorio en el caso de uso
	sessionRepo := repository.NewMongoSessionRepository(mongoDB.DB)
	sessionUseCase := session.NewSessionUseCase(sessionRepo)
	sessionController := controllers.NewSessionController(sessionUseCase)

	// Registrar rutas
	routes.RegisterRoutes(authController, songsController, sessionController)

	// Iniciar servidor
	server.StartServer()
}
