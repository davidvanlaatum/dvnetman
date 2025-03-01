package dto

import (
	"context"
	"dvnetman/pkg/auth"
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/openapi"
	"dvnetman/pkg/utils"
	"github.com/google/uuid"
)

func (c *Converter) UserToOpenAPI(ctx context.Context, mod *modal.User) (user *openapi.User, err error) {
	user = &openapi.User{
		Id:        *c.cloneUUID((*uuid.UUID)(mod.ID)),
		Email:     mod.Email,
		FirstName: mod.FirstName,
		LastName:  mod.LastName,
	}
	return
}

func (c *Converter) UpdateUserFromOpenAPI(ctx context.Context, body *openapi.User, mod *modal.User) error {
	mod.Email = body.Email
	mod.FirstName = body.FirstName
	mod.LastName = body.LastName
	mod.DisplayName = utils.JoinPtr(" ", body.FirstName, body.LastName)
	mod.ExternalID = body.ExternalID
	mod.ExternalProvider = body.ExternalProvider
	return nil
}

func (c *Converter) UserToOpenAPISearchResults(ctx context.Context, users []*modal.User) (
	[]*openapi.UserResult, error,
) {
	res := make([]*openapi.UserResult, len(users))
	for i, user := range users {
		res[i] = &openapi.UserResult{
			Id:        *(*uuid.UUID)(user.ID),
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
	}
	return res, nil
}

func (c *Converter) AuthUserToOpenAPI(ctx context.Context, u *auth.User) (user openapi.CurrentUser, err error) {
	user = openapi.CurrentUser{
		Email:            u.Email,
		DisplayName:      u.DisplayName,
		ExternalProvider: u.ExternalProvider,
		ExternalID:       u.ExternalID,
		LoggedIn:         utils.ToPtr(u.IsAuthenticated()),
		ProfileImageURL:  u.ProfileImageURL,
	}
	return
}
