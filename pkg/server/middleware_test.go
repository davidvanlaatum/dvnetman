package server

import (
	"context"
	"dvnetman/pkg/logger"
	"dvnetman/pkg/testutils"
	"dvnetman/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogContextMiddleware(t *testing.T) {
	r := require.New(t)
	ctx := testutils.GetTestContext(t)
	router := mux.NewRouter()
	router.Use(logContextMiddleware)
	router.Path("/test").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			l := logger.Ctx(r.Context())
			l.Info(r.Context()).Msg("test")
		},
	).Name("test-route")
	logs := logger.NewCollector()
	ctx = logger.NewLogger(logger.LevelTrace, logs).Context(ctx)
	srv := httptest.NewUnstartedServer(router)
	srv.Config.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}
	srv.Start()
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/test")
	r.NoError(err)
	r.Equal(http.StatusOK, resp.StatusCode)
	r.Len(logs.Logs(), 1)
	l := logs.Logs()[0]
	r.NotEmpty(l.Keys["remote"])
	r.EqualExportedValues(
		[]logger.EventData{
			{
				Time:  l.Time,
				Level: logger.LevelInfo,
				Keys: map[string]interface{}{
					"endpoint": "/test",
					"remote":   l.Keys["remote"],
					"method":   "GET",
					"url":      "/test",
					"route":    "test-route",
				},
				Message: "test",
				File:    l.File,
			},
		}, logs.Logs(),
	)
}

func TestLogContextMiddlewareNoRoute(t *testing.T) {
	r := require.New(t)
	ctx := testutils.GetTestContext(t)
	router := mux.NewRouter()
	router.Use(logContextMiddleware)
	router.NotFoundHandler = logContextMiddleware(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				logger.Info(r.Context()).Msg("test")
			},
		),
	)
	logs := logger.NewCollector()
	ctx = logger.NewLogger(logger.LevelTrace, logs).Context(ctx)
	srv := httptest.NewUnstartedServer(router)
	srv.Config.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}
	srv.Start()
	defer srv.Close()
	resp, err := srv.Client().Get(srv.URL + "/test-not-found")
	r.NoError(err)
	r.Equal(http.StatusNotFound, resp.StatusCode)
	r.Len(logs.Logs(), 1)
	l := logs.Logs()[0]
	r.NotEmpty(l.Keys["remote"])
	r.EqualExportedValues(
		[]logger.EventData{
			{
				Time:  l.Time,
				Level: logger.LevelInfo,
				Keys: map[string]interface{}{
					"remote": l.Keys["remote"],
					"method": "GET",
					"url":    "/test-not-found",
				},
				Message: "test",
				File:    l.File,
			},
		}, logs.Logs(),
	)
}

func TestAccessLogMiddleware(t *testing.T) {
	r := require.New(t)
	ctx := testutils.GetTestContext(t)
	router := mux.NewRouter()
	router.Use(logContextMiddleware)
	router.Use(accessLogMiddleware)
	router.HandleFunc(
		"/test/{id:"+utils.UUIDRegexString+"}", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		},
	).Name("test-route")
	logs := logger.NewCollector()
	ctx = logger.NewLogger(logger.LevelTrace, logs).Context(ctx)
	srv := httptest.NewUnstartedServer(router)
	srv.Config.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}
	srv.Start()
	defer srv.Close()
	id := uuid.New()
	resp, err := http.Get(srv.URL + "/test/" + id.String())
	r.NoError(err)
	r.Equal(http.StatusTeapot, resp.StatusCode)
	r.Len(logs.Logs(), 1)
	l := logs.Logs()[0]
	r.EqualExportedValues(
		[]logger.EventData{
			{
				Time:  l.Time,
				Level: logger.LevelInfo,
				Keys: map[string]interface{}{
					"endpoint": "/test/{id:uuid}",
					"remote":   l.Keys["remote"],
					"method":   "GET",
					"url":      "/test/" + id.String(),
					"route":    "test-route",
					"code":     http.StatusTeapot,
					"duration": l.Keys["duration"],
					"written":  l.Keys["written"],
				},
				Message: "request",
				File:    l.File,
			},
		}, logs.Logs(),
	)
}

func TestRecovery(t *testing.T) {
	r := require.New(t)
	ctx := testutils.GetTestContext(t)
	router := mux.NewRouter()
	router.Use(recoveryMiddleware)
	router.HandleFunc(
		"/panic", func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		},
	)
	logs := logger.NewCollector()
	ctx = logger.NewLogger(logger.LevelTrace, logs).Context(ctx)
	srv := httptest.NewUnstartedServer(router)
	srv.Config.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}
	srv.Start()
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/panic")
	r.NoError(err)
	r.Equal(http.StatusInternalServerError, resp.StatusCode)
	r.Len(logs.Logs(), 1)
	l := logs.Logs()[0]
	r.NotEmpty(l.Keys["stack"])
	r.Contains(l.File, "middleware.go")
	r.EqualExportedValues(
		[]logger.EventData{
			{
				Time:  l.Time,
				Level: logger.LevelError,
				Keys: map[string]interface{}{
					"recover": "test panic",
					"stack":   l.Keys["stack"],
				},
				Message: "Panic recovered",
				File:    l.File,
			},
		}, logs.Logs(),
	)
}
