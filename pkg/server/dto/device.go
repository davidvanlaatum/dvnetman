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
		Location:     c.getLocationRef(device.Location),
		Site:         c.getSiteRef(device.Site),
		Created:      utils.ToPtr(device.Created),
		Updated:      utils.ToPtr(device.Updated),
		Serial:       device.Serial,
		AssetTag:     device.AssetTag,
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
				Id:          *(*uuid.UUID)(device.ID),
				Name:        device.Name,
				DeviceType:  c.getDeviceTypeRef(device.DeviceType),
				Location:    c.getLocationRef(device.Location),
				Site:        c.getSiteRef(device.Site),
				Version:     device.Version,
				Status:      (*string)(device.Status),
				Created:     utils.ToPtr(device.Created),
				Updated:     utils.ToPtr(device.Updated),
				Description: device.Description,
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
	mod.Site = c.checkSiteRefExists(device.Site)
	mod.Location = c.checkLocationRefExists(device.Location)
	return c.resolveQueue(ctx)
}
