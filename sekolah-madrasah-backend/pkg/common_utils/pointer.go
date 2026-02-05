package common_utils

func ToPointer[T any](v T) *T {
	return &v
}
