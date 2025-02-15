package mongosession

import (
	"bytes"
	"context"
	"dvnetman/pkg/logger"
	"dvnetman/pkg/mongo/adapt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"net/http"
	"time"
)

type sessionData struct {
	ID             bson.ObjectID          `bson:"_id,omitempty"`
	Data           map[string]interface{} `bson:"data"`
	LastAccess     time.Time              `bson:"last_access"`
	Version        int                    `bson:"version"`
	serializedData []byte                 `bson:"-"`
}

func (s *sessionData) setSerialized() (err error) {
	s.serializedData, err = bson.Marshal(s.Data)
	return
}

func (s *sessionData) setSessionValues(values map[interface{}]interface{}) (err error) {
	var c map[string]interface{}
	if err = bson.Unmarshal(s.serializedData, &c); err == nil {
		for k, v := range c {
			values[k] = v
		}
	}
	return
}

func (s *sessionData) hasDataChangedAndUpdate(values map[interface{}]interface{}) (changed bool, err error) {
	t := map[string]interface{}{}
	for k, v := range values {
		t[k.(string)] = v
	}
	var b []byte
	if b, err = bson.Marshal(t); err != nil {
		return false, err
	}
	if bytes.Equal(s.serializedData, b) {
		return false, nil
	}
	s.Data = nil
	if err = bson.Unmarshal(b, &s.Data); err != nil {
		return true, err
	}
	s.serializedData = b
	return true, nil
}

type MongoStore struct {
	db     mongoadapt.MongoCollection
	codecs []securecookie.Codec
	//options *sessions.Options
}

func NewMongoStore(db mongoadapt.MongoCollection, codecs ...securecookie.Codec) (
	s *MongoStore, err error,
) {
	if _, err = db.Indexes().CreateOne(
		context.Background(), mongo.IndexModel{
			Keys: bson.M{
				"last_access": 1,
			},
			Options: options.Index().SetExpireAfterSeconds(3600),
		},
	); err != nil {
		return nil, errors.Wrap(err, "failed to create ttl index")
	}
	return &MongoStore{
		db:     db,
		codecs: codecs,
	}, nil
}

func (m *MongoStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	store := GetMiddleWareStorage(r.Context())
	if s := store.GetSession(name); s != nil {
		return s.Session, nil
	}
	return m.New(r, name)
}

func (m *MongoStore) New(r *http.Request, name string) (*sessions.Session, error) {
	s := sessions.NewSession(m, name)
	s.Options = &sessions.Options{
		Domain:   r.URL.Host,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   3600,
		Path:     "/",
	}
	store := GetMiddleWareStorage(r.Context())
	data := &SessionData{Session: s}
	store.AddSession(name, data)
	if c := r.CookiesNamed(name); len(c) > 0 {
		if err := securecookie.DecodeMulti(name, c[0].Value, &s.ID, m.codecs...); err == nil {
			if err = m.load(r.Context(), data); err == nil {
				s.IsNew = false
				return s, nil
			}
		} else {
			logger.Ctx(r.Context()).Error().Msgf("failed to decode session %s starting new: %+v", c[0].Value, err)
		}
	}
	s.IsNew = true
	s.ID = ""
	return s, nil
}

func (m *MongoStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) (err error) {
	l := logger.Ctx(r.Context())
	l.Debug().
		Key("session", s.ID).
		Key("name", s.Name()).
		Msgf("saving session with options %#v: %#v", s.Options, s.Values)
	if s.IsNew && s.Options.MaxAge > 0 {
		l.Trace().Msg("insert")
		err = m.insert(r.Context(), s)
	} else if !s.IsNew && s.Options.MaxAge > 0 {
		l.Trace().Msg("update")
		err = m.update(r.Context(), s)
	} else if !s.IsNew {
		l.Trace().Msg("delete")
		err = m.delete(r.Context(), s)
	} else {
		l.Trace().Msg("session is new and max age is 0, not saving")
	}
	if err == nil {
		var encoded string
		if encoded, err = securecookie.EncodeMulti(s.Name(), s.ID, m.codecs...); err != nil {
			return errors.Wrap(err, "failed to encode session")
		}
		http.SetCookie(
			w, &http.Cookie{
				Name:     s.Name(),
				Value:    encoded,
				Path:     s.Options.Path,
				Domain:   s.Options.Domain,
				MaxAge:   s.Options.MaxAge,
				Secure:   s.Options.Secure,
				HttpOnly: s.Options.HttpOnly,
			},
		)
	}
	return
}

func (m *MongoStore) insert(ctx context.Context, s *sessions.Session) (err error) {
	logger.Ctx(ctx).Debug().Key("session", s.ID).Key("name", s.Name()).Msgf(
		"inserting session %s: %#v", s.Name(), s.Values,
	)
	data := &sessionData{
		Data:       map[string]interface{}{},
		LastAccess: time.Now(),
		Version:    1,
	}
	for k, v := range s.Values {
		data.Data[k.(string)] = v
	}
	var res *mongo.InsertOneResult
	if res, err = m.db.InsertOne(ctx, data); err != nil {
		return errors.Wrap(err, "failed to insert session")
	}
	data.ID = res.InsertedID.(bson.ObjectID)
	s.ID = data.ID.Hex()
	t := map[string]interface{}{}
	if err = m.db.FindOne(ctx, bson.M{"_id": data.ID}).Decode(t); err != nil {
		return errors.Wrap(err, "failed to load session")
	}
	logger.Ctx(ctx).Debug().Key("session", s.ID).Key("name", s.Name()).Msgf("inserted session: %+v", t)
	return
}

func (m *MongoStore) update(ctx context.Context, s *sessions.Session) (err error) {
	data := GetMiddleWareStorage(ctx).GetSession(s.Name()).Data.(*sessionData)
	var res *mongo.UpdateResult
	var changed bool
	if changed, err = data.hasDataChangedAndUpdate(s.Values); err != nil {
		return errors.Wrap(err, "failed to check data changes")
	} else if !changed {
		logger.Ctx(ctx).Debug().Key("session", s.ID).Msg("no changes in session")
		now := time.Now()
		if res, err = m.db.UpdateOne(
			ctx, bson.M{"_id": data.ID}, bson.M{"$set": bson.M{"last_access": now}},
		); err != nil {
			return errors.Wrap(err, "failed to update session")
		}
		data.LastAccess = now
	} else {
		logger.Ctx(ctx).Debug().Key("session", s.ID).Msg("changes in session")
		oldVersion := data.Version
		data.Version++
		if res, err = m.db.ReplaceOne(ctx, bson.M{"_id": data.ID, "version": oldVersion}, data); err != nil {
			return errors.Wrap(err, "failed to update session")
		}
		if err = data.setSerialized(); err != nil {
			return errors.Wrap(err, "failed to set serialized data")
		}
	}
	if res.MatchedCount == 0 {
		return errors.New("session not found")
	}
	return
}

func (m *MongoStore) load(ctx context.Context, session *SessionData) (err error) {
	data := &sessionData{}
	var id bson.ObjectID
	if id, err = bson.ObjectIDFromHex(session.Session.ID); err != nil {
		return errors.Wrap(err, "failed to parse session id")
	}
	if err = m.db.FindOne(ctx, bson.M{"_id": id}).Decode(data); err == nil {
		session.Data = data
		if err = data.setSerialized(); err != nil {
			return errors.Wrap(err, "failed to set serialized data")
		}
		if err = data.setSessionValues(session.Session.Values); err != nil {
			return errors.Wrap(err, "failed to set session values")
		}
		return
	}
	return errors.Wrap(err, "failed to load session")
}

func (m *MongoStore) delete(ctx context.Context, s *sessions.Session) (err error) {
	data := GetMiddleWareStorage(ctx).GetSession(s.Name()).Data.(*sessionData)
	// not checking version here, as we are deleting the session
	if _, err = m.db.DeleteOne(ctx, bson.M{"_id": data.ID}); err != nil {
		return errors.Wrap(err, "failed to delete session")
	}
	return
}

var _ sessions.Store = (*MongoStore)(nil)
