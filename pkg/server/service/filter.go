package service

import (
	"dvnetman/pkg/mongo/modal"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type filter bson.D

func (f *filter) equalsStr(field string, value *string) {
	if value != nil && *value != "" {
		*f = append(*f, bson.E{field, *value})
	}
}

func (f *filter) equalsUUID(field string, value *uuid.UUID) {
	if value != nil && *value != uuid.Nil {
		*f = append(*f, bson.E{field, *value})
	}
}

func (f *filter) inUUID(field string, values []modal.UUID) {
	if len(values) > 0 {
		*f = append(*f, bson.E{field, bson.M{"$in": values}})
	}
}

func (f *filter) regex(field string, value *string, options string) {
	if value != nil && *value != "" {
		*f = append(*f, bson.E{field, bson.Regex{*value, options}})
	}
}
