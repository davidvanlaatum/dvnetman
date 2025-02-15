package dto

import (
	"context"
	"dvnetman/pkg/logger"
	"dvnetman/pkg/mongo/modal"
	errors2 "errors"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type queueResolverFunc func(ctx context.Context, c *Converter, resolved *int) (err error)

var queueResolverFuncs []queueResolverFunc

func addQueueResolverFunc(f queueResolverFunc) {
	queueResolverFuncs = append(queueResolverFuncs, f)
}

type queueEntry[T any] struct {
	obj          *T
	callbacks    []func(*T)
	errCallbacks []func(modal.UUID, *T) error
	resolving    bool
	errors       []error
}

type lookupQueue[T any] map[modal.UUID]*queueEntry[T]

type Converter struct {
	db                *modal.DBClient
	deviceTypeQueue   lookupQueue[modal.DeviceType]
	deviceQueue       lookupQueue[modal.Device]
	manufacturerQueue lookupQueue[modal.Manufacturer]
}

func init() {
	addQueueResolverFunc(
		func(ctx context.Context, c *Converter, resolved *int) (err error) {
			return resolveQueue(ctx, c.deviceTypeQueue, c.db.ListDeviceTypes, resolved)
		},
	)
	addQueueResolverFunc(
		func(ctx context.Context, c *Converter, resolved *int) (err error) {
			return resolveQueue(ctx, c.deviceQueue, c.db.ListDevices, resolved)
		},
	)
	addQueueResolverFunc(
		func(ctx context.Context, c *Converter, resolved *int) (err error) {
			return resolveQueue(ctx, c.manufacturerQueue, c.db.ListManufacturers, resolved)
		},
	)
}

func NewConverter(db *modal.DBClient) *Converter {
	return &Converter{
		db: db,
	}
}

func (c *Converter) cloneUUID(id *uuid.UUID) *uuid.UUID {
	x := *id
	return &x
}

func addToQueue[T any](queue *lookupQueue[T], id *modal.UUID, f func(*T)) {
	if *queue == nil {
		*queue = make(lookupQueue[T])
	}
	entry, ok := (*queue)[*id]
	if !ok {
		entry = &queueEntry[T]{}
		(*queue)[*id] = entry
	}
	if entry.obj != nil {
		f(entry.obj)
		return
	}
	entry.callbacks = append(entry.callbacks, f)
}

func addToErrQueue[T any](queue *lookupQueue[T], id *modal.UUID, f func(modal.UUID, *T) error) {
	if *queue == nil {
		*queue = make(lookupQueue[T])
	}
	entry, ok := (*queue)[*id]
	if !ok {
		entry = &queueEntry[T]{}
		(*queue)[*id] = entry
	}
	if entry.obj != nil {
		if err := f(*id, entry.obj); err != nil {
			entry.errors = append(entry.errors, err)
		}
		return
	}
	entry.errCallbacks = append(entry.errCallbacks, f)
}

type baseInterface interface {
	GetBase() *modal.Base
}

func getId(obj interface{}) modal.UUID {
	return *obj.(baseInterface).GetBase().ID
}

type listFunc[T any] func(
	ctx context.Context, filter interface{}, opts ...options.Lister[options.FindOptions],
) (
	[]T, error,
)

func callErrorCallbacks[T any](q *queueEntry[T], id modal.UUID) {
	for _, cb := range q.errCallbacks {
		if err := cb(id, q.obj); err != nil {
			q.errors = append(q.errors, err)
		}
	}
}

func resolveQueue[T any](ctx context.Context, queue lookupQueue[T], f listFunc[*T], resolved *int) (err error) {
	for {
		var ids []*modal.UUID
		count := 0
		for id, v := range queue {
			if v.obj == nil && !v.resolving {
				ids = append(ids, &id)
				v.resolving = true
				count++
				if count >= 100 {
					break
				}
			}
		}
		if len(ids) == 0 {
			break
		}
		logger.Ctx(ctx).Trace().Key("references", len(ids)).Msg("resolving references")
		var d []*T
		if d, err = f(ctx, bson.M{"id": bson.M{"$in": ids}}); err != nil {
			return
		}
		for _, dt := range d {
			*resolved++
			id := getId(dt)
			entry := queue[id]
			entry.obj = dt
			for _, cb := range entry.callbacks {
				cb(dt)
			}
			entry.callbacks = nil
			callErrorCallbacks(entry, id)
			entry.resolving = false
		}
		for id, v := range queue {
			if v.resolving {
				callErrorCallbacks(v, id)
				if len(v.errors) == 0 {
					v.errors = append(
						v.errors, errors.Errorf("failed to find reference for %T: %s", v.obj, (uuid.UUID)(id).String()),
					)
				}
			}
		}
	}
	for _, q := range queue {
		if len(q.errors) > 0 {
			return errors2.Join(q.errors...)
		}
	}
	return
}

func (c *Converter) resolveQueue(ctx context.Context) (err error) {
	for {
		var resolved int
		for _, resolverFunc := range queueResolverFuncs {
			if err = resolverFunc(ctx, c, &resolved); err != nil {
				return
			}
		}
		logger.Ctx(ctx).Trace().Key("references", resolved).Msg("resolved references")
		if resolved == 0 {
			break
		}
	}
	return
}
