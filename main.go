// /infraestructure/server/main.go
package main

import (
    "shufflerion/infrastructure/routes"
    "shufflerion/infrastructure/server"
    "shufflerion/infrastructure/controllers"
    "shufflerion/infrastructure/services"
    auth "shufflerion/modules/auth/application"
    songs "shufflerion/modules/song/application"
)

func main() {
    // Inyección de dependencias para UserService
    authService := services.NewAuthService()
    getAccessTokensUC := auth.NewGetAccessTokensUseCase(authService)
    authController := controllers.NewAuthController(getAccessTokensUC)

    // inyeccion para song controller
    songsService := services.NewSongsService()
    getRandomSongUC := songs.NewGetSongsUseCase(songsService)
    songsController := controllers.NewSongsController(getRandomSongUC)

    // Registrar rutas
    routes.RegisterRoutes(authController, songsController)

    // Iniciar servidor
    server.StartServer()
}
