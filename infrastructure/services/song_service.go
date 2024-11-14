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
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
		} `json:"track"`
	} `json:"items"`
}

type SpotifyService struct {}

func NewSongsService() domain.SongsRepository {
	return &SpotifyService{}
}

func (s *SpotifyService) GetRandomSongs(accessToken1 string, accessToken2 string) ([]domain.Song, error) {
	client := &http.Client{}

	fmt.Println("Token:", accessToken1)
	fmt.Println("Token2:", accessToken2)

	firstRequest, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/tracks?limit=1&offset=0", nil)
	if err != nil {
		return nil, fmt.Errorf("error creando la solicitud para obtener el total de pistas: %v", err)
	}
	firstRequest.Header.Set("Authorization", "Bearer BQC5twR43yJ7bURtKwbm36a2VWpMFQXK54G57Me191YxgKMR_H3YgMDb289wssy5_WsBuBqn5q_56nhf-J-CNSjYDrq7C_-aaDTq33MssIqAEvs0w8Ou5xRwdr9ZNbKmRyHTWWAtJHiVVq_STAGhFbAwhLHf2SKDZedTRapHPhzJTUgs8Urn6cQ9Y_e8llAuyXyVTtw_CPX6xNBTBrJMH0DEN_c")

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
	for i := 0; i < 50; i++ {
		offset := randSource.Intn(totalTracks) + 1
		request, err := http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/me/tracks?limit=1&offset=%d", offset), nil)
		if err != nil {
			return nil, fmt.Errorf("error creando la solicitud para obtener las pistas: %v", err)
		}
		request.Header.Set("Authorization", "Bearer BQC5twR43yJ7bURtKwbm36a2VWpMFQXK54G57Me191YxgKMR_H3YgMDb289wssy5_WsBuBqn5q_56nhf-J-CNSjYDrq7C_-aaDTq33MssIqAEvs0w8Ou5xRwdr9ZNbKmRyHTWWAtJHiVVq_STAGhFbAwhLHf2SKDZedTRapHPhzJTUgs8Urn6cQ9Y_e8llAuyXyVTtw_CPX6xNBTBrJMH0DEN_c")

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
				Url:    item.Track.ExternalURLs.Spotify,
			}
			tracks = append(tracks, track)
		}
	}

	return tracks, nil
}
