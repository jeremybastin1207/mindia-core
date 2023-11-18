package media

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jeremybastin1207/mindia-core/pkg/path"
)

const uuidRegex = "[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}"

type Path struct {
	Path string `json:"path"`
}

func NewPath(path string) Path {
	if path == "" || path[0] != '/' {
		panic("path can only be absolute")
	}
	return Path{Path: path}
}

func (p *Path) ToString() string {
	return p.Path
}

func (p *Path) Dir() string {
	if p.Extension() == "" {
		return p.Path
	}
	return strings.Replace(filepath.Dir(p.Path), "\\", "/", -1)
}

func (p *Path) Basename() string {
	base := filepath.Base(p.Path)
	return base[:len(base)-len(filepath.Ext(base))]
}

func (p *Path) BasePath() Path {
	base := filepath.Base(p.Path)
	return NewPath(
		p.Dir() + "/" + base[:len(base)-len(filepath.Ext(base))],
	)
}

func (p *Path) Uuid() string {
	uuidRegex := regexp.MustCompile(uuidRegex)
	return uuidRegex.FindString(p.Basename())
}

func (p *Path) Extension() string {
	return filepath.Ext(p.Path)
}

func (p *Path) Filename() string {
	return filepath.Base(p.Path)
}

func (p *Path) AppendSuffix(suffix string) Path {
	return NewPath(
		path.JoinPath(p.Dir(), p.Uuid()+strings.ReplaceAll(strings.ReplaceAll(suffix, "/", "-"), ".webp", "")+p.Extension()),
	)
}

func (p *Path) SetExtension(ext string) Path {
	return NewPath(
		path.JoinPath(p.Dir(), p.Uuid()+ext),
	)
}

func (p *Path) WithDir(dir string) Path {
	return NewPath(
		path.JoinPath(dir, p.Filename()),
	)
}
