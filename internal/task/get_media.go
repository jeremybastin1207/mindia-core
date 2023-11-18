package task

import (
	"github.com/jeremybastin1207/mindia-core/internal/analytics"
	"github.com/jeremybastin1207/mindia-core/internal/media"
)

type GetMediaTask struct {
	mediaStorage      media.Storer
	analyticsRecorder analytics.AnalyticsRecorder
}

func NewGetMediaTask(
	mediaStorage media.Storer,
	analyticsRecorder analytics.AnalyticsRecorder,
) GetMediaTask {
	return GetMediaTask{
		mediaStorage,
		analyticsRecorder,
	}
}

func (t *GetMediaTask) Get(path media.Path) (*media.Media, error) {
	t.analyticsRecorder.RecordMediaRequest()

	return t.mediaStorage.Get(path)
}

func (t *GetMediaTask) GetMultiple(
	path media.Path,
	offset int,
	limit int,
	sortBy string,
	asc bool,
) ([]media.Media, error) {
	return t.mediaStorage.GetMultiple(path, offset, limit, sortBy, asc)
}
