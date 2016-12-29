//
// File: favorite.go
// Date: 1/26/2014
//
// This file contains response containers for API methods that return favorites.
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

// The base structure for holding favorites from XML responses.
type Favorite struct {
	SlideshowId uint32 `xml:"slideshow_id"`
	TagText     string `xml:"tag_text"`
}

// Used when an API response contains multiple favorites.
type Favorites struct {
	Values []Favorite `xml:"favorite"`
}

// Used to check if a call to AddFavorite returns successfully. If the
// SlideshowId is non-zero after the call then it was a success.
type AddFavoriteResponse struct {
	SlideshowId uint32 `xml:"slideshowid"`
}

// Used to check if a call to CheckFavorite returns successfully. If the
// slideshow is favorited The Favorited field will be true.
type CheckFavoriteResponse struct {
	SlideshowId uint32 `xml:"slideshowid"`
	User        string `xml:"user"`
	Favorited   bool   `xml:"favorited"`
}
