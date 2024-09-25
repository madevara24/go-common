package array

// Reverse array of string
func Reverse[T comparable](arr []T) []T {
	for i := 0; i < len(arr)/2; i++ {
		j := len(arr) - i - 1
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func AppendUnique[T comparable](dest []T, src ...T) []T {
	exists := make(map[T]bool)
	// Mark existing elements from dest in the map
	for _, v := range dest {
		exists[v] = true
	}

	// Only append elements from src that don't already exist
	for _, v := range src {
		if !exists[v] {
			dest = append(dest, v)
			exists[v] = true
		}
	}

	return dest
}
