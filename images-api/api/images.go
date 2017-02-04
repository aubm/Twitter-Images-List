package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aubm/twitter-image/images-api/images"
	"github.com/aubm/twitter-image/images-api/shared"
	"golang.org/x/net/context"
)

type ImagesHandlers struct {
	Ctx    ContextProviderInterface `inject:""`
	Logger shared.LoggerInterface   `inject:""`
	Finder interface {
		Find(ctx context.Context, options images.FindOptions) (*images.FindResult, error)
	} `inject:""`
	Indexer interface {
		Index(ctx context.Context, data images.IndexRequest) error
	} `inject:""`
}

func (h *ImagesHandlers) List(w http.ResponseWriter, r *http.Request) {
	ctx := h.Ctx.New(r)

	imagesList, err := h.Finder.Find(ctx, images.FindOptions{
		Limit:      100,
		Offset:     0,
		FilterTags: r.URL.Query().Get("tags"),
	})
	if err != nil {
		h.Logger.Errorf(ctx, "Call to Manager.Find failed: %v", err)
		httpError(w, http.StatusInternalServerError, serverError)
		return
	}
	writeJSON(w, imagesList, 200)
}

func (h *ImagesHandlers) Index(w http.ResponseWriter, r *http.Request) {
	ctx := h.Ctx.New(r)

	indexRequest := images.IndexRequest{}
	if err := json.NewDecoder(r.Body).Decode(&indexRequest); err != nil {
		httpError(w, http.StatusBadRequest, invalidJSONError)
		return
	}

	if err := h.validateIndexRequest(indexRequest); err != nil {
		httpError(w, http.StatusBadRequest, newError(err.Error()))
		return
	}

	if err := h.Indexer.Index(ctx, indexRequest); err != nil {
		h.Logger.Errorf(ctx, "Call to Indexer.Index failed: %v", err)
		httpError(w, http.StatusInternalServerError, serverError)
		return
	}

	writeJSON(w, map[string]interface{}{
		"error": false,
	}, 200)
}

func (h *ImagesHandlers) validateIndexRequest(indexRequest images.IndexRequest) error {
	if indexRequest.Url == "" {
		return errors.New(`Missing "url" property`)
	}
	return nil
}
