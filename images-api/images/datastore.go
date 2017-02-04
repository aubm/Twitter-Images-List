package images

import "google.golang.org/appengine/datastore"

func buildKeyForImageID(imageID string) *datastore.Key {
	return datastore.NewKey(ctx, DATASTORE_ENTITY_NAME, id, 0, nil)
}
