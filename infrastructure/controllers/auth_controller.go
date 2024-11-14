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
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }

    var requestBody AuthControllerRequestBody

    errorDecodingBody := json.NewDecoder(r.Body).Decode(&requestBody)
    if errorDecodingBody != nil {
        http.Error(w, "Auth: Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
        return
    }

    code1 := requestBody.Code1
    code2 := requestBody.Code2

    if code1 == "" || code2 == "" {
        http.Error(w, "Both code1 and code2 parameters are required", http.StatusBadRequest)
        return
    }

    access_tokens := c.GetAccessTokensUC.Execute(code1, code2)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(access_tokens)
}