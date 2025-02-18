package controllers

import (
    "encoding/json"
    "net/http"
    authModule "shufflerion/modules/auth/application"
)

type AuthControllerRequestBody struct {
    Code1 string `json:"code1"`
    Code2 string `json:"code2"`
}

type AuthController struct {
    GetAccessTokensUC  *authModule.GetAccessTokensUseCase
}

func NewAuthController(getAccessTokensUC *authModule.GetAccessTokensUseCase) *AuthController {
    return &AuthController{GetAccessTokensUC: getAccessTokensUC}
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