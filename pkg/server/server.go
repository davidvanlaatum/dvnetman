package server

import (
	"context"
	"dvnetman/api"
	"dvnetman/pkg/auth"
	"dvnetman/pkg/cache"
	"dvnetman/pkg/config"
	"dvnetman/pkg/logger"
	"dvnetman/pkg/mongo/adapt"
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/mongo/otel"
	"dvnetman/pkg/mongo/session"
	"dvnetman/pkg/openapi"
	"dvnetman/pkg/server/service"
	"dvnetman/pkg/ymlutil"
	"dvnetman/web"
	"encoding/base64"
	"github.com/felixge/httpsnoop"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"gopkg.in/yaml.v3"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	httpServer *http.Server
	router     *mux.Router
	service    *service.Service
	config     *config.Config
	db         *mongo.Client
	store      sessions.Store
	auth       *auth.Auth
	otel       *otelServer
	client     *modal.DBClient
	cache      cache.Pool
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) startHTTP(ctx context.Context, cancel context.CancelCauseFunc) error {
	s.httpServer = &http.Server{
		Handler: s.router,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		IdleTimeout: 5 * time.Second,
		ReadTimeout: 5 * time.Second,
	}
	for _, address := range s.config.Listen {
		logger.Info(ctx).Key("address", address.Addr).Msg("Listening")
		if l, err := net.Listen("tcp", address.Addr); err != nil {
			return err
		} else {
			go func(l net.Listener) {
				if err := s.httpServer.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
					cancel(err)
				} else {
					cancel(nil)
				}
			}(l)
		}
	}
	return nil
}

func (s *Server) connectToMongo(ctx context.Context) (err error) {
	apiOptions := options.Client().
		ApplyURI(s.config.Mongo.URL).
		SetAppName("DVNetMan").
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetMonitor(otelmongo.NewMonitor(otelmongo.WithCommandAttributeDisabled(false)))
	if s.db, err = mongo.Connect(apiOptions); err != nil {
		return errors.WithMessage(err, "failed to connect to MongoDB")
	}

	if err = s.db.Ping(ctx, nil); err != nil {
		if err = errors.WithMessage(err, "failed to ping MongoDB"); err != nil {
			return
		}
	}

	s.client = modal.NewDBClient(s.getMongoDatabase())
	if err = s.client.Init(ctx); err != nil {
		return errors.WithMessage(err, "failed to initialize database client")
	}
	return
}

func (s *Server) disconnectFromMongo(ctx context.Context) {
	if s.db != nil {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		logger.Info(ctx).Msg("disconnecting from MongoDB")
		_ = s.db.Disconnect(ctx)
	}
}

func accessLogMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			m := httpsnoop.CaptureMetrics(handler, w, r)
			logger.Info(r.Context()).
				Key("code", m.Code).
				Key("duration", m.Duration).
				Key("written", m.Written).
				Msg("request")
		},
	)
}

func (s *Server) setupRouter(ctx context.Context) error {
	s.service = service.NewService(ctx, s.client, s.auth, s.cache)
	router := mux.NewRouter()
	router.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
	router.Use(handlers.ProxyHeaders)
	router.Use(otelhttp.NewMiddleware("http"))
	router.Use(traceIDHeaderMiddleware)
	router.Use(logger.Middleware())
	router.Use(mongosession.Middleware())
	router.Use(s.auth.AuthMiddleware)
	router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	router.Use(accessLogMiddleware)
	s.auth.AddRoutes(router)
	router.Methods("GET").Path("/api/openapi.yaml").HandlerFunc(s.openapiSpec).Name("OpenAPI")
	apiRouter := openapi.NewRouter(s.service)
	router.PathPrefix("/api").Handler(apiRouter)
	router.PathPrefix("/ui/").Handler(web.NewUI().GetRouter())
	router.Methods("GET").Path("/metrics").Handler(promhttp.Handler())
	_ = router.Walk(
		func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			if _, ok := route.GetHandler().(*mux.Router); ok {
				return nil
			}
			l := logger.Info(ctx)
			if methods, err := route.GetMethods(); err == nil {
				l.Key("methods", methods)
			}
			if path, err := route.GetPathTemplate(); err == nil {
				l.Key("path", path)
			}
			l.Msg("Route")
			router.NotFoundHandler = http.HandlerFunc(s.notFound)
			router.MethodNotAllowedHandler = http.HandlerFunc(s.methodNotAllowed)
			return nil
		},
	)
	s.router = router
	s.otel.attach(ctx, router)
	return nil
}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	logger.Warn(r.Context()).Key("path", r.URL.Path).Msg("Not found")
	http.Error(w, "Not found", http.StatusNotFound)
}

func (s *Server) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	logger.Warn(r.Context()).Key("path", r.URL.Path).Key("method", r.Method).Msg("Method not allowed")
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (s *Server) openapiSpec(w http.ResponseWriter, r *http.Request) {
	x := &yaml.Node{}
	if err := yaml.Unmarshal(api.OpenapiYAML, x); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	ymlutil.Walk(
		x, func(n *yaml.Node, path []string) {
			if len(path) == 3 && path[0] == "servers" && path[2] == "url" {
				n.Value = "http://" + r.Host
			}
		},
	)
	if b, err := yaml.Marshal(x); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Header().Set("Content-Length", strconv.FormatUint(uint64(len(b)), 10))
		_, _ = w.Write(b)
	}
}

func (s *Server) getMongoDatabase() mongoadapt.MongoDatabase {
	return mongoadapt.AdapterMongoDatabase(s.db.Database(s.config.Mongo.Database))
}

func (s *Server) setupAuth(ctx context.Context) (err error) {
	hashKey := s.config.Session.HashKeyBytes()
	if hashKey == nil {
		hashKey = securecookie.GenerateRandomKey(64)
		logger.Info(ctx).Msgf("Generated hash key: %v", base64.StdEncoding.EncodeToString(hashKey))
	}
	blockKey := s.config.Session.BlockKeyBytes()
	if blockKey == nil {
		blockKey = securecookie.GenerateRandomKey(32)
		logger.Info(ctx).Msgf("Generated block key: %v", base64.StdEncoding.EncodeToString(blockKey))
	}
	secureCookie := securecookie.New(hashKey, blockKey)
	if s.store, err = mongosession.NewMongoStore(
		s.getMongoDatabase().Collection("session"), secureCookie,
	); err != nil {
		return errors.WithMessage(err, "failed to create session store")
	}
	s.auth = auth.NewAuth(s.store, s.config, s.client, s.service.ErrorHandler)
	if err = s.auth.Init(); err != nil {
		return errors.WithMessage(err, "failed to initialize authentication")
	}
	return nil
}

func (s *Server) stopHTTP(ctx context.Context) {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		logger.Info(ctx).Msg("shutting down HTTP server")
		if err := s.httpServer.Shutdown(ctx); err != nil {
			logger.Error(ctx).Err(err).Msg("Failed to shutdown HTTP server")
		}
	}
}

func (s *Server) setupOtel(ctx context.Context) (err error) {
	s.otel = &otelServer{}
	if err = s.otel.setup(ctx); err != nil {
		return errors.WithMessage(err, "failed to setup OpenTelemetry")
	}
	return nil
}

func (s *Server) setupCache(ctx context.Context) (err error) {
	cfg := s.config.Cache
	if cfg == "" {
		cfg = "memory://?size=10MB&ttl=1h"
	}
	s.cache, err = cache.NewPool(ctx, cfg)
	return
}

func (s *Server) Start(ctx context.Context) (err error) {
	var c, cancel = context.WithCancelCause(ctx)
	defer cancel(nil)
	if err = s.setupOtel(c); err != nil {
		return errors.WithMessage(err, "failed to setup OpenTelemetry")
	}
	defer s.otel.shutdown(context.WithoutCancel(ctx))
	if err = s.setupCache(c); err != nil {
		return errors.WithMessage(err, "failed to setup cache")
	}
	defer func() {
		if err := s.cache.Shutdown(context.WithoutCancel(ctx)); err != nil {
			logger.Error(ctx).Err(err).Msg("failed to shutdown cache")
		}
	}()
	if err = s.connectToMongo(c); err != nil {
		return errors.WithMessage(err, "failed to connect to MongoDB")
	}
	defer s.disconnectFromMongo(context.WithoutCancel(ctx))
	if err = s.setupAuth(ctx); err != nil {
		return errors.WithMessage(err, "failed to setup authentication")
	}
	if err = s.setupRouter(ctx); err != nil {
		return errors.WithMessage(err, "failed to setup HTTP router")
	}
	if err = s.startHTTP(c, cancel); err != nil {
		return errors.WithMessage(err, "failed to start HTTP server")
	}
	defer s.stopHTTP(context.WithoutCancel(ctx))
	<-ctx.Done()
	return ctx.Err()
}
