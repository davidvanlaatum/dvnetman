// Code generated by dvnetman. DO NOT EDIT.

package openapi

import uuid "github.com/google/uuid"

type User struct {
	Email            *string    `json:"email,omitzero"`
	ExternalProvider *string    `json:"externalProvider,omitzero"`
	Id               *uuid.UUID `json:"id,omitzero"`
	Password         *string    `json:"password,omitzero"`
	Username         *string    `json:"username,omitzero"`
}
