package auth

import (
	"dvnetman/pkg/openapi"
	"github.com/pkg/errors"
	"net/http"
)

type NotLoggedInError struct {
}

func (e *NotLoggedInError) Error() string {
	return "Not logged in"
}

func init() {
	openapi.RegisterErrorConverter(
		func(err error) *openapi.Response {
			var notLoggedIn *NotLoggedInError
			if ok := errors.As(err, &notLoggedIn); ok {
				return &openapi.Response{
					Code: http.StatusUnauthorized,
					Object: openapi.APIErrorModal{
						Errors: []*openapi.ErrorMessage{
							{Code: "NOT_LOGGED_IN", Message: err.Error()},
						},
					},
				}
			}
			return nil
		},
	)
}

type NotAuthorizedError struct {
}

func (e *NotAuthorizedError) Error() string {
	return "Not authorized"
}

func init() {
	openapi.RegisterErrorConverter(
		func(err error) *openapi.Response {
			var notAuthorizedError *NotAuthorizedError
			if ok := errors.As(err, &notAuthorizedError); ok {
				return &openapi.Response{
					Code: http.StatusForbidden,
					Object: openapi.APIErrorModal{
						Errors: []*openapi.ErrorMessage{
							{Code: "NOT_AUTHORIZED", Message: err.Error()},
						},
					},
				}
			}
			return nil
		},
	)
}
