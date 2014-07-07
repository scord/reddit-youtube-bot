package main

import (
	"fmt"
	"github.com/scord/reddit-youtube-bot/reddit"
	"github.com/scord/reddit-youtube-bot/youtubebot"
	"time"
)

const APIKey = "" //add your own API Key

// start point
func main() {

	service, err := youtubebot.Initialise(APIKey)

	if err != nil {
		fmt.Println("Could not create new Youtube client")
	}

	channel, subreddit, user, pass := "", "", "", ""

	fmt.Println("Enter Channel Username:")
	fmt.Scanf("%s", &channel)
	fmt.Println("Enter Subreddit")
	fmt.Scanf("%s", &subreddit)

	fmt.Println("Enter Bot Username:")
	fmt.Scanf("%s", &user)
	fmt.Println("Enter Bot Password")
	fmt.Scanf("%s", &pass)

	latestVideo, err := youtubebot.LatestVideo(service, channel)

	if err != nil {
		fmt.Println("Could not get latest video")
	}

	for {
		time.Sleep(5 * time.Second)

		fmt.Println("Looking for new upload")

		newVideo, err := youtubebot.LatestVideo(service, channel)

		if err != nil {
			fmt.Println("Could not get latest video")
		}

		if latestVideo != newVideo {
			newVideo = latestVideo

			err := postLink(user, pass, newVideo.Title, fmt.Sprintf("https://www.youtube.com/watch?v=%s", newVideo.Id), subreddit)

			if err != nil {
				fmt.Println("Could not post video")
			}

			fmt.Printf("Posted video: %s\n", newVideo.Title)

		} else {
			fmt.Println("No new videos")
		}
	}
}

// logs into reddit and submits post
func postLink(user, pass, title, linkURL, subreddit string) error {

	session, err := reddit.Login(user, pass)

	if err != nil {
		return err
	}

	err = session.Submit(title, linkURL, subreddit)

	if err != nil {
		return err
	}

	return nil
}
