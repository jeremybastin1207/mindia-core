package task

import (
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/rs/zerolog/log"
)

type MoveMediaTask struct {
	fileStorage  media.FileStorer
	cacheStorage media.FileStorer
	mediaStorage media.Storer
}

func NewMoveMediaTask(
	fileStorage media.FileStorer,
	cacheStorage media.FileStorer, mediaStorage media.Storer,
) MoveMediaTask {
	return MoveMediaTask{
		fileStorage,
		cacheStorage,
		mediaStorage,
	}
}

func (t *MoveMediaTask) Move(src media.Path, dst media.Path) (*media.Media, error) {
	m, err := t.mediaStorage.Get(src)
	if err != nil {
		return nil, err
	}
	err = t.fileStorage.Move(src, dst)
	if err != nil {
		return nil, err
	}
	oldPath := m.Path
	m.Path = m.Path.WithDir(dst.ToString())

	movedAssets := []media.DerivedMedia{}
	for _, asset := range m.DerivedMedias {
		dst := asset.Path.WithDir(dst.ToString())
		err = t.cacheStorage.Move(asset.Path, dst)
		if err != nil {
			log.Warn().Err(err)
			continue
		}
		asset.Path = dst
		movedAssets = append(movedAssets, asset)
	}
	m.DerivedMedias = movedAssets

	err = t.mediaStorage.Delete(oldPath)
	if err != nil {
		return nil, err
	}
	err = t.mediaStorage.Save(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
