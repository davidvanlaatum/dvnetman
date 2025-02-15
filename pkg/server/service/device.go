package service

import (
	"context"
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/openapi"
	"dvnetman/pkg/server/dto"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"net/http"
)

func (s *Service) CreateDevice(ctx context.Context, opts *openapi.CreateDeviceOpts) (res *openapi.Response, err error) {
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

func (s *Service) UpdateDevice(ctx context.Context, opts *openapi.UpdateDeviceOpts) (res *openapi.Response, err error) {
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
	if err = s.db.SaveDevice(ctx, mod); err != nil {
		return
	}
	res.Code = http.StatusAccepted
	return
}

func (s *Service) GetDevice(ctx context.Context, opts *openapi.GetDeviceOpts) (res *openapi.Response, err error) {
	c := dto.NewConverter(s.db)
	var d *modal.Device
	if d, err = s.db.GetDevice(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	res = &openapi.Response{}
	if err = s.checkIfModified(opts.IfNoneMatch, opts.IfModifiedSince, d.Version, d.Updated, res); err != nil {
		return
	}
	if res.Object, err = c.DeviceToOpenAPI(ctx, d); err != nil {
		return
	}
	res.Code = http.StatusOK
	return
}

func (s *Service) ListDevices(ctx context.Context, opts *openapi.ListDevicesOpts) (res *openapi.Response, err error) {
	c := dto.NewConverter(s.db)
	var page, size int64 = 0, 10
	if opts.PerPage != nil && *opts.PerPage > 0 {
		size = int64(*opts.PerPage)
	}
	var devices []*modal.Device
	if devices, err = s.db.ListDevices(
		ctx, bson.M{}, options.Find().SetLimit(size+1).SetSkip(page*size),
	); err != nil {
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

func (s *Service) DeleteDevice(ctx context.Context, opts *openapi.DeleteDeviceOpts) (res *openapi.Response, err error) {
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
