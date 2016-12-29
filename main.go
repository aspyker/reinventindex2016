package main

import (
	"os"
	"log"
	"github.com/cbroglie/mustache"
	"fmt"
)

type Session struct {
	SessionId string
	Title string
	Abstract string
	YoutubeUrl string
	SlideshareUrl string
}

type AllSessions struct {
	Sessions []Session
}

// look at view:source and look for data-channel-external-id="XXXX"
const amazonWebServicesChannelExternalId = "UCd6MoB9NC6uYN2grvUNT-Zg"
const awsSlideShareUsername = "AmazonWebServices"

var tracks = [...]string{"ALX"}//, "ARC", "BDA", "BDM", "BAP", "CMP", "CON", "CTD", "DAT", "DSC", "DEV", "ENT",
//	"FIN", "GAM", "HLC", "IOT", "LFS", "MAC", "MAE", "MBL", "NET", "SAC", "SEC", "SVR", "STG", "WIN", "WWPS"}

func main() {
	var sessions []Session

	for _, track := range tracks {
		trackSessions := parse(track)
		//for _, session := range trackSessions {
		//	log.Printf("Session: %s, %s, %s, %s, %s", session.SessionId, session.Title, session.Abstract, session.SlideshareUrl, session.YoutubeUrl)
		//}
		sessions = append(sessions, trackSessions...)
	}

	ytDevKey := os.Getenv("YOUTUBE_DEVKEY")
	ssApiKey := os.Getenv("SLIDESHARE_APIKEY")
	ssSecret := os.Getenv("SLIDESHARE_SECRET")

	for sessionIndex, session := range sessions {
		slideShowUrl := searchOnUserForKeyword(awsSlideShareUsername, session.SessionId, ssApiKey, ssSecret)
		youtubeUrl := searchOnChannelByKeyword(amazonWebServicesChannelExternalId, session.SessionId, ytDevKey)
		(&sessions[sessionIndex]).SlideshareUrl = slideShowUrl
		(&sessions[sessionIndex]).YoutubeUrl = youtubeUrl
	}

	for _, session := range sessions {
		log.Printf("Session: %s, %s, %s, %s, %s", session.SessionId, session.Title, session.Abstract, session.SlideshareUrl, session.YoutubeUrl)
	}

	mustacheSessions := AllSessions{Sessions: sessions}

	output, err := mustache.RenderFile("resources/out.mustache", mustacheSessions)
	if err != nil {
		log.Fatalf("error running mustache %v", err)
	}
	fmt.Print(output)
}