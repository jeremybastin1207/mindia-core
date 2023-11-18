package cmd

import (
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/task"
)

type GetMediaCommand struct {
	getMediaTask task.GetMediaTask
}

func NewGetMediaCommand(getMediaTask task.GetMediaTask) GetMediaCommand {
	return GetMediaCommand{
		getMediaTask,
	}
}

func (c *GetMediaCommand) Get(path media.Path) (*media.Media, error) {
	return c.getMediaTask.Get(path)
}
