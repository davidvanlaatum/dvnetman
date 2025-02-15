package dto

import (
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/openapi"
	"github.com/google/uuid"
)

func (c *Converter) getManufacturerRef(id *modal.UUID) *openapi.ObjectReference {
	ref := &openapi.ObjectReference{Id: *(*uuid.UUID)(id)}
	cb := func(m *modal.Manufacturer) {
		ref.DisplayName = &m.Name
	}
	addToQueue(&c.manufacturerQueue, id, cb)
	return ref
}
