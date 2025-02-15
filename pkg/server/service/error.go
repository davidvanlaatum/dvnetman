package service

import (
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/openapi"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
)

var errorConverters []func(error) *openapi.Response

func RegisterErrorConverter(conv func(error) *openapi.Response) {
	errorConverters = append(errorConverters, conv)
}

func init() {
	RegisterErrorConverter(
		func(err error) *openapi.Response {
			if errors.Is(err, modal.OptimisticLockError) {
				return &openapi.Response{
					Code: http.StatusConflict,
					Object: openapi.APIErrorModal{
						Errors: []*openapi.ErrorMessage{
							{Code: "OPTIMISTIC_LOCK_ERROR", Message: "Optimistic lock error"},
						},
					},
				}
			}
			return nil
		},
	)
	//RegisterErrorConverter(
	//	func(err error) *openapi.Response {
	//		var parsingErr *openapi.ParsingError
	//		if ok := errors.As(err, &parsingErr); ok {
	//			return &openapi.Response{
	//				Code: http.StatusBadRequest,
	//				Object: openapi.APIErrorModal{
	//					Errors: []*openapi.ErrorMessage{
	//						{Code: "PARSING_ERROR", Message: err.Error()},
	//					},
	//				},
	//			}
	//		}
	//		return nil
	//	},
	//)
	//RegisterErrorConverter(
	//	func(err error) *openapi.Response {
	//		var requiredErr *openapi.RequiredError
	//		if ok := errors.As(err, &requiredErr); ok {
	//			return &openapi.Response{
	//				Code: http.StatusUnprocessableEntity,
	//				Object: openapi.APIErrorModal{
	//					Errors: []*openapi.ErrorMessage{
	//						{Code: "REQUIRED_FIELD", Message: err.Error()},
	//					},
	//				},
	//			}
	//		}
	//		return nil
	//	},
	//)
	RegisterErrorConverter(
		func(err error) *openapi.Response {
			if ok := errors.Is(err, mongo.ErrNoDocuments); ok {
				return &openapi.Response{
					Code: http.StatusNotFound,
					Object: openapi.APIErrorModal{
						Errors: []*openapi.ErrorMessage{
							{Code: "NOT_FOUND", Message: err.Error()},
						},
					},
				}
			}
			return nil
		},
	)
}
