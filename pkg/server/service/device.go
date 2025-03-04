package service

import (
	"context"
	"dvnetman/pkg/auth"
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/openapi"
	"dvnetman/pkg/server/dto"
	"dvnetman/pkg/utils"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"net/http"
)

type DeviceService struct {
	db *modal.DBClient
}

func NewDeviceService(db *modal.DBClient) *DeviceService {
	return &DeviceService{db: db}
}

func (s *DeviceService) CreateDevice(ctx context.Context, opts *openapi.CreateDeviceOpts) (
	res *openapi.Response, err error,
) {
	if err = auth.RequirePerm(ctx, auth.PermissionWrite); err != nil {
		return
	}
	c := dto.NewConverter(s.db)
	mod := &modal.Device{}
	if err = c.UpdateDeviceFromOpenAPI(ctx, opts.Body, mod); err != nil {
		return
	}
	if err = s.db.SaveDevice(ctx, mod); err != nil {
		return
	}
	res = &openapi.Response{}
	if res.Object, err = c.DeviceToOpenAPI(ctx, mod); err != nil {
		return
	}
	res.Code = http.StatusCreated
	return
}

func (s *DeviceService) UpdateDevice(ctx context.Context, opts *openapi.UpdateDeviceOpts) (
	res *openapi.Response, err error,
) {
	if err = auth.RequirePerm(ctx, auth.PermissionWrite); err != nil {
		return
	}
	c := dto.NewConverter(s.db)
	var mod *modal.Device
	if mod, err = s.db.GetDevice(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	if mod.Version != opts.Body.Version {
		err = errors.WithStack(modal.OptimisticLockError)
		return
	}
	if err = c.UpdateDeviceFromOpenAPI(ctx, opts.Body, mod); err != nil {
		return
	}
	res = &openapi.Response{}
	if err = s.db.SaveDevice(ctx, mod); err != nil {
		return
	}
	res.Code = http.StatusOK
	return
}

func (s *DeviceService) GetDevice(ctx context.Context, opts *openapi.GetDeviceOpts) (res *openapi.Response, err error) {
	if err = auth.RequirePerm(ctx, auth.PermissionRead); err != nil {
		return
	}
	c := dto.NewConverter(s.db)
	var d *modal.Device
	if d, err = s.db.GetDevice(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	res = &openapi.Response{}
	if err = checkIfModified(opts.IfNoneMatch, opts.IfModifiedSince, d.Version, d.Updated, res); err != nil {
		return
	}
	if res.Object, err = c.DeviceToOpenAPI(ctx, d); err != nil {
		return
	}
	res.Code = http.StatusOK
	return
}

func (s *DeviceService) ListDevices(ctx context.Context, opts *openapi.ListDevicesOpts) (
	res *openapi.Response, err error,
) {
	if err = auth.RequirePerm(ctx, auth.PermissionRead); err != nil {
		return
	}
	c := dto.NewConverter(s.db)
	var page, size int64 = 0, 10
	if opts.PerPage != nil && *opts.PerPage > 0 {
		size = int64(*opts.PerPage)
	}
	search := filter{}
	search.inUUID("id", utils.MapTo(opts.Body.Ids, modal.ConvertUUID))
	search.inUUID("device_type", utils.MapTo(opts.Body.DeviceType, modal.ConvertUUID))
	search.equalsStr("name", opts.Body.Name)
	search.regex("name", opts.Body.NameRegex, "i")
	search.equalsStr("status", opts.Body.Status)
	search.equalsStr("serial", opts.Body.Serial)
	search.regex("serial", opts.Body.SerialRegex, "i")
	search.equalsStr("asset_tag", opts.Body.AssetTag)
	search.regex("asset_tag", opts.Body.AssetTagRegex, "i")
	var devices []*modal.Device
	findOpts := options.Find().SetLimit(size + 1).SetSkip(page * size)
	if findOpts, err = setProjection(
		opts.Body.Fields, []string{"name", "device_type", "status"}, findOpts,
	); err != nil {
		return
	}
	if devices, err = s.db.ListDevices(ctx, search, findOpts); err != nil {
		return
	}
	rt := &openapi.DeviceSearchResults{}
	if len(devices) > int(size) {
		rt.Next = true
		devices = devices[:size]
	}
	rt.Count = len(devices)
	if rt.Items, err = c.DeviceToOpenAPISearchResults(ctx, devices); err != nil {
		return
	}
	res = &openapi.Response{
		Object: rt,
		Code:   http.StatusOK,
	}
	return
}

func (s *DeviceService) DeleteDevice(ctx context.Context, opts *openapi.DeleteDeviceOpts) (
	res *openapi.Response, err error,
) {
	if err = auth.RequirePerm(ctx, auth.PermissionWrite); err != nil {
		return
	}
	var d *modal.Device
	if d, err = s.db.GetDevice(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	if err = s.db.DeleteDevice(ctx, d); err != nil {
		return
	}
	res.Code = http.StatusNoContent
	return
}
