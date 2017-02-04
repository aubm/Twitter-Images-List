package shared

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

type HttpClientProvider struct{}

func (h *HttpClientProvider) Provide(ctx context.Context) *http.Client {
	return urlfetch.Client(ctx)
}
