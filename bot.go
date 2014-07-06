package main

import (
	"code.google.com/p/google-api-go-client/googleapi/transport"
	"code.google.com/p/google-api-go-client/youtube/v3"

	"fmt"
	"github.com/scord/reddit-youtube-bot/reddit"
	"net/http"
	"time"
)

const APIKey = "" //add your own API Key

type Video struct {
	Title string
	Id    string
}

func main() {

	client := &http.Client{Transport: &transport.APIKey{Key: APIKey}}

	service, err := youtube.New(client)

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

	newVideo := LatestVideo(service, channel)

	for {
		time.Sleep(5 * time.Second)
		fmt.Println("Looking for new upload")
		latestVideo := LatestVideo(service, channel)
		if latestVideo != newVideo {
			newVideo = latestVideo
			reddit.PostLink(user, pass, newVideo.Title, fmt.Sprintf("https://www.youtube.com/watch?v=%s", newVideo.Id), subreddit)
		} else {
			fmt.Println("No new videos")
		}
	}
}

func LatestVideo(service *youtube.Service, channelName string) (video Video) {

	call := service.Channels.List("contentDetails").ForUsername(channelName)
	response, err := call.Do()

	if err != nil {
		fmt.Println("Call failed")
	}

	channel := response.Items[0]

	playlistId := channel.ContentDetails.RelatedPlaylists.Uploads

	playlistCall := service.PlaylistItems.List("snippet").
		PlaylistId(playlistId).
		MaxResults(50)

	playlistResponse, err := playlistCall.Do()

	if err != nil {
		// The playlistItems.list method call returned an error.
		fmt.Println("Error fetching playlist items")
	}
	title := playlistResponse.Items[0].Snippet.Title
	fmt.Println(title)

	videoId := playlistResponse.Items[0].Snippet.ResourceId.VideoId

	return Video{title, videoId}

}
