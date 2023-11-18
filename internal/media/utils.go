package media

import (
	"fmt"

	"github.com/jeremybastin1207/mindia-core/pkg/path"
)

func BuildDerivedPath(m Media) Path {
	p := m.Path
	return NewPath(
		path.JoinPath(p.Dir(), fmt.Sprintf("%s-%d%s", p.Uuid(), len(m.DerivedMedias)+1, p.Extension())),
	)
}
