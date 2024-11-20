package services

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"shufflerion/modules/song/domain"
	"time"
)

type SpotifyResponse struct {
	Total int `json:"total"`
	Items []struct {
		Track struct {
			Name   string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Album struct {
				Images []struct {
					Url string `json:"url"`
				}`json:"images"`
			} `json:"album"`
			Uri string `json:"uri"`
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
		} `json:"track"`
	} `json:"items"`
}

func FetchRandomSongs(accessToken string, quantity int) ([]domain.Song, error) {
	client := &http.Client{}

	firstRequest, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/tracks?limit=1&offset=0", nil)
	if err != nil {
		return nil, fmt.Errorf("error creando la solicitud para obtener el total de pistas: %v", err)
	}
	firstRequest.Header.Set("Authorization", "Bearer "+ accessToken)

	resp, err := client.Do(firstRequest)
	if err != nil {
		return nil, fmt.Errorf("error al hacer la solicitud inicial: %v", err)
	}
	defer resp.Body.Close()

	var firstResponse SpotifyResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %v", err)
	}

	err = json.Unmarshal(body, &firstResponse)
	if err != nil {
		return nil, fmt.Errorf("error al parsear la respuesta: %v", err)
	}

	totalTracks := firstResponse.Total
	if totalTracks == 0 {
		return nil, fmt.Errorf("no se encontraron pistas en la cuenta")
	}

	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tracks []domain.Song
	for i := 0; i < quantity; i++ {
		offset := randSource.Intn(totalTracks) + 1
		request, err := http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/me/tracks?limit=1&offset=%d", offset), nil)
		if err != nil {
			return nil, fmt.Errorf("error creando la solicitud para obtener las pistas: %v", err)
		}
		request.Header.Set("Authorization", "Bearer " + accessToken)
		resp, err := client.Do(request)
		if err != nil {
			return nil, fmt.Errorf("error al hacer la solicitud de pistas: %v", err)
		}
		defer resp.Body.Close()

		var response SpotifyResponse
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error al leer la respuesta de pistas: %v", err)
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("error al parsear la respuesta de pistas: %v", err)
		}

		for _, item := range response.Items {
			track := domain.Song{
				Title:  item.Track.Name,
				Artist: item.Track.Artists[0].Name,
				Url:    item.Track.Uri,
				Image: item.Track.Album.Images[0].Url,
			}
			tracks = append(tracks, track)
		}
	}

	return tracks, nil
}
