package graphql

func convert[T, V any](slice []T, convertor func(T) V) []V {
	result := make([]V, len(slice))
	for i, v := range slice {
		result[i] = convertor(v)
	}
	return result
}
