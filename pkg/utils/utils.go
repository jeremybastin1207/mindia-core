package utils

import (
	"fmt"
	"os"
)

func ToArray[T any](mp map[string]T) []T {
	arr := []T{}
	for _, p := range mp {
		arr = append(arr, p)
	}
	return arr
}

func ExitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
