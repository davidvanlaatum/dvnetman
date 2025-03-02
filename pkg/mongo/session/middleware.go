package mongosession

import (
	"context"
	"dvnetman/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
)

type SessionData struct {
	Session *sessions.Session
	Data    interface{}
}

type MiddleWareContextStorage struct {
	sessions map[string]*SessionData
}

func (m *MiddleWareContextStorage) GetSession(name string) *SessionData {
	return m.sessions[name]
}

func (m *MiddleWareContextStorage) AddSession(name string, session *SessionData) {
	m.sessions[name] = session
}

func GetMiddleWareStorage(ctx context.Context) *MiddleWareContextStorage {
	return ctx.Value(middleWareContextKey).(*MiddleWareContextStorage)
}

type middlewareContextKey struct{}

var middleWareContextKey middlewareContextKey

func Middleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				storage := &MiddleWareContextStorage{
					sessions: make(map[string]*SessionData),
				}
				ctx := context.WithValue(r.Context(), middleWareContextKey, storage)
				req := r.WithContext(ctx)
				wr := &middlewareWriteWrapper{
					ResponseWriter: w,
					storage:        storage,
					r:              req,
				}
				next.ServeHTTP(wr, req)
			},
		)
	}
}

type middlewareWriteWrapper struct {
	http.ResponseWriter
	saveDone bool
	storage  *MiddleWareContextStorage
	r        *http.Request
}

func (m *middlewareWriteWrapper) save() {
	if m.saveDone {
		return
	}
	m.saveDone = true
	for key, sess := range m.storage.sessions {
		if err := sess.Session.Save(m.r, m.ResponseWriter); err != nil {
			logger.Error(m.r.Context()).Key("session", key).Msgf("failed to save session: %+v", err)
		}
	}
}

func (m *middlewareWriteWrapper) Write(bytes []byte) (int, error) {
	m.save()
	return m.ResponseWriter.Write(bytes)
}

func (m *middlewareWriteWrapper) WriteHeader(statusCode int) {
	m.save()
	m.ResponseWriter.WriteHeader(statusCode)
}

var _ http.ResponseWriter = (*middlewareWriteWrapper)(nil)
