package main

import (
	"context"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func SetDownloadTask(path string) chromedp.Tasks {
	return chromedp.Tasks{
		page.SetAdBlockingEnabled(true),
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorAllow).WithDownloadPath(path),
	}
}

func DownloadSong(song SearchMatch, ct *context.Context) {
	ctx, cancel := chromedp.NewContext(*ct)

	inputFieldSelector := "#layout > header > div.container.header__container > div.convert-form > div.container > div.convert-form__input-container > label > input"
	submitURLSelector := "div.format_selection > button"
	downloadButtonSelector := "#layout > header > div.container.header__container > div.convert-form > div > div.download__buttons > button"

	err := chromedp.Run(
		ctx,
		chromedp.Navigate("https://2conv.com/"),
		chromedp.WaitReady(inputFieldSelector),
		chromedp.SendKeys(inputFieldSelector, song.URL),
		chromedp.Click(submitURLSelector),
		chromedp.WaitVisible(downloadButtonSelector),
		SetDownloadTask(MakeDownloadPath(song, os.Args[1])),
		chromedp.Click(downloadButtonSelector),
	)
	HandleError(err)

	cancel()
}
