package images

import "time"

const (
	DATASTORE_ENTITY_NAME = "Image"
	SEARCH_INDEX_NAME     = "images"
)

type Image struct {
	ID          string    `datastore:"id" json:"id"`
	Url         string    `datastore:"url,noindex" json:"url"`
	Description string    `datastore:"description,noindex" json:"description"`
	Tags        []string  `datastore:"tags,noindex" json:"tags"`
	CreatedAt   time.Time `datastore:"created_at,noindex" json:"created_at"`
}
