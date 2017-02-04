package images

import (
	"fmt"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/search"
)

type Indexer struct {
}

func (i *Indexer) Index(ctx context.Context, data IndexRequest) error {
	newImage := &Image{
		ID:          uuid.NewV4().String(),
		Url:         data.Url,
		GsUrl:       data.GsUrl,
		Description: data.Description,
		Tags:        []string{},
		CreatedAt:   time.Now(),
	}
	if err := i.putToDatastore(ctx, newImage); err != nil {
		return fmt.Errorf("Failed to put the new image into datastore: %v", err)
	}
	if err := i.putToSearchIndex(ctx, newImage); err != nil {
		return fmt.Errorf("Failed to add the new image to the search index: %v", err)
	}
	return nil
}

func (i *Indexer) putToDatastore(ctx context.Context, image *Image) error {
	if _, err := datastore.Put(ctx, buildKeyForImageID(image.ID), image); err != nil {
		return fmt.Errorf("Datastore operation failed: %v", err)
	}
	return nil
}

func (i *Indexer) putToSearchIndex(ctx context.Context, image *Image) error {
	index, err := search.Open(SEARCH_INDEX_NAME)
	if err != nil {
		return fmt.Errorf("Failed to open the search index: %v", err)
	}
	if _, err := index.Put(ctx, image.ID, &searchImageEntry{
		Tags:      strings.Join(image.Tags, ", "),
		CreatedAt: image.CreatedAt,
	}); err != nil {
		return fmt.Errorf("Failed to put the image to index: %v", err)
	}
	return nil
}

type IndexRequest struct {
	GsUrl       string `json:"gs_url"`
	Url         string `json:"url"`
	Description string `json:"description"`
}

type searchImageEntry struct {
	Tags      string
	CreatedAt time.Time
}
