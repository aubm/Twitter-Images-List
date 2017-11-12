package images

import "time"

const (
	datastoreEntityName = "Image"
	searchIndexName     = "images_01"
)

type Image struct {
	ID          string    `datastore:"id" json:"id"`
	Url         string    `datastore:"url,noindex" json:"url"`
	Description string    `datastore:"description,noindex" json:"description"`
	Tags        []string  `datastore:"tags,noindex" json:"tags"`
	CreatedAt   time.Time `datastore:"created_at,noindex" json:"created_at"`
}
