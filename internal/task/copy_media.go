package task

import (
	"github.com/jeremybastin1207/mindia-core/internal/media"
)

type CopyMediaTask struct {
	fileStorage  media.FileStorer
	cacheStorage media.FileStorer
	mediaStorage media.Storer
}

func NewCopyMediaTask(
	fileStorage media.FileStorer,
	cacheStorage media.FileStorer,
	mediaStorage media.Storer,
) CopyMediaTask {
	return CopyMediaTask{
		fileStorage,
		cacheStorage,
		mediaStorage,
	}
}

func (t *CopyMediaTask) Copy(src media.Path, dst media.Path) (*media.Media, error) {
	return nil, nil
}
