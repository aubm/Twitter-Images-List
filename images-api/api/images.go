package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aubm/twitter-image/images-api/images"
	"github.com/aubm/twitter-image/images-api/shared"
	"github.com/markbates/going/defaults"
	"golang.org/x/net/context"
	"google.golang.org/appengine/taskqueue"
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

	findOptions, err := h.buildUserFindOptions(r)
	if err != nil {
		httpError(w, http.StatusBadRequest, newError(err.Error()))
		return
	}

	imagesList, err := h.Finder.Find(ctx, findOptions)
	if err != nil {
		h.Logger.Errorf(ctx, "Call to Manager.Find failed: %v", err)
		httpError(w, http.StatusInternalServerError, serverError)
		return
	}
	writeJSON(w, imagesList, 200)
}

func (h *ImagesHandlers) buildUserFindOptions(r *http.Request) (images.FindOptions, error) {
	findOptions := images.FindOptions{}
	queryParams := r.URL.Query()

	limitValue := defaults.String(queryParams.Get("limit"), "100")
	limit, err := strconv.Atoi(limitValue)
	if err != nil {
		return findOptions, errors.New(`Invalid value for parameter "limit"`)
	}

	offsetValue := defaults.String(queryParams.Get("offset"), "0")
	offset, err := strconv.Atoi(offsetValue)
	if err != nil {
		return findOptions, errors.New(`Invalid value for parameter "offset"`)
	}

	findOptions.Limit = limit
	findOptions.Offset = offset
	findOptions.FilterTags = queryParams.Get("tags")

	return findOptions, nil
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

	ok(w)
}

func (h *ImagesHandlers) QueueIndex(w http.ResponseWriter, r *http.Request) {
	ctx := h.Ctx.New(r)

	indexRequests := []images.IndexRequest{}
	if err := json.NewDecoder(r.Body).Decode(&indexRequests); err != nil {
		httpError(w, http.StatusBadRequest, invalidJSONError)
		return
	}

	tasks, err := h.buildTasksFromIndexRequests(indexRequests)
	if err != nil {
		h.Logger.Errorf(ctx, "Failed to build tasks from index requests: %v", err)
		httpError(w, http.StatusInternalServerError, serverError)
		return
	}

	if _, err := taskqueue.AddMulti(ctx, tasks, "index-image"); err != nil {
		h.Logger.Errorf(ctx, "Failed to queue the index of the image: %v", err)
		httpError(w, http.StatusInternalServerError, serverError)
		return
	}

	ok(w)
}

func (h *ImagesHandlers) buildTasksFromIndexRequests(indexRequests []images.IndexRequest) ([]*taskqueue.Task, error) {
	tasks := []*taskqueue.Task{}
	for _, indexReq := range indexRequests {
		b, err := json.Marshal(indexReq)
		if err != nil {
			return nil, fmt.Errorf("Failed to marshal index request to JSON: %v", err)
		}
		tasks = append(tasks, &taskqueue.Task{
			Path:    "/index",
			Method:  "POST",
			Payload: b,
		})
	}
	return tasks, nil
}

func (h *ImagesHandlers) validateIndexRequest(indexRequest images.IndexRequest) error {
	if indexRequest.Url == "" {
		return errors.New(`Missing "url" property`)
	}
	return nil
}
