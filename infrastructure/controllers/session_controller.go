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
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody CreateSessionRequest

	errorDecodingBody := json.NewDecoder(r.Body).Decode(&requestBody)
	if errorDecodingBody != nil {
		http.Error(w, "Error reading request body: "+errorDecodingBody.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Request Body: %+v\n", requestBody)

	if requestBody.Id == "" || requestBody.Host.Email == "" {
		http.Error(w, "sessionId and host data are mandatory", http.StatusBadRequest)
		return
	}

	existingSession, err := c.sessionUseCase.GetSession(r.Context(), requestBody.Id)
	if err == nil && existingSession != nil {
		http.Error(w, "session already exist", http.StatusConflict)
		return
	}

	session := domain.Session{
		Id:   requestBody.Id,
		Host: requestBody.Host,
	}

	err2 := c.sessionUseCase.CreateSession(r.Context(), session)
	if err2 != nil {
		http.Error(w, "Error creating session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

func (c *SessionController) GetSessionById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is mandatory", http.StatusBadRequest)
		return
	}

	session, err := c.sessionUseCase.GetSession(r.Context(), id)
	if err != nil {
		http.Error(w, "Error retrieving the session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func (c *SessionController) UpdateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody UpdateSessionRequest

	errorDecodingBody := json.NewDecoder(r.Body).Decode(&requestBody)
	if errorDecodingBody != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if requestBody.Id == "" || requestBody.Guest.Email == "" {
		http.Error(w, "sessionId and guest data are mandatory", http.StatusBadRequest)
		return
	}

	session := domain.UpdateSession{
		Id:    requestBody.Id,
		Guest: requestBody.Guest,
	}

	updatedSession, err := c.sessionUseCase.UpdateSession(r.Context(), session)
	if err != nil {
		http.Error(w, "Error updating the session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSession)
}
