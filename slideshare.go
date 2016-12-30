package main

import (
	slideshare "github.com/seanwatson/go-slideshare"
	"log"
	"strings"
	"strconv"
)

var ssSearchOptions = map[string]string{
	"sort": "latest",
}

var cachedSlideshowsToUrl = make(map[string]string)

const numCachePages = 25
const cachePageSize = 200

func searchOnUserForKeyword(userName string, keyword string, ssApiKey string, ssSecret string) string {
	cachedShowUrl, ok := cachedSlideshowsToUrl[keyword]
	if (ok) {
		return cachedShowUrl
	}
	log.Printf("cache miss on %s", keyword)

	url := ""
	ss := slideshare.SlideShare{ssApiKey, ssSecret}

	shows, err := ss.SearchSlideshows(keyword, ssSearchOptions)
	if err != nil {
		log.Fatalf("Error connecting to slideshare - %v", err)
	}

	if len(shows) != 0 {
		for _, show := range shows {
			if (strings.Contains(show.Title, "AWS re:Invent 2016") && show.Username == userName) {
				url = show.Url
				break
			}
		}
	}

	if url == "" {
		log.Printf("no slideshare data found for %s", keyword)
	}

	return url
}

func cacheSlideSharesByUsername(userName string, ssApiKey string, ssSecret string) {
	ss := slideshare.SlideShare{ssApiKey, ssSecret}
	var ssUsernameOptions = map[string]string {
		"limit": strconv.Itoa(cachePageSize),
	}

	for ii := 0; ii < numCachePages; ii++ {
		offset := ii * cachePageSize
		ssUsernameOptions["offset"] = strconv.Itoa(offset)
		// this requires a fix I need to upstream into slideshare client
		shows, err := ss.GetSlideshowsByUser(userName, ssUsernameOptions)
		if err != nil {
			log.Fatalf("Error connecting to slideshare - %v", err)
		} else if (len(shows) == 0) {
			log.Printf("finished getting slideshows at offset = %v", offset)
			return
		}

		// to speed up loading of cache (assumes in reverse cronological ordering)
		if !allShowsAreNewerThan2016(shows) {
			return
		}

		for _, show := range shows {
			//log.Printf("title = %s", show.Title)
			if (strings.Contains(show.Title, "AWS re:Invent 2016")) {
				fields := strings.FieldsFunc(show.Title, func (r rune) bool {
					return r == '(' || r == ')'
				})

				for _, field := range fields {
					for _, track := range tracks {
						trimmedField := strings.TrimSpace(field)
						if (strings.HasPrefix(trimmedField, track)) {
							log.Printf("adding %s to cache", trimmedField)
							cachedSlideshowsToUrl[trimmedField] = show.Url
						}
					}
				}
			}
		}
	}
}

func allShowsAreNewerThan2016(shows []slideshare.Slideshow) bool {
	for _, show := range shows {
		// pretty hacky date parsing, could/should have used time
		created := show.Created
		fields := strings.Split(created, "-")
		year, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Fatalf("created data malformed %s", created)
		}
		if (year < 2016) {
			return false
		}
	}
	return true
}

