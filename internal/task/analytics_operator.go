package task

import (
	"github.com/jeremybastin1207/mindia-core/internal/media"
)

type AnalyticsOperator struct {
	fileStorage  media.FileStorer
	cacheStorage media.FileStorer
}

func NewAnalyticsOperator(
	fileStorage media.FileStorer,
	cacheStorage media.FileStorer,
) AnalyticsOperator {
	return AnalyticsOperator{
		fileStorage,
		cacheStorage,
	}
}

type SpaceUsageResponse struct {
	TotalSpaceUsage int64 `json:"total_space_usage"`
	FileSpaceUsage  int64 `json:"file_space_usage"`
	CacheSpaceSuage int64 `json:"cache_space_usage"`
}

func (o *AnalyticsOperator) SpaceUsage() (*SpaceUsageResponse, error) {
	fileSpaceUsage, err := o.fileStorage.SpaceUsage()
	if err != nil {
		return nil, err
	}
	cacheSpaceUsage, err := o.cacheStorage.SpaceUsage()
	if err != nil {
		return nil, err
	}
	return &SpaceUsageResponse{
		TotalSpaceUsage: fileSpaceUsage + cacheSpaceUsage,
		FileSpaceUsage:  fileSpaceUsage,
		CacheSpaceSuage: cacheSpaceUsage,
	}, nil
}
