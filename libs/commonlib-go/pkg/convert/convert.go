package convert

// ToValue converts *T to T (returns zero value if nil)
func ToValue[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// ToPtr converts T to *T
func ToPtr[T any](v T) *T {
	return &v
}

type ConvertibleInt interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

// ConvertInt converts T to U
func ConvertInt[T ConvertibleInt, U ConvertibleInt](v T) U {
	return U(v)
}
