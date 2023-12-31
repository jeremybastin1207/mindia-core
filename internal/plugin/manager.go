package plugin

import (
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/scheduler"
	"github.com/jeremybastin1207/mindia-core/internal/transform"
)

type Plugin interface {
	Name() string
	Execute(task *scheduler.Task) (*scheduler.Task, error)
}

type PluginManager struct {
	fileStorage       media.FileStorer
	cacheStorage      media.FileStorer
	mediaStorage      media.Storer
	taskStorage       scheduler.Storer
	mediaOptimization *transform.MediaOptimization
	plugins           map[string]Plugin
}

func NewPluginManager(
	fileStorage media.FileStorer,
	cacheStorage media.FileStorer,
	mediaStorage media.Storer,
	taskStorage scheduler.Storer,
) PluginManager {
	return PluginManager{
		fileStorage,
		cacheStorage,
		mediaStorage,
		taskStorage,
		transform.NewMediaOptimization(),
		map[string]Plugin{},
	}
}

func (p *PluginManager) RegisterPlugin(pl Plugin) {
	p.plugins[pl.Name()] = pl
}

func (p *PluginManager) GetPlugin(name string) (Plugin, error) {
	return p.plugins[name], nil
}

func (p *PluginManager) GetFileStorage() media.FileStorer {
	return p.fileStorage
}

func (p *PluginManager) GetCacheStorage() media.FileStorer {
	return p.cacheStorage
}

func (p *PluginManager) GetMediaStorage() media.Storer {
	return p.mediaStorage
}

func (p *PluginManager) GetTaskStorage() scheduler.Storer {
	return p.taskStorage
}

func (p *PluginManager) GetMediaOptimization() *transform.MediaOptimization {
	return p.mediaOptimization
}
