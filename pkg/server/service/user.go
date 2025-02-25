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

func (s *Service) CreateUser(ctx context.Context, opts *openapi.CreateUserOpts) (res *openapi.Response, err error) {
	c := dto.NewConverter(s.db)
	mod := &modal.User{}
	if err = c.UpdateUserFromOpenAPI(ctx, opts.Body, mod); err != nil {
		return
	}
	if err = s.db.SaveUser(ctx, mod); err != nil {
		return
	}
	res = &openapi.Response{}
	if res.Object, err = c.UserToOpenAPI(ctx, mod); err != nil {
		return
	}
	res.Code = http.StatusCreated
	return
}

func (s *Service) UpdateUser(ctx context.Context, opts *openapi.UpdateUserOpts) (res *openapi.Response, err error) {
	c := dto.NewConverter(s.db)
	var mod *modal.User
	if mod, err = s.db.GetUser(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	if mod.Version != opts.Body.Version {
		err = errors.WithStack(modal.OptimisticLockError)
		return
	}
	if err = c.UpdateUserFromOpenAPI(ctx, opts.Body, mod); err != nil {
		return
	}
	if err = s.db.SaveUser(ctx, mod); err != nil {
		return
	}
	res.Code = http.StatusOK
	return
}

func (s *Service) DeleteUser(ctx context.Context, opts *openapi.DeleteUserOpts) (res *openapi.Response, err error) {
	var d *modal.User
	if d, err = s.db.GetUser(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	if err = s.db.DeleteUser(ctx, d); err != nil {
		return
	}
	res.Code = http.StatusNoContent
	return
}

func (s *Service) GetUser(ctx context.Context, opts *openapi.GetUserOpts) (res *openapi.Response, err error) {
	c := dto.NewConverter(s.db)
	var d *modal.User
	if d, err = s.db.GetUser(ctx, (*modal.UUID)(&opts.Id)); err != nil {
		return
	}
	res = &openapi.Response{}
	if err = s.checkIfModified(opts.IfNoneMatch, opts.IfModifiedSince, d.Version, d.Updated, res); err != nil {
		return
	}
	if res.Object, err = c.UserToOpenAPI(ctx, d); err != nil {
		return
	}
	res.Code = http.StatusOK
	return
}

func (s *Service) GetCurrentUser(ctx context.Context) (res *openapi.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ListUsers(ctx context.Context, opts *openapi.ListUsersOpts) (res *openapi.Response, err error) {
	c := dto.NewConverter(s.db)
	var page, size int64 = 0, 10
	if opts.PerPage != nil && *opts.PerPage > 0 {
		size = int64(*opts.PerPage)
	}
	search := filter{}
	search.inUUID("id", utils.MapTo(opts.Body.Ids, modal.ConvertUUID))
	search.equalsStr("email", opts.Body.Email)
	search.equalsStr("first_name", opts.Body.FirstName)
	search.regex("first_name_regex", opts.Body.FirstNameRegex, "i")
	search.equalsStr("last_name", opts.Body.LastName)
	search.regex("last_name_regex", opts.Body.LastNameRegex, "i")
	search.equalsStr("full_name", opts.Body.FullName)
	search.regex("full_name_regex", opts.Body.FullNameRegex, "i")
	search.equalsStr("username", opts.Body.Username)
	search.regex("username_regex", opts.Body.UsernameRegex, "i")
	findOpts := options.Find().SetLimit(size + 1).SetSkip(page * size)
	if findOpts, err = s.setProjection(opts.Body.Fields, []string{"model", "manufacturer"}, findOpts); err != nil {
		return
	}
	var Users []*modal.User
	if Users, err = s.db.ListUsers(ctx, search, findOpts); err != nil {
		return
	}
	rt := &openapi.UserSearchResults{}
	if len(Users) > int(size) {
		rt.Next = true
		Users = Users[:size]
	}
	rt.Count = len(Users)
	if rt.Items, err = c.UserToOpenAPISearchResults(ctx, Users); err != nil {
		return
	}
	res = &openapi.Response{
		Object: rt,
		Code:   http.StatusOK,
	}
	return
}
