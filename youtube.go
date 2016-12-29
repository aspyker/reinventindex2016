package main

import (
	"net/http"
	"log"
	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
	"fmt"
	"strings"
)

// basically does:
// https://www.googleapis.com/youtube/v3/search?part=id,snippet&channelId=UCd6MoB9NC6uYN2grvUNT-Zg&q=ALX203&key=XXXX
func searchOnChannelByKeyword(channelId string, keyword string, developerKey string) string {
	url := ""
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	call := service.Search.List("id,snippet").
		ChannelId(channelId).
		Type("video").
		Q(keyword).
		MaxResults(10)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}

	if len(response.Items) != 0 {
		for _, item := range response.Items {
			if (item.Id.Kind == "youtube#video" && strings.Contains(item.Snippet.Title, keyword)) {
				url = fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId)
				break
			}
		}
	}

	if (url == "") {
		log.Printf("no youtube data found for %s", keyword)
	}

	return url
}