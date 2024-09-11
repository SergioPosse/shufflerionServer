package routes

import (
    "net/http"
    "shufflerion/infraestructure/controllers"
)

func RegisterRoutes(songController *controllers.SongController, userController *controllers.UserController) {
    http.HandleFunc("/songs/random", songController.GetRandomSong)
    http.HandleFunc("/users/access-link", userController.GetAccessLink)
}