package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

type Config struct {
	consumerKey        string
	consumerSecret     string
	accessToken        string
	accessTokenSecret  string
	queueIndexEndpoint string
}

const (
	INDEX_REQUESTS_BUFFER_MAX_SIZE = 50
)

var (
	config              = Config{}
	tweetsToIndex       = make(chan anaconda.Tweet)
	indexRequestsBuffer = []IndexRequest{}
)

func init() {
	config.consumerKey = os.Getenv("TWITTER_CONSUMER_KEY")
	config.consumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
	config.accessToken = os.Getenv("TWITTER_ACCESS_TOKEN")
	config.accessTokenSecret = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	config.queueIndexEndpoint = os.Getenv("QUEUE_INDEX_ENDPOINT")
}

func main() {
	anaconda.SetConsumerKey(config.consumerKey)
	anaconda.SetConsumerSecret(config.consumerSecret)
	api := anaconda.NewTwitterApi(config.accessToken, config.accessTokenSecret)

	go func() {
		for {
			newTweet := <-tweetsToIndex
			indexRequestsBuffer = append(indexRequestsBuffer, IndexRequest{
				Description: newTweet.Text,
				Url:         newTweet.Entities.Media[0].Media_url,
			})
			if len(indexRequestsBuffer) > INDEX_REQUESTS_BUFFER_MAX_SIZE {
				flushIndexRequestsBuffer()
			}
		}
	}()

	fmt.Println("Start listening on tweets...")
	stream := api.PublicStreamSample(nil)
	for {
		item := <-stream.C
		switch tweet := item.(type) {
		case anaconda.Tweet:
			inspectTweet(tweet)
		}
	}
}

func flushIndexRequestsBuffer() {
	fmt.Printf("About to flush %v index requests\n", INDEX_REQUESTS_BUFFER_MAX_SIZE)
	defer func() {
		indexRequestsBuffer = []IndexRequest{}
		fmt.Println("Flush done")
	}()
	b, err := json.Marshal(indexRequestsBuffer)
	if err != nil {
		fmt.Printf("Failed to marshal index requests buffer: %v\n", err)
	}
	http.Post(config.queueIndexEndpoint, "application/json", bytes.NewReader(b))
}

func inspectTweet(tweet anaconda.Tweet) {
	tweetMedia := tweet.Entities.Media
	if len(tweetMedia) < 1 {
		return
	}

	if tweetMedia[0].Type == "photo" {
		tweetsToIndex <- tweet
	}
}

type IndexRequest struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}
