package task

import (
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/plugin"
)

type ColorizeMediaTask struct {
	pluginManager *plugin.PluginManager
}

func NewColorizeMediaTask(pluginManager *plugin.PluginManager) ColorizeMediaTask {
	return ColorizeMediaTask{
		pluginManager,
	}
}

func (t *ColorizeMediaTask) Colorize(path media.Path) (*media.Media, error) {
	plugin, err := t.pluginManager.GetPlugin("colorize")
	if err != nil {
		return nil, err
	}
	plugin.Execute(path)
	return nil, nil
}
