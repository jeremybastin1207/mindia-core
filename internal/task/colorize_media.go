package task

import (
	"fmt"

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
	p, err := t.pluginManager.GetPlugin(plugin.ColorizePluginName)
	if err != nil {
		return nil, err
	}
	t2 := plugin.NewColorizeTask(path)
	t3, err := p.Execute(&t2)
	fmt.Println(t3)
	t.pluginManager.GetTaskStorage().EnqueueTask(t3)
	return nil, err
}
