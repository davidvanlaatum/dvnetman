package service

import (
	"context"
	"dvnetman/pkg/auth"
	"dvnetman/pkg/cache"
	"dvnetman/pkg/logger"
	"dvnetman/pkg/openapi"
	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
)

func (s *Service) GetStats(ctx context.Context) (res *openapi.Response, err error) {
	if err = auth.RequirePerm(ctx, auth.PermissionRead); err != nil {
		return
	}
	var stats *openapi.Stats
	stats, err = cache.Lazy(
		ctx, s.cache, "stats", func(ctx context.Context) (stats *openapi.Stats, err error) {
			stats = &openapi.Stats{}
			set := func(c *int) func(int64, error) error {
				return func(x int64, err error) error {
					if err == nil {
						*c = int(x)
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
			return
		},
	)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("failed to get stats")
		return
	}
	res = &openapi.Response{Object: stats, Code: http.StatusOK}
	return
}
