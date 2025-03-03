package dto

import (
	"context"
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/openapi"
	"dvnetman/pkg/utils"
	"github.com/google/uuid"
	"strings"
)

func (c *Converter) LocationToOpenAPI(ctx context.Context, location *modal.Location) (
	dev *openapi.Location, err error,
) {
	dev = &openapi.Location{
		Id:          *c.cloneUUID((*uuid.UUID)(location.ID)),
		Name:        location.Name,
		Description: location.Description,
		Parent:      c.getLocationRef(location.Parent),
		Site:        c.getSiteRef(location.Site),
		Version:     location.Version,
	}
	if err = c.resolveQueue(ctx); err != nil {
		return
	}
	return
}

func (c *Converter) LocationToOpenAPISearchResults(ctx context.Context, locations []*modal.Location) (
	res []*openapi.LocationResult, err error,
) {
	res = utils.MapTo(
		locations, func(location *modal.Location) *openapi.LocationResult {
			dev := &openapi.LocationResult{
				Id:      *(*uuid.UUID)(location.ID),
				Name:    location.Name,
				Parent:  c.getLocationRef(location.Parent),
				Site:    c.getSiteRef(location.Site),
				Version: location.Version,
			}
			return dev
		},
	)
	if err = c.resolveQueue(ctx); err != nil {
		return
	}
	return
}

func (c *Converter) UpdateLocationFromOpenAPI(
	ctx context.Context, location *openapi.Location, mod *modal.Location,
) error {
	mod.Name = location.Name
	mod.Description = location.Description
	mod.Parent = c.checkLocationRefExists(location.Parent)
	mod.Site = c.checkSiteRefExists(location.Site)
	return c.resolveQueue(ctx)
}

func (c *Converter) checkLocationRefExists(ref *openapi.ObjectReference) *modal.UUID {
	if ref == nil || ref.Id == uuid.Nil {
		return nil
	}
	addToErrQueue(
		&c.locationQueue, (*modal.UUID)(&ref.Id), func(id modal.UUID, dt *modal.Location) error {
			return nil
		},
	)
	return (*modal.UUID)(&ref.Id)
}

func (c *Converter) getLocationRefNested(id *modal.UUID, path []string, cb func(path []string)) {
	addToQueue(
		&c.locationQueue, id, func(loc *modal.Location) {
			if loc.Parent == nil {
				cb(append(path, loc.Name))
			} else {
				c.getLocationRefNested(loc.Parent, append(path, loc.Name), cb)
			}
		},
	)
}

func (c *Converter) getLocationRef(id *modal.UUID) *openapi.ObjectReference {
	if id == nil {
		return nil
	}
	ref := &openapi.ObjectReference{Id: *(*uuid.UUID)(id)}
	c.getLocationRefNested(
		id, nil, func(path []string) {
			ref.DisplayName = utils.ToPtr(strings.Join(path, "/"))
		},
	)
	return ref
}
