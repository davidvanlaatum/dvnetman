package service

import (
	"context"
	"dvnetman/pkg/auth"
	"dvnetman/pkg/cache"
	"dvnetman/pkg/logger"
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/openapi"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type Service struct {
	db    *modal.DBClient
	auth  *auth.Auth
	cache cache.Cache
}

func (s *Service) ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	for _, converter := range errorConverters {
		if res := converter(err); res != nil {
			if err := res.Write(r, w); err != nil {
				logger.Error(r.Context()).Err(err).Msg("error writing error response")
				return
			}
			return
		}
	}
	logger.Error(r.Context()).Err(err).Msg("no error handler found")
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (s *Service) WriteErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	//TODO implement me
	panic("implement me")
}

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

func (s *Service) checkIfModified(
	ifNoneMatch *string, ifModifiedSince *time.Time, modVersion int, modTime time.Time, res *openapi.Response,
) (err error) {
	etag := fmt.Sprintf("w/\"%d\"", modVersion)
	lastModified := modTime.Format(time.RFC1123)
	if ifModifiedSince != nil {
		if !modTime.Truncate(time.Second).After(*ifModifiedSince) {
			err = errors.WithStack(&notModifiedError{etag: etag, lastModified: lastModified})
		}
	}
	if ifNoneMatch != nil && *ifNoneMatch != "" {
		e := fmt.Sprintf("w/\"%d\"", modVersion)
		if e == *ifNoneMatch {
			err = errors.WithStack(&notModifiedError{etag: etag, lastModified: lastModified})
		}
	}
	if res.Headers == nil {
		res.Headers = make(map[string][]string)
	}
	res.Headers["ETag"] = []string{etag}
	res.Headers["Last-Modified"] = []string{lastModified}
	return
}

func NewService(ctx context.Context, db *modal.DBClient, auth *auth.Auth, cache cache.Pool) *Service {
	return &Service{
		db:    db,
		auth:  auth,
		cache: cache,
	}
}
