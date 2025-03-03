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

func (s *Service) CreateLocation(ctx context.Context, opts *openapi.CreateLocationOpts) (
	res *openapi.Response, err error,
) {
	if err = auth.RequirePerm(ctx, auth.PermissionWrite); err != nil {
		return
	}
	c := dto.NewConverter(s.db)
	mod := &modal.Location{}
	if err = c.UpdateLocationFromOpenAPI(ctx, opts.Body, mod); err != nil {
		return
	}
	if err = s.db.SaveLocation(ctx, mod); err != nil {
		return
	}
	res = &openapi.Response{}
	if res.Object, err = c.LocationToOpenAPI(ctx, mod); err != nil {
		return
	}
	res.Code = http.StatusCreated
	return
}

func (s *Service) UpdateLocation(ctx context.Context, opts *openapi.UpdateLocationOpts) (
	res *openapi.Response, err error,
) {
	if err = auth.RequirePerm(ctx, auth.PermissionWrite); err != nil {
		return
	}
	c := dto.NewConverter(s.db)
	var mod *modal.Location
	if mod, err = s.db.GetLocation(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	if mod.Version != opts.Body.Version {
		err = errors.WithStack(modal.OptimisticLockError)
		return
	}
	if err = c.UpdateLocationFromOpenAPI(ctx, opts.Body, mod); err != nil {
		return
	}
	res = &openapi.Response{}
	if err = s.db.SaveLocation(ctx, mod); err != nil {
		return
	}
	res.Code = http.StatusOK
	return
}

func (s *Service) GetLocation(ctx context.Context, opts *openapi.GetLocationOpts) (res *openapi.Response, err error) {
	if err = auth.RequirePerm(ctx, auth.PermissionRead); err != nil {
		return
	}
	c := dto.NewConverter(s.db)
	var d *modal.Location
	if d, err = s.db.GetLocation(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	res = &openapi.Response{}
	if err = s.checkIfModified(opts.IfNoneMatch, opts.IfModifiedSince, d.Version, d.Updated, res); err != nil {
		return
	}
	if res.Object, err = c.LocationToOpenAPI(ctx, d); err != nil {
		return
	}
	res.Code = http.StatusOK
	return
}

func (s *Service) ListLocations(ctx context.Context, opts *openapi.ListLocationsOpts) (
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
	search.equalsStr("name", opts.Body.Name)
	search.regex("name", opts.Body.NameRegex, "i")
	search.equalsUUID("parent", opts.Body.Parent)
	search.equalsUUID("site", opts.Body.Site)
	var locations []*modal.Location
	findOpts := options.Find().SetLimit(size + 1).SetSkip(page * size)
	if findOpts, err = s.setProjection(
		opts.Body.Fields, []string{"name", "location_type", "status"}, findOpts,
	); err != nil {
		return
	}
	if locations, err = s.db.ListLocations(ctx, search, findOpts); err != nil {
		return
	}
	rt := &openapi.LocationSearchResults{}
	if len(locations) > int(size) {
		rt.Next = true
		locations = locations[:size]
	}
	rt.Count = len(locations)
	if rt.Items, err = c.LocationToOpenAPISearchResults(ctx, locations); err != nil {
		return
	}
	res = &openapi.Response{
		Object: rt,
		Code:   http.StatusOK,
	}
	return
}

func (s *Service) DeleteLocation(ctx context.Context, opts *openapi.DeleteLocationOpts) (
	res *openapi.Response, err error,
) {
	if err = auth.RequirePerm(ctx, auth.PermissionWrite); err != nil {
		return
	}
	var d *modal.Location
	if d, err = s.db.GetLocation(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	if err = s.db.DeleteLocation(ctx, d); err != nil {
		return
	}
	res.Code = http.StatusNoContent
	return
}
