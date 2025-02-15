package dto

import (
	"context"
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/openapi"
	"dvnetman/pkg/utils"
	"fmt"
	"github.com/google/uuid"
)

func (c *Converter) getDeviceTypeRef(id *modal.UUID) *openapi.ObjectReference {
	if id == nil {
		return nil
	}
	ref := &openapi.ObjectReference{Id: *(*uuid.UUID)(id)}
	cb := func(dt *modal.DeviceType) {
		addToQueue(
			&c.manufacturerQueue, dt.Manufacturer, func(m *modal.Manufacturer) {
				ref.DisplayName = utils.ToPtr(fmt.Sprintf("%s %s", m.Name, dt.Model))
			},
		)
	}
	addToQueue(&c.deviceTypeQueue, id, cb)
	return ref
}

func (c *Converter) checkDeviceTypeRefExists(ref *openapi.ObjectReference) *modal.UUID {
	if ref == nil {
		return nil
	}
	addToErrQueue(
		&c.deviceTypeQueue, (*modal.UUID)(&ref.Id), func(id modal.UUID, dt *modal.DeviceType) error {
			return nil
		},
	)
	return (*modal.UUID)(&ref.Id)
}

func (c *Converter) DeviceTypeToOpenAPI(ctx context.Context, dt *modal.DeviceType) (
	dev *openapi.DeviceType, err error,
) {
	dev = &openapi.DeviceType{
		Id:           c.cloneUUID((*uuid.UUID)(dt.ID)),
		Manufacturer: c.getManufacturerRef(dt.Manufacturer),
		Model:        utils.ToPtr(dt.Model),
	}
	if err = c.resolveQueue(ctx); err != nil {
		return
	}
	return
}
