package task

import (
	"time"

	"github.com/jeremybastin1207/mindia-core/internal/analytics"
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/rs/zerolog/log"
)

type StorageUsageCollector struct {
	fileStorage       media.FileStorer
	cacheStorage      media.FileStorer
	ticker            *time.Ticker
	done              chan bool
	analyticsRecorder analytics.AnalyticsRecorder
}

func NewStorageUsageCollector(
	fileStorage media.FileStorer,
	cacheStorage media.FileStorer,
	analyticsRecorder analytics.AnalyticsRecorder,
) StorageUsageCollector {
	s := StorageUsageCollector{
		fileStorage:       fileStorage,
		cacheStorage:      cacheStorage,
		ticker:            time.NewTicker(5 * time.Second),
		done:              make(chan bool),
		analyticsRecorder: analyticsRecorder,
	}
	go s.Collect()
	return s
}

func (c *StorageUsageCollector) Stop() {
	c.done <- true
}

func (c *StorageUsageCollector) Collect() {
	c.collect()

	for {
		select {
		case <-c.done:
			return
		case <-c.ticker.C:
			c.collect()
		}
	}
}

func (c *StorageUsageCollector) collect() {
	dataStorageUsage, err := c.fileStorage.SpaceUsage()
	if err != nil {
		log.Err(err)
	} else {
		c.analyticsRecorder.RecordDataStorageUsage(dataStorageUsage)
	}

	cacheStorageUsage, err := c.cacheStorage.SpaceUsage()
	if err != nil {
		log.Err(err)
	} else {
		c.analyticsRecorder.RecordCacheStorageUsage(cacheStorageUsage)
	}
}
