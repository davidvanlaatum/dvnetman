package utils

import (
	"sort"
)

func MapTo[T any, R any](in []T, f func(T) R) []R {
	if in == nil {
		return nil
	}
	out := make([]R, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}

func MapErr[T any, R any](in []T, f func(T) (R, error)) ([]R, error) {
	out := make([]R, len(in))
	for i, v := range in {
		var err error
		if out[i], err = f(v); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func MapMapTo[T any, K comparable, R any](in map[K]T, f func(key K, value T) R) map[K]R {
	if in == nil {
		return nil
	}
	out := make(map[K]R, len(in))
	for k, v := range in {
		out[k] = f(k, v)
	}
	return out
}

type SortedMapEntries[K comparable, V any] struct {
	Key   K
	Value V
}

func MapSortedByKey[K comparable, V any](
	in map[K]V, less func(a, b SortedMapEntries[K, V]) bool,
) []SortedMapEntries[K, V] {
	out := make([]SortedMapEntries[K, V], 0, len(in))
	for k, v := range in {
		out = append(out, SortedMapEntries[K, V]{Key: k, Value: v})
	}
	sort.Slice(
		out, func(i, j int) bool {
			return less(out[i], out[j])
		},
	)
	return out
}

func MapSortedByKeyString[V any](a, b SortedMapEntries[string, V]) bool {
	return a.Key < b.Key
}

func SortedMapKeys[K comparable, V any](in map[K]V, less func(a, b K) bool) []K {
	keys := make([]K, 0, len(in))
	for k := range in {
		keys = append(keys, k)
	}
	sort.Slice(
		keys, func(i, j int) bool {
			return less(keys[i], keys[j])
		},
	)
	return keys
}
