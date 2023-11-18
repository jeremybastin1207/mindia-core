package task

import (
	"github.com/jeremybastin1207/mindia-core/internal/analytics"
	"github.com/jeremybastin1207/mindia-core/internal/media"
)

type ClearCacheTask struct {
	cacheStorage      media.FileStorer
	analyticsRecorder analytics.AnalyticsRecorder
}

func NewClearCacheTask(
	cacheStorage media.FileStorer,
	analyticsRecorder analytics.AnalyticsRecorder,
) ClearCacheTask {
	return ClearCacheTask{
		cacheStorage,
		analyticsRecorder,
	}
}

func (c *ClearCacheTask) Clear(p media.Path) error {
	defer c.analyticsRecorder.RecordCacheClear()
	return c.cacheStorage.Delete(p)
}

func (c *ClearCacheTask) ClearAll() error {
	defer c.analyticsRecorder.RecordCacheClear()
	return c.cacheStorage.Delete(media.NewPath("/"))
}
