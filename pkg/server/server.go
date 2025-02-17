package server

import (
	"context"
	"dvnetman/api"
	"dvnetman/pkg/auth"
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
	"log"
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
	log        logger.Logger
	store      sessions.Store
	auth       *auth.Auth
	otel       *otelServer
}

func NewServer(cfg *config.Config, log logger.Logger) *Server {
	return &Server{
		config: cfg,
		log:    log,
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
		s.log.Info().Key("address", address.Addr).Msg("Listening")
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
		err = errors.WithMessage(err, "failed to ping MongoDB")
	}
	return
}

func (s *Server) disconnectFromMongo() {
	if s.db != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.log.Info().Msg("disconnecting from MongoDB")
		_ = s.db.Disconnect(ctx)
	}
}

func (s *Server) setupRouter() error {
	s.service = service.NewService(modal.NewDBClient(s.getMongoDatabase()))
	router := mux.NewRouter()
	router.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
	router.Use(otelhttp.NewMiddleware("http"))
	router.Use(logger.Middleware(s.log))
	router.Use(mongosession.Middleware())
	router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	router.Use(
		func(handler http.Handler) http.Handler {
			return http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					m := httpsnoop.CaptureMetrics(handler, w, r)
					logger.Ctx(r.Context()).Info().Key("code", m.Code).Key("duration", m.Duration).Msg("request")
				},
			)
		},
	)
	router.Use(handlers.CompressHandler)
	router.Use(handlers.ProxyHeaders)
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
			l := s.log.Info()
			if methods, err := route.GetMethods(); err == nil {
				l.Key("methods", methods)
			}
			if path, err := route.GetPathTemplate(); err == nil {
				l.Key("path", path)
			}
			l.Msg("Route")
			return nil
		},
	)
	router.NotFoundHandler = handlers.LoggingHandler(
		log.Writer(), http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				log.Printf("Not found: %s", r.URL.Path)
				http.Error(w, "Not found", http.StatusNotFound)
			},
		),
	)
	s.router = router
	s.otel.attach(router)
	return nil
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

func (s *Server) setupAuth() (err error) {
	hashKey := s.config.Session.HashKeyBytes()
	if hashKey == nil {
		hashKey = securecookie.GenerateRandomKey(64)
		s.log.Info().Msgf("Generated hash key: %v", base64.StdEncoding.EncodeToString(hashKey))
	}
	blockKey := s.config.Session.BlockKeyBytes()
	if blockKey == nil {
		blockKey = securecookie.GenerateRandomKey(32)
		s.log.Info().Msgf("Generated block key: %v", base64.StdEncoding.EncodeToString(blockKey))
	}
	secureCookie := securecookie.New(hashKey, blockKey)
	if s.store, err = mongosession.NewMongoStore(
		s.getMongoDatabase().Collection("session"), secureCookie,
	); err != nil {
		return errors.WithMessage(err, "failed to create session store")
	}
	s.auth = auth.NewAuth(s.log, s.store, s.config, s.service.ErrorHandler)
	if err = s.auth.Init(); err != nil {
		return errors.WithMessage(err, "failed to initialize authentication")
	}
	return nil
}

//func (s *Server) errorHandler(w http.ResponseWriter, _ *http.Request, err error, _ *openapi.ImplResponse) {
//	var res *openapi.ImplResponse
//	for _, conv := range errorConverters {
//		if res = conv(err); res != nil {
//			break
//		}
//	}
//	if res == nil {
//		s.log.Error().Msgf("Unhandled error %v (%T)", err, errors.Unwrap(err))
//		res = &openapi.ImplResponse{
//			Code: http.StatusInternalServerError,
//			Body: openapi.Error{
//				Errors: []openapi.ErrorMessage{
//					{Code: "UNKNOWN", Message: err.Error()},
//				},
//			},
//		}
//	}
//	_ = openapi.EncodeJSONResponse(res.Body, &res.Code, res.Headers, w)
//}

func (s *Server) stopHTTP() {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.log.Info().Msg("shutting down HTTP server")
		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.log.Error().Err(err).Msg("Failed to shutdown HTTP server")
		}
	}
}

func (s *Server) setupOtel(ctx context.Context) (err error) {
	s.otel = &otelServer{
		log: s.log,
	}
	if err = s.otel.setup(ctx); err != nil {
		return errors.WithMessage(err, "failed to setup OpenTelemetry")
	}
	return nil
}

func (s *Server) Start(ctx context.Context) (err error) {
	var c, cancel = context.WithCancelCause(ctx)
	defer cancel(nil)
	if err = s.setupOtel(c); err != nil {
		return errors.WithMessage(err, "failed to setup OpenTelemetry")
	}
	defer s.otel.shutdown(context.Background())
	if err = s.connectToMongo(c); err != nil {
		return errors.WithMessage(err, "failed to connect to MongoDB")
	}
	defer s.disconnectFromMongo()
	if err = s.setupAuth(); err != nil {
		return errors.WithMessage(err, "failed to setup authentication")
	}
	if err = s.setupRouter(); err != nil {
		return errors.WithMessage(err, "failed to setup HTTP router")
	}
	if err = s.startHTTP(c, cancel); err != nil {
		return errors.WithMessage(err, "failed to start HTTP server")
	}
	defer s.stopHTTP()
	<-ctx.Done()
	return ctx.Err()
}
