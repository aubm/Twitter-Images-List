package images

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/search"
	"github.com/aubm/twitter-image/images-api/shared"
	"cloud.google.com/go/vision"
	"google.golang.org/api/option"
)

type Indexer struct {
	Logger shared.LoggerInterface `inject:""`
	Config *shared.AppConfig `inject:""`
	HttpClient interface {
		Provide(ctx context.Context) *http.Client
	} `inject:""`
}

func (i *Indexer) Index(ctx context.Context, data IndexRequest) error {
	newImage := i.newImageFromIndexRequest(data)
	if err := i.annotateImageWithTags(ctx, newImage); err != nil {
		return fmt.Errorf("Failed to compute the new image tags: %v", err)
	}
	if err := i.putToDatastore(ctx, newImage); err != nil {
		return fmt.Errorf("Failed to put the new image into datastore: %v", err)
	}
	if err := i.putToSearchIndex(ctx, newImage); err != nil {
		return fmt.Errorf("Failed to add the new image to the search index: %v", err)
	}
	return nil
}

func (i *Indexer) newImageFromIndexRequest(data IndexRequest) *Image {
	return &Image{
		ID:          uuid.NewV4().String(),
		Url:         data.Url,
		Description: data.Description,
		Tags:        []string{},
		CreatedAt:   time.Now(),
	}
}

func (i *Indexer) annotateImageWithTags(ctx context.Context, image *Image) error {
	visionClient, err := vision.NewClient(ctx, option.WithAPIKey(i.Config.VisionAPIKey))
	if err != nil {
		return fmt.Errorf("Failed to instanciate the vision client: %v", err)
	}
	defer visionClient.Close()

	r, err := i.readUrl(ctx, image.Url)
	if err != nil {
		return fmt.Errorf("Failed to read image url %v: %v", image.Url, err)
	}

	src, err := vision.NewImageFromReader(r)
	if err != nil {
		return fmt.Errorf("Failed to create source image: %v", err)
	}

	annotations, err := visionClient.DetectLabels(ctx, src, 10)
	if err != nil {
		return fmt.Errorf("Failed to detect image labels: %v", err)
	}

	tags := []string{}
	for _, annotation := range annotations {
		tags = append(tags, annotation.Description)
	}

	i.Logger.Infof(ctx, "For image: %v, found tags: %v", image.Url, tags)
	image.Tags = tags

	return nil
}

func (i *Indexer) readUrl(ctx context.Context, url string) (io.ReadCloser, error) {
	resp, err := i.HttpClient.Provide(ctx).Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to download image at url %v: %v", url, err)
	}
	if resp.StatusCode < 200 && resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Failed to download image at url %v, response code is %v", url, resp.StatusCode)
	}
	return resp.Body, nil
}

func (i *Indexer) putToDatastore(ctx context.Context, image *Image) error {
	if _, err := datastore.Put(ctx, buildKeyForImageID(ctx, image.ID), image); err != nil {
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
	Url         string `json:"url"`
	Description string `json:"description"`
}

type searchImageEntry struct {
	Tags      string
	CreatedAt time.Time
}
