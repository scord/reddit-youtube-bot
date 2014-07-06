package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Session struct {
	Cookie  *http.Cookie
	Modhash string `json:"modhash"`
}

var client http.Client

func PostLink(user, pass, title, linkURL, subreddit string) error {

	session, err := Login(user, pass)

	if err != nil {
		return err
	}

	fmt.Println("Login Successful")

	err = Submit(title, linkURL, subreddit, session)

	fmt.Println(err)

	return err
}

func Submit(title string, linkURL string, subreddit string, session *Session) error {

	submitURL := fmt.Sprintf("http://www.reddit.com/api/submit")

	values := url.Values{
		"url":   {linkURL},
		"kind":  {"link"},
		"sr":    {subreddit},
		"title": {title},
		"r":     {subreddit},
		"uh":    {session.Modhash},
	}

	resp, err := Post(submitURL, values, session)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	r := &Response{}

	err = json.NewDecoder(resp.Body).Decode(r)

	if len(r.JSON.Errors) != 0 {
		var msg []string
		for _, k := range r.JSON.Errors {
			msg = append(msg, k[1])
		}
		return errors.New(strings.Join(msg, ", "))
	}

	return nil

}

type Response struct {
	JSON struct {
		Errors [][]string
		Data   struct {
			Modhash string
		}
	}
}

func Post(postURL string, postValues url.Values, session *Session) (*http.Response, error) {

	req, err := http.NewRequest("POST", postURL+"?"+postValues.Encode(), nil)

	if err != nil {
		return nil, err
	}

	if session.Cookie != nil {
		req.AddCookie(session.Cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println(resp)

	if resp.StatusCode != http.StatusOK {
		return resp, errors.New(resp.Status)
	}

	return resp, nil

}

func Login(user, pass string) (*Session, error) {

	session := &Session{}

	loginURL := fmt.Sprintf("http://www.reddit.com/api/login/%s", user)

	values := url.Values{
		"user":     {user},
		"passwd":   {pass},
		"api_type": {"json"},
	}

	resp, err := Post(loginURL, values, session)

	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "reddit_session" {
			session.Cookie = cookie
		}
	}

	r := &Response{}

	err = json.NewDecoder(resp.Body).Decode(r)

	if len(r.JSON.Errors) != 0 {
		var msg []string
		for _, k := range r.JSON.Errors {
			msg = append(msg, k[1])
		}
		return nil, errors.New(strings.Join(msg, ", "))
	}

	session.Modhash = r.JSON.Data.Modhash

	return session, nil
}
