package youtubebot

import (
	"code.google.com/p/google-api-go-client/googleapi/transport"
	"code.google.com/p/google-api-go-client/youtube/v3"
	"net/http"
)

type Video struct {
	Title string
	Id    string
}

// creates a new youtube service
func Initialise(APIKey string) (*youtube.Service, error) {
	client := &http.Client{Transport: &transport.APIKey{Key: APIKey}}
	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// returns the latest uploaded video from a channel
func LatestVideo(service *youtube.Service, channelName string) (Video, error) {

	call := service.Channels.List("contentDetails").ForUsername(channelName)
	response, err := call.Do()

	if err != nil {
		return Video{}, err
	}

	channel := response.Items[0]

	playlistId := channel.ContentDetails.RelatedPlaylists.Uploads

	playlistCall := service.PlaylistItems.List("snippet").
		PlaylistId(playlistId).
		MaxResults(50)

	playlistResponse, err := playlistCall.Do()

	if err != nil {
		return Video{}, err
	}

	title := playlistResponse.Items[0].Snippet.Title
	videoId := playlistResponse.Items[0].Snippet.ResourceId.VideoId

	return Video{title, videoId}, nil

}
