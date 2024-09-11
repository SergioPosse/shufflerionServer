// /infraestructure/server/main.go
package main

import (
    "shufflerion/infraestructure/routes"
    "shufflerion/infraestructure/server"
    "shufflerion/infraestructure/controllers"
    "shufflerion/infraestructure/services"
    "shufflerion/modules/user/application"
)

func main() {
    // Inyección de dependencias para UserService
    userService := services.NewUserService()
    generateUserAccessLinkUC := application.NewGenerateUserAccessLinkUseCase(userService)
    userController := controllers.NewUserController(generateUserAccessLinkUC)

    // Registrar rutas
    routes.RegisterRoutes(nil, userController)

    // Iniciar servidor
    server.StartServer()
}
