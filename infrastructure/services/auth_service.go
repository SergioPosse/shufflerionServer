package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	config "shufflerion/infrastructure/server/config"
	"shufflerion/modules/auth/domain"
	"strings"
)
type AuthService struct{
	config *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{config: cfg}
}

func (s *AuthService) GetAccessTokens(code1, code2 string) ([]shared.GetAccessTokensResponse, error) {
	tokens := []shared.GetAccessTokensResponse{}

	// get the first token
	token1, err := s.fetchAccessToken(code1)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo el primer access token: %v", err)
	}
	tokens = append(tokens, token1)

	// get the second token
	if(len(code2) > 20){
		token2, err := s.fetchAccessToken(code2)
		if err != nil {
			return nil, fmt.Errorf("error obteniendo el segundo access token: %v", err)
		}
		tokens = append(tokens, token2)
	} else {
		tokens = append(tokens, shared.GetAccessTokensResponse{AccessToken: "asd", RefreshToken: "asd"})
	}

	return tokens, nil
}

// fetchAccessToken run a spotify request to obtain access token
func (s *AuthService) fetchAccessToken(code string) (shared.GetAccessTokensResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", s.config.RedirectURI)
	data.Set("client_id", s.config.ClientID)
	data.Set("client_secret", s.config.ClientSecret)

	req, err := http.NewRequest("POST", s.config.APIURL, strings.NewReader(data.Encode()))
	if err != nil {
		return shared.GetAccessTokensResponse{}, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return shared.GetAccessTokensResponse{}, fmt.Errorf("error running request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return shared.GetAccessTokensResponse{}, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return shared.GetAccessTokensResponse{}, fmt.Errorf("error request, status code: %d, response: %s", resp.StatusCode, body)
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return shared.GetAccessTokensResponse{}, fmt.Errorf("error parsing response: %v", err)
	}

	accessToken, ok := response["access_token"].(string)
	refreshToken, ok2 := response["refresh_token"].(string)
	if !ok || !ok2 {
		return shared.GetAccessTokensResponse{}, fmt.Errorf("access_token o refresh_token missing in response")
	}

	return shared.GetAccessTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
