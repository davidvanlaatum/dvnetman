package service

import (
	"dvnetman/pkg/openapi"
	"github.com/pkg/errors"
	"net/http"
)

type notModifiedError struct {
	etag         string
	lastModified string
}

func (n *notModifiedError) Error() string {
	return "not modified"
}

func init() {
	RegisterErrorConverter(
		func(err error) *openapi.Response {
			var e *notModifiedError
			if errors.As(err, &e) {
				return &openapi.Response{
					Code: http.StatusNotModified,
					Headers: map[string][]string{
						"ETag":          {e.etag},
						"Last-Modified": {e.lastModified},
					},
				}
			}
			return nil
		},
	)
}
