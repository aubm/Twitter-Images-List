package main

import (
	"fmt"

	"github.com/ChimeraCoder/anaconda"
)

type Config struct {
	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string
}

func main() {
	config := Config{}

	anaconda.SetConsumerKey(config.consumerKey)
	anaconda.SetConsumerSecret(config.consumerSecret)
	api := anaconda.NewTwitterApi(config.accessToken, config.accessTokenSecret)

	stream := api.PublicStreamSample(nil)

	fmt.Println("Start listening on tweets...")
	for {
		item := <-stream.C
		switch tweet := item.(type) {
		case anaconda.Tweet:
			tweetMedia := tweet.Entities.Media
			if len(tweetMedia) > 0 {
				if tweetMedia[0].Type == "photo" {
					fmt.Println(tweetMedia[0].Media_url)
				}
			}
		}
	}
}
