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

func (o *MediaOptimization) GetSteps(contentType media.ContentType) ([]pipeline.Step, error) {
	switch contentType {
	case media.ImageJpeg, media.ImagePng:
		exifReader := NewExifReader()
		webpConverter := NewWebpConverter()
		return []pipeline.Step{
			&exifReader,
			&webpConverter,
		}, nil
	case media.VideoMp4:
		return []pipeline.Step{}, nil
	default:
		return []pipeline.Step{}, errors.New("content-type not supported")
	}
}
