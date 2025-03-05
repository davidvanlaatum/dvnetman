package modal

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var fileChecksums = []string{"md5", "sha1"}

type File struct {
	Base           `bson:",inline"`
	Name           string            `bson:"name,omitempty"`
	Size           int64             `bson:"size,omitempty"`
	ReferenceCount int64             `bson:"reference_count,omitempty"`
	ContentType    string            `bson:"content_type,omitempty"`
	Checksums      map[string]string `bson:"checksums,omitempty"`
}

func init() {
	var checksumIndexes []mongo.IndexModel
	for _, checksum := range fileChecksums {
		checksumIndexes = append(
			checksumIndexes, mongo.IndexModel{
				Keys: bson.D{{"checksums." + checksum, int32(1)}},
			},
		)
	}
	register(
		&CollectionInfo{
			Name: "files",
			Indexes: append(
				[]mongo.IndexModel{
					{
						Keys:    bson.D{{"id", int32(1)}},
						Options: options.Index().SetUnique(true),
					},
					{
						Keys: bson.D{{"name", int32(1)}},
					},
					{
						Keys: bson.D{{"size", int32(1)}},
					},
					{
						Keys: bson.D{{"reference_count", int32(1)}, {"created", int32(1)}},
					},
				},
				checksumIndexes...,
			),
		},
	)
}

func (f *File) GetCollectionName() string {
	return "files"
}
