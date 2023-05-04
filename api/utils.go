package api

import (
	"fmt"
	"strings"
)

func PathToID(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func search[T any](reader func(T) string, needle string, array []T) (*T, error) {
	for cIdx := 0; cIdx < len(array); cIdx++ {
		if reader(array[cIdx]) == needle {
			return &array[cIdx], nil
		}
	}
	return nil, fmt.Errorf("could not find object '%s' in collection", needle)
}
