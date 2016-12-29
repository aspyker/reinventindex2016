package main

import (
	"io/ioutil"
	"log"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xpath"
	"strings"
	"fmt"
)

func parse(trackAbbrev string) ([]Session) {
	var sessions []Session

	input, err := ioutil.ReadFile(fmt.Sprintf("resources/%s.html", trackAbbrev))
	if err != nil {
		log.Fatalf("error reading file %s", trackAbbrev)
	}

	xml, err := gokogiri.ParseXml(input)
	if err != nil {
		log.Fatalf("error parsing xml")
	}

	xp1 := xpath.Compile("//div[starts-with(@id, 'session_')]")
	xp2 := xpath.Compile(".//span[@class='title']")
	xp3 := xpath.Compile(".//span[@class='abbreviation']")
	xp4 := xpath.Compile(".//span[@class='abstract']")

	divs, err := xml.Search(xp1)

	if err != nil {
		log.Fatalf("error finding divs")
	}

	for _, div := range divs {
		var session Session
		titles, err2 := div.Search(xp2)
		if (err2 != nil) {
			log.Printf("%v", err2)
			session.Title = "unknown"
		} else if (len(titles) == 0) {
			log.Printf("no titles")
			session.Title = "unknown"
		} else {
			session.Title = strings.TrimSpace(titles[0].Content())
		}

		abbreviations, err2 := div.Search(xp3)
		if (err2 != nil) {
			log.Printf("%v", err2)
			session.SessionId = "UKNOWN"
		} else if (len(abbreviations) == 0) {
			log.Printf("no abbreviations")
			session.SessionId = "UKNOWN"
		} else {
			firstAbbreviation := strings.TrimSpace(abbreviations[0].Content())
			session.SessionId = strings.Split(firstAbbreviation, " ")[0]
		}
		if (strings.Contains(session.SessionId, "-R")) {
			log.Printf("skipping repeat session %s", session.SessionId)
			continue
		}

		abstracts, err2 := div.Search(xp4)
		if (err2 != nil) {
			log.Printf("%v", err2)
			session.Abstract = "unknown"
		} else if (len(abstracts) == 0) {
			log.Printf("no abbreviations")
			session.Abstract = "unknown"
		} else {
			firstAbstract := strings.TrimSpace(abstracts[0].Content())
			session.Abstract = firstAbstract
		}

		sessions = append(sessions, session)
	}

	return sessions
}
