// Code generated by dvnetman. DO NOT EDIT.

package openapi

import (
	uuid "github.com/google/uuid"
	"time"
)

type SiteResult struct {
	Created     *time.Time `json:"created,omitzero"`
	Description *string    `json:"description,omitzero"`
	Id          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Tags        []*Tag     `json:"tags,omitzero"`
	Updated     *time.Time `json:"updated,omitzero"`
	Version     int        `json:"version"`
}
