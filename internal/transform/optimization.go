package transform

import (
	"errors"

	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/pipeline"
)

type MediaOptimization struct {
}

func NewMediaOptimization() *MediaOptimization {
	return &MediaOptimization{}
}

func (o *MediaOptimization) GetSteps(contentType media.ContentType) ([]pipeline.PipelineStep, error) {
	switch contentType {
	case media.ImageJpeg, media.ImagePng:
		exifReader := NewExifReader()
		webpConverter := NewWebpConverter()
		return []pipeline.PipelineStep{
			&exifReader,
			&webpConverter,
		}, nil
	case media.VideoMp4:
		return []pipeline.PipelineStep{}, nil
	default:
		return []pipeline.PipelineStep{}, errors.New("content-type not supported")
	}
}
