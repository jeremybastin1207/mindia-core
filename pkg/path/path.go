package path

import "path/filepath"

func JoinPath(elems ...string) string {
	return filepath.ToSlash(filepath.Join(elems...))
}
