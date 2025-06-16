package utils

// BatchSlice splits a slice 'arr' into batches of size 'batchSize'.
// The last batch may be smaller if arr's length is not a multiple of batchSize.
func BatchSlice[T any](arr []T, batchSize int) [][]T {

	if batchSize <= 0 || len(arr) == 0 {
		return [][]T{}
	}

	if batchSize >= len(arr) {
		return [][]T{arr}
	}
	var batches [][]T
	for batchSize < len(arr) {
		batches = append(batches, arr[:batchSize])
		arr = arr[batchSize:]
	}
	if len(batches) > 0 {
		batches = append(batches, arr)
	}

	return batches
}
