// Code generated by dvnetman. DO NOT EDIT.

package openapi

import (
	uuid "github.com/google/uuid"
	"time"
)

type Device struct {
	AssetTag     *string          `json:"assetTag,omitzero"`
	Created      *time.Time       `json:"created,omitzero"`
	Description  *string          `json:"description,omitzero"`
	DeviceType   *ObjectReference `json:"deviceType,omitzero"`
	Id           uuid.UUID        `json:"id"`
	Location     *ObjectReference `json:"location,omitzero"`
	Name         *string          `json:"name,omitzero"`
	Ports        []*DevicePort    `json:"ports,omitzero"`
	Rack         *ObjectReference `json:"rack,omitzero"`
	RackFace     *DeviceRackFace  `json:"rackFace,omitzero"`
	RackPosition *float64         `json:"rackPosition,omitzero"`
	Serial       *string          `json:"serial,omitzero"`
	Site         *ObjectReference `json:"site,omitzero"`
	Status       *string          `json:"status,omitzero"`
	Tags         []*Tag           `json:"tags,omitzero"`
	Updated      *time.Time       `json:"updated,omitzero"`
	Version      int              `json:"version"`
}
