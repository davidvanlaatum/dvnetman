package dto

import (
	"context"
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/openapi"
	"dvnetman/pkg/utils"
	"github.com/google/uuid"
)

func (c *Converter) SiteToOpenAPI(ctx context.Context, site *modal.Site) (dev *openapi.Site, err error) {
	dev = &openapi.Site{
		Id:          *c.cloneUUID((*uuid.UUID)(site.ID)),
		Name:        site.Name,
		Description: site.Description,
		Version:     site.Version,
	}
	if err = c.resolveQueue(ctx); err != nil {
		return
	}
	return
}

func (c *Converter) SiteToOpenAPISearchResults(ctx context.Context, sites []*modal.Site) (
	res []*openapi.SiteResult, err error,
) {
	res = utils.MapTo(
		sites, func(site *modal.Site) *openapi.SiteResult {
			dev := &openapi.SiteResult{
				Id:          *(*uuid.UUID)(site.ID),
				Name:        site.Name,
				Description: site.Description,
				Version:     site.Version,
			}
			return dev
		},
	)
	if err = c.resolveQueue(ctx); err != nil {
		return
	}
	return
}

func (c *Converter) UpdateSiteFromOpenAPI(ctx context.Context, site *openapi.Site, mod *modal.Site) error {
	mod.Name = site.Name
	mod.Description = site.Description
	return c.resolveQueue(ctx)
}

func (c *Converter) checkSiteRefExists(ref *openapi.ObjectReference) *modal.UUID {
	if ref == nil || ref.Id == uuid.Nil {
		return nil
	}
	addToErrQueue(
		&c.siteQueue, (*modal.UUID)(&ref.Id), func(id modal.UUID, dt *modal.Site) error {
			return nil
		},
	)
	return (*modal.UUID)(&ref.Id)
}

func (c *Converter) getSiteRef(id *modal.UUID) *openapi.ObjectReference {
	if id == nil {
		return nil
	}
	ref := &openapi.ObjectReference{Id: *(*uuid.UUID)(id)}
	addToQueue(
		&c.siteQueue, id, func(site *modal.Site) {
			ref.DisplayName = utils.ToPtr(site.Name)
		},
	)
	return ref
}
