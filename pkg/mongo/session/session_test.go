package mongosession

import (
	mongoadapt "dvnetman/pkg/mongo/adapt"
	"dvnetman/pkg/testutils"
	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSession(t *testing.T) {
	m := testutils.StartMongo(t, testutils.GetTestContext(t))
	t.Run(
		"SaveAndGetSession", func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			ctx := testutils.GetTestContext(t)
			client := testutils.GetMongoClient(t, ctx, m)
			db := client.Database("session")
			key := securecookie.GenerateRandomKey(64)
			key2 := securecookie.GenerateRandomKey(32)
			s, err := NewMongoStore(
				mongoadapt.AdapterMongoDatabase(db).Collection(t.Name()), securecookie.New(key, key2),
			)
			r.NoError(err)
			req := httptest.NewRequestWithContext(ctx, "GET", "/", nil)

			_, err = NewMongoStore(
				mongoadapt.AdapterMongoDatabase(db).Collection(t.Name()), securecookie.New(key, key2),
			)
			r.NoError(err)

			recorder := httptest.NewRecorder()
			Middleware()(
				http.HandlerFunc(
					func(w http.ResponseWriter, req *http.Request) {
						session, err := s.Get(req, "test")
						r.NoError(err)
						session.Values["test"] = "test"
						w.WriteHeader(http.StatusOK)
					},
				),
			).ServeHTTP(recorder, req)
			req.AddCookie(recorder.Result().Cookies()[0])
			recorder = httptest.NewRecorder()
			Middleware()(
				http.HandlerFunc(
					func(w http.ResponseWriter, req *http.Request) {
						session, err := s.Get(req, "test")
						r.NoError(err)
						r.Equal("test", session.Values["test"])
						session.Values["test"] = "test2"
						w.WriteHeader(http.StatusOK)
					},
				),
			).ServeHTTP(recorder, req)
			recorder = httptest.NewRecorder()
			Middleware()(
				http.HandlerFunc(
					func(w http.ResponseWriter, req *http.Request) {
						session, err := s.Get(req, "test")
						r.NoError(err)
						r.Equal("test2", session.Values["test"])
						w.WriteHeader(http.StatusOK)
					},
				),
			).ServeHTTP(recorder, req)
		},
	)
	t.Run(
		"DeleteSession", func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			ctx := testutils.GetTestContext(t)
			client := testutils.GetMongoClient(t, ctx, m)
			db := client.Database("session")
			key := securecookie.GenerateRandomKey(64)
			key2 := securecookie.GenerateRandomKey(32)
			s, err := NewMongoStore(
				mongoadapt.AdapterMongoDatabase(db).Collection(t.Name()), securecookie.New(key, key2),
			)
			r.NoError(err)

			req := httptest.NewRequestWithContext(ctx, "GET", "/", nil)
			recorder := httptest.NewRecorder()
			var sessionId string
			Middleware()(
				http.HandlerFunc(
					func(w http.ResponseWriter, req *http.Request) {
						session, err := s.Get(req, "test")
						session.Values["test"] = "test"
						r.NoError(err)
						w.WriteHeader(http.StatusOK)
						sessionId = session.ID
					},
				),
			).ServeHTTP(recorder, req)
			req.AddCookie(recorder.Result().Cookies()[0])
			recorder = httptest.NewRecorder()
			Middleware()(
				http.HandlerFunc(
					func(w http.ResponseWriter, req *http.Request) {
						session, err := s.Get(req, "test")
						r.NoError(err)
						r.Equal(sessionId, session.ID)
						session.Options.MaxAge = -1
						w.WriteHeader(http.StatusOK)
					},
				),
			).ServeHTTP(recorder, req)
			recorder = httptest.NewRecorder()
			Middleware()(
				http.HandlerFunc(
					func(w http.ResponseWriter, req *http.Request) {
						session, err := s.Get(req, "test")
						r.NoError(err)
						r.True(session.IsNew)
						r.NotEqual(sessionId, session.ID)
						w.WriteHeader(http.StatusOK)
					},
				),
			).ServeHTTP(recorder, req)
		},
	)
}
