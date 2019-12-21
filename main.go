package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/crosscode-nl/partition"
)

func MakeDownloadPath(track SearchMatch, rootDir string) string {
	name := fmt.Sprintf("%v ~ %v.mp3", track.Title, track.Authors)
	// cwd, err := os.Getwd()
	// HandleError(err)
	return path.Join(rootDir, name)
}

func Scrape(ps *PlaylistRes, ctx *context.Context, out chan SearchMatch) {
	partition.ToFunc(len(ps.Items), 5, func(l, h int) {
		var wg sync.WaitGroup
		for _, t := range ps.Items[l:h] {
			wg.Add(1)
			go func(t Track) {
				s := FindSong(ctx, t, out)
				DownloadSong(s, ctx)
				wg.Done()
			}(t.Track)
		}

		wg.Wait()
	})
}

func main() {
	os.Setenv("__CLIENT_ID__", "90e95e868e05453790dcae54d7dd451c")
	os.Setenv("__CLIENT_SECRET__", "e82b199394ba433687a15b4da6b351d6")

	url := "https://open.spotify.com/playlist/0b0vzbWWMUUqzYPB07CLm3?si=R8UugCFdQUm-NAJYlSknzA"
	token := GetAccessToken()
	id := GetPlaylistID(url)
	out := make(chan SearchMatch)

	ps := RequestPlaylist(id, token.AccessToken)
	ops := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
	c, _ := chromedp.NewExecAllocator(context.Background(), ops...)

	ctx, cancel := chromedp.NewContext(c)
	defer cancel()
	c2, cc2 := chromedp.NewContext(ctx)
	defer cc2()
	go chromedp.Run(c2, chromedp.Navigate("https://2conv.com/"))
	time.Sleep(5 * time.Second)
	chromedp.Run(ctx)

	Scrape(&ps, &ctx, out)

	for s := range out {
		fmt.Println(s)
	}

}
