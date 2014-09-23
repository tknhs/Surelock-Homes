package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/go-oauth/oauth"
)

type Tweet struct {
	Text      string
	TimeStamp string `json:"timestamp_ms"`
	User      struct {
		ScreenName string `json:"screen_name"`
	}
}

var oauthClient = oauth.Client{}

func TwitterInit(configTwitter TwitterConfig) *oauth.Credentials {
	oauthClient.Credentials.Token = configTwitter.ConsumerKey
	oauthClient.Credentials.Secret = configTwitter.ConsumerSecret
	accessToken := configTwitter.AccessToken
	accessTokenSecret := configTwitter.AccessTokenSecret

	var token *oauth.Credentials
	token = &oauth.Credentials{accessToken, accessTokenSecret}

	return token
}

func TwitterPost(token *oauth.Credentials, twText string) error {
	twUrl := "https://api.twitter.com/1.1/statuses/update.json"
	twStatus := []string{twText, strconv.Itoa(int(time.Now().Unix()))}
	twParam := make(url.Values)
	twParam.Set("status", strings.Join(twStatus, ","))

	oauthClient.SignParam(token, "POST", twUrl, twParam)
	res, err := http.PostForm(twUrl, url.Values(twParam))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(strconv.Itoa(res.StatusCode))
	}
	log.Println(twText)

	defer res.Body.Close()
	return nil
}

func TwitterStreaming(twitterTimestamp chan string, token *oauth.Credentials, serverAccount string) {
	// Kill this function after 5 minutes
	timer := time.NewTimer(time.Second * 300)
	go func() {
		<-timer.C
		log.Println("[kill] streaming twtter")
		twitterTimestamp <- "8000000000"
	}()

	// start twitterstreaming
	twUrl := "https://userstream.twitter.com/1.1/user.json"
	twParam := make(url.Values)

	oauthClient.SignParam(token, "GET", twUrl, twParam)
	twUrl = twUrl + "?" + twParam.Encode()
	res, err := http.Get(twUrl)
	if err != nil {
		log.Fatalf("failed to get a tweet\n", err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("failed to get a tweet\n", res.StatusCode)
	}

	buf := bufio.NewReader(res.Body)
	var last []byte
	for {
		b, _, err := buf.ReadLine()
		last = append(last, b...)
		var tweets [1]Tweet
		err = json.Unmarshal(last, &tweets[0])
		if err != nil {
			continue
		}
		last = []byte{}

		for i := len(tweets) - 1; i >= 0; i-- {
			user := tweets[i].User.ScreenName
			text := tweets[i].Text
			text = strings.Split(text, ",")[0]
			ts := tweets[i].TimeStamp
			if err != nil {
				continue
			}

			if user == serverAccount && text == "open" {
				timer.Stop()
				res.Body.Close()
				twitterTimestamp <- ts
			}
		}
	}
}
