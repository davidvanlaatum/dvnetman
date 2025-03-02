package cache

import (
	"bytes"
	"context"
	"dvnetman/pkg/logger"
	"encoding/gob"
	"fmt"
	"github.com/inhies/go-bytesize"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
	"net/url"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

var OTelStatusKey = attribute.Key("cache.status")
var OTelValueTypeKey = attribute.Key("cache.value_type")

type memoryCacheEntry struct {
	value  []byte
	expiry time.Time
}

type MemoryCache struct {
	cache  map[string]memoryCacheEntry
	mu     sync.RWMutex
	limit  bytesize.ByteSize
	size   bytesize.ByteSize
	ttl    time.Duration
	c      chan struct{}
	tracer trace.Tracer
}

func (m *MemoryCache) startSpan(ctx context.Context, attr ...attribute.KeyValue) (context.Context, trace.Span) {
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	s := strings.Split(f.Name(), ".")
	op := strings.ToUpper(s[len(s)-1])
	attr = append(attr, semconv.DBSystemKey.String("memory"), semconv.DBOperationName(op))
	return m.tracer.Start(ctx, "cache."+op, trace.WithAttributes(attr...), trace.WithSpanKind(trace.SpanKindInternal))
}

func (m *MemoryCache) Flush(ctx context.Context) error {
	_, s := m.startSpan(ctx)
	defer s.End()
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache = map[string]memoryCacheEntry{}
	return nil
}

func (m *MemoryCache) New(ctx context.Context) (Cache, error) {
	return m, nil
}

func (m *MemoryCache) Shutdown(ctx context.Context) error {
	close(m.c)
	return nil
}

func (m *MemoryCache) Set(ctx context.Context, key string, value interface{}) (err error) {
	_, s := m.startSpan(ctx, semconv.DBCollectionName(key), OTelValueTypeKey.String(fmt.Sprintf("%T", value)))
	defer s.End()
	m.mu.Lock()
	defer m.mu.Unlock()
	b := &bytes.Buffer{}
	e := gob.NewEncoder(b)
	if err = e.Encode(value); err != nil {
		return
	}
	m.cache[key] = memoryCacheEntry{value: b.Bytes(), expiry: time.Now().Add(m.ttl)}
	m.size += bytesize.ByteSize(b.Len())
	logger.Trace(ctx).Key("key", key).Msg("cache set")
	if m.size > m.limit {
		select {
		case m.c <- struct{}{}:
		default:
		}
	}
	return
}

func (m *MemoryCache) Get(ctx context.Context, key string, value interface{}) (ok bool, err error) {
	ctx, s := m.startSpan(ctx, semconv.DBCollectionName(key), OTelValueTypeKey.String(fmt.Sprintf("%T", value)))
	defer s.End()
	l := logger.Trace(ctx).Key("key", key)
	defer l.Msg("cache get")
	m.mu.RLock()
	defer m.mu.RUnlock()
	var e memoryCacheEntry
	if e, ok = m.cache[key]; ok && time.Now().Before(e.expiry) {
		l.Key("cache", "hit")
		s.SetAttributes(OTelStatusKey.String("hit"))
		d := gob.NewDecoder(bytes.NewBuffer(e.value))
		if err = d.Decode(value); err != nil {
			s.RecordError(err)
			s.SetStatus(codes.Error, err.Error())
			return false, errors.Wrap(err, "decode")
		}
		return true, nil
	} else if ok {
		l.Key("cache", "expired")
		s.SetAttributes(OTelStatusKey.String("expired"))
		go func() {
			_ = m.Delete(ctx, key)
		}()
		ok = false
	} else {
		l.Key("cache", "miss")
		s.SetAttributes(OTelStatusKey.String("miss"))
	}
	return
}

func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	_, s := m.startSpan(ctx, semconv.DBCollectionName(key))
	defer s.End()
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.cache, key)
	logger.Trace(ctx).Key("key", key).Msg("cache delete")
	return nil
}

func (m *MemoryCache) Close(ctx context.Context) error {
	return nil
}

func (m *MemoryCache) expire() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.cache {
		if time.Now().After(v.expiry) {
			delete(m.cache, k)
			m.size -= bytesize.ByteSize(len(v.value))
		}
	}
	if m.limit > 0 && m.size > m.limit {
		type t struct {
			expiry time.Time
			key    string
		}
		var s []t
		for k, v := range m.cache {
			s = append(s, t{expiry: v.expiry, key: k})
		}
		sort.Slice(
			s, func(i, j int) bool {
				return s[i].expiry.Before(s[j].expiry)
			},
		)
		for m.size > m.limit {
			k := s[0].key
			v := m.cache[k]
			delete(m.cache, k)
			m.size -= bytesize.ByteSize(len(v.value))
			s = s[1:]
		}
	}
}

func (m *MemoryCache) cleanup(ctx context.Context) {
	t := time.NewTicker(m.ttl)
	defer func() {
		_ = m.Flush(ctx)
	}()
	defer t.Stop()
	for {
		select {
		case <-t.C:
			m.expire()
		case _, ok := <-m.c:
			if !ok {
				return
			}
			m.expire()
		case <-ctx.Done():
			return
		}
	}
}

var _ Pool = (*MemoryCache)(nil)

type MemoryDriver struct {
}

func (m *MemoryDriver) New(ctx context.Context, config *url.URL) (_ Cache, err error) {
	traceProvider := otel.GetTracerProvider()
	c := &MemoryCache{
		cache:  map[string]memoryCacheEntry{},
		c:      make(chan struct{}, 1),
		tracer: traceProvider.Tracer("cache"),
	}
	if c.ttl, err = time.ParseDuration(config.Query().Get("ttl")); err != nil {
		return nil, errors.Wrap(err, "parse ttl")
	}
	if x := config.Query().Get("limit"); x != "" {
		if c.limit, err = bytesize.Parse(x); err != nil {
			return nil, errors.Wrap(err, "parse limit")
		}
	}
	go c.cleanup(ctx)
	return c, nil
}

func (m *MemoryDriver) Name() string {
	return "memory"
}

func init() {
	Register(&MemoryDriver{})
}
