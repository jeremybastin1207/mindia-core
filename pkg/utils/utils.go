package utils

func ToArray[T any](mp map[string]T) []T {
	arr := []T{}
	for _, p := range mp {
		arr = append(arr, p)
	}
	return arr
}
