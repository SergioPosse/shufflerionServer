package services

import (
		"shufflerion/modules/auth/domain"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
)

type AuthService struct{}

func NewAuthService() *AuthService {
    return &AuthService{}
}

// GenerateUserAccessLink genera los access tokens a partir de los códigos
func (s *AuthService) GetAccessTokens(code1 string, code2 string) ([]shared.GetAccessTokensResponse, error) {
    // Spotify API endpoint
    apiURL := "https://accounts.spotify.com/api/token"

    // Datos para la solicitud (se usa code1 aquí para el primer fetch)
    data := url.Values{}
    data.Set("grant_type", "authorization_code")
    data.Set("code", code1) // Primer código para el primer access_token
    data.Set("redirect_uri", "http://localhost:3000/callback")
    data.Set("client_id", "335ea7b32dd24009bd0529ba85f0f8cc")
    data.Set("client_secret", "b482ee9f0aa4408da21b224d59c2d445")

    // Crear la solicitud POST para obtener el primer access_token
    req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
    if err != nil {
        return nil, fmt.Errorf("error creando la solicitud: %v", err)
    }

    // Setear los headers
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Realizar la solicitud
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error al hacer la solicitud: %v", err)
    }
    defer resp.Body.Close()

    // Leer la respuesta
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error al leer la respuesta: %v", err)
    }

    // Verificar si la solicitud fue exitosa
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("error en la solicitud, status code: %d, respuesta: %s", resp.StatusCode, body)
    }

    // Parsear el JSON para extraer el access_token
    var response map[string]interface{}
    err = json.Unmarshal(body, &response)
    if err != nil {
        return nil, fmt.Errorf("error al parsear el cuerpo de la respuesta: %v", err)
    }

    accessToken1, ok := response["access_token"].(string)
		refreshToken1, ok2 := response["refresh_token"].(string)
    if !ok || !ok2 {
        return nil, fmt.Errorf("no se encontró el access_token o refresh_token en la respuesta")
    }

    // data.Set("code", code2) // Cambiar el código a code2

    // req, err = http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
    // if err != nil {
    //     return nil, fmt.Errorf("error creando la solicitud para el segundo código: %v", err)
    // }

    // resp, err = client.Do(req)
    // if err != nil {
    //     return nil, fmt.Errorf("error al hacer la solicitud para el segundo código: %v", err)
    // }
    // defer resp.Body.Close()

    // // Leer la respuesta para el segundo access_token
    // body, err = io.ReadAll(resp.Body)
    // if err != nil {
    //     return nil, fmt.Errorf("error al leer la respuesta del segundo código: %v", err)
    // }

    // if resp.StatusCode != http.StatusOK {
    //     return nil, fmt.Errorf("error en la solicitud para el segundo código, status code: %d, respuesta: %s", resp.StatusCode, body)
    // }

    // // Parsear el JSON para el segundo access_token
    // err = json.Unmarshal(body, &response)
    // if err != nil {
    //     return nil, fmt.Errorf("error al parsear la respuesta del segundo código: %v", err)
    // }

    // accessToken2, ok := response["access_token"].(string)
    // if !ok {
    //     return nil, fmt.Errorf("no se encontró el access_token en la respuesta del segundo código")
    // }

    // Retornar ambos access_tokens
		authResponse := []shared.GetAccessTokensResponse{}
		authResponse = append(authResponse, shared.GetAccessTokensResponse{AccessToken: accessToken1, RefreshToken: refreshToken1})
		authResponse = append(authResponse, shared.GetAccessTokensResponse{AccessToken: "accessToken2", RefreshToken: "refreshToken2"})
    return authResponse, nil
}