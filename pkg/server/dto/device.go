package dto

import (
	"context"
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/openapi"
	"dvnetman/pkg/utils"
	"github.com/google/uuid"
)

func (c *Converter) DeviceToOpenAPI(ctx context.Context, device *modal.Device) (dev *openapi.Device, err error) {
	dev = &openapi.Device{
		Id:           *c.cloneUUID((*uuid.UUID)(device.ID)),
		Name:         device.Name,
		Description:  device.Description,
		RackPosition: utils.ConvertPtr(device.RackPosition, func(x int) float64 { return float64(x) / 2 }),
		RackFace:     (*openapi.DeviceRackFace)(device.RackFace),
		Version:      device.Version,
		Status:       (*string)(device.Status),
		DeviceType:   c.getDeviceTypeRef(device.DeviceType),
	}
	if err = c.resolveQueue(ctx); err != nil {
		return
	}
	return
}

func (c *Converter) DeviceToOpenAPISearchResults(ctx context.Context, devices []*modal.Device) (
	res []*openapi.DeviceResult, err error,
) {
	res = utils.MapTo(
		devices, func(device *modal.Device) *openapi.DeviceResult {
			dev := &openapi.DeviceResult{
				Id:         *(*uuid.UUID)(device.ID),
				Name:       device.Name,
				DeviceType: c.getDeviceTypeRef(device.DeviceType),
				Version:    device.Version,
			}
			return dev
		},
	)
	if err = c.resolveQueue(ctx); err != nil {
		return
	}
	return
}

func (c *Converter) UpdateDeviceFromOpenAPI(ctx context.Context, device *openapi.Device, mod *modal.Device) error {
	mod.Name = device.Name
	mod.Description = device.Description
	mod.RackPosition = utils.ConvertPtr(device.RackPosition, func(x float64) int { return int(x * 2) })
	mod.DeviceType = c.checkDeviceTypeRefExists(device.DeviceType)
	return c.resolveQueue(ctx)
}
