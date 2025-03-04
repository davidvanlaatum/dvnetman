package service

import (
	"dvnetman/pkg/openapi"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"slices"
	"time"
)

func setProjection(
	fields []string, allowed []string, builder *options.FindOptionsBuilder,
) (_ *options.FindOptionsBuilder, err error) {
	if len(fields) == 0 {
		return builder, nil
	}
	projection := bson.D{{"id", 1}, {"version", 1}}
	for _, field := range fields {
		if !slices.Contains(allowed, field) {
			return nil, errors.Errorf("field %s is not allowed", field)
		}
		projection = append(projection, bson.E{field, 1})
	}
	return builder.SetProjection(projection), nil
}

func checkIfModified(
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
