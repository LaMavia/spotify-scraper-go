package main

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type SearchMatch struct {
	Title      string
	URL        string
	DurationMS int
	Authors    string
}

func NormalizeString(s string) string {
	return strings.ToLower(s)
}

func GetTrackArtists(track Track) string {
	var output string
	n := len(track.Album.Artists) - 1

	for i, a := range track.Album.Artists {
		output += NormalizeString(a.Name)
		if i != n {
			output += " "
		}
	}

	return output
}

func MakeSearchURL(track Track) string {
	query := fmt.Sprintf("\"%v\" %v", GetTrackArtists(track), NormalizeString(track.Name))
	return fmt.Sprintf("https://www.youtube.com/results?search_query=%v&sp=EgIQAQ%%253D%%253D", url.QueryEscape(query))
}

func CalcCos(a, b string) float64 {
	a, b = strings.ToLower(a), strings.ToLower(b)
	wordsA, wordsB := strings.Fields(a), strings.Fields(b)
	words := Unique(append(wordsA, wordsB...))
	vecA, vecB := make([]int, len(words)), make([]int, len(words))

	for i, w := range words {
		vecA[i] = Count(wordsA, w)
		vecB[i] = Count(wordsB, w)
	}

	dotProd := 0
	sumA, sumB := 0, 0
	for i := range words {
		dotProd += vecA[i] * vecB[i]
		sumA += vecA[i] * vecA[i]
		sumB += vecB[i] * vecB[i]
	}

	return float64(dotProd) / (math.Sqrt(float64(sumA)) * math.Sqrt(float64(sumB)))
}

var isDig = regexp.MustCompile("\\d+")

// CalcDuration : calculates the duration of search result
func CalcDuration(s string) int {
	wds := strings.Fields(s)
	out := 0

	for _, w := range wds {
		if isDig.MatchString(w) {
			d, err := strconv.Atoi(w)
			HandleError(err)

			// fmt.Printf("%v : %v => %v\n", w, d, out)

			out = out*60 + d
		}
	}

	return out * 1000
}

func Score(track Track, match SearchMatch) float64 {
	cos := CalcCos(track.Name, match.Title)
	dDuration := float64(Abs(track.DurationMs - match.DurationMS))

	return cos / dDuration
}

func FindSong(contx *context.Context, track Track, out chan SearchMatch) SearchMatch {
	ctx, cancel := chromedp.NewContext(*contx)

	// var t string
	var titles []*cdp.Node
	var durations []*cdp.Node
	var match SearchMatch
	lastScore := 0.0
	initialMatch := false

	err := chromedp.Run(ctx,
		chromedp.Navigate(MakeSearchURL(track)),
		chromedp.WaitVisible("a#video-title"),
		chromedp.Nodes("a#video-title", &titles, chromedp.ByQueryAll),
		chromedp.Nodes("ytd-thumbnail-overlay-time-status-renderer > span.style-scope.ytd-thumbnail-overlay-time-status-renderer", &durations, chromedp.ByQueryAll),
	)
	HandleError(err)

	for i, x := range titles {
		if i >= len(durations) {
			continue
		}

		result := SearchMatch{
			Title:      NormalizeString(x.AttributeValue("title")),
			DurationMS: CalcDuration(durations[i].AttributeValue("aria-label")),
			URL:        x.AttributeValue("href"),
			Authors:    "",
		}

		score := Score(track, result)

		if !initialMatch || Score(track, result) > lastScore {
			match = result
			lastScore = score
			initialMatch = true
		}
	}

	match.Authors = GetTrackArtists(track)
	cancel()
	return match
}
