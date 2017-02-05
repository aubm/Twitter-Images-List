package images

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func buildKeyForImageID(ctx context.Context, imageID string) *datastore.Key {
	return datastore.NewKey(ctx, DATASTORE_ENTITY_NAME, imageID, 0, nil)
}
