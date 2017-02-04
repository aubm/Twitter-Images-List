package images

import (
	"google.golang.org/appengine/datastore"
	"golang.org/x/net/context"
)

func buildKeyForImageID(ctx context.Context, imageID string) *datastore.Key {
	return datastore.NewKey(ctx, DATASTORE_ENTITY_NAME, imageID, 0, nil)
}
