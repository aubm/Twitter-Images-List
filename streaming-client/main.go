package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	twitter "github.com/ChimeraCoder/anaconda"
)

type Config struct {
	consumerKey        string
	consumerSecret     string
	accessToken        string
	accessTokenSecret  string
	queueIndexEndpoint string
}

const (
	indexRequestsBufferMaxSize = 50
)

var (
	config              = Config{}
	tweetsToIndex       = make(chan twitter.Tweet)
	indexRequestsBuffer = make([]IndexRequest, 0)
)

func init() {
	config.consumerKey = os.Getenv("TWITTER_CONSUMER_KEY")
	config.consumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
	config.accessToken = os.Getenv("TWITTER_ACCESS_TOKEN")
	config.accessTokenSecret = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	config.queueIndexEndpoint = os.Getenv("QUEUE_INDEX_ENDPOINT")
}

func main() {
	twitter.SetConsumerKey(config.consumerKey)
	twitter.SetConsumerSecret(config.consumerSecret)
	twitterClient := twitter.NewTwitterApi(config.accessToken, config.accessTokenSecret)

	go listenToTweetsToIndex()

	log.Print("start listening on tweets...")
	stream := twitterClient.PublicStreamSample(nil)
	for {
		item := <-stream.C
		switch tweet := item.(type) {
		case twitter.Tweet:
			inspectTweet(tweet)
		}
	}
}

func listenToTweetsToIndex() {
	for {
		newTweet := <-tweetsToIndex
		indexRequestsBuffer = append(indexRequestsBuffer, IndexRequest{
			Description: newTweet.Text,
			Url:         newTweet.Entities.Media[0].Media_url,
		})
		if len(indexRequestsBuffer) > indexRequestsBufferMaxSize {
			flushIndexRequestsBuffer()
		}
	}
}

func flushIndexRequestsBuffer() {
	log.Printf("about to flush %v index requests", indexRequestsBufferMaxSize)
	defer func() {
		indexRequestsBuffer = []IndexRequest{}
		log.Print("flush done")
	}()
	b, err := json.Marshal(indexRequestsBuffer)
	if err != nil {
		log.Printf("failed to marshal index requests buffer: %v", err)
	}
	http.Post(config.queueIndexEndpoint, "application/json", bytes.NewReader(b))
}

func inspectTweet(tweet twitter.Tweet) {
	tweetMedia := tweet.Entities.Media
	if len(tweetMedia) < 1 {
		return
	}

	if tweetMedia[0].Type == "photo" {
		log.Printf("indexing new tweet %v", tweet.Id)
		tweetsToIndex <- tweet
	}
}

type IndexRequest struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}
