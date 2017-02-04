package shared

import (
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/net/context"
	"google.golang.org/appengine/file"
)

const CONFIG_FILE_NAME = "configuration.json"

type AppConfig struct {
	VisionAPIKey string `json:"vision_api_key"`
}

type ConfigurationLoader struct {
	appConfig *AppConfig
	StorageService interface {
		ReadFile(ctx context.Context, bucket, filename string) (io.ReadCloser, error)
	} `inject:""`
	Logger LoggerInterface `inject:""`
}

func (l *ConfigurationLoader) Load(ctx context.Context) *AppConfig {
	if l.appConfig != nil {
		return l.appConfig
	}

	bucketName, err := file.DefaultBucketName(ctx)
	if err != nil {
		panic(fmt.Errorf("Failed to get the default bucket name: %v", err))
	}

	r, err := l.StorageService.ReadFile(ctx, bucketName, CONFIG_FILE_NAME)
	if err != nil {
		panic(fmt.Errorf("Failed to read configuration file: %v", err))
	}
	defer r.Close()

	l.appConfig = &AppConfig{}
	if err := json.NewDecoder(r).Decode(l.appConfig); err != nil {
		panic(fmt.Errorf("Failed to parse configuration file: %v", err))
	}

	return l.appConfig
}
