package config

import (
	"fmt"
	"os"
)

type Config struct {
	APIURL       string
	APIURL_GET_SONGS string
	APIURL_ADD_USER string
	RedirectURI  string
	ClientID     string
	ClientSecret string
}

// NewConfig carga las variables de entorno y devuelve la configuración
func NewConfig() (*Config, error) {
	cfg := &Config{
		APIURL:       os.Getenv("SPOTIFY_API_URL_TOKEN"),
		APIURL_GET_SONGS: os.Getenv("SPOTIFY_API_URL_GET_SONGS"),
		APIURL_ADD_USER: os.Getenv("SPOTIFY_API_URL_APP_AD_USER"),
		RedirectURI:  os.Getenv("SPOTIFY_REDIRECT_URI"),
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}

	// Validar que todas las variables estén configuradas
	if cfg.APIURL == "" || cfg.RedirectURI == "" || cfg.ClientID == "" || cfg.ClientSecret == "" || cfg.APIURL_ADD_USER == "" || cfg.APIURL_GET_SONGS == ""{
		return nil, fmt.Errorf("error: faltan variables de entorno necesarias")
	}

	return cfg, nil
}