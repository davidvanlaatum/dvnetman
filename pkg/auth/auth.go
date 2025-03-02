package auth

import (
	"bytes"
	"context"
	"dvnetman/pkg/config"
	"dvnetman/pkg/logger"
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/utils"
	"encoding/gob"
	"github.com/automattic/go-gravatar"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"time"
)

const (
	sessionKey              = "session"
	returnToSessionKey      = "returnTo"
	oauthProviderSessionKey = "oauthProvider"
	userSessionKey          = "user"
)

type Auth struct {
	store        sessions.Store
	cfg          *config.Config
	db           *modal.DBClient
	errorHandler func(w http.ResponseWriter, r *http.Request, err error)
}

func NewAuth(
	store sessions.Store, cfg *config.Config, db *modal.DBClient,
	errorHandler func(w http.ResponseWriter, r *http.Request, err error),
) *Auth {
	return &Auth{
		store:        store,
		cfg:          cfg,
		db:           db,
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
		u = u.ResolveReference(&url.URL{Path: "auth/" + cfg.Provider + "/callback"})
		var p *openidConnect.Provider
		if p, err = openidConnect.NewNamed(
			cfg.Provider, cfg.OpenIDConnect.ClientID, cfg.OpenIDConnect.ClientSecret, u.String(),
			cfg.OpenIDConnect.AutoDiscoveryURL, cfg.OpenIDConnect.Scopes...,
		); err != nil {
			err = errors.Wrap(err, "failed to create OpenID Connect provider")
		}
		p.HTTPClient = &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
		p.SetName(cfg.Provider)
		provider = p
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
	router.Methods("GET").Path("/auth/logout").HandlerFunc(a.logout).Name("logout")
	router.Methods("GET").Path("/auth/{provider}").HandlerFunc(a.oAuthBeginHandler).Name("auth")
	router.Methods("GET").Path("/auth/{provider}/callback").HandlerFunc(a.oAuthCallbackHandler).Name("authCallback")
}

func (a *Auth) logout(w http.ResponseWriter, r *http.Request) {
	logger.Info(r.Context()).Msg("Logging out")
	if err := gothic.Logout(w, r); err != nil {
		a.errorHandler(w, r, errors.Wrap(err, "failed to logout"))
		return
	}
	if s, err := a.store.Get(r, sessionKey); err != nil {
		a.errorHandler(w, r, errors.Wrap(err, "failed to get session"))
		return
	} else {
		delete(s.Values, userSessionKey)
	}
	logger.Info(r.Context()).Msgf("headers %+v", w.Header())
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (a *Auth) AuthProviders() []Provider {
	var providers []Provider
	for _, provider := range a.cfg.Auth {
		providers = append(
			providers, Provider{
				Provider:            provider.Provider,
				DisplayName:         provider.DisplayName,
				LoginURL:            "/auth/" + provider.Provider,
				LoginButtonImageURL: provider.LoginButtonURL,
			},
		)
	}
	return providers
}

func (a *Auth) oAuthBeginHandler(w http.ResponseWriter, r *http.Request) {
	if s, err := a.store.Get(r, sessionKey); err != nil {
		a.errorHandler(w, r, errors.Wrap(err, "failed to get session"))
		return
	} else {
		s.Values[returnToSessionKey] = r.Header.Get("Referer")
		if s.Values[oauthProviderSessionKey], err = gothic.GetProviderName(r); err != nil {
			a.errorHandler(w, r, errors.Wrap(err, "failed to get provider name"))
			return
		}
		logger.Debug(r.Context()).Key("returnTo", s.Values[returnToSessionKey]).Msg("Saved returnTo")
	}
	gothic.BeginAuthHandler(w, r)
}

func (a *Auth) lookupUser(ctx context.Context, gothUser goth.User) (user *User, err error) {
	var u *modal.User
	if u, err = a.db.GetUserByExternalID(ctx, gothUser.Provider, gothUser.UserID); err != nil && !errors.Is(
		err, mongo.ErrNoDocuments,
	) {
		return nil, errors.Wrap(err, "failed to get user by external ID")
	}
	if u == nil {
		u = &modal.User{
			ExternalProvider: utils.ToPtr(gothUser.Provider),
			ExternalID:       utils.ToPtr(gothUser.UserID),
		}
	}
	if gothUser.Email != "" {
		u.Email = utils.ToPtr(gothUser.Email)
	}
	if gothUser.FirstName != "" {
		u.FirstName = utils.ToPtr(gothUser.FirstName)
	}
	if gothUser.LastName != "" {
		u.LastName = utils.ToPtr(gothUser.LastName)
	}
	if gothUser.Name != "" {
		u.DisplayName = utils.ToPtr(gothUser.Name)
	} else if gothUser.NickName != "" {
		u.DisplayName = utils.ToPtr(gothUser.NickName)
	} else if u.DisplayName == nil {
		u.DisplayName = utils.JoinPtr(" ", u.FirstName, u.LastName)
	}
	if err = a.db.SaveUser(ctx, u); err != nil {
		return nil, errors.Wrap(err, "failed to save user")
	}
	user = &User{
		ID:               uuid.UUID(*u.ID),
		DisplayName:      u.DisplayName,
		Email:            u.Email,
		ExternalID:       u.ExternalID,
		ExternalProvider: u.ExternalProvider,
		OAuthUser:        &gothUser,
	}
	if u.Email != nil {
		grav := gravatar.NewGravatarFromEmail(*u.Email)
		grav.Default = "mp"
		user.ProfileImageURL = utils.ToPtr(grav.GetURL())
	}
	return
}

func (a *Auth) setUserInSession(r *http.Request, u *User) (err error) {
	var s *sessions.Session
	if s, err = a.store.Get(r, sessionKey); err != nil {
		err = errors.Wrap(err, "failed to get session")
		return
	}
	b := &bytes.Buffer{}
	e := gob.NewEncoder(b)
	if err = e.Encode(u); err != nil {
		err = errors.Wrap(err, "failed to encode user")
		return
	}
	s.Values[userSessionKey] = b.String()
	return
}

func (a *Auth) oAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	var oAuthUser goth.User
	var err error
	if oAuthUser, err = gothic.CompleteUserAuth(w, r); err != nil {
		a.errorHandler(w, r, errors.Wrap(err, "failed to complete user auth"))
		return
	}
	var u *User
	if u, err = a.lookupUser(r.Context(), oAuthUser); err != nil {
		a.errorHandler(w, r, errors.Wrap(err, "failed to lookup user"))
		return
	}
	if err = a.setUserInSession(r, u); err != nil {
		a.errorHandler(w, r, errors.Wrap(err, "failed to set user in session"))
		return
	}
	var s *sessions.Session
	if s, err = a.store.Get(r, sessionKey); err != nil {
		a.errorHandler(w, r, errors.Wrap(err, "failed to get session"))
		return
	}
	if returnTo, ok := s.Values[returnToSessionKey].(string); ok {
		logger.Debug(r.Context()).Key("returnTo", returnTo).Msg("Redirecting")
		delete(s.Values, returnToSessionKey)
		http.Redirect(w, r, returnTo, http.StatusTemporaryRedirect)
	} else {
		logger.Debug(r.Context()).Msg("Redirecting to /")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

type userContextKeyStruct struct{}

var userContextKey = userContextKeyStruct{}

func UserFromContext(ctx context.Context) *User {
	if x := ctx.Value(userContextKey); x != nil {
		return x.(*User)
	}
	return nil
}

func (a *Auth) GetUser(r *http.Request) (user *User, err error) {
	if user = UserFromContext(r.Context()); user != nil {
		return
	}
	var s *sessions.Session
	if s, err = a.store.Get(r, sessionKey); err != nil {
		err = errors.Wrap(err, "failed to get session")
		return
	}
	if s.Values[userSessionKey] == nil {
		logger.Debug(r.Context()).Msg("No user in session")
		user = &User{
			ID: anonymousUserId,
		}
		return
	}
	b := bytes.NewBufferString(s.Values[userSessionKey].(string))
	d := gob.NewDecoder(b)
	user = &User{}
	if err = d.Decode(user); err != nil {
		err = errors.Wrap(err, "failed to decode user")
	}
	if user.IsAuthenticated() && user.OAuthUser.ExpiresAt.Before(time.Now().Add(time.Minute)) {
		logger.Debug(r.Context()).Msg("Time to refresh token")
		var provider goth.Provider
		if provider, err = goth.GetProvider(user.OAuthUser.Provider); err != nil {
			return nil, errors.Wrap(err, "failed to get provider")
		}
		var token *oauth2.Token
		if token, err = provider.RefreshToken(user.OAuthUser.RefreshToken); err != nil {
			return nil, errors.Wrap(err, "failed to refresh token")
		}
		user.OAuthUser.AccessToken = token.AccessToken
		user.OAuthUser.RefreshToken = token.RefreshToken
		user.OAuthUser.ExpiresAt = token.Expiry
		if err = a.setUserInSession(r, user); err != nil {
			return nil, errors.Wrap(err, "failed to set user in session")
		}
	} else {
		logger.Debug(r.Context()).Msgf("user token expires in %v", time.Until(user.OAuthUser.ExpiresAt))
	}
	return
}

func (a *Auth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if u, err := a.GetUser(r); err != nil {
				logger.Debug(r.Context()).Err(err).Msg("No user in session")
			} else if u != nil {
				ctx := context.WithValue(r.Context(), userContextKey, u)
				if u.IsAuthenticated() {
					ctx = logger.Ctx(r.Context()).SubLogger().
						Key("user", u.ID).
						Logger().
						Context(ctx)
				}
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		},
	)
}

func RequirePerm(ctx context.Context, perm Permission) (err error) {
	if u := UserFromContext(ctx); u == nil || !u.IsAuthenticated() {
		err = errors.WithStack(&NotLoggedInError{})
	}
	return
}
