package images

import (
	"fmt"

	"github.com/aubm/twitter-image/images-api/shared"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/search"
)

type Finder struct {
	Logger shared.LoggerInterface `inject:""`
}

func (f *Finder) Find(ctx context.Context, options FindOptions) (*FindResult, error) {
	it, err := f.performSearch(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform the search operation: %v", err)
	}

	keys, err := f.extractKeysFromIterator(ctx, it)
	if err != nil {
		return nil, fmt.Errorf("Failed to extract datastore keys from search result: %v", err)
	}

	imagesList, err := f.getImagesListFromDatastore(ctx, keys)
	if err != nil {
		return nil, fmt.Errorf("Failed to get images list from datastore: %v", err)
	}

	return &FindResult{
		Items:      imagesList,
		Count:      len(imagesList),
		Offset:     options.Offset,
	}, nil
}

func (f *Finder) performSearch(ctx context.Context, options FindOptions) (*search.Iterator, error) {
	index, err := search.Open(SEARCH_INDEX_NAME)
	if err != nil {
		return nil, fmt.Errorf("Failed to open the search index %v: %v", SEARCH_INDEX_NAME, err)
	}

	query := ""
	if options.FilterTags != "" {
		query = fmt.Sprintf("Tags: (%v)", options.FilterTags)
	}

	return index.Search(ctx, query, &search.SearchOptions{
		Limit:   options.Limit,
		Offset:  options.Offset,
		IDsOnly: true,
	}), nil
}

func (f *Finder) extractKeysFromIterator(ctx context.Context, it *search.Iterator) ([]*datastore.Key, error) {
	keys := []*datastore.Key{}
	for {
		id, err := it.Next(nil)
		if err != nil {
			if err == search.Done {
				break
			}
			return nil, fmt.Errorf("Failed to read next search entry: %v", err)
		}
		keys = append(keys, buildKeyForImageID(ctx, id))
	}
	return keys, nil
}

func (f *Finder) getImagesListFromDatastore(ctx context.Context, keys []*datastore.Key) ([]Image, error) {
	imagesList := make([]Image, len(keys))
	if err := datastore.GetMulti(ctx, keys, imagesList); err != nil {
		return nil, fmt.Errorf("Failed to get multiple keys from datastore: %v", err)
	}
	return imagesList, nil
}

type FindOptions struct {
	Limit      int
	Offset     int
	FilterTags string
}

type FindResult struct {
	Items       []Image `json:"items"`
	MaxPerPage  int     `json:"max_per_page"`
	Offset      int     `json:"offset"`
	Count       int     `json:"count"`
}
