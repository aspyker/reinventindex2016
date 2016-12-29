package main

import (
	"net/http"
	"log"
	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
	"fmt"
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

	//log.Printf("about to search youtube for \"%s\"", keyword)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}
	//log.Printf("youtube results len %v", len(response.Items))

	if len(response.Items) != 0 {
		for _, item := range response.Items {
			switch item.Id.Kind {
			case "youtube#video":
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