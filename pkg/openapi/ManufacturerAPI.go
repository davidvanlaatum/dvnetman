// Code generated by dvnetman. DO NOT EDIT.

package openapi

import (
	"context"
	utils "dvnetman/pkg/utils"
	"encoding/json"
	uuid "github.com/google/uuid"
	mux "github.com/gorilla/mux"
	errors "github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

type ManufacturerAPI interface {
	CreateManufacturer(ctx context.Context, opts *CreateManufacturerOpts) (res *Response, err error)
	DeleteManufacturer(ctx context.Context, opts *DeleteManufacturerOpts) (res *Response, err error)
	GetManufacturer(ctx context.Context, opts *GetManufacturerOpts) (res *Response, err error)
	ListManufacturers(ctx context.Context, opts *ListManufacturersOpts) (res *Response, err error)
	UpdateManufacturer(ctx context.Context, opts *UpdateManufacturerOpts) (res *Response, err error)
}
type CreateManufacturerOpts struct {
	Body *Manufacturer
}
type ListManufacturersOpts struct {
	Page    *int
	PerPage *int
	Sort    *string
	Body    *ManufacturerSearchBody
}
type DeleteManufacturerOpts struct {
	Id uuid.UUID
}
type GetManufacturerOpts struct {
	Id              uuid.UUID
	IfNoneMatch     *string
	IfModifiedSince *time.Time
}
type UpdateManufacturerOpts struct {
	Id   uuid.UUID
	Body *Manufacturer
}
type ManufacturerAPIHandler struct {
	service ManufacturerAPI
	errors  ErrorHandler
}

func (h *ManufacturerAPIHandler) CreateManufacturer(w http.ResponseWriter, r *http.Request) {
	var res *Response
	var err error
	opts := &CreateManufacturerOpts{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&opts.Body); err != nil {
		h.errors.ErrorHandler(w, r, errors.WithStack(NewBodyParamError(err)))
		return
	}
	if decoder.More() {
		h.errors.ErrorHandler(w, r, errors.WithStack(NewBodyParamError(errors.New("unexpected data after body"))))
		return
	}
	if res, err = h.service.CreateManufacturer(r.Context(), opts); err != nil {
		h.errors.ErrorHandler(w, r, err)
	} else if res == nil {
		h.errors.ErrorHandler(w, r, errors.Errorf("no response returned"))
	} else if err = res.Write(r, w); err != nil {
		h.errors.WriteErrorHandler(w, r, err)
	}
}
func (h *ManufacturerAPIHandler) ListManufacturers(w http.ResponseWriter, r *http.Request) {
	var res *Response
	var err error
	opts := &ListManufacturersOpts{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&opts.Body); err != nil {
		h.errors.ErrorHandler(w, r, errors.WithStack(NewBodyParamError(err)))
		return
	}
	if decoder.More() {
		h.errors.ErrorHandler(w, r, errors.WithStack(NewBodyParamError(errors.New("unexpected data after body"))))
		return
	}
	for k, v := range r.URL.Query() {
		switch k {
		case "page":
			var x int
			if x, err = strconv.Atoi(v[0]); err != nil {
				h.errors.ErrorHandler(w, r, errors.WithStack(NewQueryParamError("page", err)))
				return
			}
			opts.Page = utils.ToPtr(x)
		case "per_page":
			var x int
			if x, err = strconv.Atoi(v[0]); err != nil {
				h.errors.ErrorHandler(w, r, errors.WithStack(NewQueryParamError("per_page", err)))
				return
			}
			opts.PerPage = utils.ToPtr(x)
		case "sort":
			opts.Sort = utils.ToPtr(v[0])
		}
	}
	if res, err = h.service.ListManufacturers(r.Context(), opts); err != nil {
		h.errors.ErrorHandler(w, r, err)
	} else if res == nil {
		h.errors.ErrorHandler(w, r, errors.Errorf("no response returned"))
	} else if err = res.Write(r, w); err != nil {
		h.errors.WriteErrorHandler(w, r, err)
	}
}
func (h *ManufacturerAPIHandler) DeleteManufacturer(w http.ResponseWriter, r *http.Request) {
	var res *Response
	var err error
	opts := &DeleteManufacturerOpts{}
	vars := mux.Vars(r)
	if opts.Id, err = uuid.Parse(vars["id"]); err != nil {
		h.errors.ErrorHandler(w, r, errors.WithStack(NewPathParamError("Id", err)))
		return
	}
	if res, err = h.service.DeleteManufacturer(r.Context(), opts); err != nil {
		h.errors.ErrorHandler(w, r, err)
	} else if res == nil {
		h.errors.ErrorHandler(w, r, errors.Errorf("no response returned"))
	} else if err = res.Write(r, w); err != nil {
		h.errors.WriteErrorHandler(w, r, err)
	}
}
func (h *ManufacturerAPIHandler) GetManufacturer(w http.ResponseWriter, r *http.Request) {
	var res *Response
	var err error
	opts := &GetManufacturerOpts{}
	vars := mux.Vars(r)
	if opts.Id, err = uuid.Parse(vars["id"]); err != nil {
		h.errors.ErrorHandler(w, r, errors.WithStack(NewPathParamError("Id", err)))
		return
	}
	for k, v := range r.Header {
		switch k {
		case "If-None-Match":
			opts.IfNoneMatch = utils.ToPtr(v[0])
		case "If-Modified-Since":
			var t time.Time
			if t, err = time.Parse(time.RFC1123, v[0]); err != nil {
				h.errors.ErrorHandler(w, r, errors.WithStack(NewQueryParamError("If-Modified-Since", err)))
				return
			}
			opts.IfModifiedSince = utils.ToPtr(t)
		}
	}
	if res, err = h.service.GetManufacturer(r.Context(), opts); err != nil {
		h.errors.ErrorHandler(w, r, err)
	} else if res == nil {
		h.errors.ErrorHandler(w, r, errors.Errorf("no response returned"))
	} else if err = res.Write(r, w); err != nil {
		h.errors.WriteErrorHandler(w, r, err)
	}
}
func (h *ManufacturerAPIHandler) UpdateManufacturer(w http.ResponseWriter, r *http.Request) {
	var res *Response
	var err error
	opts := &UpdateManufacturerOpts{}
	vars := mux.Vars(r)
	if opts.Id, err = uuid.Parse(vars["id"]); err != nil {
		h.errors.ErrorHandler(w, r, errors.WithStack(NewPathParamError("Id", err)))
		return
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&opts.Body); err != nil {
		h.errors.ErrorHandler(w, r, errors.WithStack(NewBodyParamError(err)))
		return
	}
	if decoder.More() {
		h.errors.ErrorHandler(w, r, errors.WithStack(NewBodyParamError(errors.New("unexpected data after body"))))
		return
	}
	if res, err = h.service.UpdateManufacturer(r.Context(), opts); err != nil {
		h.errors.ErrorHandler(w, r, err)
	} else if res == nil {
		h.errors.ErrorHandler(w, r, errors.Errorf("no response returned"))
	} else if err = res.Write(r, w); err != nil {
		h.errors.WriteErrorHandler(w, r, err)
	}
}
func AttachManufacturerAPI(service ManufacturerAPI, errors ErrorHandler, router *mux.Router) {
	handler := &ManufacturerAPIHandler{
		errors:  errors,
		service: service,
	}
	router.Methods(http.MethodPost).Path("/api/v1/manufacturer").Name("CreateManufacturer").HandlerFunc(handler.CreateManufacturer)
	router.Methods(http.MethodPost).Path("/api/v1/manufacturer/search").Name("ListManufacturers").HandlerFunc(handler.ListManufacturers)
	router.Methods(http.MethodDelete).Path("/api/v1/manufacturer/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}").Name("DeleteManufacturer").HandlerFunc(handler.DeleteManufacturer)
	router.Methods(http.MethodGet).Path("/api/v1/manufacturer/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}").Name("GetManufacturer").HandlerFunc(handler.GetManufacturer)
	router.Methods(http.MethodPut).Path("/api/v1/manufacturer/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}").Name("UpdateManufacturer").HandlerFunc(handler.UpdateManufacturer)
}
