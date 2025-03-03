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

func (s *Service) CreateSite(ctx context.Context, opts *openapi.CreateSiteOpts) (
	res *openapi.Response, err error,
) {
	if err = auth.RequirePerm(ctx, auth.PermissionWrite); err != nil {
		return
	}
	c := dto.NewConverter(s.db)
	mod := &modal.Site{}
	if err = c.UpdateSiteFromOpenAPI(ctx, opts.Body, mod); err != nil {
		return
	}
	if err = s.db.SaveSite(ctx, mod); err != nil {
		return
	}
	res = &openapi.Response{}
	if res.Object, err = c.SiteToOpenAPI(ctx, mod); err != nil {
		return
	}
	res.Code = http.StatusCreated
	return
}

func (s *Service) UpdateSite(ctx context.Context, opts *openapi.UpdateSiteOpts) (
	res *openapi.Response, err error,
) {
	if err = auth.RequirePerm(ctx, auth.PermissionWrite); err != nil {
		return
	}
	c := dto.NewConverter(s.db)
	var mod *modal.Site
	if mod, err = s.db.GetSite(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	if mod.Version != opts.Body.Version {
		err = errors.WithStack(modal.OptimisticLockError)
		return
	}
	if err = c.UpdateSiteFromOpenAPI(ctx, opts.Body, mod); err != nil {
		return
	}
	res = &openapi.Response{}
	if err = s.db.SaveSite(ctx, mod); err != nil {
		return
	}
	res.Code = http.StatusOK
	return
}

func (s *Service) GetSite(ctx context.Context, opts *openapi.GetSiteOpts) (res *openapi.Response, err error) {
	if err = auth.RequirePerm(ctx, auth.PermissionRead); err != nil {
		return
	}
	c := dto.NewConverter(s.db)
	var d *modal.Site
	if d, err = s.db.GetSite(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	res = &openapi.Response{}
	if err = s.checkIfModified(opts.IfNoneMatch, opts.IfModifiedSince, d.Version, d.Updated, res); err != nil {
		return
	}
	if res.Object, err = c.SiteToOpenAPI(ctx, d); err != nil {
		return
	}
	res.Code = http.StatusOK
	return
}

func (s *Service) ListSites(ctx context.Context, opts *openapi.ListSitesOpts) (
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
	var sites []*modal.Site
	findOpts := options.Find().SetLimit(size + 1).SetSkip(page * size)
	if findOpts, err = s.setProjection(
		opts.Body.Fields, []string{"name", "site_type", "status"}, findOpts,
	); err != nil {
		return
	}
	if sites, err = s.db.ListSites(ctx, search, findOpts); err != nil {
		return
	}
	rt := &openapi.SiteSearchResults{}
	if len(sites) > int(size) {
		rt.Next = true
		sites = sites[:size]
	}
	rt.Count = len(sites)
	if rt.Items, err = c.SiteToOpenAPISearchResults(ctx, sites); err != nil {
		return
	}
	res = &openapi.Response{
		Object: rt,
		Code:   http.StatusOK,
	}
	return
}

func (s *Service) DeleteSite(ctx context.Context, opts *openapi.DeleteSiteOpts) (
	res *openapi.Response, err error,
) {
	if err = auth.RequirePerm(ctx, auth.PermissionWrite); err != nil {
		return
	}
	var d *modal.Site
	if d, err = s.db.GetSite(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	if err = s.db.DeleteSite(ctx, d); err != nil {
		return
	}
	res.Code = http.StatusNoContent
	return
}
