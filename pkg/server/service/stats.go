package service

import (
	"context"
	"dvnetman/pkg/auth"
	"dvnetman/pkg/openapi"
	"dvnetman/pkg/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
)

func (s *Service) GetStats(ctx context.Context) (res *openapi.Response, err error) {
	if err = auth.RequirePerm(ctx, auth.PermissionRead); err != nil {
		return
	}
	stats := &openapi.Stats{}
	set := func(c **int) func(int64, error) error {
		return func(x int64, err error) error {
			if err == nil {
				*c = utils.ToPtr(int(x))
			}
			return err
		}
	}
	if err = set(&stats.DeviceCount)(s.db.CountDevices(ctx, bson.D{})); err != nil {
		return
	}
	if err = set(&stats.DeviceTypeCount)(s.db.CountDeviceTypes(ctx, bson.D{})); err != nil {
		return
	}
	if err = set(&stats.ManufacturerCount)(s.db.CountManufacturers(ctx, bson.D{})); err != nil {
		return
	}
	if err = set(&stats.UserCount)(s.db.CountUsers(ctx, bson.D{})); err != nil {
		return
	}
	res = &openapi.Response{Object: stats, Code: http.StatusOK}
	return
}
