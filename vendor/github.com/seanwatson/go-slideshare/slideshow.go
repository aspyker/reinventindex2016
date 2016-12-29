//
// File: slideshow.go
// Date: 1/26/2014
//
// This file contains the response containers for API calls that return
// slideshows.
//
// The MIT License (MIT)
//
// Copyright (c) 2014 Sean Watson
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package slideshare

// Base structure for holding Slideshows from XML responses.
type Slideshow struct {
	Id                uint32 `xml:"ID"`
	Title             string `xml:"Title"`
	Description       string `xml:"Description"`
	Username          string `xml:"Username"`
	Status            uint8  `xml:"Status"`
	Url               string `xml:"URL"`
	ThumbnailUrl      string `xml:"ThumbnailURL"`
	ThumbnailSize     string `xml:"ThumbnailSize"`
	ThumbnailSmallUrl string `xml:"ThumbnailSmallURL"`
	Embed             string `xml:"Embed"`
	Created           string `xml:"Created"`
	Updated           string `xml:"Updated"`
	Language          string `xml:"Language"`
	Format            string `xml:"Format"`
	Download          bool   `xml:"Download"`
	DownloadUrl       string `xml:"DownloadUrl"`
	SlideshowType     uint8  `xml:"SlideshowType"`
	InContest         bool   `xml:"InContest"`
}

// Used for when multiple Slideshows are returned in a single response.
type Slideshows struct {
	Values []Slideshow `xml:"Slideshow"`
}

// Used to check if a call to EditSlideshow responded succesfully.
type EditSlideshowResponse struct {
	Edited uint32 `xml:"slideshowid"`
}

// Used to check if a call to DeleteSlideshow responded succesfully.
type DeleteSlideshowResponse struct {
	Deleted uint32 `xml:"slideshowid"`
}

// Used to check if a call to UploadSlideshow responded succesfully.
type UploadSlideshowResponse struct {
	Uploaded uint32 `xml:"slideshowid"`
}
