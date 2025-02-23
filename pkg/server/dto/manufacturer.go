package dto

import (
	"context"
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/openapi"
	"dvnetman/pkg/utils"
	"github.com/google/uuid"
)

func (c *Converter) checkManufacturerRefExists(ref *openapi.ObjectReference) *modal.UUID {
	if ref == nil {
		return nil
	}
	addToErrQueue(
		&c.manufacturerQueue, (*modal.UUID)(&ref.Id), func(id modal.UUID, dt *modal.Manufacturer) error {
			return nil
		},
	)
	return (*modal.UUID)(&ref.Id)
}

func (c *Converter) getManufacturerRef(id *modal.UUID) *openapi.ObjectReference {
	if id == nil {
		return nil
	}
	ref := &openapi.ObjectReference{Id: *(*uuid.UUID)(id)}
	cb := func(m *modal.Manufacturer) {
		ref.DisplayName = utils.ToPtr(m.Name)
	}
	addToQueue(&c.manufacturerQueue, id, cb)
	return ref
}

func (c *Converter) ManufacturerToOpenAPI(ctx context.Context, mod *modal.Manufacturer) (
	man *openapi.Manufacturer, err error,
) {
	man = &openapi.Manufacturer{
		Id:      *c.cloneUUID((*uuid.UUID)(mod.ID)),
		Name:    mod.Name,
		Created: utils.ToPtr(mod.Created),
		Updated: utils.ToPtr(mod.Updated),
		Version: mod.Version,
		//Tags: mod.Tags, TODO
	}
	if err = c.resolveQueue(ctx); err != nil {
		return
	}
	return
}

func (c *Converter) ManufacturerToOpenAPISearchResults(
	ctx context.Context, manufacturers []*modal.Manufacturer,
) (res []*openapi.ManufacturerResult, err error) {
	res = utils.MapTo(
		manufacturers, func(manufacturer *modal.Manufacturer) *openapi.ManufacturerResult {
			return &openapi.ManufacturerResult{
				Id:      *(*uuid.UUID)(manufacturer.ID),
				Created: utils.ToPtr(manufacturer.Created),
				Updated: utils.ToPtr(manufacturer.Updated),
				Name:    manufacturer.Name,
				Version: manufacturer.Version,
				//Tags:    nil, TODO
			}
		},
	)
	if err = c.resolveQueue(ctx); err != nil {
		return
	}
	return
}

func (c *Converter) UpdateManufacturerFromOpenAPI(
	ctx context.Context, body *openapi.Manufacturer, mod *modal.Manufacturer,
) error {
	mod.Name = body.Name
	//mod.Tags = body.Tags TODO
	return c.resolveQueue(ctx)
}
