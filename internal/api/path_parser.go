package api

import (
	"fmt"
	"strings"
)

func parsePath(path string) (string, string) {
	var (
		transformations string
		image_path      string
	)
	for _, part := range strings.Split(path, "/") {
		if strings.Contains(part, "c_") || strings.Contains(part, "t_") {
			transformations = fmt.Sprintf("%s/%s", transformations, part)
		} else {
			image_path = fmt.Sprintf("%s/%s", image_path, part)
		}
	}
	transformations = strings.TrimPrefix(transformations, "/")
	return transformations, image_path
}
