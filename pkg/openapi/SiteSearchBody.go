// Code generated by dvnetman. DO NOT EDIT.

package openapi

import uuid "github.com/google/uuid"

type SiteSearchBody struct {
	Fields    []string    `json:"fields,omitzero"`
	Ids       []uuid.UUID `json:"ids,omitzero"`
	Name      *string     `json:"name,omitzero"`
	NameRegex *string     `json:"nameRegex,omitzero"`
}
