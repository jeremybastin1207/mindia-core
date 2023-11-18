package analytics

type AnalyticsRecorder interface {
	RecordBandwithUsage(bytesLength int)
	RecordMediaRequest()
	RecordMediaDelete()
	RecordDataStorageUsage(bytesUsage int64)
	RecordCacheStorageUsage(bytesUsage int64)
	RecordTaskCreation()
	RecordCacheClear()
}
