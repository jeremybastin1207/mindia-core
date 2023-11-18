package task

import (
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/pipeline"
	"github.com/jeremybastin1207/mindia-core/internal/transform"
)

type TagMediaTask struct {
	fileStorage  media.FileStorer
	mediaStorage media.Storer
	tagger       transform.GoogleTagger
}

func NewTagMediaTask(fileStorage media.FileStorer, mediaStorage media.Storer) TagMediaTask {
	tagger := transform.NewGoogleTagger(transform.GoogleTaggerConfig{})

	return TagMediaTask{
		fileStorage,
		mediaStorage,
		tagger,
	}
}

func (t *TagMediaTask) Tag(path media.Path) (*media.Media, error) {
	media, err := t.mediaStorage.Get(path)
	if err != nil {
		return nil, err
	}

	downloadResult, err := t.fileStorage.Download(path)
	if err != nil {
		return nil, err
	}

	source := pipeline.NewSource(pipeline.SourceConfig{
		Getter: func(ctx pipeline.PipelineCtx) (pipeline.PipelineCtx, error) {
			ctx.Path = path
			ctx.Buffer = &pipeline.Buffer{
				Reader: downloadResult.Body,
			}
			return ctx, nil
		},
	})

	sinker := pipeline.NewSinker(pipeline.SinkerConfig{
		Sinker: func(ctx pipeline.PipelineCtx) error {
			media.Tags = ctx.Tags
			return t.mediaStorage.Save(media)
		},
	})

	p := pipeline.NewPipeline(pipeline.PipelineConfig{
		Source: &source,
		Sinker: &sinker,
		Steps:  []pipeline.Step{&t.tagger},
	})
	_, err = p.Execute()
	if err != nil {
		return nil, err
	}

	m, err := t.mediaStorage.Get(path)
	if err != nil {
		return nil, err
	}

	return m, nil
}
