package shared

import (
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

type StorageService struct {
}

func (s *StorageService) ReadFile(ctx context.Context, bucket, filename string) (io.ReadCloser, error) {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to instanciate the storage client: %v", err)
	}

	r, err := storageClient.Bucket(bucket).Object(filename).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the file %v in bucket %v: %v", filename, bucket, err)
	}

	return r, nil
}
