package main

import (
	"fmt"
	"os"
)

func main() {
	os.Setenv("__CLIENT_ID__", "90e95e868e05453790dcae54d7dd451c")
	os.Setenv("__CLIENT_SECRET__", "e82b199394ba433687a15b4da6b351d6")

	url := "https://open.spotify.com/playlist/2qsfiBNlB3HYOsOmXTk1aS?si=Qh2IPldOTYKO1hoMYUc6NA"
	token := GetAccessToken()
	id := GetPlaylistID(url)

	playlist := RequestPlaylist(id, token.AccessToken)
	fmt.Printf("%v", playlist)
}
