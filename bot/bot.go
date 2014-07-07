package bot

import (
	"errors"
	"fmt"
	"github.com/scord/goreddit"
	"github.com/scord/reddit-youtube-bot/bot/youtubebot"
	"time"
)

// start point
func Run(channel, subreddit, user, pass, apiKey string) error {

	if channel == "" || subreddit == "" || user == "" || pass == "" {
		return errors.New("Missing details")
	} else if apiKey == "" {
		return errors.New("Invalid API Key")
	}

	service, err := youtubebot.Initialise(apiKey)

	if err != nil {
		errors.New("Could not create new Youtube client")
	}

	latestVideo, err := youtubebot.LatestVideo(service, channel)

	if err != nil {
		errors.New("Could not get latest video")
	}

	for {
		time.Sleep(5 * time.Second)

		fmt.Println("looking")

		newVideo, err := youtubebot.LatestVideo(service, channel)

		if err != nil {
			errors.New("Could not get latest video")
		}

		if latestVideo != newVideo {
			newVideo = latestVideo

			err := postLink(user, pass, newVideo.Title, fmt.Sprintf("https://www.youtube.com/watch?v=%s", newVideo.ID), subreddit)

			if err != nil {
				errors.New("Could not post video")
			}
		}
	}
}

// logs into reddit and submits post
func postLink(user, pass, title, linkURL, subreddit string) error {

	session, err := goreddit.Login(user, pass)

	if err != nil {
		return err
	}

	err = session.Submit(title, linkURL, subreddit)

	if err != nil {
		return err
	}

	return nil
}
