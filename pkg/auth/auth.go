package auth

import (
	"dvnetman/pkg/config"
	"dvnetman/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

const returnToSessionKey = "returnTo"

type Auth struct {
	log          logger.Logger
	store        sessions.Store
	cfg          *config.Config
	errorHandler func(w http.ResponseWriter, r *http.Request, err error)
}

func NewAuth(
	log logger.Logger, store sessions.Store, cfg *config.Config,
	errorHandler func(w http.ResponseWriter, r *http.Request, err error),
) *Auth {
	return &Auth{
		log:          log,
		store:        store,
		cfg:          cfg,
		errorHandler: errorHandler,
	}
}

func (a *Auth) initProvider(cfg *config.AuthConfig) (provider goth.Provider, err error) {
	if cfg.OpenIDConnect != nil {
		var u *url.URL
		if u, err = a.cfg.GetURL(); err != nil {
			err = errors.Wrap(err, "failed to get server URL")
			return
		}
		u = u.ResolveReference(&url.URL{Path: "auth/" + cfg.OpenIDConnect.Provider + "/callback"})
		if provider, err = openidConnect.New(
			cfg.OpenIDConnect.ClientID, cfg.OpenIDConnect.ClientSecret, u.String(),
			cfg.OpenIDConnect.AutoDiscoveryURL, cfg.OpenIDConnect.Scopes...,
		); err != nil {
			err = errors.Wrap(err, "failed to create OpenID Connect provider")
		}
		provider.SetName(cfg.OpenIDConnect.Provider)
	}
	return
}

func (a *Auth) Init() (err error) {
	gothic.Store = a.store
	var providers []goth.Provider
	for i, c := range a.cfg.Auth {
		var provider goth.Provider
		if provider, err = a.initProvider(&c); err != nil {
			return errors.WithMessagef(err, "failed to initialize auth provider %v", i)
		}
		providers = append(providers, provider)
	}
	goth.UseProviders(providers...)
	return
}

func (a *Auth) AddRoutes(router *mux.Router) {
	router.Methods("GET").Path("/auth/{provider}").HandlerFunc(a.oAuthBeginHandler).Name("auth")
	router.Methods("GET").Path("/auth/{provider}/callback").HandlerFunc(a.oAuthCallbackHandler).Name("authCallback")
}

func (a *Auth) oAuthBeginHandler(w http.ResponseWriter, r *http.Request) {
	if s, err := a.store.Get(r, "session"); err != nil {
		a.errorHandler(w, r, errors.Wrap(err, "failed to get session"))
		return
	} else {
		s.Values["returnTo"] = r.Header.Get("Referer")
		s.Values["oauthProvider"] = mux.Vars(r)["provider"]
		logger.Ctx(r.Context()).Debug().Key("returnTo", s.Values["returnTo"]).Msg("Saved returnTo")
	}
	gothic.BeginAuthHandler(w, r)
}

func (a *Auth) oAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	var user goth.User
	var err error
	if user, err = gothic.CompleteUserAuth(w, r); err != nil {
		a.errorHandler(w, r, errors.Wrap(err, "failed to complete user auth"))
		return
	}
	logger.Ctx(r.Context()).Debug().Key("user", user).Msg("User")
	var s *sessions.Session
	if s, err = a.store.Get(r, "session"); err != nil {
		a.errorHandler(w, r, errors.Wrap(err, "failed to get session"))
		return
	}
	s.Values["oAuthUser"] = user
	if returnTo, ok := s.Values[returnToSessionKey].(string); ok {
		logger.Ctx(r.Context()).Debug().Key("returnTo", returnTo).Msg("Redirecting")
		delete(s.Values, returnToSessionKey)
		http.Redirect(w, r, returnTo, http.StatusTemporaryRedirect)
	} else {
		logger.Ctx(r.Context()).Debug().Msg("Redirecting to /")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}
