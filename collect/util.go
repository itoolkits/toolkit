// collect util

package collect

// ChunkArrays chuck array
func ChunkArrays[T any](slice []T, size int) [][]T {
	if len(slice) == 0 || size <= 0 {
		return nil
	}

	expectedLen := (len(slice) + size - 1) / size

	result := make([][]T, 0, expectedLen)

	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		result = append(result, slice[i:end])
	}
	return result
}
