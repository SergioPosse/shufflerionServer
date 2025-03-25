package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	config "shufflerion/infrastructure/server/config"
	"shufflerion/modules/session/domain"
	domainSong "shufflerion/modules/song/domain"
	"strconv"
	"time"
)

type SpotifyResponse struct {
	Total int `json:"total"`
	Items []struct {
		Track struct {
			Name    string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Duration int `json:"duration_ms"`
			Explicit bool `json:"explicit`
			Album struct {
				Images []struct {
					Url string `json:"url"`
				} `json:"images"`
			} `json:"album"`
			Uri          string `json:"uri"`
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
		} `json:"track"`
	} `json:"items"`
}

type SpotifyService struct {
	config *config.Config
}

func NewSpotifyService(cfg *config.Config) domain.SpotifyService {
	return &SpotifyService{config: cfg}
}

func (s *SpotifyService) FetchRandomSongs(accessToken string, quantity int) ([]domainSong.Song, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s0", s.config.APIURL_GET_SONGS)

	firstRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for get tracks: %v", err)
	}
	firstRequest.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(firstRequest)
	if err != nil {
		return nil, fmt.Errorf("error in initial request: %v", err)
	}
	defer resp.Body.Close()

	var firstResponse SpotifyResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	err = json.Unmarshal(body, &firstResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	totalTracks := firstResponse.Total
	if totalTracks == 0 {
		return nil, fmt.Errorf("there is no tracks in the account")
	}

	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tracks []domainSong.Song
	for i := 0; i < quantity; i++ {
		offset := randSource.Intn(totalTracks) + 1
		url := s.config.APIURL_GET_SONGS + strconv.Itoa(offset)

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request for get tracks: %v", err)
		}
		request.Header.Set("Authorization", "Bearer "+accessToken)
		resp, err := client.Do(request)
		if err != nil {
			return nil, fmt.Errorf("error in getting tracks request: %v", err)
		}
		defer resp.Body.Close()

		var response SpotifyResponse
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading tracks response: %v", err)
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("error parsing tracks response: %v", err)
		}

		for _, item := range response.Items {
			track := domainSong.Song{
				Title:  item.Track.Name,
				Artist: item.Track.Artists[0].Name,
				Url:    item.Track.Uri,
				Image:  item.Track.Album.Images[0].Url,
				Duration: item.Track.Duration,
				Explicit: item.Track.Explicit,
			}
			tracks = append(tracks, track)
		}
	}

	return tracks, nil
}

// not used method but probably used in production mode when spotify dev team give me feedback and approve my app
func (s *SpotifyService) AddUserToApp(accessToken string, email string) (bool, error) {
	client := &http.Client{}

	clientId := s.config.ClientID
	if clientId == "" {
		return false, fmt.Errorf("SPOTIFY_CLIENT_ID no está definido en las variables de entorno")
	}

	payload := map[string]string{
		"clientId": clientId,
		"email":    email,
		"name":     email,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("error serializing JSON: %v", err)
	}

	url := s.config.APIURL_ADD_USER + clientId + "/users"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("error creating http request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error in request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return false, fmt.Errorf("error in API response: %s", resp.Status)
	}

	return true, nil
}
