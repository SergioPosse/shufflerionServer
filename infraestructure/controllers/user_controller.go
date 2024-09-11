// /infrastructure/controllers/user_controller.go
package controllers

import (
    "encoding/json"
    "net/http"
    "shufflerion/modules/user/application"
)

type UserController struct {
    GenerateUserAccessLinkUC *application.GenerateUserAccessLinkUseCase
}

func NewUserController(generateUserAccessLinkUC *application.GenerateUserAccessLinkUseCase) *UserController {
    return &UserController{GenerateUserAccessLinkUC: generateUserAccessLinkUC}
}

// Para manejar el GET
func (c *UserController) GetAccessLink(w http.ResponseWriter, r *http.Request) {
	accessLink := c.GenerateUserAccessLinkUC.Execute()

	// Envía el enlace como respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_link": accessLink})
}

// Para manejar el POST
func (c *UserController) GenerateUserAccessLink(w http.ResponseWriter, r *http.Request) {
    // userID := r.URL.Query().Get("userID")
    // if userID == "" {
    //     http.Error(w, "Missing userID", http.StatusBadRequest)
    //     return
    // }

    // accessLink := c.GenerateUserAccessLinkUC.Execute()

    // w.Header().Set("Content-Type", "application/json")
    // json.NewEncoder(w).Encode(map[string]string{"access_link": accessLink})
}