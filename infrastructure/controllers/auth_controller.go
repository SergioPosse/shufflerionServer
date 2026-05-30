package controllers

import (
    "encoding/json"
    "net/http"
    authModule "shufflerion/modules/auth/application"
    "shufflerion/infrastructure/services"
)

type AuthControllerRequestBody struct {
    Code1 string `json:"code1"`
    Code2 string `json:"code2"`
}

type RefreshTokenRequestBody struct {
    RefreshToken string `json:"refresh_token"`
}

type AuthController struct {
    GetAccessTokensUC *authModule.GetAccessTokensUseCase
    AuthService       *services.AuthService
}

func NewAuthController(getAccessTokensUC *authModule.GetAccessTokensUseCase, authService *services.AuthService) *AuthController {
    return &AuthController{GetAccessTokensUC: getAccessTokensUC, AuthService: authService}
}

func (c *AuthController) GetAccessTokens(w http.ResponseWriter, r *http.Request) {

    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var requestBody AuthControllerRequestBody

    errorDecodingBody := json.NewDecoder(r.Body).Decode(&requestBody)
    if errorDecodingBody != nil {
        http.Error(w, "Auth: Error reading request body", http.StatusBadRequest)
        return
    }

    code1 := requestBody.Code1
    code2 := requestBody.Code2

    if code1 == "" {
        http.Error(w, "code1 is required", http.StatusBadRequest)
        return
    }

    access_tokens := c.GetAccessTokensUC.Execute(code1, code2)

    if access_tokens == nil {
        http.Error(w, "Failed to retrieve access tokens", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(access_tokens)
}

func (c *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var body RefreshTokenRequestBody
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "error reading request body", http.StatusBadRequest)
        return
    }
    if body.RefreshToken == "" {
        http.Error(w, "refresh_token is required", http.StatusBadRequest)
        return
    }

    tokens, err := c.AuthService.RefreshAccessToken(body.RefreshToken)
    if err != nil {
        http.Error(w, "failed to refresh token: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tokens)
}