package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

// TokenResponse @doc
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// Artist @doc
type Artist struct {
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	}
	Href  string `json:"href"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Ttype string `json:"type"`
	URI   string `json:"uri"`
}

// Image @doc
type Image struct {
	Height int64  `json:"height"`
	Width  int64  `json:"width"`
	URL    string `json:"url"`
}

// Album @doc
type Album struct {
	Artists []Artist `json:"artists"`
	Images  []Image  `json:"images"`
}

// Track @doc
type Track struct {
	Name       string `json:"name"`
	DurationMs int64  `json:"duration_ms"`
	Album      Album  `json:"album"`
}

// PlaylistRes @doc
type PlaylistRes struct {
	Items []struct {
		Track Track `json:"track"`
	} `json:"items"`
}

// GetIDRex @doc
var GetIDRex = regexp.MustCompile("playlist\\/(\\w+)\\??")

// GetPlaylistID @doc
func GetPlaylistID(s string) string {
	// Check if is not a URL
	matched, err := regexp.MatchString("^https?:\\/\\/.+", s)
	if HandleError(err) && !matched {
		return s
	}

	id := GetIDRex.FindStringSubmatch(s)[1]
	return id
}

// GetAccessToken @doc
func GetAccessToken() TokenResponse {
	clientID, clientSecret := os.Getenv("__CLIENT_ID__"), os.Getenv("__CLIENT_SECRET__")
	b64 := base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", clientID, clientSecret)))
	auth := fmt.Sprintf("Basic %v", b64)

	// Make a request
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token?grant_type=client_credentials", nil)
	HandleError(err)

	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	cl := &http.Client{}
	res, err := cl.Do(req)
	HandleError(err)

	// Decode the response
	decoder := json.NewDecoder(res.Body)
	var tokenRes TokenResponse
	erro := decoder.Decode(&tokenRes)
	HandleError(erro)

	return tokenRes
}

// MakePlaylistURI @doc
func MakePlaylistURI(id string) string {
	uri := fmt.Sprintf("https://api.spotify.com/v1/playlists/%v/tracks?market=ES&fields=items.track(album(images%%2C%%20artists(name%%2C%%20href))%%2C%%20name%%2C%%20duration_ms)", id)

	return uri
}

// RequestPlaylist @doc
func RequestPlaylist(id string, accessToken string) PlaylistRes {
	req, err := http.NewRequest("GET", MakePlaylistURI(id), nil)
	HandleError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", accessToken))

	cl := &http.Client{}
	res, err := cl.Do(req)
	HandleError(err)

	var playlist PlaylistRes
	dec := json.NewDecoder(res.Body)

	erro := dec.Decode(&playlist)
	HandleError(erro)

	return playlist
}
