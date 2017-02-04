package app

import (
	"fmt"
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

	if err := inject.Populate(
		imagesHandlers, context, imagesFinder, imagesIndexer, logger,
	); err != nil {
		panic(fmt.Errorf("Failed to populate application graph: %v", err))
	}

	router := mux.NewRouter()
	router.HandleFunc("/find", imagesHandlers.List).Methods("GET")
	router.HandleFunc("/index", imagesHandlers.Index).Methods("POST")
	http.Handle("/", router)
}
