package utils

// ToPtr is a wrapper that allows converting a literal to a pointer or create a pointer to a shallow copy of a value.
func ToPtr[T any](v T) *T {
	return &v
}

func FromPtr[T any](v *T) (r T) {
	if v != nil {
		r = *v
	}
	return
}

// ConvertPtr converts a pointer to a value using the provided function only if it is non-nil.
func ConvertPtr[T, R any](in *T, f func(T) R) *R {
	if in == nil {
		return nil
	}
	return ToPtr(f(*in))
}
