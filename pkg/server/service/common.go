package service

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"slices"
)

func (s *Service) setProjection(
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
