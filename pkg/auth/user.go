package auth

import (
	"github.com/google/uuid"
	"github.com/markbates/goth"
)

var (
	anonymousUserId = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	systemUserId    = uuid.MustParse("00000000-0000-0000-0000-000000000001")
)

type User struct {
	ID               uuid.UUID
	Email            *string
	DisplayName      *string
	ExternalProvider *string
	ExternalID       *string
	ProfileImageURL  *string
	OAuthUser        *goth.User
}

func (u *User) IsAnonymous() bool {
	return u.ID == anonymousUserId
}

func (u *User) IsSystem() bool {
	return u.ID == systemUserId
}

func (u *User) IsAuthenticated() bool {
	return !u.IsAnonymous() && !u.IsSystem()
}
