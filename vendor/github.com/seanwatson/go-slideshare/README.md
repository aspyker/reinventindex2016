go-slideshare
=============

Golang wrapper for the SlideShare API

Installation
------------

If you don't already have one, you will need an API key from SlideShare.

Run the below command to download the library

`go get github.com/seanwatson/go-slideshare`

Examples
--------

Get a slideshow with it's ID

    package main

    import (
        "fmt"
        "github.com/seanwatson/go-slideshare"
    )

    func main() {
        ss := slideshare.SlideShare{"<your_api_key>", "<your_shared_secret>"}
        slideshow, err := ss.GetSlideshow(13343768)
        if err != nil {
            fmt.Println("Error getting slideshow")
            fmt.Println(err)
        }
        fmt.Println(slideshow.Title)
        fmt.Println(slideshow.Description)
    }


