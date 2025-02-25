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

func (s *Service) CreateManufacturer(ctx context.Context, opts *openapi.CreateManufacturerOpts) (
	res *openapi.Response, err error,
) {
	c := dto.NewConverter(s.db)
	mod := &modal.Manufacturer{}
	if err = c.UpdateManufacturerFromOpenAPI(ctx, opts.Body, mod); err != nil {
		return
	}
	if err = s.db.SaveManufacturer(ctx, mod); err != nil {
		return
	}
	res = &openapi.Response{}
	if res.Object, err = c.ManufacturerToOpenAPI(ctx, mod); err != nil {
		return
	}
	res.Code = http.StatusCreated
	return
}

func (s *Service) UpdateManufacturer(ctx context.Context, opts *openapi.UpdateManufacturerOpts) (
	res *openapi.Response, err error,
) {
	c := dto.NewConverter(s.db)
	var mod *modal.Manufacturer
	if mod, err = s.db.GetManufacturer(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	if mod.Version != opts.Body.Version {
		err = errors.WithStack(modal.OptimisticLockError)
		return
	}
	if err = c.UpdateManufacturerFromOpenAPI(ctx, opts.Body, mod); err != nil {
		return
	}
	if err = s.db.SaveManufacturer(ctx, mod); err != nil {
		return
	}
	res.Code = http.StatusOK
	return
}

func (s *Service) GetManufacturer(ctx context.Context, opts *openapi.GetManufacturerOpts) (
	res *openapi.Response, err error,
) {
	c := dto.NewConverter(s.db)
	var d *modal.Manufacturer
	if d, err = s.db.GetManufacturer(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	res = &openapi.Response{}
	if err = s.checkIfModified(opts.IfNoneMatch, opts.IfModifiedSince, d.Version, d.Updated, res); err != nil {
		return
	}
	if res.Object, err = c.ManufacturerToOpenAPI(ctx, d); err != nil {
		return
	}
	res.Code = http.StatusOK
	return
}

func (s *Service) ListManufacturers(ctx context.Context, opts *openapi.ListManufacturersOpts) (
	res *openapi.Response, err error,
) {
	c := dto.NewConverter(s.db)
	var page, size int64 = 0, 10
	if opts.PerPage != nil && *opts.PerPage > 0 {
		size = int64(*opts.PerPage)
	}
	var Manufacturers []*modal.Manufacturer
	if Manufacturers, err = s.db.ListManufacturers(
		ctx, bson.M{}, options.Find().SetLimit(size+1).SetSkip(page*size),
	); err != nil {
		return
	}
	rt := &openapi.ManufacturerSearchResults{}
	if len(Manufacturers) > int(size) {
		rt.Next = true
		Manufacturers = Manufacturers[:size]
	}
	rt.Count = len(Manufacturers)
	if rt.Items, err = c.ManufacturerToOpenAPISearchResults(ctx, Manufacturers); err != nil {
		return
	}
	res = &openapi.Response{
		Object: rt,
		Code:   http.StatusOK,
	}
	return
}

func (s *Service) DeleteManufacturer(ctx context.Context, opts *openapi.DeleteManufacturerOpts) (
	res *openapi.Response, err error,
) {
	var d *modal.Manufacturer
	if d, err = s.db.GetManufacturer(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	if err = s.db.DeleteManufacturer(ctx, d); err != nil {
		return
	}
	res.Code = http.StatusNoContent
	return
}
