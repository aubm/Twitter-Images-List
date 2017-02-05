package app

import (
	"fmt"
	"os"
	"net/http"

	"github.com/aubm/twitter-image/images-api/api"
	"github.com/aubm/twitter-image/images-api/images"
	"github.com/aubm/twitter-image/images-api/shared"
	"github.com/facebookgo/inject"
	"github.com/gorilla/mux"
)

func init() {
	imagesHandlers := &api.ImagesHandlers{}
	imagesFinder := &images.Finder{}
	imagesIndexer := &images.Indexer{}
	context := &api.ContextProvider{}
	logger := &shared.Logger{}
	httpClientProvider := &shared.HttpClientProvider{}
	config := initConfig()

	if err := inject.Populate(
		imagesHandlers, context, imagesFinder, imagesIndexer, logger, httpClientProvider, config,
	); err != nil {
		panic(fmt.Errorf("Failed to populate application graph: %v", err))
	}

	router := mux.NewRouter()
	router.HandleFunc("/", imagesHandlers.List).Methods("GET")
	router.HandleFunc("/index", imagesHandlers.Index).Methods("POST")
	router.HandleFunc("/queue-index", imagesHandlers.QueueIndex).Methods("POST")
	http.Handle("/", router)
}
func initConfig() *shared.AppConfig {
	return &shared.AppConfig{
		VisionAPIKey: os.Getenv("VISION_API_KEY"),
	}
}
