package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shufflerion/modules/session/application"
	"shufflerion/modules/session/domain"
)

type CreateSessionRequest struct {
	Id   string      `json:"id"`
	Host domain.User `json:"host"`
}

type UpdateSessionRequest struct {
	Id    string      `json:"id"`
	Guest domain.User `json:"guest"`
}

type SessionController struct {
	sessionUseCase *application.SessionUseCase
}

func NewSessionController(sessionUseCase *application.SessionUseCase) *SessionController {
	return &SessionController{sessionUseCase}
}

func (c *SessionController) CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var requestBody CreateSessionRequest

	// Decodificar el JSON ANTES de validar
	errorDecodingBody := json.NewDecoder(r.Body).Decode(&requestBody)
	if errorDecodingBody != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud: "+errorDecodingBody.Error(), http.StatusBadRequest)
		return
	}

	// Imprimir para depuración (ahora sí tiene datos)
	fmt.Printf("Request Body: %+v\n", requestBody)

	// Validar que los campos no estén vacíos
	if requestBody.Id == "" || requestBody.Host.Email == "" {
		http.Error(w, "Se requiere el ID de la sesión y el usuario anfitrión", http.StatusBadRequest)
		return
	}

	// 💡 **Verificar si ya existe una sesión con ese ID**
	existingSession, err := c.sessionUseCase.GetSession(r.Context(), requestBody.Id)
	if err == nil && existingSession != nil {
		http.Error(w, "Ya existe una sesión con este ID", http.StatusConflict)
		return
	}

	session := domain.CreateSession{
		Id:   requestBody.Id,
		Host: requestBody.Host,
	}

	err2 := c.sessionUseCase.CreateSession(r.Context(), session)
	if err2 != nil {
		http.Error(w, "Error al crear la sesión: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

func (c *SessionController) GetSessionById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Se requiere el parámetro 'id'", http.StatusBadRequest)
		return
	}

	session, err := c.sessionUseCase.GetSession(r.Context(), id)
	if err != nil {
		http.Error(w, "Error al obtener la sesión: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func (c *SessionController) UpdateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var requestBody UpdateSessionRequest

	errorDecodingBody := json.NewDecoder(r.Body).Decode(&requestBody)
	if errorDecodingBody != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if requestBody.Id == "" || requestBody.Guest.Email == "" {
		http.Error(w, "Se requiere el ID de la sesión y el usuario invitado", http.StatusBadRequest)
		return
	}

	session := domain.UpdateSession{
		Id:    requestBody.Id,
		Guest: requestBody.Guest, // Agregar el Guest como un array
	}

	updatedSession, err := c.sessionUseCase.UpdateSession(r.Context(), session)
	if err != nil {
		http.Error(w, "Error al actualizar la sesión: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSession)
}
