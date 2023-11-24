package gqlresolver

func sliceToPointers[T any](v []T) []*T {
	if len(v) == 0 {
		return nil
	}
	result := make([]*T, len(v))
	for i := range v {
		result[i] = &v[i]
	}
	return result
}
