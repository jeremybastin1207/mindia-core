package task

import (
	"github.com/jeremybastin1207/mindia-core/internal/analytics"
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/rs/zerolog/log"
)

type DeleteMediaTask struct {
	fileStorage       media.FileStorer
	cacheStorage      media.FileStorer
	mediaStorage      media.Storer
	analyticsRecorder analytics.AnalyticsRecorder
}

func NewDeleteMediaTask(
	fileStorage media.FileStorer,
	cacheStorage media.FileStorer,
	mediaStorage media.Storer,
	analyticsRecorder analytics.AnalyticsRecorder,
) DeleteMediaTask {
	return DeleteMediaTask{
		fileStorage,
		cacheStorage,
		mediaStorage,
		analyticsRecorder,
	}
}

func (t *DeleteMediaTask) Delete(path media.Path) error {
	defer t.analyticsRecorder.RecordMediaDelete()

	media, err := t.mediaStorage.Get(path)
	if err != nil {
		return err
	}
	err = t.fileStorage.Delete(path)
	if err != nil {
		log.Warn().Err(err)
	}
	for _, dm := range media.DerivedMedias {
		err := t.cacheStorage.Delete(dm.Path)
		if err != nil {
			log.Warn().Err(err)
		}
	}
	return t.mediaStorage.Delete(media.Path)
}

func (t *DeleteMediaTask) DeleteMultiple(paths []media.Path) error {
	for _, path := range paths {
		err := t.Delete(path)
		if err != nil {
			return err
		}
	}
	return nil
}
