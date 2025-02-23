package service

import (
	"context"
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/openapi"
	"dvnetman/pkg/server/dto"
	"dvnetman/pkg/utils"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"net/http"
)

func (s *Service) CreateDeviceType(ctx context.Context, opts *openapi.CreateDeviceTypeOpts) (
	res *openapi.Response, err error,
) {
	c := dto.NewConverter(s.db)
	mod := &modal.DeviceType{}
	if err = c.UpdateDeviceTypeFromOpenAPI(ctx, opts.Body, mod); err != nil {
		return
	}
	if err = s.db.SaveDeviceType(ctx, mod); err != nil {
		return
	}
	res = &openapi.Response{}
	if res.Object, err = c.DeviceTypeToOpenAPI(ctx, mod); err != nil {
		return
	}
	res.Code = http.StatusCreated
	return
}

func (s *Service) UpdateDeviceType(ctx context.Context, opts *openapi.UpdateDeviceTypeOpts) (
	res *openapi.Response, err error,
) {
	c := dto.NewConverter(s.db)
	var mod *modal.DeviceType
	if mod, err = s.db.GetDeviceType(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	if mod.Version != opts.Body.Version {
		err = errors.WithStack(modal.OptimisticLockError)
		return
	}
	if err = c.UpdateDeviceTypeFromOpenAPI(ctx, opts.Body, mod); err != nil {
		return
	}
	if err = s.db.SaveDeviceType(ctx, mod); err != nil {
		return
	}
	res.Code = http.StatusAccepted
	return
}

func (s *Service) GetDeviceType(ctx context.Context, opts *openapi.GetDeviceTypeOpts) (
	res *openapi.Response, err error,
) {
	c := dto.NewConverter(s.db)
	var d *modal.DeviceType
	if d, err = s.db.GetDeviceType(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	res = &openapi.Response{}
	if err = s.checkIfModified(opts.IfNoneMatch, opts.IfModifiedSince, d.Version, d.Updated, res); err != nil {
		return
	}
	if res.Object, err = c.DeviceTypeToOpenAPI(ctx, d); err != nil {
		return
	}
	res.Code = http.StatusOK
	return
}

func (s *Service) ListDeviceTypes(ctx context.Context, opts *openapi.ListDeviceTypesOpts) (
	res *openapi.Response, err error,
) {
	c := dto.NewConverter(s.db)
	var page, size int64 = 0, 10
	if opts.PerPage != nil && *opts.PerPage > 0 {
		size = int64(*opts.PerPage)
	}
	search := filter{}
	search.equalsStr("model", opts.Model)
	search.regex("model", opts.ModelRegex, "i")
	search.inUUID("id", utils.MapTo(opts.Ids, modal.ConvertUUID))
	findOpts := options.Find().SetLimit(size + 1).SetSkip(page * size)
	if findOpts, err = s.setProjection(opts.Fields, []string{"model", "manufacturer"}, findOpts); err != nil {
		return
	}
	var DeviceTypes []*modal.DeviceType
	if DeviceTypes, err = s.db.ListDeviceTypes(ctx, search, findOpts); err != nil {
		return
	}
	rt := &openapi.DeviceTypeSearchResults{}
	if len(DeviceTypes) > int(size) {
		rt.Next = true
		DeviceTypes = DeviceTypes[:size]
	}
	rt.Count = len(DeviceTypes)
	if rt.Items, err = c.DeviceTypeToOpenAPISearchResults(ctx, DeviceTypes); err != nil {
		return
	}
	res = &openapi.Response{
		Object: rt,
		Code:   http.StatusOK,
	}
	return
}

func (s *Service) DeleteDeviceType(ctx context.Context, opts *openapi.DeleteDeviceTypeOpts) (
	res *openapi.Response, err error,
) {
	var d *modal.DeviceType
	if d, err = s.db.GetDeviceType(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	if err = s.db.DeleteDeviceType(ctx, d); err != nil {
		return
	}
	res.Code = http.StatusNoContent
	return
}
