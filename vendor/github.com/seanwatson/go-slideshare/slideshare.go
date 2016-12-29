//
// File: slideshare.go
// Date: 1/26/2014
//
// This file contains the main SlideShare struct and all of the API methods.
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

//TODO: tests
//TODO: check api errors, make real errors
//TODO: examples

// Package slideshare simplifies making calls to SlideShare.net's API.
package slideshare

import (
	"crypto/sha1"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Base URL for all API calls
var baseUrl = "https://www.slideshare.net/api/2"

// Parses the XML response from a call to the SlideShare API and populates the
// container with the result. The appropriate container for the API method
// called should be passed for all response fields to be populated correctly.
//
// Args:
//   resp: The response from the HTTP GET call.
//   container: A struct that can hold the data represented in the XML
//              response.
func parseFromHttpResponse(resp *http.Response, container interface{}) error {
	//TODO: check the HEAD
	responseBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	return xml.Unmarshal([]byte(responseBody), container)
}

// The base SlideShare structure on which all API methods can be called.
//
// Values:
//   ApiKey: The developer's API key. Obtained by request from SlideShare.
//   SharedSecret: The developer's shared secret. Obtained by request from
//                 SlideShare.
type SlideShare struct {
	ApiKey       string // Developer key provided by SlideShare
	SharedSecret string // Developer secret provided by SlideShare
}

// Generates a SlideShare API URL given the method and a list of arguments
// to include in the URL. This automatically generates the timestamp and hash
// as well as adds the API key to the URL.
//
// Args:
//   method: The api method to call
//   args: A map of argument names to values to include in the URL
func (s *SlideShare) getUrl(method string, args map[string]string) string {
	values := url.Values{}
	for k, v := range args {
		values.Set(k, v)
	}
	values.Set("api_key", s.ApiKey)
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	values.Set("ts", timestamp)
	hash := sha1.New()
	io.WriteString(hash, s.SharedSecret+timestamp)
	values.Set("hash", fmt.Sprintf("%x", hash.Sum(nil)))
	return baseUrl + "/" + method + "?" + values.Encode()
}

// Gets a slideshow given the id. This does not support getting slideshows by
// url.
//
// Args:
//   id: The id of the slideshow to get info about.
//   optArgs: Optional. A map of parameter names to values to add to the call.
//              username: username to use with authenticated lookup.
//              password: password for the given username. Required with
//                        username.
//              exclude_tags: set to "Y" to exclude tags from the response.
//              detailed: set to "Y" to get a more detailed response.
//              get_transcript: set to "Y" to include a transcript in the
//                              response.
func (s *SlideShare) GetSlideshow(id int, optArgs ...map[string]string) (
	Slideshow, error) {
	args := map[string]string{
		"slideshow_id": strconv.Itoa(id),
	}
	if len(optArgs) > 0 {
		for k, v := range optArgs[0] {
			args[k] = v
		}
	}
	url := s.getUrl("get_slideshow", args)
	resp, err := http.Get(url)
	if err != nil {
		return Slideshow{}, err
	}
	slideshow := Slideshow{}
	err = parseFromHttpResponse(resp, &slideshow)
	return slideshow, err
}

// Gets slideshows associated with a tag.
//
// Args:
//   tag: The tag to lookup shows for.
//   limit: The maximum number of slideshows to return.
//   offset: The number of shows to skip over from the starts of the tag's list
func (s *SlideShare) GetSlideshowsByTag(tag string, limit int,
	offset int) ([]Slideshow, error) {
	args := map[string]string{
		"tag":    tag,
		"limit":  strconv.Itoa(limit),
		"offset": strconv.Itoa(offset),
	}
	url := s.getUrl("get_slideshows_by_tag", args)
	resp, err := http.Get(url)
	if err != nil {
		return []Slideshow{}, err
	}
	slideshows := Slideshows{}
	err = parseFromHttpResponse(resp, &slideshows)
	return slideshows.Values, err
}

// Gets slideshows associated with a group.
//
// Args:
//   groupName: The group to lookup slideshows for.
//   limit: The maximum number of slideshows to return.
//   offset: The number of shows to skip over from the start of the group's
//           list
func (s *SlideShare) GetSlideshowsByGroup(groupName string, limit int,
	offset int) ([]Slideshow, error) {
	args := map[string]string{
		"group_name": groupName,
		"limit":      strconv.Itoa(limit),
		"offset":     strconv.Itoa(offset),
	}
	url := s.getUrl("get_slideshows_by_group", args)
	resp, err := http.Get(url)
	if err != nil {
		return []Slideshow{}, err
	}
	slideshows := Slideshows{}
	err = parseFromHttpResponse(resp, &slideshows)
	return slideshows.Values, err
}

// Gets the slideshows uploaded by a user.
//
// Args:
//   usernameFor: The username to get slideshows from.
//   optArgs: Optional. A map of parameter names to values to add to the call.
//              username: Username to use with authenticated lookup.
//              password: Password for the given username. Required with
//                        username.
//              limit: The maximum number of shows to return.
//              offset: The number of shows to skip over.
//              detailed: Set to "Y" to get a more detailed response.
//              get_unconverted: Set to "Y" to get unconverted shows.
func (s *SlideShare) GetSlideshowsByUser(usernameFor string,
	optArgs ...map[string]string) ([]Slideshow, error) {
	args := map[string]string{
		"username_for": usernameFor,
	}
	if len(optArgs) > 0 {
		for k, v := range optArgs[0] {
			args[k] = v
		}
	}
	url := s.getUrl("get_slideshows_by_tag", args)
	resp, err := http.Get(url)
	if err != nil {
		return []Slideshow{}, err
	}
	slideshows := Slideshows{}
	err = parseFromHttpResponse(resp, &slideshows)
	return slideshows.Values, err
}

// Search slideshows.
//
// Args:
//   q: The search query string.
//   optArgs: Optional. A map of parameter names to values to add to the call.
//              page: The page number to return results from.
//              items_per_page: The number of items to return per page.
//              lang: Langauge to restrict results to. ('es', 'pt', 'fr', etc.)
//              sort: How to sort the results. (relavent, mostdownloaded, latest)
//              upload_date: Time restriction. (any, week, month, year)
//              what: set to "tag" to search for tags.
//              download: set to "1" to restrict to shows that can be downloaded.
//              fileformat: Restrict to format. (all, pdf, ppt, ...)
//              file_type: Restrict to file type. (all, presentations, documents)
//              cc: set to "1" to restrict to creative commons licensed shows.
//              cc_adapt: set to "1" to restrict to adaptation creative commons.
//              cc_commercial: set to "1" to restrict to commercial creative commons.
//              detailed: Set to "Y" to get a more detailed response.
//              get_transcript: Set to "Y" to get a transcript of the show in
//                              responses.
func (s *SlideShare) SearchSlideshows(q string, optArgs ...map[string]string) (
	[]Slideshow, error) {
	args := map[string]string{
		"q": q,
	}
	if len(optArgs) > 0 {
		for k, v := range optArgs[0] {
			args[k] = v
		}
	}
	url := s.getUrl("search_slideshows", args)
	resp, err := http.Get(url)
	if err != nil {
		return []Slideshow{}, err
	}
	slideshows := Slideshows{}
	err = parseFromHttpResponse(resp, &slideshows)
	return slideshows.Values, err
}

// Edit the slideshow with the given ID.
//
// Args:
//   username: The username of the slideshow owner.
//   password: The password for the owner.
//   slideshowId: The ID of the slideshow to edit.
//   optArgs: Optional. A map of parameter names to values to add to the call.
//              slideshow_title: Title of the slideshow.
//              slideshow_tags: Comma seperated list of tags to search for.
//              make_slideshow_private: Set to "Y" to make the show private.
//              generate_secret_url: Set to "Y" to return a secret URL.
//              allow_embeds: Set to "Y" to allow embeds of this show.
//              share_with_contacts: Set to "Y" to let contacts view the show.
func (s *SlideShare) EditSlideshow(username, password, slideshowId string,
	optArgs ...map[string]string) error {
	args := map[string]string{
		"username":     username,
		"password":     password,
		"slideshow_id": slideshowId,
	}
	if len(optArgs) > 0 {
		for k, v := range optArgs[0] {
			args[k] = v
		}
	}
	url := s.getUrl("edit_slideshow", args)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	editSlideshowResponse := EditSlideshowResponse{}
	err = parseFromHttpResponse(resp, &editSlideshowResponse)
	if err != nil {
		return err
	}
	if editSlideshowResponse.Edited == 0 {
		return errors.New("Slideshow was not edited.")
	}
	return nil
}

// Deletes the slideshow with ID from the user's account.
//
// Args:
//   username: username of the slideshow owner.
//   password: password for the account.
//   slideshowId: the ID of the slideshow to delete.
func (s *SlideShare) DeleteSlideshow(username, password, slideshowId string) error {
	args := map[string]string{
		"username":     username,
		"password":     password,
		"slideshow_id": slideshowId,
	}
	url := s.getUrl("delete_slideshow", args)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	deleteSlideshowResponse := DeleteSlideshowResponse{}
	err = parseFromHttpResponse(resp, &deleteSlideshowResponse)
	if err != nil {
		return err
	}
	if deleteSlideshowResponse.Deleted == 0 {
		return errors.New("Slideshow was not deleted.")
	}
	return nil
}

// Upload a slideshow to the user's account.
//
// Args:
//   username: The username of the slideshow owner.
//   password: The password for the owner.
//   slideshowTitle: The title for the slideshow.
//   uploadUrl: A URL to the slideshow file. (http://example.com/show.ppt)
//   optArgs: Optional. A map of parameter names to values to add to the call.
//              make_src_public: Set to "Y" to allow downloads.
//              slideshow_description: A description of the slideshow.
//              slideshow_tags: Comma seperated list of tags to search for.
//              make_slideshow_private: Set to "Y" to make the show private.
//              generate_secret_url: Set to "Y" to return a secret URL.
//              allow_embeds: Set to "Y" to allow embeds of this show.
//              share_with_contacts: Set to "Y" to let contacts view the show.
func (s *SlideShare) UploadSlideshow(username, password, slideshowTitle,
	uploadUrl string, optArgs ...map[string]string) error {
	args := map[string]string{
		"username":        username,
		"password":        password,
		"slideshow_title": slideshowTitle,
		"upload_url":      uploadUrl,
	}
	if len(optArgs) > 0 {
		for k, v := range optArgs[0] {
			args[k] = v
		}
	}
	url := s.getUrl("upload_slideshow", args)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	uploadSlideshowResponse := UploadSlideshowResponse{}
	err = parseFromHttpResponse(resp, &uploadSlideshowResponse)
	if err != nil {
		return err
	}
	if uploadSlideshowResponse.Uploaded == 0 {
		return errors.New("Could not upload slideshow.")
	}
	return nil
}

// Gets the groups associated with a user's account. Optionally a username
// and password may be provided to look up private group memberships for
// a given account. Otherwise username and password should be passed as empty
// string.
//
// Args:
//   usernameFor: The username to look up groups for.
//   username: Optional. Username to use for authenticated lookup.
//   password: Optional. Password for "username" account.
func (s *SlideShare) GetUserGroups(usernameFor, username,
	password string) ([]Group, error) {
	args := map[string]string{
		"username_for": usernameFor,
	}
	if username != "" {
		args["username"] = username
		args["password"] = password
	}
	url := s.getUrl("get_user_groups", args)
	resp, err := http.Get(url)
	if err != nil {
		return []Group{}, err
	}
	groups := Groups{}
	err = parseFromHttpResponse(resp, &groups)
	return groups.Values, err
}

// Gets the favorite slideshows associated with a user's account.
//
// Args:
//   usernameFor: The username to look up favorites for.
func (s *SlideShare) GetUserFavorites(usernameFor string) ([]Favorite, error) {
	args := map[string]string{
		"username_for": usernameFor,
	}
	url := s.getUrl("get_user_favorites", args)
	resp, err := http.Get(url)
	if err != nil {
		return []Favorite{}, err
	}
	favorites := Favorites{}
	err = parseFromHttpResponse(resp, &favorites)
	return favorites.Values, err
}

// Gets the contacts associated with a user's account.
//
// Args:
//   usernameFor: The username to look up contacts for.
//   limit: The maximum number of contacts to return.
//   offset: The number of contacts to skip over from the start of the user's
//           contact list
func (s *SlideShare) GetUserContacts(usernameFor string, limit int,
	offset int) ([]Contact, error) {
	args := map[string]string{
		"username_for": usernameFor,
		"limit":        strconv.Itoa(limit),
		"offset":       strconv.Itoa(offset),
	}
	url := s.getUrl("get_user_contacts", args)
	resp, err := http.Get(url)
	if err != nil {
		return []Contact{}, err
	}
	contacts := Contacts{}
	err = parseFromHttpResponse(resp, &contacts)
	return contacts.Values, err
}

// Gets the tags associated with a user's account.
//
// Args:
//   username: The username to look up favorite for.
//   password: Password for the user.
func (s *SlideShare) GetUserTags(username, password string) ([]Tag, error) {
	args := map[string]string{
		"username": username,
		"password": password,
	}
	url := s.getUrl("get_user_tags", args)
	resp, err := http.Get(url)
	if err != nil {
		return []Tag{}, err
	}
	tags := Tags{}
	err = parseFromHttpResponse(resp, &tags)
	return tags.Values, err
}

// Adds the slideshow with 'slideshowId' to the user's favorties.
//
// Args:
//   username: The username to set the favorite for.
//   password: Password for the user.
//   slideshowId: ID of the slideshow to add.
//
// Returns:
//   nil on success, an error if a problem occured.
func (s *SlideShare) AddFavorite(username string, password string,
	slideshowId int) error {
	args := map[string]string{
		"username":     username,
		"password":     password,
		"slideshow_id": strconv.Itoa(slideshowId),
	}
	url := s.getUrl("add_favorite", args)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	addFavoriteResponse := AddFavoriteResponse{}
	err = parseFromHttpResponse(resp, &addFavoriteResponse)
	if err != nil {
		return err
	}
	// The slideshow ID is returned in the response on success so it will be
	// non-zero if successful.
	if addFavoriteResponse.SlideshowId == 0 {
		return errors.New("Could not add favorite")
	}
	return nil
}

// Checks whether the user has favorited the slideshow with 'slideshowId'.
//
// Args:
//   username: The username to look up favorite for.
//   password: Password for the user.
//   slideshowId: ID of the slideshow to check.
//
// Returns:
//   true if the show is favorited by the user, false otherwise.
func (s *SlideShare) CheckFavorite(username string, password string,
	slideshowId int) (bool, error) {
	args := map[string]string{
		"username":     username,
		"password":     password,
		"slideshow_id": strconv.Itoa(slideshowId),
	}
	url := s.getUrl("check_favorite", args)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	checkFavoriteResponse := CheckFavoriteResponse{}
	err = parseFromHttpResponse(resp, &checkFavoriteResponse)
	return checkFavoriteResponse.Favorited, err
}
