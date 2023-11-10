package collection

// Map applies the given function to each element of slice a, returning slice of
// results
func Map[T any, M any](a []T, fn func(T) M) []M {
	anySlice := make([]M, 0)
	for _, elem := range a {
		anySlice = append(anySlice, fn(elem))
	}

	return anySlice
}
