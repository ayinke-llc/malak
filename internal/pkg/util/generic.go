package util

func Ref[T any](value T) *T {
	return &value
}

func DeRef[T any](ptr *T) T {
	var zeroValue T

	if ptr == nil {
		return zeroValue
	}

	return *ptr
}
