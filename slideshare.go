package main

import (
	slideshare "github.com/seanwatson/go-slideshare"
	"log"
	"strings"
)

var ssSearchOptions = map[string]string{
	"sort": "latest",
}


func searchOnUserForKeyword(userName string, keyword string, ssApiKey string, ssSecret string) string {
	url := ""
	ss := slideshare.SlideShare{ssApiKey, ssSecret}

	shows, err := ss.SearchSlideshows(keyword, ssSearchOptions)
	if err != nil {
		log.Fatalf("Error connecting to slideshare - %v", err)
	} else if (len(shows) == 0) {
		log.Printf("no slideshare data found for %s", keyword)
	} else {
		for _, show := range shows {
			if (strings.Contains(show.Title, "AWS re:Invent 2016") && show.Username == userName) {
				url = show.Url
				break
			}
		}
	}
	return url
}

