package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type redditSession struct {
	client  *http.Client
	modhash string `json:"modhash"`
}

// Submit submits a post to reddit
func (session *redditSession) Submit(title string, linkURL string, subreddit string) error {

	submitURL := "http://www.reddit.com/api/submit"

	values := url.Values{
		"url":   {linkURL},
		"kind":  {"link"},
		"sr":    {subreddit},
		"title": {title},
		"r":     {subreddit},
		"uh":    {session.modhash},
	}

	_, err := session.postRequest(submitURL, values)

	if err != nil {
		return err
	}

	return nil
}

type response struct {
	JSON struct {
		Errors [][]string
		Data   struct {
			Modhash string
		}
	}
}

// sends a general post request to reddit, updating the session with any new cookies or modhash
func (session *redditSession) postRequest(postURL string, postValues url.Values) (*http.Response, error) {

	resp, err := session.client.PostForm(postURL, postValues)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, errors.New(resp.Status)
	}

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	r := &response{}

	err = json.NewDecoder(resp.Body).Decode(r)

	if err != nil {
		return nil, err
	}

	if len(r.JSON.Errors) != 0 {
		var msg []string
		for _, k := range r.JSON.Errors {
			msg = append(msg, k[1])
		}
		return nil, errors.New(strings.Join(msg, ", "))
	}

	session.modhash = r.JSON.Data.Modhash

	return resp, nil

}

// Login creates a new reddit session which contains cookies and the modhash
func Login(user, pass string) (*redditSession, error) {

	session := &redditSession{}

	cookieJar, _ := cookiejar.New(nil)

	session.client = &http.Client{
		Jar: cookieJar,
	}

	loginURL := fmt.Sprintf("http://www.reddit.com/api/login/%s", user)

	values := url.Values{
		"user":     {user},
		"passwd":   {pass},
		"api_type": {"json"},
	}

	_, err := session.postRequest(loginURL, values)

	if err != nil {
		return nil, err
	}

	return session, nil
}
